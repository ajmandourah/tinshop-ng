package utils_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/DblK/tinshop/utils"
)

var _ = Describe("Bytes", func() {
	Describe("Itob", func() {
		It("Test with 0", func() {
			Expect(utils.Itob(0)).To(Equal([]uint8{0, 0, 0, 0, 0, 0, 0, 0}))
		})
		It("Test with 42", func() {
			Expect(utils.Itob(42)).To(Equal([]uint8{0, 0, 0, 0, 0, 0, 0, 42}))
		})
	})
	Describe("ByteToUint64", func() {
		It("Test with empty byte", func() {
			Expect(utils.ByteToUint64([]byte{})).To(Equal(uint64(0)))
		})
		It("Test with 42 in byte", func() {
			Expect(utils.ByteToUint64([]byte{0, 0, 0, 0, 0, 0, 0, 42})).To(Equal(uint64(42)))
		})
	})
	Describe("ByteToMap", func() {
		It("Test with empty byte", func() {
			res, err := utils.ByteToMap([]byte{})
			Expect(err).To(BeNil())
			Expect(res).To(HaveLen(0))
		})
		It("Test with 42 in byte", func() {
			type test struct {
				Value int `json:"visit,omitempty"`
			}
			newTest := &test{Value: 42}
			buf, _ := json.Marshal(newTest)
			res, err := utils.ByteToMap(buf)
			Expect(err).To(BeNil())
			Expect(res).To(HaveLen(1))
			Expect(res["visit"]).To(Equal(float64(42)))
		})
	})
})
