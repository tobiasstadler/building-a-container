package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

func main() {
	if os.Args[1] == "run" {
		cmd := exec.Command(os.Args[0], append([]string{"doRun"}, os.Args[2:]...)...)

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		cmd.SysProcAttr = &syscall.SysProcAttr{
			Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
			Unshareflags: syscall.CLONE_NEWNS,
		}

		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error running the %s command: %s\n", os.Args[0], err)
			os.Exit(1)
		}
	} else if os.Args[1] == "doRun" {
		initCmd()

		err := syscall.Exec(os.Args[2], os.Args[3:], os.Environ())
		if err != nil {
			fmt.Printf("Error running the %s command: %s\n", os.Args[2], err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func initCmd() {
	rand.Seed(time.Now().UnixNano())

	containerName := "container" + strconv.Itoa(rand.Int())

	err := syscall.Sethostname([]byte(containerName))
	if err != nil {
		fmt.Printf("Error setting hostname to %s: %s\n", containerName, err)
		os.Exit(1)
	}

	err = os.MkdirAll(containerName+"/upper", 0777)
	if err != nil {
		fmt.Printf("Error creating directory %s: %s\n", containerName+"/upper", err)
		os.Exit(1)
	}

	err = os.MkdirAll(containerName+"/merged", 0777)
	if err != nil {
		fmt.Printf("Error creating directory %s: %s\n", containerName+"/merged", err)
		os.Exit(1)
	}

	err = os.MkdirAll(containerName+"/work", 0777)
	if err != nil {
		fmt.Printf("Error creating directory %s: %s\n", containerName+"/work", err)
		os.Exit(1)
	}

	err = syscall.Mount("none", containerName+"/merged", "overlay", 0, "lowerdir=ubuntu,upperdir="+containerName+"/upper,workdir="+containerName+"/work")
	if err != nil {
		fmt.Printf("Error mounting %s to %s: %s\n", "ubuntu", containerName+"/merged", err)
		os.Exit(1)
	}

	err = os.MkdirAll(containerName+"/merged/.oldroot", 0777)
	if err != nil {
		fmt.Printf("Error creating directory %s: %s\n", containerName+"/merged/.oldroot", err)
		os.Exit(1)
	}

	err = syscall.PivotRoot(containerName+"/merged", containerName+"/merged/.oldroot")
	if err != nil {
		fmt.Printf("Error changing root to %s: %s\n", containerName+"/merged", err)
		os.Exit(1)
	}

	err = os.Chdir("/")
	if err != nil {
		fmt.Printf("Error changing directory to %s: %s\n", "/", err)
		os.Exit(1)
	}

	err = syscall.Unmount("/.oldroot", syscall.MNT_DETACH)
	if err != nil {
		fmt.Printf("Error unmounting %s: %s\n", "/.oldroot", err)
		os.Exit(1)
	}

	err = os.RemoveAll("/.oldroot")
	if err != nil {
		fmt.Printf("Error deleting directory %s: %s\n", "/.oldroot", err)
		os.Exit(1)
	}

	err = syscall.Mount("none", "/proc", "proc", 0, "")
	if err != nil {
		fmt.Printf("Error mounting %s to %s: %s\n", "proc", "/proc", err)
		os.Exit(1)
	}
}
