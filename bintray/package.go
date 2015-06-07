package bintray

import (
	"sort"

	"github.com/hashicorp/go-version"
)

type Package struct {
	RawVersions      []string `json:"versions"`
	RawLatestVersion string   `json:"latest_version"`
}

func (bintrayPackage *Package) VersionsSince(minVersion *version.Version) []*version.Version {
	versionsToReturn := []*version.Version{}

	if minVersion == nil {
		return append(versionsToReturn, bintrayPackage.latestVersion())
	}

	for _, rawVersion := range bintrayPackage.RawVersions {
		v, _ := version.NewVersion(rawVersion)
		if v.GreaterThan(minVersion) {
			versionsToReturn = append(versionsToReturn, v)
		}
	}

	sort.Sort(sort.Reverse(version.Collection(versionsToReturn)))
	return versionsToReturn
}

func (bintrayPackage *Package) latestVersion() *version.Version {
	latestVersion, _ := version.NewVersion(bintrayPackage.RawLatestVersion)
	return latestVersion
}
