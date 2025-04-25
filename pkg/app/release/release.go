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
	"github.com/intel-sandbox/kube-score/pkg/utils"
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
	// if opts.Distribution == "rke" {
	// 	kProvider = provider.RkeProvider{
	// 		RepoURL:  common.RKEReporURL,
	// 		DBClient: &dbclient,
	// 		GHClient: &ghclient,
	// 	}
	// } else {
	kProvider = provider.OSSProvider{
		RepoURL:  common.K8sRepoUrl,
		DBClient: &dbclient,
		GHClient: &ghclient,
	}

	if opts.ListVersions {
		listReleases(ctx, kProvider, opts)
	} else if opts.ShowReport {
		r, _ := evalRelease(ctx, kProvider, opts)
		// format vuln report
		utils.PrintVulnerabilityReport(r)
	} else {
		evalRelease(ctx, kProvider, opts)
	}
	return nil
}

func evalRelease(ctx context.Context, provider provider.KubeProvider, opts *common.ReleaseCmdOpts) ([][]vulns.VulnerabilityReport, error) {
	r := reports.ReleaseReport{}
	r.ReleaseTag = opts.Version
	url := ""
	var found bool
	var topk int
	if opts.Component != "" {
		url, topk, found = GetGitHubSourceByComponent(opts, opts.Component) // read from config
		if !found {
			return nil, errors.New("unsupported component supplied")
		}
	}
	rmd, err := provider.GetReleaseMeta(ctx, opts.Version, url)
	if err != nil {
		fmt.Printf("version validation failed: %s\n", opts.Version)
		return nil, errors.New("version validation failed")
	}
	fmt.Printf("Release version = %s\nRelease time: %v\n", opts.Version, rmd.CreatedAt.String())

	rimgs := []reports.ImageReport{}
	if opts.Component != "" {
		// Third-Party
		rimgs, err = provider.GetReleaseVersions(ctx, opts.Version, url, topk)
		if err != nil {
			fmt.Printf("failed to discover releases for: %s\n", opts.Version)
			return nil, errors.New("release discovery failed")
		}
	} else {
		// Kubernetes
		rimgs, err = provider.GetReleaseImages(ctx, opts.Version)
		if err != nil {
			fmt.Printf("failed to discover images for: %s\n", opts.Version)
			return nil, errors.New("image discovery failed")
		}
	}
	scanner := vulns.TrivyScanner{}
	vulnList := [][]vulns.VulnerabilityReport{}
	for idx, release := range rimgs {
		vData := []vulns.VulnerabilityReport{}
		vSummary := &reports.VulnerabilityData{}
		if opts.Component != "" {
			// ThirdParty
			vSummary, vData, err = scanner.ScanRepo(ctx, release.URL, opts.OutputFilePath)
			if err != nil {
				fmt.Printf("error scanning repo: %s\n", release.URL)
				fmt.Printf("error: %s\n", err)
				return nil, err
			}
		} else {
			// Kubernetes
			vSummary, vData, err = scanner.ScanImage(ctx, release.URL, opts.OutputFilePath)
			if err != nil {
				fmt.Printf("error scanning image: %s\n", release.URL)
				continue
			}
		}

		vulnList = append(vulnList, vData)
		rimgs[idx].Vulnerabilities = *vSummary
	}
	r.Images = rimgs
	utils.PrintReleaseImages(r)
	return vulnList, nil
}

func listReleases(ctx context.Context, provider provider.KubeProvider, opts *common.ReleaseCmdOpts) error {
	releases, err := provider.GetReleases(ctx)
	if err != nil {
		fmt.Printf("error listing releases\n")
		return errors.New("error listing releases")
	}
	utils.PrintReleaseList(releases)
	return nil
}

func GetGitHubSourceByComponent(opts *common.ReleaseCmdOpts, component string) (string, int, bool) {
	for _, c := range opts.Config.ConfigSpec.ConfigYAML.ThirdPartyComponentPolicy {
		if c.ComponentName == component {
			return c.GitHubSource, c.TopK, true
		}
	}
	return "", 0, false
}
