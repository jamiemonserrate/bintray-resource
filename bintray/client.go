package bintray

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const URL string = "https://api.bintray.com"

type Client struct {
	url         string
	subjectName string
	repoName    string
}

type Package struct {
	LatestVersion string   `json:"latest_version"`
	Versions      []string `json:"versions"`
}

type BintrayClient interface {
	GetPackage(string) Package
}

func NewClient(bintrayURL, subjectName, repoName string) *Client {
	return &Client{url: bintrayURL, subjectName: subjectName, repoName: repoName}
}

func (client *Client) GetPackage(packageName string) Package {
	var bintrayPackage Package
	response, _ := http.Get(client.getPackageURL(packageName))
	json.NewDecoder(response.Body).Decode(&bintrayPackage)
	return bintrayPackage
}

func (client *Client) getPackageURL(packageName string) string {
	return fmt.Sprintf("%s/packages/%s/%s/%s", client.url, client.subjectName, client.repoName, packageName)
}
