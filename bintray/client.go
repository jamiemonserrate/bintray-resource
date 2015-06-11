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
	UploadPackage(string, string, string) error
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
	// req, err := http.NewRequest("GET", client.inPackageURL(packageName, version), nil)
	if err != nil {
		panic(err)
	}
	// apikey := "9dd0d7a78b11e773ef4dbc389cf36c1cfe536ebc"
	// req.SetBasicAuth("jamiemonserrate", apikey)
	// hclient := &http.Client{}
	// response, err := hclient.Do(req)
	defer response.Body.Close()

	_, err = io.Copy(downloadedFile, response.Body)

	if err != nil {
		panic(err)
	}
}

func (client *Client) DeleteVersion(packageName, version string) {
	req, err := http.NewRequest("DELETE", client.deleteVersionURL(packageName, version), nil)
	if err != nil {
		panic(err)
	}
	apikey := "9dd0d7a78b11e773ef4dbc389cf36c1cfe536ebc"
	req.SetBasicAuth("jamiemonserrate", apikey)
	hclient := &http.Client{}
	_, err = hclient.Do(req)
	// defer response.Body.Close()

	if err != nil {
		panic(err)
	}
}

func (client *Client) UploadPackage(packageName, from, version string) error {
	fullPath, _ := filepath.Abs(from)
	file, err := os.Open(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		return err
	}

	req, err := client.newRequestWithReader("PUT", client.outPackageURL(packageName, version), file, fi.Size())
	if err != nil {
		return err
	}
	err = client.execute(req)
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

func (c *Client) newRequestWithReader(method, urlStr string, requestReader io.Reader, requestLength int64) (*http.Request, error) {
	req, err := http.NewRequest(method, urlStr, requestReader)
	if err != nil {
		return nil, err
	}
	if requestLength > 0 {
		req.ContentLength = int64(requestLength)
	}
	return req, nil
}

func (c *Client) execute(req *http.Request) error {
	client := &http.Client{}

	apikey := "9dd0d7a78b11e773ef4dbc389cf36c1cfe536ebc"
	req.SetBasicAuth("jamiemonserrate", apikey)
	_, err := client.Do(req)
	if err != nil {
		return err
	}
	return nil
}
