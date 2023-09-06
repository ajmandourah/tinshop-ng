package gamescollection_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	collection "github.com/DblK/tinshop/gamescollection"
	"github.com/DblK/tinshop/mock_repository"
	"github.com/DblK/tinshop/repository"
)

var _ = Describe("Collection", func() {
	var (
		myMockConfig   *mock_repository.MockConfig
		ctrl           *gomock.Controller
		testCollection repository.Collection
	)
	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})
	JustBeforeEach(func() {
		myMockConfig = mock_repository.NewMockConfig(ctrl)

		myMockConfig.EXPECT().
			RootShop().
			Return("http://tinshop.example.com").
			AnyTimes()
		myMockConfig.EXPECT().
			WelcomeMessage().
			Return("Welcome to testing shop!").
			AnyTimes()
		myMockConfig.EXPECT().
			NoWelcomeMessage().
			Return(false).
			AnyTimes()

		testCollection = collection.New(myMockConfig)
	})
	It("Return list of games", func() {
		games := testCollection.Games()

		Expect(games.Files).To(HaveLen(0))
	})
	Describe("AddNewGames", func() {
		JustBeforeEach(func() {
			testCollection.ResetGamesCollection()
		})
		It("Add an empty array", func() {
			newGames := make([]repository.FileDesc, 0)
			testCollection.AddNewGames(newGames)

			games := testCollection.Games()
			Expect(games.Files).To(HaveLen(0))
			Expect(games.Titledb).To(HaveLen(0))
		})
		Context("No TitleDB", func() {
			// TODO: Need to add tests here!
		})
		Context("With TitleDB", func() {
			JustBeforeEach(func() {
				customDB := make(map[string]repository.TitleDBEntry)
				custom1 := repository.TitleDBEntry{ // Base
					ID:              "010034500641A000",
					Languages:       []string{"FR", "EN", "US"},
					Name:            "Attack on Titan 2",
					Region:          "US",
					NumberOfPlayers: 1,
					IconURL:         "http://fake.icon.url",
				}
				custom2 := repository.TitleDBEntry{ // Update
					ID:      "010034500641A800",
					Version: 917504,
				}
				custom3 := repository.TitleDBEntry{ // DLC
					ID:      "010034500641B001",
					Name:    "Additional Episode, \"A Sudden Rain\"",
					Version: 131072,
				}
				custom4 := repository.TitleDBEntry{ // Base (No Region)
					ID:              "0100574002AF4000",
					Languages:       []string{"FR", "EN", "US"},
					Name:            "ONE PIECE: Unlimited World Red Deluxe Edition",
					NumberOfPlayers: 1,
				}
				custom5 := repository.TitleDBEntry{ // Base (No info)
					ID: "010034501225C000",
				}
				customDB["010034500641A000"] = custom1
				customDB["010034500641A800"] = custom2
				customDB["010034500641B001"] = custom3
				customDB["0100574002AF4000"] = custom4
				customDB["010034501225C000"] = custom5
				myMockConfig = mock_repository.NewMockConfig(ctrl)
				myMockConfig.EXPECT().
					CustomDB().
					Return(customDB).
					AnyTimes()
				myMockConfig.EXPECT().
					BannedTheme().
					Return(nil).
					AnyTimes()
				myMockConfig.EXPECT().
					RootShop().
					Return("http://tinshop.example.com").
					AnyTimes()
				myMockConfig.EXPECT().
					WelcomeMessage().
					Return("Welcome to testing shop!").
					AnyTimes()
				myMockConfig.EXPECT().
					NoWelcomeMessage().
					Return(false).
					AnyTimes()

				testCollection.OnConfigUpdate(myMockConfig)
			})
			It("Add a base game", func() {
				newGames := make([]repository.FileDesc, 0)
				newFile := repository.FileDesc{
					Size:      42,
					Path:      "/here/is/my/game",
					GameID:    "010034500641A000",
					GameInfo:  "[010034500641A000][v0].nsp",
					Extension: "nsp",
					HostType:  repository.LocalFile,
				}
				newGames = append(newGames, newFile)
				testCollection.AddNewGames(newGames)

				games := testCollection.Games()
				Expect(games.Files).To(HaveLen(1))
				Expect(games.Titledb).To(HaveLen(1))
				Expect(games.Files[0].URL).To(Equal("http://tinshop.example.com/games/010034500641A000#[010034500641A000] Attack on Titan 2 (US) [BASE].nsp"))
			})
			It("Add a base game (without Region)", func() {
				newGames := make([]repository.FileDesc, 0)
				newFile := repository.FileDesc{
					Size:      42,
					Path:      "/here/is/my/game",
					GameID:    "0100574002AF4000",
					GameInfo:  "[0100574002AF4000][v0].nsp",
					Extension: "nsp",
					HostType:  repository.LocalFile,
				}
				newGames = append(newGames, newFile)
				testCollection.AddNewGames(newGames)

				games := testCollection.Games()
				Expect(games.Files).To(HaveLen(1))
				Expect(games.Titledb).To(HaveLen(1))
				Expect(games.Files[0].URL).To(Equal("http://tinshop.example.com/games/0100574002AF4000#[0100574002AF4000] ONE PIECE: Unlimited World Red Deluxe Edition [BASE].nsp"))
			})
			It("Add a base game (without any information)", func() {
				newGames := make([]repository.FileDesc, 0)
				newFile := repository.FileDesc{
					Size:      42,
					Path:      "/here/is/my/game",
					GameID:    "010034501225C000",
					GameInfo:  "[010034501225C000][v0].nsp",
					Extension: "nsp",
					HostType:  repository.LocalFile,
				}
				newGames = append(newGames, newFile)
				testCollection.AddNewGames(newGames)

				games := testCollection.Games()
				Expect(games.Files).To(HaveLen(1))
				Expect(games.Titledb).To(HaveLen(1))
				Expect(games.Files[0].URL).To(Equal("http://tinshop.example.com/games/010034501225C000#[010034501225C000] [BASE].nsp"))
			})
			It("Add a DLC game", func() {
				newGames := make([]repository.FileDesc, 0)
				newFile := repository.FileDesc{
					Size:      42,
					Path:      "/here/is/my/game",
					GameID:    "010034500641B001",
					GameInfo:  "[010034500641B001][v0].nsp",
					Extension: "nsp",
					HostType:  repository.LocalFile,
				}
				newGames = append(newGames, newFile)
				testCollection.AddNewGames(newGames)

				games := testCollection.Games()
				Expect(games.Files).To(HaveLen(1))
				Expect(games.Titledb).To(HaveLen(1))
				Expect(games.Files[0].URL).To(Equal("http://tinshop.example.com/games/010034500641B001#[010034500641B001] Attack on Titan 2 (US) - Additional Episode, \"A Sudden Rain\" [DLC].nsp"))
			})
			It("Add an UPDATE game", func() {
				newGames := make([]repository.FileDesc, 0)
				newFile := repository.FileDesc{
					Size:      42,
					Path:      "/here/is/my/game",
					GameID:    "010034500641A800",
					GameInfo:  "[010034500641A800][v0].nsp",
					Extension: "nsp",
					HostType:  repository.LocalFile,
				}
				newGames = append(newGames, newFile)
				testCollection.AddNewGames(newGames)

				games := testCollection.Games()
				Expect(games.Files).To(HaveLen(1))
				Expect(games.Titledb).To(HaveLen(1))
				Expect(games.Files[0].URL).To(Equal("http://tinshop.example.com/games/010034500641A800#[010034500641A800] Attack on Titan 2 (US) [v917504][UPD].nsp"))
			})
			It("Add a duplicate game", func() {
				newGames := make([]repository.FileDesc, 0)
				newFile1 := repository.FileDesc{
					Size:      42,
					Path:      "/here/is/my/game",
					GameID:    "010034500641A000",
					GameInfo:  "[010034500641A000][v0].nsp",
					Extension: "nsp",
					HostType:  repository.LocalFile,
				}
				newFile2 := repository.FileDesc{
					Size:      43,
					Path:      "/here/is/my/game",
					GameID:    "010034500641A000",
					GameInfo:  "[010034500641A000][v0].nsp",
					Extension: "nsp",
					HostType:  repository.LocalFile,
				}
				newGames = append(newGames, newFile1)
				newGames = append(newGames, newFile2)
				testCollection.AddNewGames(newGames)

				games := testCollection.Games()
				Expect(games.Files).To(HaveLen(1))
				Expect(games.Titledb).To(HaveLen(1))
				Expect(games.Files[0].URL).To(Equal("http://tinshop.example.com/games/010034500641A000#[010034500641A000] Attack on Titan 2 (US) [BASE].nsp"))
			})
			It("Add a duplicate game (with different path)", func() {
				newGames := make([]repository.FileDesc, 0)
				newFile1 := repository.FileDesc{
					Size:      42,
					Path:      "/here/is/my/game1",
					GameID:    "010034500641A000",
					GameInfo:  "[010034500641A000][v0].nsp",
					Extension: "nsp",
					HostType:  repository.LocalFile,
				}
				newFile2 := repository.FileDesc{
					Size:      43,
					Path:      "/here/is/my/game2",
					GameID:    "010034500641A000",
					GameInfo:  "[010034500641A000][v0].nsp",
					Extension: "nsp",
					HostType:  repository.LocalFile,
				}
				newGames = append(newGames, newFile1)
				newGames = append(newGames, newFile2)
				testCollection.AddNewGames(newGames)

				games := testCollection.Games()
				Expect(games.Files).To(HaveLen(1))
				Expect(games.Titledb).To(HaveLen(1))
				Expect(games.Files[0].URL).To(Equal("http://tinshop.example.com/games/010034500641A000#[010034500641A000] Attack on Titan 2 (US) [BASE].nsp"))
			})
		})
	})
	Describe("RemoveGame", func() {
		JustBeforeEach(func() {
			testCollection.ResetGamesCollection()
		})
		It("Removing existing ID", func() {
			newGames := make([]repository.FileDesc, 0)
			newFile := repository.FileDesc{
				Size:     42,
				Path:     "/here/is/my/game",
				GameID:   "0000000000000001",
				GameInfo: "[0000000000000001][v0].nsp",
				HostType: repository.LocalFile,
			}
			newGames = append(newGames, newFile)
			testCollection.AddNewGames(newGames)

			Expect(testCollection.Games().Files).To(HaveLen(1))
			testCollection.RemoveGame("0000000000000001")
			Expect(testCollection.Games().Files).To(HaveLen(0))
		})
		It("Removing not existing ID", func() {
			newGames := make([]repository.FileDesc, 0)
			newFile := repository.FileDesc{
				Size:     42,
				Path:     "/here/is/my/game",
				GameID:   "0000000000000001",
				GameInfo: "[0000000000000001][v0].nsp",
				HostType: repository.LocalFile,
			}
			newGames = append(newGames, newFile)
			testCollection.AddNewGames(newGames)

			Expect(testCollection.Games().Files).To(HaveLen(1))
			testCollection.RemoveGame("0000000000000002")
			Expect(testCollection.Games().Files).To(HaveLen(1))
		})
	})
	Describe("Filter", func() {
		var (
			myMockConfig *mock_repository.MockConfig
			ctrl         *gomock.Controller
		)
		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
		})
		JustBeforeEach(func() {
			myMockConfig = mock_repository.NewMockConfig(ctrl)
			customDB := make(map[string]repository.TitleDBEntry)
			custom1 := repository.TitleDBEntry{
				ID:              "0000000000000001",
				Languages:       []string{"FR", "EN", "US"},
				NumberOfPlayers: 1,
			}
			customDB["0000000000000001"] = custom1
			custom2 := repository.TitleDBEntry{
				ID:              "0000000000000002",
				Languages:       []string{"JP"},
				NumberOfPlayers: 2,
			}
			customDB["0000000000000001"] = custom1
			customDB["0000000000000002"] = custom2

			myMockConfig.EXPECT().
				Host().
				Return("tinshop.example.com").
				AnyTimes()
			myMockConfig.EXPECT().
				CustomDB().
				Return(customDB).
				AnyTimes()
			myMockConfig.EXPECT().
				BannedTheme().
				Return(nil).
				AnyTimes()
			myMockConfig.EXPECT().
				RootShop().
				Return("http://tinshop.example.com").
				AnyTimes()
			myMockConfig.EXPECT().
				WelcomeMessage().
				Return("Welcome to testing shop!").
				AnyTimes()
			myMockConfig.EXPECT().
				NoWelcomeMessage().
				Return(false).
				AnyTimes()

			testCollection.OnConfigUpdate(myMockConfig)

			newGames := make([]repository.FileDesc, 0)
			newFile1 := repository.FileDesc{
				Size:     1,
				Path:     "/here/is/my/game",
				GameID:   "0000000000000001",
				GameInfo: "[0000000000000001][v0].nsp",
				HostType: repository.LocalFile,
			}
			newFile2 := repository.FileDesc{
				Size:     22,
				Path:     "/here/is/my/game",
				GameID:   "0000000000000002",
				GameInfo: "[0000000000000002][v0].nsp",
				HostType: repository.LocalFile,
			}
			newGames = append(newGames, newFile1)
			newGames = append(newGames, newFile2)
			testCollection.AddNewGames(newGames)
		})
		It("Filtering world", func() {
			filteredGames := testCollection.Filter("WORLD")
			Expect(len(filteredGames.Titledb)).To(Equal(2))
			Expect(filteredGames.Titledb["0000000000000001"]).NotTo(BeNil())
			Expect(filteredGames.Titledb["0000000000000002"]).NotTo(BeNil())
			Expect(len(filteredGames.Files)).To(Equal(2))
		})
		It("Filtering US", func() {
			filteredGames := testCollection.Filter("US")
			Expect(len(filteredGames.Titledb)).To(Equal(1))
			Expect(filteredGames.Titledb["0000000000000001"]).NotTo(BeNil())
			_, ok := filteredGames.Titledb["0000000000000002"]
			Expect(ok).To(BeFalse())
			Expect(len(filteredGames.Files)).To(Equal(1))
		})
		It("Filtering non existing language entry (HK)", func() {
			filteredGames := testCollection.Filter("HK")
			Expect(len(filteredGames.Titledb)).To(Equal(0))
			Expect(len(filteredGames.Files)).To(Equal(0))
		})
		It("Filtering multi", func() {
			filteredGames := testCollection.Filter("MULTI")
			Expect(len(filteredGames.Titledb)).To(Equal(1))
			_, ok := filteredGames.Titledb["0000000000000001"]
			Expect(ok).To(BeFalse())
			Expect(filteredGames.Titledb["0000000000000002"]).NotTo(BeNil())
			Expect(len(filteredGames.Files)).To(Equal(1))
		})
	})
	Describe("CountGames", func() {
		It("Test with empty collection", func() {
			Expect(testCollection.CountGames()).To(Equal(0))
		})
		It("Test one game in collection", func() {
			newGames := make([]repository.FileDesc, 0)
			newFile := repository.FileDesc{
				Size:     42,
				Path:     "/here/is/my/game",
				GameID:   "0000000000000001",
				GameInfo: "[0000000000000001][v0].nsp",
				HostType: repository.LocalFile,
			}
			customDB := make(map[string]repository.TitleDBEntry)
			newEntry := &repository.TitleDBEntry{
				IconURL: "https://example.com",
			}
			customDB["0000000000000001"] = *newEntry

			myMockConfig.EXPECT().
				CustomDB().
				Return(customDB).
				AnyTimes()
			myMockConfig.EXPECT().
				BannedTheme().
				Return(nil).
				AnyTimes()

			testCollection.OnConfigUpdate(myMockConfig)

			newGames = append(newGames, newFile)
			testCollection.AddNewGames(newGames)

			Expect(testCollection.CountGames()).To(Equal(1))
		})
	})
	Describe("GetKey", func() {
		var (
			myMockConfig *mock_repository.MockConfig
			ctrl         *gomock.Controller
		)
		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
		})
		JustBeforeEach(func() {
			myMockConfig = mock_repository.NewMockConfig(ctrl)
			customDB := make(map[string]repository.TitleDBEntry)
			custom1 := repository.TitleDBEntry{
				ID:  "0000000000000001",
				Key: "My Key",
			}
			customDB["0000000000000001"] = custom1
			custom2 := repository.TitleDBEntry{
				ID:  "0000000000000002",
				Key: "",
			}
			customDB["0000000000000001"] = custom1
			customDB["0000000000000002"] = custom2

			myMockConfig.EXPECT().
				CustomDB().
				Return(customDB).
				AnyTimes()
			myMockConfig.EXPECT().
				BannedTheme().
				Return(nil).
				AnyTimes()
			myMockConfig.EXPECT().
				WelcomeMessage().
				Return("Welcome to testing shop!").
				AnyTimes()
			myMockConfig.EXPECT().
				NoWelcomeMessage().
				Return(false).
				AnyTimes()

			testCollection.OnConfigUpdate(myMockConfig)
		})
		It("Retrieving existing Key", func() {
			key, err := testCollection.GetKey("0000000000000001")
			Expect(err).To(BeNil())
			Expect(key).NotTo(BeEmpty())
			Expect(key).To(Equal("My Key"))
		})
		It("Retrieving not existing Key", func() {
			key, err := testCollection.GetKey("0000000000000002")
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(Equal("TitleDBKey for game 0000000000000002 is not found"))
			Expect(key).To(BeEmpty())
		})
	})
})
