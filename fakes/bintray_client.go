package fakes

import "github.com/jamiemonserrate/bintray-resource/bintray"

type BintrayClient struct {
	PackageNameRequested    string
	VersionsToReturn        []string
	VersionRequested        string
	DestinationDirRequested string
	FileToBeUploaded        string
	VersionToBeUploaded     string
	ErrorToBeReturned       error
}

func (fakeBintrayClient *BintrayClient) GetPackage(packageName string) (*bintray.Package, error) {
	if fakeBintrayClient.ErrorToBeReturned != nil {
		return nil, fakeBintrayClient.ErrorToBeReturned
	}

	fakeBintrayClient.PackageNameRequested = packageName
	return &bintray.Package{RawVersions: fakeBintrayClient.VersionsToReturn}, nil
}

func (fakeBintrayClient *BintrayClient) DownloadPackage(packageName, version, destinationDir string) error {
	if fakeBintrayClient.ErrorToBeReturned != nil {
		return fakeBintrayClient.ErrorToBeReturned
	}

	fakeBintrayClient.PackageNameRequested = packageName
	fakeBintrayClient.VersionRequested = version
	fakeBintrayClient.DestinationDirRequested = destinationDir
	return nil
}

func (fakeBintrayClient *BintrayClient) UploadPackage(packageName, fileToBeUploaded, version string) error {
	if fakeBintrayClient.ErrorToBeReturned != nil {
		return fakeBintrayClient.ErrorToBeReturned
	}

	fakeBintrayClient.PackageNameRequested = packageName
	fakeBintrayClient.FileToBeUploaded = fileToBeUploaded
	fakeBintrayClient.VersionToBeUploaded = version
	return nil
}
