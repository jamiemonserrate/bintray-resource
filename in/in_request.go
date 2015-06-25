package in

import (
	"github.com/hashicorp/go-version"
	"github.com/jamiemonserrate/bintray-resource/bintrayresource"
)

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
	if _, err := version.NewVersion(inRequest.RawVersion.Number); err != nil {
		return false, err.Error()
	}

	return true, ""
}
