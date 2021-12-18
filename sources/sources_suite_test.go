package sources_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSources(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sources Suite")
}
