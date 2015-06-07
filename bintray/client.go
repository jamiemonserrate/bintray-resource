package bintray

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

const URL string = "https://api.bintray.com"

type Client struct {
	url         string
	subjectName string
	repoName    string
}

type BintrayClient interface {
	GetPackage(string) Package
	DownloadPackage(string, string, string)
}

func NewClient(bintrayURL, subjectName, repoName string) *Client {
	return &Client{url: bintrayURL, subjectName: subjectName, repoName: repoName}
}

func (client *Client) GetPackage(packageName string) Package {
	var bintrayPackage Package

	response, _ := http.Get(client.getPackageURL(packageName))
	defer response.Body.Close()

	json.NewDecoder(response.Body).Decode(&bintrayPackage)
	return bintrayPackage
}

func (client *Client) DownloadPackage(packageName, version, destinationDir string) {
	downloadedFile, err := os.Create(filepath.Join(destinationDir, packageName))
	if err != nil {
		panic(err)
	}
	defer downloadedFile.Close()
	response, err := http.Get(client.inPackageURL(packageName, version))
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	_, err = io.Copy(downloadedFile, response.Body)

	if err != nil {
		panic(err)
	}
}

func (client *Client) getPackageURL(packageName string) string {
	getPackagePath := path.Join("packages", client.subjectName, client.repoName, packageName)
	return fmt.Sprintf("%s/%s", client.url, getPackagePath)
}

func (client *Client) inPackageURL(packageName, version string) string {
	downloadPackagePath := path.Join(client.subjectName, client.repoName, version, packageName)
	return fmt.Sprintf("%s/%s", client.url, downloadPackagePath)
}
