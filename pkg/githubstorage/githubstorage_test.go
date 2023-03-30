package githubstorage_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/dathan/go-openai-prompt-git-save/pkg/githubstorage"
)

var _ = Describe("Githubstorage", func() {
	var input string
	BeforeEach(func() {

		input = "Hello"

	})
	Describe("Saving to github", func() {
		It("should save to github", func() {
			err := githubstorage.SaveInput(input, "Initial commit")
			Expect(err).To(BeNil())
		})
	})

})
