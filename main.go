package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func main() {
	//default
	root := os.Getenv("HOME")
	dir := "dhnt"
	base := filepath.Join(root, dir)

	if len(os.Args) > 1 {
		base = os.Args[1]
		root, dir = filepath.Split(base)
	}

	url := "https://github.com/dhnt/dhnt.git"
	branch := runtime.GOOS

	//
	if _, err := os.Stat(root); os.IsNotExist(err) {
		fmt.Println(err)
		os.Exit(1)

	}
	if _, err := os.Stat(base); err == nil {
		fmt.Printf("Directory %v exists. M3 already installed?\n", base)
		os.Exit(1)

	}
	if err := os.Chdir(root); err != nil {
		fmt.Println(err)
		os.Exit(1)

	}

	// clone
	fmt.Printf("Cloning M3 branch %v into %v\n", branch, dir)

	err := cloneRepo(url, branch, dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// install
	fmt.Println("Installing M3 ...")
	install(base)

	fmt.Println("M3 installation done!")
	os.Exit(0)
}

func cloneRepo(url, branch, directory string) error {
	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL:               url,
		ReferenceName:     plumbing.NewBranchReferenceName(branch),
		SingleBranch:      true,
		Depth:             1,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Progress:          os.Stdout,
	})
	if err != nil {
		return err
	}

	// retrieving the branch being pointed by HEAD
	ref, err := r.Head()
	if err != nil {
		return err
	}
	// retrieving the commit object
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return err
	}

	fmt.Println(commit)
	return nil
}

func install(base string) {
	m3d := filepath.Join(base, "go/bin/m3d")
	//tray := filepath.Join(base, "go/bin/systray")
	execute(base, "sudo", m3d, "install", "--base", base)
	execute(base, "sudo", m3d, "start")
	// execute(base, tray)
}

// execute sets up env and runs file
func execute(base, binary string, args ...string) error {
	cmd := exec.Command(binary, args...)
	//
	cmd.Env = DefaultEnviron(base)

	//
	cmdOut, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("error creating stdout", err)
		return err
	}
	scanOut := bufio.NewScanner(cmdOut)
	go func() {
		for scanOut.Scan() {
			fmt.Println(">", scanOut.Text())
		}
	}()

	cmdErr, err := cmd.StderrPipe()
	scanErr := bufio.NewScanner(cmdErr)
	go func() {
		for scanErr.Scan() {
			fmt.Println(">>", scanErr.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Println("error starting cmd", err)
		return err
	}

	return cmd.Wait()
}
