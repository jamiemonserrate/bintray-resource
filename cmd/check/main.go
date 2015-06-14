package main

import (
	"encoding/json"
	"os"

	"github.com/jamiemonserrate/bintray-resource/bintray"
	"github.com/jamiemonserrate/bintray-resource/check"
)

func main() {
	checkRequest := decodeJSONFrom(os.Stdin)

	checkCommand := check.NewCheckCommand(bintray.NewClient(
		bintray.APIURL,
		checkRequest.Source.SubjectName,
		checkRequest.Source.RepoName,
		checkRequest.Source.Username,
		checkRequest.Source.APIKey))

	checkResponse, _ := checkCommand.Execute(checkRequest)

	writeToStdout(checkResponse)
}

func decodeJSONFrom(request *os.File) check.CheckRequest {
	checkRequest := check.CheckRequest{}
	json.NewDecoder(request).Decode(&checkRequest)
	return checkRequest
}

func writeToStdout(response check.CheckResponse) {
	json.NewEncoder(os.Stdout).Encode(response)
}
