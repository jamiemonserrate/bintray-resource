package main

import (
	"encoding/json"
	"os"

	"github.com/jamiemonserrate/bintray-resource/bintray"
	"github.com/jamiemonserrate/bintray-resource/bintrayresource"
	"github.com/jamiemonserrate/bintray-resource/out"
)

func main() {
	outRequest := decodeJSONFrom(os.Stdin)

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

	writeToStdout(outResponse)
}

func decodeJSONFrom(request *os.File) out.OutRequest {
	outRequest := out.OutRequest{}
	json.NewDecoder(request).Decode(&outRequest)
	return outRequest
}

func writeToStdout(response *out.OutResponse) {
	json.NewEncoder(os.Stdout).Encode(response)
}
