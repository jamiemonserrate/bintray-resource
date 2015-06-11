package main

import (
	"encoding/json"
	"os"

	"github.com/jamiemonserrate/bintray-resource/bintray"
	"github.com/jamiemonserrate/bintray-resource/out"
)

func main() {
	outRequest := decodeJSONFrom(os.Stdin)

	outCommand := out.NewOutCommand(bintray.NewClient(
		"https://api.bintray.com",
		"jamiemonserrate",
		"jamie-concourse"))

	outResponse := outCommand.Execute(outRequest)

	writeToStdout(outResponse)
}

func decodeJSONFrom(request *os.File) out.OutRequest {
	outRequest := out.OutRequest{}
	json.NewDecoder(request).Decode(&outRequest)
	return outRequest
}

func writeToStdout(response out.OutResponse) {
	json.NewEncoder(os.Stdout).Encode(response)
}
