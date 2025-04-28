package provider

import (
	"context"

	"github.com/intel-sandbox/kube-score/pkg/reports"
)

type KubeProvider interface {
	GetReleases(context.Context) ([]reports.ReleaseMD, error)
	GetReleaseMeta(context.Context, string, string) (reports.ReleaseMD, error)
	GetReleaseAssets(context.Context, string, string) ([]reports.ReleaseAsset, error)
	GetReleaseImages(context.Context, string) ([]reports.ImageReport, error)
	GetReleaseVersions(context.Context, string, string, int) ([]reports.ImageReport, error)
}
