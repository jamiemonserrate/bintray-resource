package out

import "github.com/jamiemonserrate/bintray-resource"

type OutRequest struct {
	Source      bintrayresource.Source `json:"source"`
	From        string                 `json:"from"`
	VersionFile string                 `json:"version_file"`
}

type OutResponse struct {
	Version bintrayresource.Version
}
