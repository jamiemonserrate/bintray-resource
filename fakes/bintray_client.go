package fakes

import "github.com/jamiemonserrate/bintray-resource/bintray"

type BintrayClient struct {
	PackageNameRequested string
	VersionsToReturn     []string
}

func (fakeBintrayClient *BintrayClient) GetPackage(packageName string) bintray.Package {
	fakeBintrayClient.PackageNameRequested = packageName
	return bintray.Package{RawVersions: fakeBintrayClient.VersionsToReturn}
}
