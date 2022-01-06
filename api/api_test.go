package api_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/DblK/tinshop/api"
	"github.com/DblK/tinshop/repository"
)

var _ = Describe("Api", func() {
	var (
		myAPI  repository.API
		writer *httptest.ResponseRecorder
	)
	BeforeEach(func() {
		myAPI = api.New()
	})
	Describe("Stats", func() {
		It("Test with empty stats", func() {
			emptyStats := &repository.StatsSummary{}
			writer = httptest.NewRecorder()

			myAPI.Stats(writer, *emptyStats)
			Expect(writer.Code).To(Equal(http.StatusOK))
			Expect(writer.Body.String()).To(Equal("{}"))
		})
		It("Test with some stats", func() {
			emptyStats := &repository.StatsSummary{
				Visit: 42,
			}
			writer = httptest.NewRecorder()

			myAPI.Stats(writer, *emptyStats)
			Expect(writer.Code).To(Equal(http.StatusOK))
			Expect(writer.Body.String()).To(Equal("{\"visit\":42}"))
		})
	})
})
