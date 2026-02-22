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
	"github.com/rh-ecosystem-edge/eco-gotests/tests/rhwa/nhc-operator/internal/nhcparams"
)

var _ = Describe(
	"NHC tests",
	Ordered,
	ContinueOnFailure,
	Label(nhcparams.Label), func() {
		BeforeAll(func() {
			By("Get NHC deployment object")
			nhcDeployment, err := deployment.Pull(
				APIClient, nhcparams.OperatorDeploymentName, rhwaparams.RhwaOperatorNs)
			Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("Failed to get NHC deployment %s", err))

			By("Verify NHC deployment is Ready")
			Expect(nhcDeployment.IsReady(rhwaparams.DefaultTimeout)).To(BeTrue(), "NHC deployment is not Ready")
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
