package out

import "github.com/jamiemonserrate/bintray-resource/bintrayresource"

type OutRequest struct {
	Source      bintrayresource.Source `json:"source"`
	From        string                 `json:"from"`
	VersionFile string                 `json:"version_file"`
}

type OutResponse struct {
	Version bintrayresource.Version
}

func (outRequest *OutRequest) IsValid() bool {
	return outRequest.Source.IsValid() && outRequest.From != "" && outRequest.VersionFile != ""
}
