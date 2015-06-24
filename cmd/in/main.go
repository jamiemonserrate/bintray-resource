package main

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/jamiemonserrate/bintray-resource/bintray"
	"github.com/jamiemonserrate/bintray-resource/bintrayresource"
	"github.com/jamiemonserrate/bintray-resource/in"
)

func main() {
	if len(os.Args) < 2 {
		bintrayresource.LogErrorAndExit(errors.New("Please specify destination directory`"))
	}

	inRequest, err := decodeJSONFrom(os.Stdin)
	if err != nil {
		bintrayresource.LogErrorAndExit(err)
	}

	inCommand := in.NewInCommand(bintray.NewClient(
		bintray.DownloadURL,
		inRequest.Source.SubjectName,
		inRequest.Source.RepoName,
		inRequest.Source.Username,
		inRequest.Source.APIKey))

	destinationDir := os.Args[1]
	inResponse, err := inCommand.Execute(inRequest, destinationDir)
	if err != nil {
		bintrayresource.LogErrorAndExit(err)
	}

	err = writeToStdout(inResponse)
	if err != nil {
		bintrayresource.LogErrorAndExit(err)
	}
}

func decodeJSONFrom(request *os.File) (in.InRequest, error) {
	inRequest := in.InRequest{}
	err := json.NewDecoder(request).Decode(&inRequest)
	return inRequest, err
}

func writeToStdout(response *in.InResponse) error {
	return json.NewEncoder(os.Stdout).Encode(response)
}
