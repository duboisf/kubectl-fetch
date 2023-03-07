package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
	"sync"
	"time"
)

// Fetcher is an interface for cmd.Plugin
type Fetcher interface {
	Fetch(ctx context.Context) ([]string, error)
}

// Starter is an interface for terminal.UI
type Starter interface {
	Start(ctx context.Context, wg *sync.WaitGroup)
}

type Stdout interface {
	io.Writer
	Stat() (fs.FileInfo, error)
}

type Cmd struct {
	plugin        Fetcher
	stderr        io.Writer
	stdout        Stdout
	ui            Starter
	UIStopTimeout time.Duration
}

func NewCmd(plugin Fetcher, stdout Stdout, stderr io.Writer, ui Starter) (*Cmd, error) {
	return &Cmd{
		plugin:        plugin,
		stderr:        stderr,
		stdout:        stdout,
		ui:            ui,
		UIStopTimeout: 500 * time.Millisecond,
	}, nil
}

func (c *Cmd) Run(ctx context.Context) error {
	wg := &sync.WaitGroup{}
	fileInfo, err := c.stdout.Stat()
	if err != nil {
		return err
	}
	uiCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	// Only start UI if we are connected to a TTY
	if fileInfo.Mode()&os.ModeCharDevice != 0 {
		wg.Add(1)
		go c.ui.Start(uiCtx, wg)
	}
	resources, err := c.plugin.Fetch(ctx)
	cancel()
	c.waitForUI(wg)
	if err != nil {
		return err
	}
	if len(resources) == 0 {
		fmt.Fprintln(c.stderr, "No resources found.")
		return nil
	}
	bufferedStdout := bufio.NewWriter(c.stdout)
	bufferedStdout.WriteString(strings.Join(resources, "\n") + "\n")
	return bufferedStdout.Flush()
}

func (c *Cmd) waitForUI(wg *sync.WaitGroup) {
	stopped := make(chan struct{})
	go func() {
		defer close(stopped)
		wg.Wait()
	}()
	<-stopped
}
