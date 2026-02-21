package mdr

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
	"github.com/rh-ecosystem-edge/eco-gotests/tests/rhwa/mdr-operator/internal/mdrparams"
	_ "github.com/rh-ecosystem-edge/eco-gotests/tests/rhwa/mdr-operator/tests"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

var (
	_, currentFile, _, _ = runtime.Caller(0)
	// Initialize the namespace builder using the standard RHWA test namespace name (dast-tests)
	testNS = namespace.NewBuilder(APIClient, rhwaparams.TestNamespaceName)
)

func TestMDR(t *testing.T) {
	_, reporterConfig := GinkgoConfiguration()
	reporterConfig.JUnitReport = RHWAConfig.GetJunitReportPath(currentFile)

	RegisterFailHandler(Fail)
	RunSpecs(t, "MDR", Label(mdrparams.Labels...), reporterConfig)
}

// Added BeforeSuite to create the namespace required for DAST/Trivy scans
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

var _ = JustAfterEach(func() {
	reporter.ReportIfFailed(
		CurrentSpecReport(), currentFile, mdrparams.ReporterNamespacesToDump, mdrparams.ReporterCRDsToDump)
})

var _ = ReportAfterSuite("", func(report Report) {
	reportxml.Create(
		report, RHWAConfig.GetReportPath(), RHWAConfig.TCPrefix)
})

// Added AfterSuite to clean up the namespace after tests finish
var _ = AfterSuite(func() {
	By("Deleting test namespace")
	err := testNS.DeleteAndWait(rhwaparams.DefaultTimeout)
	Expect(err).ToNot(HaveOccurred(), "error to delete test namespace")
})