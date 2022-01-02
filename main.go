package main

import (
	"fmt"
	"log"

	"net/http"

	"github.com/go-playground/webhooks/v6/github"
)

const (
	path = "/webhooks"
)

func main() {
	hook, _ := github.New(github.Options.Secret("heyuheng1.22.3"))

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, github.ReleaseEvent, github.PushEvent, github.WorkflowRunEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				fmt.Println("Invalid EventType")
				// ok event wasn;t one of the ones asked to be parsed
			}
		}
		switch payload.(type) {
		case github.ReleasePayload:
			release := payload.(github.ReleasePayload)
			// Do whatever you want from here...
			fmt.Printf("%+v", release)
		case github.PushPayload:
			// push := payload.(github.PushPayload)
			// fmt.Printf("%+v", push)
		case github.WorkflowRunPayload:
			workflow := payload.(github.WorkflowRunPayload)
			if workflow.Action == "completed" && workflow.WorkflowRun.Conclusion == "success" {
				switch workflow.WorkflowRun.Name {
				case "Docker":
					log.Println("Receive Docker workflow finished.")
					go deployNewContainer(stableContainerName, stableImageName)
				case "Docker-nightly":
					log.Println("Receive Docker workflow finished.")
					go deployNewContainer(nightlyContainerName, nightlyImageName)
				}
			}
		}
	})
	http.ListenAndServe(":3000", nil)
}
