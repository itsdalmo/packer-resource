package main

import (
	"encoding/json"
	"github.com/itsdalmo/packer-resource/src/models"
	"github.com/itsdalmo/packer-resource/src/out"
	"log"
	"os"
)

func main() {
	var request models.PutRequest
	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		log.Fatalf("failed to unmarshal request: %s", err)
	}

	if len(os.Args) < 2 {
		log.Fatalf("missing arguments")
	}
	sourceDir := os.Args[1]

	response, err := out.Run(request, sourceDir)
	if err != nil {
		log.Fatalf("put failed: %s", err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		log.Fatalf("failed to marshal response: %s", err)
	}
}
