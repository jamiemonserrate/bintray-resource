package main

import (
	"encoding/json"
	"os"

	"github.com/jamiemonserrate/bintray-resource/bintray"
	"github.com/jamiemonserrate/bintray-resource/in"
)

func main() {
	inRequest := decodeJSONFrom(os.Stdin)

	inCommand := in.NewInCommand(bintray.NewClient(
		bintray.DownloadURL,
		inRequest.Source.SubjectName,
		inRequest.Source.RepoName,
		inRequest.Source.Username,
		inRequest.Source.APIKey))

	destinationDir := os.Args[1]
	inResponse := inCommand.Execute(inRequest, destinationDir)

	writeToStdout(inResponse)
}

func decodeJSONFrom(request *os.File) in.InRequest {
	inRequest := in.InRequest{}
	json.NewDecoder(request).Decode(&inRequest)
	return inRequest
}

func writeToStdout(response in.InResponse) {
	json.NewEncoder(os.Stdout).Encode(response)
}
