package config_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dblk/tinshop/config"
	"github.com/dblk/tinshop/mock_repository"
)

var _ = Describe("Config", func() {
	It("Ensure to be able to set a RootShop", func() {
		config.GetConfig().SetRootShop("http://tinshop.example.com")
		cfg := config.GetConfig()

		Expect(cfg.RootShop()).To(Equal("http://tinshop.example.com"))
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
				config.ComputeDefaultValues(myMockConfig)

				Expect(testRootShop).To(Equal("http://tinshop.example.com:8080"))
			})
		})
	})
})
