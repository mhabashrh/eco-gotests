package tests

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/deployment"
	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/pod"
	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/reportxml"

	. "github.com/rh-ecosystem-edge/eco-gotests/tests/rhwa/internal/rhwainittools"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/rhwa/internal/rhwaparams"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/rhwa/sbr-operator/internal/sbrparams"
)

var _ = Describe(
	"SBR tests",
	Ordered,
	ContinueOnFailure,
	Label(sbrparams.Label), func() {
		BeforeAll(func() {
			By("Get SBR deployment object")
			sbrDeployment, err := deployment.Pull(
				APIClient, sbrparams.OperatorDeploymentName, rhwaparams.RhwaOperatorNs)
			Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("Failed to get SBR deployment %s", err))

			By("Verify SBR deployment is Ready")
			Expect(sbrDeployment.IsReady(rhwaparams.DefaultTimeout)).To(BeTrue(), "SBR deployment is not Ready")
		})
		It("Verify Node Maintenance Operator pod is running", reportxml.ID("46315"), func() {
			_, err := pod.WaitForAllPodsInNamespaceRunning(
				APIClient,
				rhwaparams.RhwaOperatorNs,
				rhwaparams.DefaultTimeout,
			)
			Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("Pod is not ready %s", err))
		})
	})
