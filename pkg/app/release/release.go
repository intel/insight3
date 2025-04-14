package release

import (
	"context"
	"errors"
	"fmt"

	"github.com/intel-sandbox/kube-score/pkg/actions/vulns"
	"github.com/intel-sandbox/kube-score/pkg/clients/db"
	"github.com/intel-sandbox/kube-score/pkg/clients/ghclient"
	"github.com/intel-sandbox/kube-score/pkg/common"
	"github.com/intel-sandbox/kube-score/pkg/provider"
	"github.com/intel-sandbox/kube-score/pkg/reports"
)

func Start(opts *common.ReleaseCmdOpts) error {
	ctx := context.Background()

	ghclient := ghclient.GHClient{}
	if err := ghclient.Setup(ctx, opts.Config.GitHub.APIKey); err != nil {
		return fmt.Errorf("error connecting to github")
	}

	dbclient := db.RedisClient{}
	// if err := dbclient.NewClient(ctx, opts.Config.Database); err != nil {
	// 	return fmt.Errorf("error connecting to redis db")
	// }

	var kProvider provider.KubeProvider
	if opts.Distribution == "rke" {
		kProvider = provider.RkeProvider{
			RepoURL:  common.RKEReporURL,
			DBClient: &dbclient,
			GHClient: &ghclient,
		}
	} else {
		kProvider = provider.OSSProvider{
			RepoURL:  common.K8sRepoUrl,
			DBClient: &dbclient,
			GHClient: &ghclient,
		}
	}

	if opts.ListVersions {
		listReleases(ctx, kProvider, opts)
	} else {
		evalRelease(ctx, kProvider, opts)
	}
	// fmt.Printf("opts recieved: %v\n", opts)
	return nil
}

func evalRelease(ctx context.Context, provider provider.KubeProvider, opts *common.ReleaseCmdOpts) error {
	r := reports.ReleaseReport{}
	r.ReleaseTag = opts.Version

	rmd, err := provider.GetReleaseMeta(ctx, opts.Version)
	if err != nil {
		fmt.Printf("version validation failed: %s\n", opts.Version)
		return errors.New("version validation failed")
	}
	fmt.Printf("Kubernetes version = %s\nRelease time: %v\n", opts.Version, rmd.CreatedAt.String())

	rimgs, err := provider.GetReleaseImages(ctx, opts.Version)
	if err != nil {
		fmt.Printf("failed to discover images for: %s\n", opts.Version)
		return errors.New("image discovery failed")
	}

	scanner := vulns.TrivyScanner{}
	// scanner.Init(opts.Config.VulnerabilityScannner)

	for idx, img := range rimgs {
		vData, err := scanner.ScanImage(ctx, img.URL)
		if err != nil {
			fmt.Printf("error scanning image: %s\n", img.URL)
			continue
		}
		rimgs[idx].Vulnerabilities = *vData
	}
	r.Images = rimgs
	reports.PrintReleaseImages(r)
	return nil
}

func listReleases(ctx context.Context, provider provider.KubeProvider, opts *common.ReleaseCmdOpts) error {
	releases, err := provider.GetReleases(ctx)
	if err != nil {
		fmt.Printf("error listing releases\n")
		return errors.New("error listing releases")
	}
	reports.PrintReleaseList(releases)
	return nil
}
