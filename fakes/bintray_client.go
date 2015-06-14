package fakes

import "github.com/jamiemonserrate/bintray-resource/bintray"

type BintrayClient struct {
	PackageNameRequested    string
	VersionsToReturn        []string
	VersionRequested        string
	DestinationDirRequested string
	FileToBeUploaded        string
	VersionToBeUploaded     string
}

func (fakeBintrayClient *BintrayClient) GetPackage(packageName string) bintray.Package {
	fakeBintrayClient.PackageNameRequested = packageName
	return bintray.Package{RawVersions: fakeBintrayClient.VersionsToReturn}
}

func (fakeBintrayClient *BintrayClient) DownloadPackage(packageName, version, destinationDir string) {
	fakeBintrayClient.PackageNameRequested = packageName
	fakeBintrayClient.VersionRequested = version
	fakeBintrayClient.DestinationDirRequested = destinationDir
}

func (fakeBintrayClient *BintrayClient) UploadPackage(packageName, fileToBeUploaded, version string) error {
	fakeBintrayClient.PackageNameRequested = packageName
	fakeBintrayClient.FileToBeUploaded = fileToBeUploaded
	fakeBintrayClient.VersionToBeUploaded = version
	return nil
}
