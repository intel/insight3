package vulns

import (
	"context"

	"github.com/intel-sandbox/kube-score/pkg/common"
	"github.com/intel-sandbox/kube-score/pkg/reports"
)

type Scanner interface {
	Init(common.ScannerConfig) error
	ScanImage(context.Context, string) (reports.VulnerabilityData, error)
}
