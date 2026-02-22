package nhc

import (
	"runtime"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
		"github.com/rh-ecosystem-edge/eco-goinfra/pkg/namespace"
	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/reportxml"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/internal/params"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/internal/reporter"
	. "github.com/rh-ecosystem-edge/eco-gotests/tests/rhwa/internal/rhwainittools"
		"github.com/rh-ecosystem-edge/eco-gotests/tests/rhwa/internal/rhwaparams"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/rhwa/nhc-operator/internal/nhcparams"
	_ "github.com/rh-ecosystem-edge/eco-gotests/tests/rhwa/nhc-operator/tests"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

var (
	_, currentFile, _, _ = runtime.Caller(0)
	testNS               = namespace.NewBuilder(APIClient, rhwaparams.TestNamespaceName)
)


func TestNHC(t *testing.T) {
	_, reporterConfig := GinkgoConfiguration()
	reporterConfig.JUnitReport = RHWAConfig.GetJunitReportPath(currentFile)

	RegisterFailHandler(Fail)
	RunSpecs(t, "NHC", Label(nhcparams.Labels...), reporterConfig)
}

var _ = JustAfterEach(func() {
	reporter.ReportIfFailed(
		CurrentSpecReport(), currentFile, nhcparams.ReporterNamespacesToDump, nhcparams.ReporterCRDsToDump)
})

var _ = ReportAfterSuite("", func(report Report) {
	reportxml.Create(
		report, RHWAConfig.GetReportPath(), RHWAConfig.TCPrefix)
})

var _ = BeforeSuite(func() {
	By("Creating test namespace with privileged labels")
	for key, value := range params.PrivilegedNSLabels {
		testNS.WithLabel(key, value)
	}
	_, err := testNS.Create()

	if err != nil && !apierrors.IsAlreadyExists(err) {
		Expect(err).ToNot(HaveOccurred(), "error to create test namespace")
	}
})


