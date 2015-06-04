package check

type Version struct {
	Number string `json:"number"`
}

type CheckRequest struct {
	Version Version `json:"version"`
}

type CheckResponse []Version
