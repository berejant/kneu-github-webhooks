package main

import (
	"fmt"
	"github.com/go-playground/webhooks/v6/github"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

const (
	path = "/webhooks"
)

type Config struct {
	listenAddr   string
	githubSecret string
	projectDir   string
	deployScript string
}

func main() {
	var err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := Config{
		listenAddr:   os.Getenv("LISTEN_ADDR"),
		githubSecret: os.Getenv("GITHUB_SECRET"),
		projectDir:   os.Getenv("PROJECT_DIR"),
		deployScript: os.Getenv("DEPLOY_SCRIPT"),
	}

	if config.githubSecret == "" {
		log.Fatal("Error no Github secret in .env file")
	}

	hook, _ := github.New(github.Options.Secret(config.githubSecret))

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, github.PushEvent, github.PingEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				// ok event wasn't one of the ones asked to be parsed
				return
			} else {
				fmt.Println("Webhook request error: " + err.Error())
			}
		}
		switch payload.(type) {
		case github.PingPayload:
			go handlePingEvent(config, payload.(github.PingPayload))

		case github.PushPayload:
			go handlePushEvent(config, payload.(github.PushPayload))
		}
	})

	listenAddr := config.listenAddr
	if listenAddr == "" {
		listenAddr = ":3000"
	}

	fmt.Printf("%s Start listen as %s \n", currentDatetime(), listenAddr)
	http.ListenAndServe(listenAddr, nil)
}

func handlePingEvent(config Config, payload github.PingPayload) {
	fmt.Printf("%s: Receive ping event %s \n", currentDatetime(), payload.Repository.FullName)
}

func handlePushEvent(config Config, payload github.PushPayload) {
	fmt.Printf("%s: Receive push event %s \n", currentDatetime(), payload.Repository.FullName)

	deployScript := config.deployScript
	var cmd *exec.Cmd
	if config.deployScript == "" {
		cmd = exec.Command("git", "pull")
	} else {
		cmd = exec.Command(deployScript)
	}

	cmd.Dir = config.projectDir
	out, err := cmd.CombinedOutput()

	if len(out) != 0 {
		fmt.Printf("Output: \n%s\n", out)
	}

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("------------------------------")
}

func currentDatetime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
