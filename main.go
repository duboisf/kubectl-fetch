package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/duboisf/kubectl-fetch/internal/cmd"
	"github.com/duboisf/kubectl-fetch/internal/pkg/kubectl"
	"github.com/duboisf/kubectl-fetch/internal/pkg/terminal"
)

func main() {
	err := Main()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func Main() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	tput := terminal.NewTPut(exec.Command)

	foregroundColor, _ := tput.Query("setaf 4")
	backgroundColor, _ := tput.Query("setab 0")
	resetColor, _ := tput.Query("sgr0")
	progressBar := terminal.NewProgressBar(foregroundColor, backgroundColor, resetColor)
	spinner := terminal.NewSpinner(100*time.Millisecond)
	tui := terminal.NewUI(progressBar, spinner, tput, os.Stderr)
	kubeClient := kubectl.New(exec.CommandContext)
	opts, err := cmd.GetOptions(os.Args[1:])
	if err != nil {
		return err
	}
	plugin, err := cmd.NewPlugin(kubeClient, opts, tui)
	if err != nil {
		return err
	}
	cmd, err := cmd.NewCmd(plugin, os.Stdout, os.Stderr, tui)
	if err != nil {
		return err
	}
	return cmd.Run(ctx)
}
