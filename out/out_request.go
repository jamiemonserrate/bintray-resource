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

func (outRequest *OutRequest) IsValid() (bool, string) {
	if isValid, errMssg := outRequest.Source.IsValid(); !isValid {
		return false, errMssg
	}
	if outRequest.From == "" {
		return false, "Please specify the From"
	}
	if outRequest.VersionFile == "" {
		return false, "Please specify the VersionFile"
	}

	return true, ""
}
