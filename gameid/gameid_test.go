package gameid_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ajmandourah/tinshop/gameid"
)

var _ = Describe("Gameid", func() {
	It("Ensure creation of gameID", func() {
		game := gameid.New("shortID", "fullID", "extension")

		Expect(game.ShortID()).To(Equal("shortID"))
		Expect(game.FullID()).To(Equal("fullID"))
		Expect(game.Extension()).To(Equal("extension"))
	})
})
