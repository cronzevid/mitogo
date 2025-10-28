package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
)

const (
	defaultPort = 22
)

func main() {
	host := flag.String("host", "", "target host")
	bastion := flag.String("bastion", "", "optional bastion host")
	file := flag.String("file", "", "Go file to execute remotely")
	overwrite := flag.Bool("overwrite", false, "overwrite local binary in payloads dir")
	currentUser, _ := user.Current()
	userFlag := flag.String("user", currentUser.Username, "SSH user for connection")

	flag.Parse()

	if *host == "" || *file == "" {
		log.Fatalf("Usage: %s -host <target> [-bastion <jump>] -file <path> [-user <user>] [-overwrite]", os.Args[0])
	}

	sshClient, err := ConnectSSH(*host, *bastion, *userFlag, defaultPort)
	if err != nil {
		log.Fatalf("SSH connect: %v", err)
	}
	defer sshClient.Close()

	arch, err := RemoteArch(sshClient)
	if err != nil {
		log.Fatalf("get arch: %v", err)
	}
	fmt.Printf("[+] Remote arch: %s\n", arch)

	binPath, err := BuildForArchWithOverwrite(*file, arch, *overwrite)
	if err != nil {
		log.Fatalf("build: %v", err)
	}

	if err := UploadAndRun(sshClient, binPath); err != nil {
		log.Fatalf("remote exec: %v", err)
	}
}
