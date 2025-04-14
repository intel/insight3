package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/intel-sandbox/kube-score/pkg/app/recommend"
	"github.com/intel-sandbox/kube-score/pkg/common"
	"github.com/intel-sandbox/kube-score/pkg/config"
	"github.com/intel-sandbox/kube-score/pkg/logging"

	"github.com/spf13/cobra"
)

var recOpts common.RecommendCmdOpts

var recommendCmd = &cobra.Command{
	Use:   "recommend",
	Short: "recommend::kube-score",
	Long:  `This command recommends kubernetes release upgrade options`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := logging.WithLogger(context.Background())
		logger := logging.FromContext(ctx)

		err := validateRecommendFlags(args)
		if err != nil {
			logger.Errorf("unable to validate flags: %v\n", err)
			_ = cmd.Help()
			os.Exit(1)
		}
		recommend.Start(&recOpts)
	},
}

func validateRecommendFlags(args []string) error {
	if recOpts.CurrentVersion == "" {
		return fmt.Errorf("missing input parameters `version`")
	}

	if cfg := config.ConfigParser(&recOpts.ConfigFilepath); cfg == nil {
		return fmt.Errorf("error parsing config file")
	} else {
		recOpts.Config = *cfg
	}
	return nil
}

func init() {
	recommendCmd.PersistentFlags().StringVar(&recOpts.CurrentVersion, "version", "", "kubernetes release version")
	recommendCmd.PersistentFlags().StringVar(&recOpts.OutputFormat, "output", "stdout", "output format (stdout, json) (default: stdout)")
	recommendCmd.PersistentFlags().StringVar(&recOpts.ConfigFilepath, "config", ".kube_score_use.yaml", "kube-scopre config file (default: .kube_score.yaml")
	_ = releaseCmd.MarkPersistentFlagRequired("version")
}
