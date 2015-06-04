package check

type Version struct {
	Path string `json:"path"`
}

type CheckResponse []Version
