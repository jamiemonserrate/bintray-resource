package fakes

import "github.com/jamiemonserrate/bintray-resource/bintray"

type BintrayClient struct {
	PackageNameRequested  string
	VersionsToReturn      []string
	LatestVersionToReturn string
}

func (fakeBintrayClient *BintrayClient) GetPackage(packageName string) bintray.Package {
	fakeBintrayClient.PackageNameRequested = packageName
	return bintray.Package{LatestVersion: fakeBintrayClient.LatestVersionToReturn,
		Versions: fakeBintrayClient.VersionsToReturn}
}
