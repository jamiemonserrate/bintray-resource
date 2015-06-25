package bintray

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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

type ErrorResponse struct {
	Message string `json:"message"`
}

type BintrayClient interface {
	GetPackage(packageName string) (*Package, error)
	DownloadPackage(packageName, version, destinationDir string) error
	UploadPackage(packageName, from, version string) error
	InPackageURL(packageName, version string) string
}

func NewClient(bintrayURL, subjectName, repoName, username, password string) *Client {
	return &Client{url: bintrayURL, subjectName: subjectName, repoName: repoName,
		username: username, password: password}
}

func (client *Client) GetPackage(packageName string) (*Package, error) {
	response, err := http.Get(client.getPackageURL(packageName))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, client.parseErrorFrom(response)
	}

	var bintrayPackage *Package
	responseBytes, _ := ioutil.ReadAll(response.Body)
	if err := json.NewDecoder(bytes.NewReader(responseBytes)).Decode(&bintrayPackage); err != nil || bintrayPackage.RawLatestVersion == "" {
		errorResponse := ErrorResponse{}
		if err = json.NewDecoder(bytes.NewReader(responseBytes)).Decode(&errorResponse); err != nil {
			return nil, errors.New(string(responseBytes))
		}
		return nil, errors.New(errorResponse.Message)
	}

	return bintrayPackage, nil
}

func (client *Client) DownloadPackage(packageName, version, destinationDir string) error {
	downloadedFile, err := os.Create(filepath.Join(destinationDir, packageName))
	if err != nil {
		return err
	}

	defer downloadedFile.Close()
	response, err := http.Get(client.InPackageURL(packageName, version))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return client.parseErrorFrom(response)
	}
	_, err = io.Copy(downloadedFile, response.Body)
	return err
}

func (client *Client) DeleteVersion(packageName, version string) error {
	request, err := http.NewRequest("DELETE", client.deleteVersionURL(packageName, version), nil)
	if err != nil {
		return err
	}
	request.SetBasicAuth(client.username, client.password)
	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return client.parseErrorFrom(response)
	}
	return nil
}

func (client *Client) UploadPackage(packageName, from, version string) error {
	file, err := os.Open(from)
	if err != nil {
		return err
	}
	defer file.Close()
	fileStat, err := file.Stat()
	if err != nil {
		return err
	}

	request, err := http.NewRequest("PUT", client.outPackageURL(packageName, version), file)
	if err != nil {
		return err
	}
	request.ContentLength = int64(fileStat.Size())
	request.SetBasicAuth(client.username, client.password)
	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusCreated {
		return client.parseErrorFrom(response)
	}
	return nil
}

func (client *Client) getPackageURL(packageName string) string {
	getPackagePath := path.Join("packages", client.subjectName, client.repoName, packageName)
	return fmt.Sprintf("%s/%s", client.url, getPackagePath)
}

func (client *Client) InPackageURL(packageName, version string) string {
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

func (client *Client) parseErrorFrom(response *http.Response) error {
	errorResponse := ErrorResponse{}
	errorString, _ := ioutil.ReadAll(response.Body)
	if err := json.NewDecoder(bytes.NewReader(errorString)).Decode(&errorResponse); err != nil {
		return errors.New(string(errorString))
	}
	return errors.New(errorResponse.Message)
}
