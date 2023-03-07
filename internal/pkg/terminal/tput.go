package terminal

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Cmd is an interface for exec.Cmd
type Cmd interface {
	Output() ([]byte, error)
	StdinPipe() (io.WriteCloser, error)
}

type NewCommand[C Cmd] func(name string, args ...string) C

type TPut[C Cmd] struct {
	newCommand NewCommand[C]
}

func NewTPut[C Cmd](newCommand NewCommand[C]) *TPut[C] {
	return &TPut[C]{newCommand: newCommand}
}

func (t *TPut[C]) Query(capnames ...string) (string, error) {
	cmd := t.newCommand("tput", "-S")
	rawStdin, err := cmd.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("could get stdin pipe to tput proces: %w", err)
	}
	err = writeCapnamesToStdin(rawStdin, capnames)
	if err != nil {
		return "", fmt.Errorf("could not write capabilities to tput's stdin: %w", err)
	}
	stdout, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("tput returned an error: %w", err)
	}
	return strings.TrimSpace(string(stdout)), nil
}

func (t *TPut[C]) QueryInt(capname string) (int, error) {
	output, err := t.Query(capname)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(output)
}

func writeCapnamesToStdin(w io.WriteCloser, capnames []string) error {
	defer w.Close()
	stdin := bufio.NewWriter(w)
	defer stdin.Flush()
	for _, capname := range capnames {
		_, err := stdin.WriteString(capname + "\n")
		if err != nil {
			return err
		}
	}
	return stdin.Flush()
}
