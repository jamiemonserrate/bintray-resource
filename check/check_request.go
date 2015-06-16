package check

import (
	"github.com/hashicorp/go-version"
	"github.com/jamiemonserrate/bintray-resource/bintrayresource"
)

type CheckRequest struct {
	Source     bintrayresource.Source  `json:"source"`
	RawVersion bintrayresource.Version `json:"version"`
}

func (checkRequest *CheckRequest) Version() *version.Version {
	convertedVersion, _ := version.NewVersion(checkRequest.RawVersion.Number)
	return convertedVersion
}

type CheckResponse []bintrayresource.Version
