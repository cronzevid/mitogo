package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
	"path/filepath"
	"time"
)

func privateKey() ssh.AuthMethod {
	key, err := os.ReadFile(filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa"))
	if err != nil {
		panic(err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		panic(err)
	}
	return ssh.PublicKeys(signer)
}

func ConnectSSH(target, bastion, user string, port int) (*ssh.Client, error) {
	cfg := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{privateKey()},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}
	addr := fmt.Sprintf("%s:%d", target, port)
	if bastion == "" {
		return ssh.Dial("tcp", addr, cfg)
	}

	bastionClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", bastion, port), cfg)
	if err != nil {
		return nil, fmt.Errorf("bastion connect: %w", err)
	}
	conn, err := bastionClient.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("proxy dial: %w", err)
	}
	ncc, chans, reqs, err := ssh.NewClientConn(conn, addr, cfg)
	if err != nil {
		return nil, fmt.Errorf("target client conn: %w", err)
	}
	return ssh.NewClient(ncc, chans, reqs), nil
}

func RemoteArch(c *ssh.Client) (string, error) {
	sess, err := c.NewSession()
	if err != nil {
		return "", err
	}
	defer sess.Close()
	out, err := sess.CombinedOutput("uname -m")
	if err != nil {
		return "", err
	}
	return string(out), nil
}

