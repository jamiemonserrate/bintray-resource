package main

import (
	"encoding/json"
	"os"

	"github.com/jamiemonserrate/bintray-resource/bintray"
	"github.com/jamiemonserrate/bintray-resource/bintrayresource"
	"github.com/jamiemonserrate/bintray-resource/check"
)

func main() {
	checkRequest, err := decodeJSONFrom(os.Stdin)
	if err != nil {
		bintrayresource.LogErrorAndExit(err)
	}

	checkCommand := check.NewCheckCommand(bintray.NewClient(
		bintray.APIURL,
		checkRequest.Source.SubjectName,
		checkRequest.Source.RepoName,
		checkRequest.Source.Username,
		checkRequest.Source.APIKey))

	checkResponse, err := checkCommand.Execute(checkRequest)
	if err != nil {
		bintrayresource.LogErrorAndExit(err)
	}

	err = writeToStdout(checkResponse)
	if err != nil {
		bintrayresource.LogErrorAndExit(err)
	}
}

func decodeJSONFrom(request *os.File) (check.CheckRequest, error) {
	checkRequest := check.CheckRequest{}
	err := json.NewDecoder(request).Decode(&checkRequest)
	return checkRequest, err
}

func writeToStdout(response check.CheckResponse) error {
	return json.NewEncoder(os.Stdout).Encode(response)
}
