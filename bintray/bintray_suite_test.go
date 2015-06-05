package bintray_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestBintray(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bintray Suite")
}
