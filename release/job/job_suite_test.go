package job_test

import (
	. "github.com/cloudfoundry/bosh-init/internal/github.com/onsi/ginkgo"
	. "github.com/cloudfoundry/bosh-init/internal/github.com/onsi/gomega"
	"testing"
)

func TestRelease(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Release Job Suite")
}
