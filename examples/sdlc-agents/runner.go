package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type RunResult struct {
	Language string
	Stdout   string
	Stderr   string
	ExitCode int
}

func RunInDocker(ctx context.Context, language, code string) (*RunResult, error) {
	switch language {
	case "python":
		return runPython(ctx, code)
	default:
		return nil, fmt.Errorf("unsupported language: %s", language)
	}
}

func runPython(ctx context.Context, code string) (*RunResult, error) {
	dir, err := os.MkdirTemp("", "agent-run-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)

	srcPath := filepath.Join(dir, "main.py")
	if err := os.WriteFile(srcPath, []byte(code), 0644); err != nil {
		return nil, err
	}

	args := []string{
		"run", "--rm",
		"-v", fmt.Sprintf("%s:/app", dir),
		"-w", "/app",
		"python:3.11-alpine",
		"python", "main.py", "10", // sample input for test
	}

	cmd := exec.CommandContext(ctx, "docker", args...)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = cmd.Run()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}

	return &RunResult{
		Language: "python",
		Stdout:   outBuf.String(),
		Stderr:   errBuf.String(),
		ExitCode: exitCode,
	}, nil
}
