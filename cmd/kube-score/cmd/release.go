package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/intel-sandbox/kube-score/pkg/app/release"
	"github.com/intel-sandbox/kube-score/pkg/common"
	"github.com/intel-sandbox/kube-score/pkg/config"
	"github.com/intel-sandbox/kube-score/pkg/logging"

	"github.com/spf13/cobra"
)

var opts common.ReleaseCmdOpts

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "release::kube-score",
	Long:  `This command evaluates a given kubernetes release`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := logging.WithLogger(context.Background())
		logger := logging.FromContext(ctx)

		err := validateFlags(args)
		if err != nil {
			logger.Errorf("unable to validate flags: %v\n", err)
			_ = cmd.Help()
			os.Exit(1)
		}

		release.Start(&opts)
	},
}

func validateFlags(args []string) error {

	if opts.Version == "" && !opts.ListVersions {
		return fmt.Errorf("missing input parameters, select `list` or `version` option ")
	}

	if opts.Distribution != "k8s" && opts.Distribution != "rke" {
		return fmt.Errorf("un-supported distribution %s ", opts.Distribution)
	}
	if cfg := config.ConfigParser(&opts.ConfigFilepath); cfg == nil {
		return fmt.Errorf("error parsing config file")
	} else {
		opts.Config = *cfg
	}
	return nil
}

func init() {
	releaseCmd.PersistentFlags().StringVar(&opts.Version, "version", "", "kubernetes release version")
	releaseCmd.PersistentFlags().StringVar(&opts.OutputFormat, "output", "stdout", "output format (stdout, json) (default: stdout)")
	releaseCmd.PersistentFlags().StringVar(&opts.ConfigFilepath, "config", ".kube_score.yaml", "kube-scopre config file (default: .kube_score.yaml")
	releaseCmd.PersistentFlags().BoolVar(&opts.ListVersions, "list", false, "list recent kubernetes releases (recent 20)")
	releaseCmd.PersistentFlags().StringVar(&opts.Distribution, "dist", "k8s", "kubernetes distribution (k8s, rke) (default: k8s)")
	// _ = releaseCmd.MarkPersistentFlagRequired("version")
}
