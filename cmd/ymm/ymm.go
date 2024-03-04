package main

import (
	"os"

	"github.com/crowdigit/ymm/internal"
	"github.com/crowdigit/ymm/pkg/exec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "ymm",
	Short: "Download MP3 files from Youtube and manage files",
}

var singleCmd = &cobra.Command{
	Use:   "single",
	Short: "Download a MP3 file from URL",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ci := exec.NewCommandProvider()
		return internal.DownloadSingle(
			ci,
			viper.GetString("jq.path"),
			viper.GetStringSlice("jq.args"),
			args[1],
		)
	},
}

func init() {
	rootCmd.AddCommand(singleCmd)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
