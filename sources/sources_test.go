package sources_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/DblK/tinshop/repository"
	"github.com/DblK/tinshop/sources"
)

var _ = Describe("Sources", func() {
	var allSources repository.Sources
	BeforeEach(func() {
		allSources = sources.New(nil)
	})
	It("Return list of game files", func() {
		files := allSources.GetFiles()

		Expect(len(files)).To(Equal(0))
	})
	// It("Should add files into the current list", func() {
	// 	newFile := repository.FileDesc{Size: 42, Path: "/somewhere/here"}
	// 	newfiles := make([]repository.FileDesc, 0)
	// 	newfiles = append(newfiles, newFile)
	// 	sources.AddFiles(newfiles)

	// 	Expect(len(sources.GetFiles())).To(Equal(1))
	// 	firstFile := sources.GetFiles()[0]
	// 	Expect(firstFile.GameID).To(BeEmpty())
	// 	Expect(firstFile.GameInfo).To(BeEmpty())
	// 	Expect(firstFile.HostType).To(BeEmpty())
	// 	Expect(firstFile.Path).To(Equal("/somewhere/here"))
	// 	Expect(firstFile.Size).To(Equal(int64(42)))

	// })
})
