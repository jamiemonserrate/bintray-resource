package main

import (
	"encoding/json"
	"os"

	"github.com/jamiemonserrate/bintray-resource/bintray"
	"github.com/jamiemonserrate/bintray-resource/bintrayresource"
	"github.com/jamiemonserrate/bintray-resource/out"
)

func main() {
	outRequest, err := decodeJSONFrom(os.Stdin)
	if err != nil {
		bintrayresource.LogErrorAndExit(err)
	}

	outCommand := out.NewOutCommand(bintray.NewClient(
		bintray.APIURL,
		outRequest.Source.SubjectName,
		outRequest.Source.RepoName,
		outRequest.Source.Username,
		outRequest.Source.APIKey))

	outResponse, err := outCommand.Execute(outRequest)
	if err != nil {
		bintrayresource.LogErrorAndExit(err)
	}

	err = writeToStdout(outResponse)
	if err != nil {
		bintrayresource.LogErrorAndExit(err)
	}
}

func decodeJSONFrom(request *os.File) (out.OutRequest, error) {
	outRequest := out.OutRequest{}
	err := json.NewDecoder(request).Decode(&outRequest)
	return outRequest, err
}

func writeToStdout(response *out.OutResponse) error {
	return json.NewEncoder(os.Stdout).Encode(response)
}
