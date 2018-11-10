package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func brodcastURL(url string) {
	fmt.Printf("this is the url will be sent: %s", url)
}

type requestPayload struct {
	HostName string `json:"host_name"`
}

func logRequestPayload(requestionPayload requestPayload) {
	log.Printf("Host name: %s\n", requestionPayload.HostName)
}

func parseRequest(request *http.Request) *requestPayload {
	requestPayload := &requestPayload{request.Host}
	return requestPayload
}

func getRequestBody(request *http.Request) []byte {
	// Read body to buffer
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		panic(err)
	}

	request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return body
}

func processRequest(r *http.Request, body []byte) {
	fmt.Println(string(body[:]))
	fmt.Println(r.URL.Host)
}

func serverHandler(w http.ResponseWriter, r *http.Request) {
	requestPayload := parseRequest(r)
	logRequestPayload(*requestPayload)
	body := getRequestBody(r)
	processRequest(r, body)
	httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   requestPayload.HostName,
	}).ServeHTTP(w, r)
}

func main() {
	http.HandleFunc("/", serverHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
