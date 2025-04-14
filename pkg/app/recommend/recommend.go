package recommend

import (
	"context"
	"fmt"

	"github.com/hako/durafmt"
	"github.com/intel-sandbox/kube-score/pkg/clients/ghclient"
	"github.com/intel-sandbox/kube-score/pkg/common"
	"github.com/intel-sandbox/kube-score/pkg/reports"
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
	laterReleases, _ := ghclient.GetAllReleasesGreaterThan(ctx, common.K8sRepoUrl, opts.CurrentVersion)
	// for _, lr := range laterReleases {
	// 	fmt.Printf("%s\n", lr.Tag)
	// }
	if len(laterReleases) == 0 {
		fmt.Printf("you are at the latest available version\n")
		return nil
	}

	currentReleaseTime, _ := ghclient.GetReleaseTimestamp(ctx, common.K8sRepoUrl, opts.CurrentVersion)
	latestReleaseTime, _ := ghclient.GetReleaseTimestamp(ctx, common.K8sRepoUrl, laterReleases[0].Tag)

	report.CurrentRelease = opts.CurrentVersion
	report.RecommendedRelease = laterReleases[0].Tag
	durationStr := latestReleaseTime.Sub(currentReleaseTime).String()
	d, _ := durafmt.ParseString(durationStr)
	report.ReleaseLagTime = d.String()
	report.ReleaseLagSpace = len(laterReleases)

	reports.PrintRecommendationReport(report)

	return nil
}
