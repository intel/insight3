package provider

import (
	"context"

	"github.com/intel-sandbox/kube-score/pkg/actions/imageutils"
	"github.com/intel-sandbox/kube-score/pkg/clients/db"
	"github.com/intel-sandbox/kube-score/pkg/clients/ghclient"
	"github.com/intel-sandbox/kube-score/pkg/reports"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	"k8s.io/kubernetes/cmd/kubeadm/app/images"
)

type OSSProvider struct {
	RepoURL  string
	GHClient *ghclient.GHClient
	DBClient *db.RedisClient
}

func (oss OSSProvider) GetReleases(ctx context.Context) ([]reports.ReleaseMD, error) {
	return oss.GHClient.GetAllReleases(ctx, oss.RepoURL, 20)
}

func (oss OSSProvider) GetReleaseMeta(ctx context.Context, version string) (reports.ReleaseMD, error) {
	return oss.GHClient.GetRelease(ctx, oss.RepoURL, version)
}

func (oss OSSProvider) GetReleaseAssets(ctx context.Context, version, name string) ([]reports.ReleaseAsset, error) {

	return nil, nil
}

func (oss OSSProvider) GetReleaseImages(ctx context.Context, version string) ([]reports.ImageReport, error) {
	imgs := []reports.ImageReport{}
	cfg := kubeadmapi.ClusterConfiguration{}
	// externalcfg := &kubeadmapiv1.ClusterConfiguration{}
	// kubeadmscheme.Scheme.Default(externalcfg)
	cfg.KubernetesVersion = version
	cfg.Etcd = kubeadmapi.Etcd{}
	cfg.ImageRepository = "registry.k8s.io"
	cfg.Etcd.Local = &kubeadmapi.LocalEtcd{}
	imglist := images.GetControlPlaneImages(&cfg)
	for _, i := range imglist {
		imgs = append(imgs, reports.ImageReport{
			URL:       i,
			CreatedAt: imageutils.GetImageBuildTime(i),
			Digest:    imageutils.GetImageDigest(i),
		})
	}

	return imgs, nil
}
