package in

import "github.com/jamiemonserrate/bintray-resource/bintrayresource"

type InRequest struct {
	Source     bintrayresource.Source  `json:"source"`
	RawVersion bintrayresource.Version `json:"version"`
}

type InResponse struct {
	Version bintrayresource.Version
}
