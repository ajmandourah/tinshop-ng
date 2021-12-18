package gamescollection_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGamescollection(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gamescollection Suite")
}
