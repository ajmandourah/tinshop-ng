package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTinshop(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tinshop Suite")
}
