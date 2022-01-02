package gameid_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGameid(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gameid Suite")
}
