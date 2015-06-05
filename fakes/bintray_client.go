package fakes

import "github.com/jamiemonserrate/bintray-resource/bintray"

type BintrayClient struct {
	PackageNameInvokedWith string
	VersionsToReturn       []string
	LatestVersionToReturn  string
}

func (fakeBintrayClient BintrayClient) GetPackage(packageName string) bintray.Package {
	fakeBintrayClient.PackageNameInvokedWith = packageName
	return bintray.Package{LatestVersion: fakeBintrayClient.LatestVersionToReturn,
		Versions: fakeBintrayClient.VersionsToReturn}
}
