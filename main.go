package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"net/http"
	"net/url"

	"github.com/go-playground/webhooks/v6/github"
)

const (
	path = "/webhooks"
)

var (
	message, htmlURL string
)

// RequestInfo 请求字段的结构体
type RequestInfo struct {
	URL     string
	Cookies []*http.Cookie
	Params  map[string][]string
}

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

			resp, err := PostWithParams(
				RequestInfo{
					URL: strings.Replace(workflow.WorkflowRun.Repository.CommitsURL, "{/sha}", workflow.WorkflowRun.HeadSha, -1),
				},
			)
			if err != nil {
				log.Println(err.Error())
				return
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println(err.Error())
				return
			}
			defer resp.Body.Close()
			resMap := make(map[string]interface{})
			if err := json.Unmarshal(body, &resMap); err != nil {
				log.Println(err.Error())
				return
			}
			message = resMap["message"].(string)
			htmlURL = resMap["html_url"].(string)
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

// PostWithParams 发送带Cookie Params的POST请求
func PostWithParams(info RequestInfo) (resp *http.Response, err error) {
	params := url.Values{}
	for key, values := range info.Params {
		for index := range values {
			params.Add(key, values[index])
		}
	}
	params.Set("timestamp", fmt.Sprintf("%d", time.Now().UnixNano()))

	// body := ioutil.NopCloser(strings.NewReader(params.Encode()))
	resp, err = http.PostForm(info.URL+"?timestamp="+fmt.Sprint(time.Now().UnixNano()), params)
	if err != nil {
		log.Println(err.Error())
		return
	}

	return
}
