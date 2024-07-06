package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stdout, "$ ")
		cmdLine, _ := reader.ReadString('\n')
		cmd, args := getCmdAndArgs(cmdLine)
		cmd = strings.TrimSpace(cmd)
		switch cmd {
		case "type":
			if len(args) == 0 {
				fmt.Println("type: missing argument")
				continue
			}
			switch args[0] {
			case "exit", "echo", "type", "pwd", "cd":
				fmt.Printf("%s is a shell builtin\n", args[0])
			default:
				path, found := findExecutable(args[0])
				if found {
					fmt.Printf("%s is %s\n", args[0], path)
				} else {
					fmt.Printf("%s: not found\n", args[0])
				}
			}
		case "exit":
			os.Exit(0)
		case "echo":
			fmt.Printf("%s\n", strings.Join(args, " "))
		case "pwd":
			pwd()
		case "cd":
			if len(args) == 0 {
				fmt.Println("cd: missing argument")
			} else {
				cd(args[0])
			}
		default:
			path, found := findExecutable(cmd)
			if found {
				runExternalCommand(path, args)
			} else {
				fmt.Printf("%s: not found\n", cmd)
			}
		}
	}
}

func getCmdAndArgs(cmd string) (string, []string) {
	l := strings.Fields(cmd)
	if len(l) == 0 {
		return "", []string{}
	}
	return l[0], l[1:]
}

func findExecutable(name string) (string, bool) {
	pathEnv := os.Getenv("PATH")
	paths := strings.Split(pathEnv, ":")

	for _, dir := range paths {
		fullPath := filepath.Join(dir, name)
		if fileExistsAndExecutable(fullPath) {
			return fullPath, true
		}
	}
	return "", false
}

func fileExistsAndExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir() && (info.Mode()&0111 != 0)
}

func runExternalCommand(path string, args []string) {
	cmd := exec.Command(path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error running command: %s\n", err)
	}
}

func pwd() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current working directory: %s\n", err)
		return
	}
	fmt.Println(dir)
}

func cd(path string) {
	err := os.Chdir(path)
	if err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", path)
	}
}

