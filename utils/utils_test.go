package utils_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ajmandourah/tinshop/repository"
	"github.com/ajmandourah/tinshop/utils"
)

var _ = Describe("Utils", func() {
	Describe("ExtractGameId", func() {
		Context("Should succeed", func() {
			It("Nicely separated groups", func() {
				game := utils.ExtractGameID("Paw Patrol Mighty Pups Save Adventure Bay [01001F201121E800][v131072] (1.58 GB).nsz")

				Expect(game.Extension()).To(Equal("nsz"))
				Expect(game.ShortID()).To(Equal("01001F201121E800"))
				Expect(game.FullID()).To(Equal("[01001F201121E800][v131072].nsz"))
			})
			It("Make upper of Game Id", func() {
				game := utils.ExtractGameID("Game [01001f201121e800][v131072] (1.58 GB).nsz")

				Expect(game.Extension()).To(Equal("nsz"))
				Expect(game.ShortID()).To(Equal("01001F201121E800"))
				Expect(game.FullID()).To(Equal("[01001F201121E800][v131072].nsz"))
			})
			It("Should only take interesting part", func() {
				game := utils.ExtractGameID("Luigi’s Mansion 3 [Luigi’s Mansion 3 Multiplayer Pack 1][0100DCA0064A7001][US][v131072].nsp")

				Expect(game.Extension()).To(Equal("nsp"))
				Expect(game.ShortID()).To(Equal("0100DCA0064A7001"))
				Expect(game.FullID()).To(Equal("[0100DCA0064A7001][v131072].nsp"))
			})
			It("Group tied with parenthesis group", func() {
				game := utils.ExtractGameID("Paw Patrol Mighty Pups Save Adventure Bay [01001F201121E800][v131072](1.58 GB).nsz")

				Expect(game.Extension()).To(Equal("nsz"))
				Expect(game.ShortID()).To(Equal("01001F201121E800"))
				Expect(game.FullID()).To(Equal("[01001F201121E800][v131072].nsz"))
			})
			It("Nice filename with nsp file", func() {
				game := utils.ExtractGameID("Super Mario Odyssey [0100000000010000][v0].nsp")

				Expect(game.Extension()).To(Equal("nsp"))
				Expect(game.ShortID()).To(Equal("0100000000010000"))
				Expect(game.FullID()).To(Equal("[0100000000010000][v0].nsp"))
			})
			It("Nice separated DLC information", func() {
				game := utils.ExtractGameID("The Legend of Zelda Breath of the Wild [DLC Pack 1 The Master Trials] [01007EF00011F001][v196608].nsp")

				Expect(game.Extension()).To(Equal("nsp"))
				Expect(game.ShortID()).To(Equal("01007EF00011F001"))
				Expect(game.FullID()).To(Equal("[01007EF00011F001][v196608].nsp"))
			})
			It("Tied DLC info to game id and version", func() {
				game := utils.ExtractGameID("Fake - The Legend of Zelda Breath of the Wild [DLC Pack 1 The Master Trials][01007EF00011F001][v196608].nsp")

				Expect(game.Extension()).To(Equal("nsp"))
				Expect(game.ShortID()).To(Equal("01007EF00011F001"))
				Expect(game.FullID()).To(Equal("[01007EF00011F001][v196608].nsp"))
			})
			It("Tied DLC info with no space to game id and version", func() {
				game := utils.ExtractGameID("Fake - The Legend of Zelda Breath of the Wild [DLCPack1TheMasterTrials][01007EF00011F001][v196608].nsp")

				Expect(game.Extension()).To(Equal("nsp"))
				Expect(game.ShortID()).To(Equal("01007EF00011F001"))
				Expect(game.FullID()).To(Equal("[01007EF00011F001][v196608].nsp"))
			})
			It("Game inside sub directory", func() {
				game := utils.ExtractGameID("Fake - My Directory/Fake - [0100152000022800][v655360].nsz")

				Expect(game.Extension()).To(Equal("nsz"))
				Expect(game.ShortID()).To(Equal("0100152000022800"))
				Expect(game.FullID()).To(Equal("[0100152000022800][v655360].nsz"))
			})
		})
		Context("Should Fail", func() {
			It("Test with not size valid game id", func() {
				game := utils.ExtractGameID("Fake - My Game [NSP]/Fake - My Own Game [1231231][v0].nsz")

				Expect(game.Extension()).To(BeEmpty())
				Expect(game.ShortID()).To(BeEmpty())
				Expect(game.FullID()).To(BeEmpty())
			})
			It("Test with bad number of version", func() {
				game := utils.ExtractGameID("Fake - My Game [NSP]/Fake - My Own Game [0100152000022800][0].nsz")

				Expect(game.Extension()).To(BeEmpty())
				Expect(game.ShortID()).To(BeEmpty())
				Expect(game.FullID()).To(BeEmpty())
			})
			It("Test with no game id no version", func() {
				game := utils.ExtractGameID("Fake - Bad name.txt")

				Expect(game.Extension()).To(BeEmpty())
				Expect(game.ShortID()).To(BeEmpty())
				Expect(game.FullID()).To(BeEmpty())
			})
			It("Test with double extension", func() {
				game := utils.ExtractGameID("Fake - Bad name.old.txt")

				Expect(game.Extension()).To(BeEmpty())
				Expect(game.ShortID()).To(BeEmpty())
				Expect(game.FullID()).To(BeEmpty())
			})
		})
	})
	Describe("RemoveFileDesc", func() {
		It("With empty source", func() {
			source := make([]repository.FileDesc, 0)
			res := utils.RemoveFileDesc(source, 0)
			Expect(res).To(HaveLen(0))
		})
		It("With two elements source", func() {
			source := make([]repository.FileDesc, 0)
			file1 := repository.FileDesc{Size: 1}
			file2 := repository.FileDesc{Size: 2}
			source = append(source, file1)
			source = append(source, file2)
			res := utils.RemoveFileDesc(source, 0)
			Expect(res).To(HaveLen(1))
			Expect(res[0].Size).To(Equal(int64(2)))
		})
		It("With negative index", func() {
			source := make([]repository.FileDesc, 0)
			file1 := repository.FileDesc{Size: 1}
			file2 := repository.FileDesc{Size: 2}
			source = append(source, file1)
			source = append(source, file2)
			res := utils.RemoveFileDesc(source, -2)
			Expect(res).To(HaveLen(2))
			Expect(res[0].Size).To(Equal(int64(1)))
			Expect(res[1].Size).To(Equal(int64(2)))

		})
		It("With out of bound index", func() {
			source := make([]repository.FileDesc, 0)
			file1 := repository.FileDesc{Size: 1}
			file2 := repository.FileDesc{Size: 2}
			source = append(source, file1)
			source = append(source, file2)
			res := utils.RemoveFileDesc(source, 4)
			Expect(res).To(HaveLen(2))
			Expect(res[0].Size).To(Equal(int64(1)))
			Expect(res[1].Size).To(Equal(int64(2)))

		})
	})
	Describe("RemoveGameFile", func() {
		It("With empty source", func() {
			source := make([]repository.GameFileType, 0)
			res := utils.RemoveGameFile(source, 0)
			Expect(res).To(HaveLen(0))
		})
		It("With two elements source", func() {
			source := make([]repository.GameFileType, 0)
			file1 := repository.GameFileType{Size: 1}
			file2 := repository.GameFileType{Size: 2}
			source = append(source, file1)
			source = append(source, file2)
			res := utils.RemoveGameFile(source, 0)
			Expect(res).To(HaveLen(1))
			Expect(res[0].Size).To(Equal(int64(2)))
		})
		It("With negative index", func() {
			source := make([]repository.GameFileType, 0)
			file1 := repository.GameFileType{Size: 1}
			file2 := repository.GameFileType{Size: 2}
			source = append(source, file1)
			source = append(source, file2)
			res := utils.RemoveGameFile(source, -2)
			Expect(res).To(HaveLen(2))
			Expect(res[0].Size).To(Equal(int64(1)))
			Expect(res[1].Size).To(Equal(int64(2)))
		})
		It("With out of bound index", func() {
			source := make([]repository.GameFileType, 0)
			file1 := repository.GameFileType{Size: 1}
			file2 := repository.GameFileType{Size: 2}
			source = append(source, file1)
			source = append(source, file2)
			res := utils.RemoveGameFile(source, 4)
			Expect(res).To(HaveLen(2))
			Expect(res[0].Size).To(Equal(int64(1)))
			Expect(res[1].Size).To(Equal(int64(2)))

		})
	})
	Describe("IsValidFilter", func() {
		It("Test multi", func() {
			Expect(utils.IsValidFilter("multi")).To(BeTrue())
		})
		It("Test multi (with pending slash)", func() {
			Expect(utils.IsValidFilter("multi/")).To(BeFalse())
		})
		It("Test multi (multi case)", func() {
			Expect(utils.IsValidFilter("muLtI")).To(BeTrue())
		})
		It("Test multi (multi case)", func() {
			Expect(utils.IsValidFilter("muLtI/")).To(BeFalse())
		})
		It("Test world", func() {
			Expect(utils.IsValidFilter("world")).To(BeTrue())
		})
		It("Test world (with pending slash)", func() {
			Expect(utils.IsValidFilter("world/")).To(BeFalse())
		})
		It("Test fr (lowercase)", func() {
			Expect(utils.IsValidFilter("fr")).To(BeTrue())
		})
		It("Test fr (lowercase & with pending slash)", func() {
			Expect(utils.IsValidFilter("fr/")).To(BeFalse())
		})
		It("Test FR (uppercase)", func() {
			Expect(utils.IsValidFilter("FR")).To(BeTrue())
		})
		It("Test FR (multi case)", func() {
			Expect(utils.IsValidFilter("fR")).To(BeTrue())
		})
		It("Test GG (Not existing country)", func() {
			Expect(utils.IsValidFilter("GG")).To(BeFalse())
		})
		It("Test superpath (Not valid entry)", func() {
			Expect(utils.IsValidFilter("superpath")).To(BeFalse())
		})
	})
	Describe("GetIPFromRequest", func() {
		It("Test with ip", func() {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = "10.0.0.10"
			Expect(utils.GetIPFromRequest(req)).To(Equal("10.0.0.10"))
		})
		It("Test with ip and X-Forwarded-For", func() {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = "10.0.0.10"
			req.Header.Set("X-Forwarded-For", "1.1.1.1")
			Expect(utils.GetIPFromRequest(req)).To(Equal("1.1.1.1"))
		})
	})
	Describe("Search", func() {
		It("Test with not found value", func() {
			myTab := make([]string, 0)
			myTab = append(myTab, "test")
			idxMyTab := utils.Search(len(myTab), func(index int) bool {
				return myTab[index] == "dblk"
			})
			Expect(idxMyTab).To(Equal(-1))
		})
		It("Test with found value", func() {
			myTab := make([]string, 0)
			myTab = append(myTab, "dblk")
			myTab = append(myTab, "test")
			idxMyTab := utils.Search(len(myTab), func(index int) bool {
				return myTab[index] == "test"
			})
			Expect(idxMyTab).To(Equal(1))
		})
	})
})
