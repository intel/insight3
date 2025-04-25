package provider

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-github/github"
	"github.com/intel-sandbox/kube-score/pkg/actions/imageutils"
	"github.com/intel-sandbox/kube-score/pkg/clients/db"
	"github.com/intel-sandbox/kube-score/pkg/clients/ghclient"
	"github.com/intel-sandbox/kube-score/pkg/reports"
	"github.com/intel-sandbox/kube-score/pkg/utils"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	"k8s.io/kubernetes/cmd/kubeadm/app/images"
)

var releaseFilers = []string{
	"rc",
	"alpha",
	"beta",
}

type OSSProvider struct {
	RepoURL  string
	GHClient *ghclient.GHClient
	DBClient *db.RedisClient
}

func (oss OSSProvider) GetReleases(ctx context.Context) ([]reports.ReleaseMD, error) {
	return oss.GHClient.GetAllReleases(ctx, oss.RepoURL, 20)
}

func (oss OSSProvider) GetReleaseMeta(ctx context.Context, version string, url string) (reports.ReleaseMD, error) {
	repoURL := oss.RepoURL
	if url != "" {
		repoURL = url
	}
	return oss.GHClient.GetRelease(ctx, repoURL, version)
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

func (oss OSSProvider) GetReleaseVersions(ctx context.Context, version string, repoURL string, topk int) ([]reports.ImageReport, error) {
	releases := []reports.ImageReport{}

	repoOwner := ""
	if repoURL != "" {
		repoOwner = utils.ParseRepositoryOwner(repoURL)
	}
	//fmt.Printf("Repo url: %s\n", repoURL)
	repo := utils.ParseRepositoryName(repoURL)
	pageCount := 0
	for {
		pageCount++
		rlist, ghres, err := oss.GHClient.ClientV3.Repositories.ListReleases(ctx, repoOwner, repo, &github.ListOptions{Page: pageCount, PerPage: 100})
		if err != nil || ghres.StatusCode != 200 || len(rlist) == 0 {
			// logger.Error(err, "error making github api calls")
			if ghres != nil && ghres.StatusCode == http.StatusForbidden {
				fmt.Printf("un-expected return code: %s", ghres.StatusCode)
				return releases, err
			}
			break
		}
		for _, r := range rlist {
			ignoreRelease := false
			for _, f := range releaseFilers {
				if strings.Contains(r.GetTagName(), f) {
					ignoreRelease = true
				}
			}
			if !ignoreRelease {
				releases = append(releases, reports.ImageReport{
					URL:       r.GetHTMLURL(),
					CreatedAt: r.GetCreatedAt().String(),
				})
			}
		}
	}
	// trim results to latest 20
	res := utils.SortAndTrimReports(releases, topk)
	return res, nil
}
