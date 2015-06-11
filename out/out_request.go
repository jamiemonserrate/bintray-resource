package out

type Source struct {
	SubjectName string `json:"subject_name"`
	RepoName    string `json:"repo_name"`
	PackageName string `json:"package_name"`
}

type Version struct {
	Number string `json:"number"`
}

type OutRequest struct {
	Source Source `json:"source"`
	From   string `json:"from"`
}

type OutResponse struct {
	Version Version
}
