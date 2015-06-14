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

const DownloadURL string = "https://dl.bintray.com"
const APIURL string = "https://api.bintray.com"

type Client struct {
	url         string
	subjectName string
	repoName    string
	username    string
	password    string
}

type BintrayClient interface {
	GetPackage(string) Package
	DownloadPackage(string, string, string)
	UploadPackage(string, string, string) error
}

func NewClient(bintrayURL, subjectName, repoName, username, password string) *Client {
	return &Client{url: bintrayURL, subjectName: subjectName, repoName: repoName,
		username: username, password: password}
}

func (client *Client) GetPackage(packageName string) Package {
	var bintrayPackage Package
	response, _ := http.Get(client.getPackageURL(packageName))
	defer response.Body.Close()
	json.NewDecoder(response.Body).Decode(&bintrayPackage)
	return bintrayPackage
}

func (client *Client) DownloadPackage(packageName, version, destinationDir string) {
	downloadedFile, _ := os.Create(filepath.Join(destinationDir, packageName))
	defer downloadedFile.Close()
	response, _ := http.Get(client.inPackageURL(packageName, version))
	defer response.Body.Close()
	io.Copy(downloadedFile, response.Body)
}

func (client *Client) DeleteVersion(packageName, version string) error {
	req, _ := http.NewRequest("DELETE", client.deleteVersionURL(packageName, version), nil)
	req.SetBasicAuth(client.username, client.password)
	c := &http.Client{}
	_, err := c.Do(req)
	return err
}

func (client *Client) UploadPackage(packageName, from, version string) error {
	file, _ := os.Open(from)
	defer file.Close()
	fileStat, _ := file.Stat()

	req, _ := http.NewRequest("PUT", client.outPackageURL(packageName, version), file)
	req.ContentLength = int64(fileStat.Size())
	req.SetBasicAuth(client.username, client.password)
	c := &http.Client{}
	_, err := c.Do(req)
	return err
}

func (client *Client) getPackageURL(packageName string) string {
	getPackagePath := path.Join("packages", client.subjectName, client.repoName, packageName)
	return fmt.Sprintf("%s/%s", client.url, getPackagePath)
}

func (client *Client) inPackageURL(packageName, version string) string {
	downloadPackagePath := path.Join(client.subjectName, client.repoName, version, packageName)
	return fmt.Sprintf("%s/%s", client.url, downloadPackagePath)
}

func (client *Client) deleteVersionURL(packageName, version string) string {
	deleteVersionPath := path.Join("packages", client.subjectName, client.repoName, packageName, "versions", version)
	return fmt.Sprintf("%s/%s", client.url, deleteVersionPath)
}

func (client *Client) outPackageURL(packageName, version string) string {
	uploadPackagePath := path.Join("content", client.subjectName, client.repoName, packageName, version, version+"/"+packageName)
	return fmt.Sprintf("%s/%s?publish=1", client.url, uploadPackagePath)
}
