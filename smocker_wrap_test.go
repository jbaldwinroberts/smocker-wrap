package smocker_wrap_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	sw "smocker-wrap"
)

// TODO mock out smocker

var _ = Describe("SmockerWrap", func() {
	Context("Reset", func() {
		When("Reset is called, with force = true", func() {
			It("clears the mocks and the history of calls", func() {
				err := sw.Reset(true)
				Expect(err).To(BeNil())
			})
		})
	})
	
	Context("AddMock", func() {
		When("AddMock is called", func() {
			It("adds the mock", func() {
				request := sw.MockRequest{
					Method: "GET",
					Path:   "/hello/world",
				}

				response := &sw.MockResponse{
					Status: 200,
					Body:   `{"message": "Hello, World!"}`,
				}
				
				err := sw.AddMock(true, "some-session", request, response)
				Expect(err).To(BeNil())

				result, err := sw.VerifyMocks("some-session")
				Expect(err).To(BeNil())
				Expect(result.Mocks.AllUsed).To(BeFalse())
				Expect(result.Mocks.Verified).To(BeFalse())
			})
		})
	})
})
