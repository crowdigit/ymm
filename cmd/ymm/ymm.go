package main

import (
	"fmt"
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

		var jqConf internal.ExecConfig
		if err := viper.UnmarshalKey("jq", &jqConf); err != nil {
			return fmt.Errorf("failed to unmarshal config: %w", err)
		}

		var ytConf internal.ExecConfig
		if err := viper.UnmarshalKey("yt", &ytConf); err != nil {
			return fmt.Errorf("failed to unmarshal config: %w", err)
		}

		return internal.DownloadSingle(ci, jqConf, ytConf, args[0])
	},
}

func init() {
	viper.SetDefault("yt.path", "yt-dlp")
	viper.SetDefault("yt.args", []string{"--dump-json"})
	viper.SetDefault("jq.path", "jaq")
	viper.SetDefault("jq.args", []string{"--slurp", "."})

	rootCmd.AddCommand(singleCmd)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
