package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/go-playground/webhooks/v6/github"
)

const (
	path = "/webhooks"
)

var (

	// GitRes  git请求返回的结果
	GitRes apiRes

	AuthStr string
)

// RequestInfo 请求字段的结构体
type RequestInfo struct {
	URL     string
	Cookies []*http.Cookie
	Params  map[string][]string
}
type apiRes struct {
	Commit struct {
		Message string `json:"message"`
	}
	HTMLURL     string `json:"html_url"`
	CommentsURL string `json:"comments_url"`
	WorkFlowRun struct {
		HTMLURL string `json:"html_url"`
		Actor   struct {
			Login string `json:"login"`
		} `json:"actor"`
	} `json:"workflow_run"`
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
			resp, err := GetReq(strings.Replace(workflow.WorkflowRun.Repository.CommitsURL, "{/sha}", "/"+workflow.WorkflowRun.HeadSha, -1))
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
			if err := json.Unmarshal(body, &GitRes); err != nil {
				log.Println(err.Error())
				return
			}

			if workflow.Action == "completed" && workflow.WorkflowRun.Conclusion == "success" {
				switch workflow.WorkflowRun.Name {
				case "Docker":
					log.Println("Receive Docker workflow finished.")
					go deployNewContainer(stableContainerName, stableImageNameDocker)
				case "Docker-nightly":
					log.Println("Receive Docker workflow finished.")
					go deployNewContainer(nightlyContainerName, nightlyImageNameDocker)
				}
			}
		}
	})
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// GetReq 获取请求
//  @param info 传入的参数、url、cookie信息
//  @return resp
//  @return err
func GetReq(URL string) (resp *http.Response, err error) {
	//创建client
	resp, err = http.Get(URL)
	if err != nil {
		return
	}
	return
}

func init() {
	// 获取镜像Token
	username := os.Getenv("DOCKER_USERNAME")
	password := os.Getenv("DOCKER_PASSWORD")
	if username == "" || password == "" {
		log.Fatal("DOCKER_USERNAME or DOCKER_PASSWORD is empty")
	}
	a := types.AuthConfig{
		Username: username,
		Password: password,
	}
	encodeJSON, err := json.Marshal(&a)
	if err != nil {
		log.Fatal(err.Error())
	}
	AuthStr = base64.StdEncoding.EncodeToString(encodeJSON)
}
