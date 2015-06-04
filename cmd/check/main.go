package main

import (
	"encoding/json"
	"os"

	"github.com/jamiemonserrate/bintray-resource/check"

)

func main(){
	json.NewEncoder(os.Stdout).Encode(check.CheckResponse{})	
}
