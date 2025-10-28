package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const payloadsDir = "payloads"

func BuildForArch(srcPath, arch string) (string, error) {
	arch = strings.TrimSpace(arch)
	goArch := map[string]string{
		"x86_64":  "amd64",
		"aarch64": "arm64",
		"armv7l":  "arm",
	}[arch]
	if goArch == "" {
		return "", fmt.Errorf("unsupported arch: %s", arch)
	}

	if err := os.MkdirAll(payloadsDir, 0755); err != nil {
		return "", fmt.Errorf("create payloads dir: %w", err)
	}

	base := filepath.Base(srcPath)
	binaryName := fmt.Sprintf("%s_%s", base, goArch)
	binaryPath := filepath.Join(payloadsDir, binaryName)
	if fi, err := os.Stat(binaryPath); err == nil && fi.Mode().IsRegular() {
		fmt.Printf("[+] Using existing binary from %s\n", binaryPath)
		return binaryPath, nil
	}

	fmt.Printf("[+] Compiling %s for arch %s...\n", srcPath, goArch)
	cmd := exec.Command("go", "build",
		"-ldflags", "-s -w",
		"-o", binaryPath, srcPath)
	cmd.Env = append(os.Environ(),
		"GOOS=linux",
		fmt.Sprintf("GOARCH=%s", goArch),
		"CGO_ENABLED=0",
	)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		_ = os.Remove(binaryPath)
		return "", err
	}

	return binaryPath, nil
}

