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

type Metadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (s *Source) IsValid() (bool, string) {
	if s.SubjectName == "" {
		return false, "Please specify the SubjectName"
	}
	if s.RepoName == "" {
		return false, "Please specify the RepoName"
	}
	if s.PackageName == "" {
		return false, "Please specify the PackageName"
	}
	if s.Username == "" {
		return false, "Please specify the Username"
	}
	if s.APIKey == "" {
		return false, "Please specify the APIKey"
	}
	return true, ""
}
