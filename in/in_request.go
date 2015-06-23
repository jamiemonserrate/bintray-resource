package in

import "github.com/jamiemonserrate/bintray-resource/bintrayresource"

type InRequest struct {
	Source     bintrayresource.Source  `json:"source"`
	RawVersion bintrayresource.Version `json:"version"`
}

type InResponse struct {
	Version bintrayresource.Version
}

func (inRequest *InRequest) IsValid() (bool, string) {
	if isValid, errMssg := inRequest.Source.IsValid(); !isValid {
		return false, errMssg
	}
	if inRequest.RawVersion.Number == "" {
		return false, "Please specify the Version"
	}

	return true, ""
}
