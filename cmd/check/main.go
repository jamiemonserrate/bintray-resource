package main

import (
	"encoding/json"
	"os"
	"net/http"

	"github.com/jamiemonserrate/bintray-resource/check"

)

type VersionResponse struct {
	LatestVersion string `json:"latest_version"`
	Versions []string `json:"versions"`
}

func main(){
	checkRequest := check.CheckRequest{}
	json.NewDecoder(os.Stdin).Decode(&checkRequest)
	r, _ := http.Get("https://api.bintray.com/packages/jamiemonserrate/jamie-concourse/cf-artifactory")
	var response VersionResponse 
	json.NewDecoder(r.Body).Decode(&response)
	responseToReturn := check.CheckResponse{}
	if checkRequest.Version.Number == response.LatestVersion {
		json.NewEncoder(os.Stdout).Encode(responseToReturn)	
	}else {
		responseToReturn = append(responseToReturn, check.Version{Number: response.Versions[0]})
		json.NewEncoder(os.Stdout).Encode(responseToReturn)
	}
}
