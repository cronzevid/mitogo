package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

func UploadAndRun(c *ssh.Client, localPath string) error {
	session, err := c.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	remoteFile := fmt.Sprintf("/tmp/%s", filepath.Base(localPath))

	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()

		f, _ := os.Open(localPath)
		defer f.Close()
		info, _ := f.Stat()
		fmt.Fprintf(w, "C0755 %d %s\n", info.Size(), filepath.Base(localPath))
		io.Copy(w, f)
		fmt.Fprint(w, "\x00")
	}()

	if err := session.Run(fmt.Sprintf("scp -t /tmp")); err != nil {
		return fmt.Errorf("upload failed: %w", err)
	}

	return runWithWrapper(c, remoteFile)
}

func runWithWrapper(c *ssh.Client, remoteFile string) error {
	sess, err := c.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()

	sess.Stdout = os.Stdout
	sess.Stderr = os.Stderr

	wrapper := fmt.Sprintf(
		"set -e; chmod +x %[1]s; echo '[+] Executing'; %[1]s; rc=$?; rm -f %[1]s; exit $rc",
		remoteFile)

	if err := sess.Run(wrapper); err != nil {
		return fmt.Errorf("remote run failed: %w", err)
	}
	return nil
}

