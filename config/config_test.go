package config_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ajmandourah/tinshop-ng/config"
	"github.com/ajmandourah/tinshop-ng/mock_repository"
	"github.com/ajmandourah/tinshop-ng/repository"
)

var _ = Describe("Config", func() {
	var testConfig repository.Config
	BeforeEach(func() {
		testConfig = config.New()
	})

	It("Ensure to be able to set a RootShop", func() {
		testConfig.SetRootShop("http://tinshop.example.com")

		Expect(testConfig.RootShop()).To(Equal("http://tinshop.example.com"))
	})
	Context("ComputeDefaultValues", func() {
		var (
			myMockConfig *mock_repository.MockConfig
			ctrl         *gomock.Controller
		)
		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
		})
		JustBeforeEach(func() {
			myMockConfig = mock_repository.NewMockConfig(ctrl)

			myMockConfig.EXPECT().
				Host().
				Return("tinshop.example.com").
				AnyTimes()
			myMockConfig.EXPECT().
				ShopTitle().
				Return("Tinshop!").
				AnyTimes()
			myMockConfig.EXPECT().
				SetShopTemplateData(gomock.Any()).
				Return().
				AnyTimes()
		})
		Describe("ComputeDefaultValues", func() {

			It("Should append the default port(3000)", func() {
				var testRootShop string
				myMockConfig.EXPECT().
					Protocol().
					Return("http").
					AnyTimes()
				myMockConfig.EXPECT().
					SetRootShop(gomock.Any()).
					Return().
					Do(func(rootShop string) {
						testRootShop = rootShop
					}).
					AnyTimes()
				myMockConfig.EXPECT().
					Port().
					Return(0).
					AnyTimes()
				myMockConfig.EXPECT().
					ReverseProxy().
					Return(false).
					AnyTimes()
				config.ComputeDefaultValues(myMockConfig)

				Expect(testRootShop).To(Equal("http://tinshop.example.com:3000"))
			})
			It("Should not set port for https/443 config", func() {
				var testRootShop string
				myMockConfig.EXPECT().
					SetRootShop(gomock.Any()).
					Return().
					Do(func(rootShop string) {
						testRootShop = rootShop
					}).
					AnyTimes()
				myMockConfig.EXPECT().
					Protocol().
					Return("https").
					AnyTimes()
				myMockConfig.EXPECT().
					Port().
					Return(443).
					AnyTimes()
				myMockConfig.EXPECT().
					ReverseProxy().
					Return(false).
					AnyTimes()
				config.ComputeDefaultValues(myMockConfig)

				Expect(testRootShop).To(Equal("https://tinshop.example.com"))
			})
			It("Should not set port for http/80 config", func() {
				var testRootShop string
				myMockConfig.EXPECT().
					SetRootShop(gomock.Any()).
					Return().
					Do(func(rootShop string) {
						testRootShop = rootShop
					}).
					AnyTimes()
				myMockConfig.EXPECT().
					Protocol().
					Return("http").
					AnyTimes()
				myMockConfig.EXPECT().
					Port().
					Return(80).
					AnyTimes()
				myMockConfig.EXPECT().
					ReverseProxy().
					Return(false).
					AnyTimes()
				config.ComputeDefaultValues(myMockConfig)

				Expect(testRootShop).To(Equal("http://tinshop.example.com"))
			})
			It("Should set port non standard port", func() {
				var testRootShop string
				myMockConfig.EXPECT().
					SetRootShop(gomock.Any()).
					Return().
					Do(func(rootShop string) {
						testRootShop = rootShop
					}).
					AnyTimes()
				myMockConfig.EXPECT().
					Protocol().
					Return("http").
					AnyTimes()
				myMockConfig.EXPECT().
					Port().
					Return(8080).
					AnyTimes()
				myMockConfig.EXPECT().
					ReverseProxy().
					Return(false).
					AnyTimes()
				config.ComputeDefaultValues(myMockConfig)

				Expect(testRootShop).To(Equal("http://tinshop.example.com:8080"))
			})
			It("Should not set port if reverse proxy", func() {
				var testRootShop string
				myMockConfig.EXPECT().
					Protocol().
					Return("http").
					AnyTimes()
				myMockConfig.EXPECT().
					SetRootShop(gomock.Any()).
					Return().
					Do(func(rootShop string) {
						testRootShop = rootShop
					}).
					AnyTimes()
				myMockConfig.EXPECT().
					Port().
					Return(0).
					AnyTimes()
				myMockConfig.EXPECT().
					ReverseProxy().
					Return(true).
					AnyTimes()
				config.ComputeDefaultValues(myMockConfig)

				Expect(testRootShop).To(Equal("http://tinshop.example.com"))
			})
		})
	})
	Context("Security for Blacklist/Whitelist tests", func() {
		var myConfig config.Configuration

		BeforeEach(func() {
			myConfig = config.Configuration{}
		})

		Describe("Blacklist tests", func() { //nolint:dupl
			It("With empty blacklist", func() {
				Expect(myConfig.IsBlacklisted("me")).To(BeFalse())
				Expect(myConfig.IsWhitelisted("me")).To(BeTrue())
			})
			It("With a blacklist", func() {
				var blacklist = make([]string, 0)
				blacklist = append(blacklist, "me")

				myConfig.Security.Blacklist = blacklist
				Expect(myConfig.IsBlacklisted("me")).To(BeTrue())
				Expect(myConfig.IsWhitelisted("me")).To(BeFalse())
			})
			It("With a blacklist on other person", func() {
				var blacklist = make([]string, 0)
				blacklist = append(blacklist, "someoneElse")

				myConfig.Security.Blacklist = blacklist
				Expect(myConfig.IsBlacklisted("me")).To(BeFalse())
				Expect(myConfig.IsWhitelisted("me")).To(BeTrue())
			})
		})
		Describe("Whitelist tests", func() { //nolint:dupl
			It("With empty whitelist", func() {
				Expect(myConfig.IsWhitelisted("me")).To(BeTrue())
				Expect(myConfig.IsBlacklisted("me")).To(BeFalse())
			})
			It("With a whitelist", func() {
				var whitelist = make([]string, 0)
				whitelist = append(whitelist, "me")

				myConfig.Security.Whitelist = whitelist
				Expect(myConfig.IsWhitelisted("me")).To(BeTrue())
				Expect(myConfig.IsBlacklisted("me")).To(BeFalse())
			})
			It("With a whitelist for someone else", func() {
				var whitelist = make([]string, 0)
				whitelist = append(whitelist, "someoneElse")

				myConfig.Security.Whitelist = whitelist
				Expect(myConfig.IsWhitelisted("me")).To(BeFalse())
				Expect(myConfig.IsBlacklisted("me")).To(BeTrue())
			})
		})
		Describe("Mix of blacklist/whitelist", func() {
			It("With a blacklist and someone else in whitelist", func() {
				var blacklist = make([]string, 0)
				blacklist = append(blacklist, "me")
				var whitelist = make([]string, 0)
				whitelist = append(whitelist, "someoneElse")

				myConfig.Security.Blacklist = blacklist
				myConfig.Security.Whitelist = whitelist
				Expect(myConfig.IsBlacklisted("me")).To(BeTrue())
				Expect(myConfig.IsWhitelisted("me")).To(BeFalse())
			})
			It("With a blacklist another person and someone else in whitelist", func() {
				var blacklist = make([]string, 0)
				blacklist = append(blacklist, "anotherPerson")
				var whitelist = make([]string, 0)
				whitelist = append(whitelist, "someoneElse")

				myConfig.Security.Blacklist = blacklist
				myConfig.Security.Whitelist = whitelist
				Expect(myConfig.IsBlacklisted("me")).To(BeTrue())
				Expect(myConfig.IsWhitelisted("me")).To(BeFalse())
			})
		})
	})
	Context("Security for theme", func() {
		var myConfig config.Configuration

		BeforeEach(func() {
			myConfig = config.Configuration{}
		})

		Describe("IsBannedTheme", func() {
			It("should not be banned if empty config", func() {
				Expect(myConfig.IsBannedTheme("myTheme")).To(BeFalse())
			})
			It("should not be banned if no corresponding config", func() {
				var bannedThemes = make([]string, 0)
				bannedThemes = append(bannedThemes, "banned")
				myConfig.Security.BannedTheme = bannedThemes
				Expect(myConfig.IsBannedTheme("myTheme")).To(BeFalse())
			})
			It("should not be banned if no corresponding config", func() {
				var bannedThemes = make([]string, 0)
				bannedThemes = append(bannedThemes, "myTheme")
				myConfig.Security.BannedTheme = bannedThemes
				Expect(myConfig.IsBannedTheme("myTheme")).To(BeTrue())
			})
		})
	})
	Describe("Protocol", func() {
		var myConfig config.Configuration

		BeforeEach(func() {
			myConfig = config.Configuration{}
		})

		It("Test with empty object", func() {
			Expect(myConfig.Protocol()).To(BeEmpty())
		})
		It("Test with a value", func() {
			myConfig.ShopProtocol = "https"
			Expect(myConfig.Protocol()).To(Equal("https"))
		})
	})
	Describe("Host", func() {
		var myConfig config.Configuration

		BeforeEach(func() {
			myConfig = config.Configuration{}
		})

		It("Test with empty object", func() {
			Expect(myConfig.Host()).To(BeEmpty())
		})
		It("Test with a value", func() {
			myConfig.ShopHost = "tinshop.example.com"
			Expect(myConfig.Host()).To(Equal("tinshop.example.com"))
		})
	})
	Describe("WelcomeMessage", func() {
		var myConfig config.Configuration

		BeforeEach(func() {
			myConfig = config.Configuration{}
		})

		It("Test with empty object", func() {
			Expect(myConfig.WelcomeMessage()).To(BeEmpty())
		})
		It("Test with a value", func() {
			myConfig.ShopWelcomeMessage = "We are testing it!"
			Expect(myConfig.WelcomeMessage()).To(Equal("We are testing it!"))
		})
		It("Test with a empty value value", func() {
			myConfig.ShopWelcomeMessage = ""
			Expect(myConfig.WelcomeMessage()).To(BeEmpty())
		})
	})
	Describe("Port", func() {
		var myConfig config.Configuration

		BeforeEach(func() {
			myConfig = config.Configuration{}
		})

		It("Test with empty object", func() {
			Expect(myConfig.Port()).To(Equal(0))
		})
		It("Test with a value", func() {
			myConfig.ShopPort = 12345
			Expect(myConfig.Port()).To(Equal(12345))
		})
	})
	Describe("ReverseProxy", func() {
		var myConfig config.Configuration

		BeforeEach(func() {
			myConfig = config.Configuration{}
		})

		It("Test with empty object", func() {
			Expect(myConfig.ReverseProxy()).To(BeFalse())
		})
		It("Test with a value", func() {
			myConfig.Proxy = true
			Expect(myConfig.ReverseProxy()).To(BeTrue())
		})
	})
	Describe("ShopTitle", func() {
		var myConfig config.Configuration

		BeforeEach(func() {
			myConfig = config.Configuration{}
		})

		It("Test with empty object", func() {
			Expect(myConfig.ShopTitle()).To(BeEmpty())
		})
		It("Test with a value", func() {
			myConfig.Name = "Tinshop"
			Expect(myConfig.ShopTitle()).To(Equal("Tinshop"))
		})
	})
	Describe("DebugNfs", func() {
		var myConfig config.Configuration

		BeforeEach(func() {
			myConfig = config.Configuration{}
		})

		It("Test with empty object", func() {
			Expect(myConfig.DebugNfs()).To(BeFalse())
		})
		It("Test with a value", func() {
			myConfig.Debug.Nfs = true
			Expect(myConfig.DebugNfs()).To(BeTrue())
		})
	})
	Describe("VerifyNSP", func() {
		var myConfig config.Configuration

		BeforeEach(func() {
			myConfig = config.Configuration{}
		})

		It("Test with empty object", func() {
			Expect(myConfig.VerifyNSP()).To(BeFalse())
		})
		It("Test with a value", func() {
			myConfig.NSP.CheckVerified = true
			Expect(myConfig.VerifyNSP()).To(BeTrue())
		})
	})
	Describe("DebugNoSecurity", func() {
		var myConfig config.Configuration

		BeforeEach(func() {
			myConfig = config.Configuration{}
		})

		It("Test with empty object", func() {
			Expect(myConfig.DebugNoSecurity()).To(BeFalse())
		})
		It("Test with a value", func() {
			myConfig.Debug.NoSecurity = true
			Expect(myConfig.DebugNoSecurity()).To(BeTrue())
		})
	})
	Describe("DebugTicket", func() {
		var myConfig config.Configuration

		BeforeEach(func() {
			myConfig = config.Configuration{}
		})

		It("Test with empty object", func() {
			Expect(myConfig.DebugTicket()).To(BeFalse())
		})
		It("Test with a value", func() {
			myConfig.Debug.Ticket = true
			Expect(myConfig.DebugTicket()).To(BeTrue())
		})
	})
	Describe("BannedTheme", func() {
		var myConfig config.Configuration

		BeforeEach(func() {
			myConfig = config.Configuration{}
		})

		It("Test with empty object", func() {
			Expect(myConfig.BannedTheme()).To(HaveLen(0))
		})
		It("Test with a value", func() {
			myConfig.Security.BannedTheme = make([]string, 0)
			myConfig.Security.BannedTheme = append(myConfig.Security.BannedTheme, "Banned")
			Expect(myConfig.BannedTheme()).To(HaveLen(1))
			Expect(myConfig.BannedTheme()[0]).To(Equal("Banned"))
		})
	})
})
