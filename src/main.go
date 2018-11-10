package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/akamensky/argparse"
	"github.com/gobwas/glob"
	"github.com/m-motawea/aggregator"
)

type requestPayload struct {
	HostName string `json:"host_name"`
}

type server struct {
	HostNames []string
	target    *string
	scheme    string
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
	aggregator.RegsiterCall(*r, body, []byte{}, 200)
}

func knownURL(urls []string, host string) bool {
	for i := 0; i < len(urls); i++ {
		g := glob.MustCompile(urls[i])
		if g.Match(host) {
			return true
		}
	}
	return false
}

func (s *server) serverHandler(w http.ResponseWriter, r *http.Request) {
	requestPayload := parseRequest(r)
	logRequestPayload(*requestPayload)
	if knownURL(s.HostNames, r.Host) {
		var body = []byte{}
		if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
			body = getRequestBody(r)
		}
		processRequest(r, body)
	}
	var host string
	if s.target != nil {
		host = *s.target
	} else {
		host = requestPayload.HostName
	}
	httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: s.scheme,
		Host:   host,
	}).ServeHTTP(w, r)
}

func main() {
	parser := argparse.NewParser("SAI", "Simple API inspector")
	hostsList := parser.List("H", "hostname", &argparse.Options{Required: true, Help: "Enter URL list with globs"})
	target := parser.String("t", "target", &argparse.Options{Required: false, Help: "Target for the requests (will proxy all requests to this endpoint"})
	port := parser.Int("p", "port", &argparse.Options{Required: false, Help: "port", Default: 8080})
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		panic(parser.Usage(err))
	}
	var host *string
	scheme := "http"
	if target != nil {
		u, err := url.Parse(*target)
		if err != nil {
			panic("target must be a valid url")
		}
		host = &u.Host
		scheme = u.Scheme
	}
	s := server{*hostsList, host, scheme}
	http.HandleFunc("/", s.serverHandler)
	go func() {
		for {
			time.Sleep(time.Second * 10)
			fmt.Println(string(aggregator.GetRaml()))
		}
	}()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
