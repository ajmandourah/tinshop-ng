package main_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	main "github.com/DblK/tinshop"
	"github.com/DblK/tinshop/mock_repository"
	"github.com/DblK/tinshop/repository"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Security", func() {
	Describe("TinfoilMiddleware", func() {
		var (
			req              *http.Request
			handler          http.Handler
			writer           *httptest.ResponseRecorder
			myMockCollection *mock_repository.MockCollection
			myMockSources    *mock_repository.MockSources
			myMockConfig     *mock_repository.MockConfig
			ctrl             *gomock.Controller
			myShop           *main.TinShop
		)

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			myMockCollection = mock_repository.NewMockCollection(ctrl)
			myMockSources = mock_repository.NewMockSources(ctrl)
			myMockConfig = mock_repository.NewMockConfig(ctrl)
			myShop = &main.TinShop{}
		})

		JustBeforeEach(func() {
			myShop.Shop = repository.Shop{}
			myShop.Shop.Config = myMockConfig
			myShop.Shop.Collection = myMockCollection
			myShop.Shop.Sources = myMockSources
		})

		Context("No security", func() {
			BeforeEach(func() {
				r := mux.NewRouter()
				r.Use(myShop.TinfoilMiddleware)
				r.HandleFunc("/", myShop.HomeHandler)
				handler = r
				req = httptest.NewRequest(http.MethodGet, "/", nil)
				writer = httptest.NewRecorder()
			})
			It("without any headers", func() {
				emptyCollection := &repository.GameType{}

				myMockCollection.EXPECT().
					Games().
					Return(*emptyCollection).
					AnyTimes()

				myMockConfig.EXPECT().
					DebugNoSecurity().
					Return(true).
					AnyTimes()

				handler.ServeHTTP(writer, req)

				Expect(writer.Code).To(Equal(http.StatusOK))

				var list repository.GameType
				err := json.NewDecoder(writer.Body).Decode(&list)

				Expect(err).To(BeNil())
				Expect(list.Files).To(HaveLen(0))
				Expect(list.ThemeBlackList).To(BeNil())
				Expect(list.Success).To(BeEmpty())
				Expect(list.Titledb).To(HaveLen(0))
			})
		})
		Context("With security", func() {
			BeforeEach(func() {
				r := mux.NewRouter()
				r.Use(myShop.TinfoilMiddleware)
				r.HandleFunc("/", myShop.HomeHandler)
				r.HandleFunc("/{filter}", myShop.HomeHandler)  // Testing purpose
				r.HandleFunc("/{filter}/", myShop.HomeHandler) // Testing purpose
				handler = r
			})
			DescribeTable("test for blacklisted switch", func(path string, valid bool) {
				req = httptest.NewRequest(http.MethodGet, "/"+path, nil)
				writer = httptest.NewRecorder()

				emptyCollection := &repository.GameType{}

				myMockCollection.EXPECT().
					Games().
					Return(*emptyCollection).
					AnyTimes()

				myMockConfig.EXPECT().
					DebugNoSecurity().
					Return(false).
					AnyTimes()

				myMockConfig.EXPECT().
					IsBlacklisted(gomock.Any()).
					Return(true).
					AnyTimes()

				shopTemplateData := &repository.ShopTemplate{
					ShopTitle: "Unit Test",
				}
				myMockConfig.EXPECT().
					ShopTemplateData().
					Return(*shopTemplateData).
					AnyTimes()

				handler.ServeHTTP(writer, req)

				if !valid {
					Expect(writer.Code).To(Equal(http.StatusOK))
					var list repository.GameType
					err := json.NewDecoder(writer.Body).Decode(&list)
					Expect(err).To(BeNil())
					return
				}

				Expect(writer.Code).To(Equal(http.StatusOK))

				var list repository.GameType
				err := json.NewDecoder(writer.Body).Decode(&list)
				Expect(err).NotTo(BeNil())
			},
				Entry("Root path", "", true),
				Entry("'world' path", "world", true),
				Entry("'world/' path", "world/", true),
				Entry("'multi' path", "multi", true),
				Entry("'multi/' path", "multi/", true),
				Entry("'fr' path", "fr", true),
				Entry("'fr/' path", "fr/", true),
				Entry("'dblk/' path", "dblk/", false),
			)
			DescribeTable("test for banned theme switch", func(path string, valid bool) {
				req = httptest.NewRequest(http.MethodGet, "/"+path, nil)
				writer = httptest.NewRecorder()

				emptyCollection := &repository.GameType{}

				myMockCollection.EXPECT().
					Games().
					Return(*emptyCollection).
					AnyTimes()

				myMockConfig.EXPECT().
					DebugNoSecurity().
					Return(false).
					AnyTimes()

				myMockConfig.EXPECT().
					IsBlacklisted(gomock.Any()).
					Return(false).
					AnyTimes()

				myMockConfig.EXPECT().
					IsBannedTheme(gomock.Any()).
					Return(true).
					AnyTimes()

				shopTemplateData := &repository.ShopTemplate{
					ShopTitle: "Unit Test",
				}
				myMockConfig.EXPECT().
					ShopTemplateData().
					Return(*shopTemplateData).
					AnyTimes()

				handler.ServeHTTP(writer, req)

				if !valid {
					Expect(writer.Code).To(Equal(http.StatusOK))
					var list repository.GameType
					err := json.NewDecoder(writer.Body).Decode(&list)
					Expect(err).To(BeNil())
					return
				}

				Expect(writer.Code).To(Equal(http.StatusOK))

				var list repository.GameType
				err := json.NewDecoder(writer.Body).Decode(&list)
				Expect(err).NotTo(BeNil())
			},
				Entry("Root path", "", true),
				Entry("'world' path", "world", true),
				Entry("'world/' path", "world/", true),
				Entry("'multi' path", "multi", true),
				Entry("'multi/' path", "multi/", true),
				Entry("'fr' path", "fr", true),
				Entry("'fr/' path", "fr/", true),
				Entry("'dblk/' path", "dblk/", false),
			)
			DescribeTable("test for an existing user agent", func(path string, valid bool) {
				req = httptest.NewRequest(http.MethodGet, "/"+path, nil)
				writer = httptest.NewRecorder()

				emptyCollection := &repository.GameType{}

				myMockCollection.EXPECT().
					Games().
					Return(*emptyCollection).
					AnyTimes()

				myMockConfig.EXPECT().
					DebugNoSecurity().
					Return(false).
					AnyTimes()

				myMockConfig.EXPECT().
					IsBlacklisted(gomock.Any()).
					Return(false).
					AnyTimes()

				myMockConfig.EXPECT().
					IsBannedTheme(gomock.Any()).
					Return(false).
					AnyTimes()

				shopTemplateData := &repository.ShopTemplate{
					ShopTitle: "Unit Test",
				}
				myMockConfig.EXPECT().
					ShopTemplateData().
					Return(*shopTemplateData).
					AnyTimes()

				req.Header.Set("User-Agent", "Tinshop testing!")
				handler.ServeHTTP(writer, req)

				if !valid {
					Expect(writer.Code).To(Equal(http.StatusOK))
					var list repository.GameType
					err := json.NewDecoder(writer.Body).Decode(&list)
					Expect(err).To(BeNil())
					return
				}

				Expect(writer.Code).To(Equal(http.StatusOK))

				var list repository.GameType
				err := json.NewDecoder(writer.Body).Decode(&list)
				Expect(err).NotTo(BeNil())
			},
				Entry("Root path", "", true),
				Entry("'world' path", "world", true),
				Entry("'world/' path", "world/", true),
				Entry("'multi' path", "multi", true),
				Entry("'multi/' path", "multi/", true),
				Entry("'fr' path", "fr", true),
				Entry("'fr/' path", "fr/", true),
				Entry("'dblk/' path", "dblk/", false),
			)
			DescribeTable("test for with missing mandatory headers", func(path string, valid bool) {
				req = httptest.NewRequest(http.MethodGet, "/"+path, nil)
				writer = httptest.NewRecorder()

				emptyCollection := &repository.GameType{}

				myMockCollection.EXPECT().
					Games().
					Return(*emptyCollection).
					AnyTimes()

				myMockConfig.EXPECT().
					DebugNoSecurity().
					Return(false).
					AnyTimes()

				myMockConfig.EXPECT().
					IsBlacklisted(gomock.Any()).
					Return(false).
					AnyTimes()

				myMockConfig.EXPECT().
					IsBannedTheme(gomock.Any()).
					Return(false).
					AnyTimes()

				shopTemplateData := &repository.ShopTemplate{
					ShopTitle: "Unit Test",
				}
				myMockConfig.EXPECT().
					ShopTemplateData().
					Return(*shopTemplateData).
					AnyTimes()

				handler.ServeHTTP(writer, req)

				if !valid {
					Expect(writer.Code).To(Equal(http.StatusOK))
					var list repository.GameType
					err := json.NewDecoder(writer.Body).Decode(&list)
					Expect(err).To(BeNil())
					return
				}

				Expect(writer.Code).To(Equal(http.StatusOK))

				var list repository.GameType
				err := json.NewDecoder(writer.Body).Decode(&list)
				Expect(err).NotTo(BeNil())
			},
				Entry("Root path", "", true),
				Entry("'world' path", "world", true),
				Entry("'world/' path", "world/", true),
				Entry("'multi' path", "multi", true),
				Entry("'multi/' path", "multi/", true),
				Entry("'fr' path", "fr", true),
				Entry("'fr/' path", "fr/", true),
				Entry("'dblk/' path", "dblk/", false),
			)
			DescribeTable("test with all mandatory headers", func(path string, valid bool) {
				req = httptest.NewRequest(http.MethodGet, "/"+path, nil)
				writer = httptest.NewRecorder()

				emptyCollection := &repository.GameType{}

				myMockCollection.EXPECT().
					Games().
					Return(*emptyCollection).
					AnyTimes()

				myMockConfig.EXPECT().
					DebugNoSecurity().
					Return(false).
					AnyTimes()

				myMockConfig.EXPECT().
					IsBlacklisted(gomock.Any()).
					Return(false).
					AnyTimes()

				myMockConfig.EXPECT().
					IsBannedTheme(gomock.Any()).
					Return(false).
					AnyTimes()

				shopTemplateData := &repository.ShopTemplate{
					ShopTitle: "Unit Test",
				}
				myMockConfig.EXPECT().
					ShopTemplateData().
					Return(*shopTemplateData).
					AnyTimes()

				req.Header.Set("Theme", "Tinshop testing!")
				req.Header.Set("Uid", "Tinshop testing!")
				req.Header.Set("Version", "13.0")
				req.Header.Set("Language", "FR")
				req.Header.Set("Hauth", "XX")
				req.Header.Set("Uauth", "XX")
				handler.ServeHTTP(writer, req)

				if !valid {
					Expect(writer.Code).To(Equal(http.StatusOK))
					var list repository.GameType
					err := json.NewDecoder(writer.Body).Decode(&list)
					Expect(err).To(BeNil())
					return
				}

				Expect(writer.Code).To(Equal(http.StatusOK))

				var list repository.GameType
				err := json.NewDecoder(writer.Body).Decode(&list)

				Expect(err).To(BeNil())
				Expect(list.Files).To(HaveLen(0))
				Expect(list.ThemeBlackList).To(BeNil())
				Expect(list.Success).To(BeEmpty())
				Expect(list.Titledb).To(HaveLen(0))
			},
				Entry("Root path", "", true),
				Entry("'world' path", "world", true),
				Entry("'world/' path", "world/", true),
				Entry("'multi' path", "multi", true),
				Entry("'multi/' path", "multi/", true),
				Entry("'fr' path", "fr", true),
				Entry("'fr/' path", "fr/", true),
				Entry("'dblk/' path", "dblk/", false),
			)
		})
	})
})
