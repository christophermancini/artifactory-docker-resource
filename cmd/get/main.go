package main

import (
	"log"
	"os"
	"path/filepath"

	resource "github.com/digitalocean/artifactory-docker-resource"
	rlog "github.com/digitalocean/concourse-resource-library/log"
	jlog "github.com/jfrog/jfrog-client-go/utils/log"
)

func main() {
	input := rlog.WriteStdin()
	defer rlog.Close()

	jlog.SetLogger(jlog.NewLogger(jlog.DEBUG, log.Writer()))

	log.Println("input:", input)

	var request resource.GetRequest
	err := request.Read(input)
	if err != nil {
		log.Fatalf("failed to read request input: %s", err)
	}

	err = request.Source.Validate()
	if err != nil {
		log.Fatalf("invalid source config: %s", err)
	}

	if len(os.Args) < 2 {
		log.Fatalf("missing arguments")
	}
	dir := os.Args[1]

	response, err := resource.Get(request, dir)
	if err != nil {
		log.Fatalf("failed to perform check: %s", err)
	}

	// write metadata to output dir
	os.MkdirAll(filepath.Join(dir, "resource"), os.ModePerm)
	err = response.Metadata.ToFiles(filepath.Join(dir, "resource"))
	if err != nil {
		log.Fatalf("failed to write metadata.json: %s", err)
	}
	log.Println("metadata written to:", filepath.Join(dir, "resource"))

	err = response.Write()
	if err != nil {
		log.Fatalf("failed to write response to stdout: %s", err)
	}

	log.Println("get complete")
}
