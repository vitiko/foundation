package channel_transfer_task_executor_only_tx

import (
	"encoding/json"
	"testing"

	"github.com/anoideaopen/foundation/test/integration/cmn/runner"
	"github.com/hyperledger/fabric/integration/nwo"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestChannelTransfer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Channel transfer with task executor only tx Suite")
}

var (
	buildServer *nwo.BuildServer
	components  *nwo.Components
)

var _ = SynchronizedBeforeSuite(func() []byte {
	nwo.RequiredImages = []string{nwo.CCEnvDefaultImage, runner.RedisDefaultImage}

	buildServer = nwo.NewBuildServer()
	buildServer.Serve()

	components = buildServer.Components()
	payload, err := json.Marshal(components)
	Expect(err).NotTo(HaveOccurred())

	return payload
}, func(payload []byte) {
	err := json.Unmarshal(payload, &components)
	Expect(err).NotTo(HaveOccurred())
})

var _ = SynchronizedAfterSuite(func() {
}, func() {
	buildServer.Shutdown()
})