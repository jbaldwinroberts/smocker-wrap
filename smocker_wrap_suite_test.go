package smocker_wrap_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSmockerWrap(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SmockerWrap Suite")
}
