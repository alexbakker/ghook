package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	logger "log"
	"net/http"
	"os"
	"strings"

	"github.com/alexbakker/ghook"
	"github.com/google/go-github/v32/github"
)

type Config struct {
	Secret   *string         `json:"secret"`
	Handlers []HandlerConfig `json:"handlers"`
}

type HandlerConfig struct {
	Repo    string `json:"repo"`
	Ref     string `json:"ref"`
	Command string `json:"command"`
}

var (
	log = logger.New(os.Stderr, "", 0)

	config   Config
	queue    = new(TaskQueue)
	addr     = flag.String("addr", "127.0.0.1:8080", "address to listen on")
	filename = flag.String("config", "config.json", "the filename of the configuration file")
)

func main() {
	flag.Parse()

	// parse the configuration file
	bytes, err := ioutil.ReadFile(*filename)
	if err != nil {
		log.Fatalf("error loading config: %s", err)
	}
	if err = json.Unmarshal(bytes, &config); err != nil {
		log.Fatalf("error parsing config: %s", err)
	}

	// verify the config
	if config.Secret == nil {
		log.Fatalf("error: secret not set")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	worker := Worker{Queue: queue}
	go worker.Run(ctx)

	hook := ghook.New([]byte(*config.Secret), handleEvent)
	log.Fatalf("error: %s", http.ListenAndServe(*addr, hook))
}

func handleEvent(event *ghook.Event) error {
	log.Printf("handling %s event %s\n", event.Name, event.GUID)

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
	for _, handler := range config.Handlers {
		if handler.Repo == *event.Repo.FullName &&
			strings.Contains(*event.Ref, handler.Ref) &&
			handler.Command != "" {
			queue.Push(&Task{Command: handler.Command})
		}
	}

	return nil
}
