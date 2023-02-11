package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/duboisf/kubectl-getall/internal/cmd"
)

func main() {
	kubectlPlugin := cmd.New()
	resources, err := kubectlPlugin.GetAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	if len(resources) == 0 {
		fmt.Fprintf(os.Stderr, "No resources found!\n")
	} else {
		fmt.Fprintln(os.Stdout, strings.Join(resources, "\n"))
	}
}
