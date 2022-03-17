package cmd

import (
	"github.com/alknopfler/tidy-mirror/config"
	"github.com/alknopfler/tidy-mirror/pkg/registry"
	"github.com/spf13/cobra"
	"os"
)

func NewMirrorOlm() *cobra.Command {
	var r *registry.Registry
	var kubeconfig, configPath string
	cmd := &cobra.Command{
		Use:   "olm",
		Short: "Mirroring the OLM operators to the registry deployed based on mode (hub or spoke)",
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := config.NewConfig(kubeconfig, configPath)
			if err != nil {
				return err
			}
			r = registry.NewRegistry(conf)
			r.WritePullSecretBaseToTempFile(r.GetPullSecretBase())
			defer os.Remove(r.PullSecretTempFile)
			return r.RunMirrorOlm()
		},
	}
	flags := cmd.Flags()
	// Read the config flag directly into the struct, so it's immediately available.
	flags.StringVar(&kubeconfig, "kubeconfig", "", "Path to kubeconfig file")
	flags.StringVar(&configPath, "config-file", "", "Path to the config file")
	//TODO add flag to get spoke name or ALL to deploy the registry
	return cmd
}
