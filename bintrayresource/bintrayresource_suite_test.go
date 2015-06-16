package bintrayresource_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestBintrayresource(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bintrayresource Suite")
}
