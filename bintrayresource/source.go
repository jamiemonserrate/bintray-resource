package bintrayresource

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

func (s *Source) IsValid() bool {
	if s.SubjectName == "" || s.RepoName == "" || s.PackageName == "" || s.Username == "" || s.APIKey == "" {
		return false
	}
	return true
}
