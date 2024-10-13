package main_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	main "github.com/ajmandourah/tinshop-ng"
	"github.com/ajmandourah/tinshop-ng/mock_repository"
	"github.com/ajmandourah/tinshop-ng/repository"
)

var _ = Describe("Main", func() {
	Describe("HomeHandler", func() {
		var (
			req              *http.Request
			handler          http.Handler
			writer           *httptest.ResponseRecorder
			myMockCollection *mock_repository.MockCollection
			myMockSources    *mock_repository.MockSources
			myMockConfig     *mock_repository.MockConfig
			myMockStats      *mock_repository.MockStats
			ctrl             *gomock.Controller
			myShop           *main.TinShop
		)

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			myMockCollection = mock_repository.NewMockCollection(ctrl)
			myMockSources = mock_repository.NewMockSources(ctrl)
			myMockConfig = mock_repository.NewMockConfig(ctrl)
			myMockStats = mock_repository.NewMockStats(ctrl)
			myShop = &main.TinShop{}
		})

		JustBeforeEach(func() {
			myShop.Shop = repository.Shop{}
			myShop.Shop.Config = myMockConfig
			myShop.Shop.Collection = myMockCollection
			myShop.Shop.Sources = myMockSources
			myShop.Shop.Stats = myMockStats
		})

		Context("With empty collection", func() {
			BeforeEach(func() {
				handler = http.HandlerFunc(myShop.HomeHandler)
				req = httptest.NewRequest(http.MethodGet, "/", nil)
				writer = httptest.NewRecorder()
			})
			It("Verify response without data", func() {
				myShop.Shop.Collection = nil
				handler.ServeHTTP(writer, req)
				Expect(writer.Code).To(Equal(http.StatusNotFound))
			})
			It("Verify empty response", func() {
				emptyCollection := &repository.GameType{}

				myMockCollection.EXPECT().
					Games().
					Return(*emptyCollection).
					AnyTimes()

				// tinshop.ResetTinshop(myShop)
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
		Context("With collection", func() {
			JustBeforeEach(func() {
				handler = http.HandlerFunc(myShop.HomeHandler)
				req = httptest.NewRequest(http.MethodGet, "/", nil)
				writer = httptest.NewRecorder()
				// tinshop.ResetTinshop(myShop)
			})
			It("Verify status code", func() {
				smallCollection := &repository.GameType{}
				smallCollection.Files = make([]repository.GameFileType, 0)
				oneFile := &repository.GameFileType{
					Size: 42,
					URL:  "http://test.tinshop.io",
				}
				smallCollection.Files = append(smallCollection.Files, *oneFile)
				smallCollection.Success = "Welcome to your own shop!"
				smallCollection.Titledb = make(map[string]repository.TitleDBEntry)
				oneEntry := &repository.TitleDBEntry{
					ID: "0000000000000001",
				}
				smallCollection.Titledb["0000000000000001"] = *oneEntry

				myMockCollection.EXPECT().
					Games().
					Return(*smallCollection).
					AnyTimes()

				handler.ServeHTTP(writer, req)
				Expect(writer.Code).To(Equal(http.StatusOK))

				var list repository.GameType
				err := json.NewDecoder(writer.Body).Decode(&list)

				Expect(err).To(BeNil())
				Expect(list.Files).To(HaveLen(1))
				Expect(list.ThemeBlackList).To(BeNil())
				Expect(list.Success).To(Equal("Welcome to your own shop!"))
				Expect(list.Titledb).To(HaveLen(1))
			})
		})
	})
	Describe("FilteringHandler", func() {
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

		Context("With empty collection", func() {
			Context("Filtering with  path", func() {
				BeforeEach(func() {
					r := mux.NewRouter()
					r.HandleFunc("/{filter}", myShop.FilteringHandler)
					r.HandleFunc("/{filter}/", myShop.FilteringHandler)
					handler = r
				})
				DescribeTable("Verify response without data", func(path string, valid bool) {
					req = httptest.NewRequest(http.MethodGet, "/"+path, nil)
					writer = httptest.NewRecorder()
					myShop.Shop.Collection = nil
					handler.ServeHTTP(writer, req)
					if !valid {
						Expect(writer.Code).To(Equal(http.StatusNotAcceptable))
						return
					}
					Expect(writer.Code).To(Equal(http.StatusNotFound))
				},
					Entry("with path 'world'", "world", true),
					Entry("with path 'world/'", "world/", true),
					Entry("with path 'multi'", "multi", true),
					Entry("with path 'multi/'", "multi/", true),
					Entry("with path 'fr'", "fr", true),
					Entry("with path 'fr/'", "fr/", true),
					Entry("with path 'us'", "us", true),
					Entry("with path 'us/'", "us/", true),
					Entry("with path 'dblk'", "dblk", false),
				)
				DescribeTable("Verify empty response", func(path string, valid bool) {
					req = httptest.NewRequest(http.MethodGet, "/"+path, nil)
					writer = httptest.NewRecorder()
					emptyCollection := &repository.GameType{}

					myMockCollection.EXPECT().
						Games().
						Return(*emptyCollection).
						AnyTimes()
					myMockCollection.EXPECT().
						Filter(gomock.Any()).
						Return(*emptyCollection).
						AnyTimes()

					handler.ServeHTTP(writer, req)
					if !valid {
						Expect(writer.Code).To(Equal(http.StatusNotAcceptable))
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
					Entry("with path 'world'", "world", true),
					Entry("with path 'world/'", "world/", true),
					Entry("with path 'multi'", "multi", true),
					Entry("with path 'multi/'", "multi/", true),
					Entry("with path 'fr'", "fr", true),
					Entry("with path 'fr/'", "fr/", true),
					Entry("with path 'us'", "us", true),
					Entry("with path 'us/'", "us/", true),
					Entry("with path 'dblk'", "dblk", false),
				)
			})
		})
		Context("With collection", func() {
			Context("Filtering with path", func() {
				JustBeforeEach(func() {
					r := mux.NewRouter()
					r.HandleFunc("/{filter}", myShop.FilteringHandler)
					r.HandleFunc("/{filter}/", myShop.FilteringHandler)
					handler = r
					req = httptest.NewRequest(http.MethodGet, "/world", nil)
					writer = httptest.NewRecorder()
				})
				DescribeTable("Verify status code", func(path string, valid bool) {
					req = httptest.NewRequest(http.MethodGet, "/"+path, nil)
					writer = httptest.NewRecorder()
					smallCollection := &repository.GameType{}
					smallCollection.Files = make([]repository.GameFileType, 0)
					oneFile := &repository.GameFileType{
						Size: 42,
						URL:  "http://test.tinshop.io",
					}
					smallCollection.Files = append(smallCollection.Files, *oneFile)
					smallCollection.Success = "Welcome to your own shop!"
					smallCollection.Titledb = make(map[string]repository.TitleDBEntry)
					oneEntry := &repository.TitleDBEntry{
						ID: "0000000000000001",
					}
					smallCollection.Titledb["0000000000000001"] = *oneEntry

					myMockCollection.EXPECT().
						Filter(gomock.Any()).
						Return(*smallCollection).
						AnyTimes()

					handler.ServeHTTP(writer, req)
					if !valid {
						Expect(writer.Code).To(Equal(http.StatusNotAcceptable))
						return
					}
					Expect(writer.Code).To(Equal(http.StatusOK))

					var list repository.GameType
					err := json.NewDecoder(writer.Body).Decode(&list)

					Expect(err).To(BeNil())
					Expect(list.Files).To(HaveLen(1))
					Expect(list.ThemeBlackList).To(BeNil())
					Expect(list.Success).To(Equal("Welcome to your own shop!"))
					Expect(list.Titledb).To(HaveLen(1))
				},
					Entry("with path 'world'", "world", true),
					Entry("with path 'world/'", "world/", true),
					Entry("with path 'multi'", "multi", true),
					Entry("with path 'multi/'", "multi/", true),
					Entry("with path 'fr'", "fr", true),
					Entry("with path 'fr/'", "fr/", true),
					Entry("with path 'us'", "us", true),
					Entry("with path 'us/'", "us/", true),
					Entry("with path 'dblk'", "dblk", false),
				)
			})
		})
	})
	Describe("TinfoilMiddleware", func() {
		var (
			req              *http.Request
			handler          http.Handler
			writer           *httptest.ResponseRecorder
			myMockCollection *mock_repository.MockCollection
			myMockSources    *mock_repository.MockSources
			myMockConfig     *mock_repository.MockConfig
			myMockStats      *mock_repository.MockStats
			ctrl             *gomock.Controller
			myShop           *main.TinShop
		)

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			myMockCollection = mock_repository.NewMockCollection(ctrl)
			myMockSources = mock_repository.NewMockSources(ctrl)
			myMockConfig = mock_repository.NewMockConfig(ctrl)
			myMockStats = mock_repository.NewMockStats(ctrl)
			myShop = &main.TinShop{}
		})

		JustBeforeEach(func() {
			myShop.Shop = repository.Shop{}
			myShop.Shop.Config = myMockConfig
			myShop.Shop.Collection = myMockCollection
			myShop.Shop.Sources = myMockSources
			myShop.Shop.Stats = myMockStats
		})
		Context("Not handled endpoint", func() {
			BeforeEach(func() {
				r := mux.NewRouter()
				r.Use(myShop.StatsMiddleware)
				r.HandleFunc("/api/{endpoint}", myShop.HomeHandler) // Testing purpose
				handler = r
			})

			It("Test with the api endpoint", func() {
				req = httptest.NewRequest(http.MethodGet, "/api/stats", nil)
				writer = httptest.NewRecorder()

				emptyCollection := &repository.GameType{}

				myMockCollection.EXPECT().
					Games().
					Return(*emptyCollection).
					AnyTimes()

				myMockSources.EXPECT().
					HasGame(gomock.Any()).
					Return(true).
					Times(0)
				myMockStats.EXPECT().
					ListVisit(gomock.Any()).
					Return(nil).
					Times(0)
				myMockStats.EXPECT().
					DownloadAsked(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(0)

				handler.ServeHTTP(writer, req)
				Expect(writer.Code).To(Equal(http.StatusOK))
			})
		})
		Context("Games endpoint", func() {
			BeforeEach(func() {
				r := mux.NewRouter()
				r.Use(myShop.StatsMiddleware)
				r.HandleFunc("/games/{game}", myShop.HomeHandler) // Testing purpose
				handler = r
			})

			It("Test with a not found game", func() {
				req = httptest.NewRequest(http.MethodGet, "/games/notFound", nil)
				writer = httptest.NewRecorder()

				emptyCollection := &repository.GameType{}

				myMockCollection.EXPECT().
					Games().
					Return(*emptyCollection).
					AnyTimes()

				myMockSources.EXPECT().
					HasGame("notFound").
					Return(false).
					Times(1)
				myMockStats.EXPECT().
					ListVisit(gomock.Any()).
					Return(nil).
					Times(0)
				myMockStats.EXPECT().
					DownloadAsked(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(0)

				handler.ServeHTTP(writer, req)
				Expect(writer.Code).To(Equal(http.StatusOK))
			})
			It("Test with a found game", func() {
				req = httptest.NewRequest(http.MethodGet, "/games/existingGame", nil)
				req.RemoteAddr = "10.0.0.10"
				writer = httptest.NewRecorder()

				emptyCollection := &repository.GameType{}

				myMockCollection.EXPECT().
					Games().
					Return(*emptyCollection).
					AnyTimes()

				myMockSources.EXPECT().
					HasGame("existingGame").
					Return(true).
					Times(1)
				myMockStats.EXPECT().
					ListVisit(gomock.Any()).
					Return(nil).
					Times(0)
				myMockStats.EXPECT().
					DownloadAsked("10.0.0.10", "existingGame").
					Return(nil).
					Times(1)

				handler.ServeHTTP(writer, req)
				Expect(writer.Code).To(Equal(http.StatusOK))
			})
		})
		Context("Listing endpoint", func() {
			BeforeEach(func() {
				r := mux.NewRouter()
				r.Use(myShop.StatsMiddleware)
				r.HandleFunc("/", myShop.HomeHandler)
				r.HandleFunc("/{filter}", myShop.HomeHandler)  // Testing purpose
				r.HandleFunc("/{filter}/", myShop.HomeHandler) // Testing purpose
				handler = r
			})

			It("Test with root endpoint", func() {
				req = httptest.NewRequest(http.MethodGet, "/", nil)
				writer = httptest.NewRecorder()

				emptyCollection := &repository.GameType{}

				myMockCollection.EXPECT().
					Games().
					Return(*emptyCollection).
					AnyTimes()

				myMockSources.EXPECT().
					HasGame(gomock.Any()).
					Return(false).
					Times(0)
				myMockStats.EXPECT().
					ListVisit(gomock.Any()).
					Return(nil).
					Times(1)
				myMockStats.EXPECT().
					DownloadAsked(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(0)

				handler.ServeHTTP(writer, req)
				Expect(writer.Code).To(Equal(http.StatusOK))
			})
			It("Test with a filter endpoint", func() {
				req = httptest.NewRequest(http.MethodGet, "/FR", nil)
				writer = httptest.NewRecorder()

				emptyCollection := &repository.GameType{}

				myMockCollection.EXPECT().
					Games().
					Return(*emptyCollection).
					AnyTimes()

				myMockSources.EXPECT().
					HasGame(gomock.Any()).
					Return(true).
					Times(0)
				myMockStats.EXPECT().
					ListVisit(gomock.Any()).
					Return(nil).
					Times(1)
				myMockStats.EXPECT().
					DownloadAsked(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(0)

				handler.ServeHTTP(writer, req)
				Expect(writer.Code).To(Equal(http.StatusOK))
			})
		})
	})
})
