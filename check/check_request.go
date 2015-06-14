package check

import "github.com/hashicorp/go-version"

type Source struct {
	SubjectName string `json:"subject_name"`
	RepoName    string `json:"repo_name"`
	PackageName string `json:"package_name"`
	Username    string `json:"username"`
	APIKey      string `json:"api_key"`
}

type Version struct {
	Number string `json:"number"`
}

type CheckRequest struct {
	Source     Source  `json:"source"`
	RawVersion Version `json:"version"`
}

func (checkRequest *CheckRequest) Version() *version.Version {
	convertedVersion, _ := version.NewVersion(checkRequest.RawVersion.Number)
	return convertedVersion
}

type CheckResponse []Version
