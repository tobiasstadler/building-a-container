package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command(os.Args[1], os.Args[2:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error running the %s command: %s\n", os.Args[1], err)
		os.Exit(1)
	}
}
