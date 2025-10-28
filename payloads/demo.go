package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	// header
	fmt.Println("payload: start")
	fmt.Printf("pid=%d argv=%v\n", os.Getpid(), os.Args)
	fmt.Printf("env: %s=%q\n", "PAYLOAD_TAG", os.Getenv("PAYLOAD_TAG"))

	// demonstrate streaming stdout
	for i := 1; i <= 5; i++ {
		fmt.Printf("stdout: tick %d\n", i)
		time.Sleep(200 * time.Millisecond)
	}

	// demonstrate streaming stderr
	fmt.Fprintln(os.Stderr, "stderr: this is an error-like message (not fatal)")

	// perform a tiny computation and print
	sum := 0
	for i := 0; i < 1000; i++ {
		sum += i
	}
	fmt.Printf("computed sum: %d\n", sum)

	// exit control: if PAYLOAD_FAIL=1, exit non-zero
	if os.Getenv("PAYLOAD_FAIL") == "1" {
		fmt.Fprintln(os.Stderr, "exiting non-zero due to PAYLOAD_FAIL=1")
		os.Exit(42)
	}

	fmt.Println("payload: done")
}
