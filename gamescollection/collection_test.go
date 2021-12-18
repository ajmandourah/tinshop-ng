package gamescollection_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	collection "github.com/dblk/tinshop/gamescollection"
)

var _ = Describe("Collection", func() {
	It("Return list of games", func() {
		games := collection.Games()

		Expect(len(games.Files)).To(Equal(0))
	})
})
