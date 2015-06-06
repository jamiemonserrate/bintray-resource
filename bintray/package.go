package bintray

import (
	"sort"

	"github.com/hashicorp/go-version"
)

type Package struct {
	RawVersions []string `json:"versions"`
}

func (bintrayPackage *Package) VersionsSince(minVersion *version.Version) []*version.Version {
	versionsToReturn := []*version.Version{}
	for _, rawVersion := range bintrayPackage.RawVersions {
		v, _ := version.NewVersion(rawVersion)
		if v.GreaterThan(minVersion) {
			versionsToReturn = append(versionsToReturn, v)
		}
	}

	sort.Sort(sort.Reverse(version.Collection(versionsToReturn)))
	return versionsToReturn
}
