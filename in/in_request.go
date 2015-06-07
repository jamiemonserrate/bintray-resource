package in

type Source struct {
	SubjectName string `json:"subject_name"`
	RepoName    string `json:"repo_name"`
	PackageName string `json:"package_name"`
}

type Version struct {
	Number string `json:"number"`
}

type InRequest struct {
	Source     Source  `json:"source"`
	RawVersion Version `json:"version"`
}

type InResponse struct {
	Version Version
}
