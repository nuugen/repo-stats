package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
)

type Response struct {
	Status  int
	Message string
	Data    interface{}
}

func newResponse(status int, message string, data interface{}) *Response {
	return &Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func writeResponse(res http.ResponseWriter, statusCode int, message string, data interface{}) error {
	res.WriteHeader(statusCode)
	httpResponse := newResponse(statusCode, message, data)
	err := json.NewEncoder(res).Encode(httpResponse)
	return err
}

func fetchIssueCount(client github.Client, owner, repo string) (int, error) {
	repository, _, err := client.Repositories.Get(context.Background(), owner, repo)
	if err != nil {
		return 0, err
	}
	return *repository.OpenIssuesCount, nil
}

func issueCountHandler(client github.Client, res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	count, err := fetchIssueCount(client, params["owner"], params["repo"])
	if err != nil {
		writeResponse(res, http.StatusBadRequest, "Invalid repository specified", nil)
		return
	}
	writeResponse(res, http.StatusOK, "", count)
}

func main() {
	HOST := ":6699"
	ghClient := github.NewClient(nil)
	r := mux.NewRouter()
	r.HandleFunc("/issue-count/{owner}/{repo}", func(w http.ResponseWriter, req *http.Request) {
		issueCountHandler(*ghClient, w, req)
	})
	fmt.Printf("Listening on http://%s", HOST)
	log.Fatal(http.ListenAndServe(HOST, r))
}
