package recommend

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hako/durafmt"
	"github.com/intel-sandbox/kube-score/pkg/clients/ghclient"
	"github.com/intel-sandbox/kube-score/pkg/common"
	"github.com/intel-sandbox/kube-score/pkg/reports"
	"github.com/intel-sandbox/kube-score/pkg/utils"
)

func Start(opts *common.RecommendCmdOpts) error {
	ctx := context.Background()
	ghclient := ghclient.GHClient{}
	if err := ghclient.Setup(ctx, opts.Config.GitHub.APIKey); err != nil {
		return fmt.Errorf("error connecting to github")
	}
	recommendation(ctx, &ghclient, opts)

	return nil
}

func recommendation(ctx context.Context, ghclient *ghclient.GHClient, opts *common.RecommendCmdOpts) error {
	report := reports.RecommendationReport{}
	repoUrl := common.K8sRepoUrl
	var found bool
	if opts.Component != "" {
		repoUrl, _, found = GetGitHubSourceByComponent(opts, opts.Component) // read from config
		if !found {
			return errors.New("unsupported component supplied")
		}
	}

	laterReleases, _ := ghclient.GetAllReleasesGreaterThan(ctx, repoUrl, opts.CurrentVersion)
	if len(laterReleases) == 0 {
		fmt.Printf("You are at the latest available version\n")
		return nil
	}

	currentReleaseTime, _ := ghclient.GetReleaseTimestamp(ctx, repoUrl, opts.CurrentVersion)

	var bestCandidate string
	var bestCandidateScore float64 = -1
	var bestCandidateReleaseTime time.Time

	// Determine score for each upgrade version
	for _, release := range laterReleases {
		score, skip := ghclient.ScoreUpgradeCandidate(ctx, opts.CurrentVersion, release.Tag)
		if skip {
			continue
		}
		if score > bestCandidateScore {
			bestCandidate = release.Tag
			bestCandidateScore = score
			bestCandidateReleaseTime, _ = ghclient.GetReleaseTimestamp(ctx, common.K8sRepoUrl, bestCandidate)
		}
	}

	if bestCandidate == "" {
		fmt.Println("No suitable upgrade candidate found.")
		return nil
	}

	// Fill report
	report.CurrentRelease = opts.CurrentVersion
	report.RecommendedRelease = bestCandidate
	d, _ := durafmt.ParseString(bestCandidateReleaseTime.Sub(currentReleaseTime).String())
	report.ReleaseLagTime = d.String()
	report.ReleaseLagSpace = ghclient.VersionDistance(opts.CurrentVersion, bestCandidate)
	report.LatestRelease = laterReleases[0].Tag
	utils.PrintRecommendationReport(report)

	return nil
}

func GetGitHubSourceByComponent(opts *common.RecommendCmdOpts, component string) (string, int, bool) {
	for _, c := range opts.Config.ConfigSpec.ConfigYAML.ThirdPartyComponentPolicy {
		if c.ComponentName == component {
			return c.GitHubSource, c.TopK, true
		}
	}
	return "", 0, false
}
