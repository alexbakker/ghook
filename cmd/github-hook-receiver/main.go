package main

import (
	"encoding/json"
	"errors"
	"flag"
	logger "log"
	"net/http"
	"os"
	"os/exec"

	"github.com/alexbakker/ghook"
	"github.com/google/go-github/github"
)

var (
	log       = logger.New(os.Stderr, "", 0)
	secretEnv = "HOOK_SECRET"

	addr    = flag.String("addr", "127.0.0.1:8080", "address to listen on")
	branch  = flag.String("branch", "master", "the branch this hook should apply for (all if omitted)")
	command = flag.String("cmd", "", "command to execute")
)

func main() {
	flag.Parse()

	secret := os.Getenv(secretEnv)
	if secret == "" {
		log.Fatalf("error: %s not set", secretEnv)
	}

	hook := ghook.New([]byte(secret), handleEvent)
	log.Fatalf("error: %s", http.ListenAndServe(*addr, hook))
}

func handleEvent(event *ghook.Event) error {
	if event.Name == "ping" {
		return nil
	}

	if event.Name != "push" {
		return errors.New("unhandled event")
	}

	var info github.PushEvent
	if err := json.Unmarshal(event.Payload, &info); err != nil {
		return err
	}

	return handlePush(&info)
}

func handlePush(event *github.PushEvent) error {
	if *command != "" {
		cmd := exec.Command("/bin/sh", "-c", *command)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
