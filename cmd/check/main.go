package main

import (
	"encoding/json"
	"os"

	"github.com/jamiemonserrate/bintray-resource/bintray"
	"github.com/jamiemonserrate/bintray-resource/check"
)

func main() {
	checkRequest := decodeJSONFrom(os.Stdin)

	checkResponse := versionDiffs(checkRequest)

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

func versionDiffs(checkRequest check.CheckRequest) check.CheckResponse {
	client := bintray.NewClient("https://api.bintray.com", "jamiemonserrate", "jamie-concourse")
	bintrayPackage := client.GetPackage("cf-artifactory")

	response := check.CheckResponse{}
	if checkRequest.Version.Number != bintrayPackage.LatestVersion {
		response = append(response, check.Version{Number: bintrayPackage.Versions[0]})
	}
	return response
}
