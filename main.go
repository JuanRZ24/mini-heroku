package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "uso: mh deploy <repo>")
		os.Exit(1)
	}

	comando := os.Args[1]

	switch comando {
	case "deploy":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "error: falta el repo")
			os.Exit(1)
		}
		repo := os.Args[2]
		ruta, err := os.MkdirTemp("", "mh-*")

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer os.RemoveAll(ruta)
		cmd := exec.Command("git", "clone", "--depth", "1", repo, ruta)
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		cmd = exec.Command("git", "rev-parse", "HEAD")
		cmd.Dir = ruta
		salida, err := cmd.Output()
		salida = []byte(strings.TrimSpace(string(salida)))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		nombreRepo := path.Base(repo)
		nombreRepo = strings.TrimSuffix(nombreRepo, ".git")

		nombreRepo = strings.ToLower(nombreRepo)

		imagen := fmt.Sprintf("mh/%s:%s", nombreRepo, salida)

		cmdDocker := exec.Command("docker", "build", "-t", imagen, ruta)
		cmdDocker.Stdout = os.Stdout
		cmdDocker.Stderr = os.Stderr
		err = cmdDocker.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "error: comando desconocido %q\n", comando)
		os.Exit(1)
	}

}
