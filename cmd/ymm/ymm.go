package main

import (
	"fmt"
	"os"

	"github.com/crowdigit/exec"
	"github.com/crowdigit/ymm/internal"
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

		var ytConf internal.ExecConfig
		if err := viper.UnmarshalKey("command.metadata.youtube", &ytConf); err != nil {
			return fmt.Errorf("failed to unmarshal config: %w", err)
		}

		var jqConf internal.ExecConfig
		if err := viper.UnmarshalKey("command.metadata.json", &jqConf); err != nil {
			return fmt.Errorf("failed to unmarshal config: %w", err)
		}

		return internal.DownloadSingle(ci, ytConf, jqConf, args[0])
	},
}

func init() {
	viper.SetDefault("command.metadata.youtube.path", "yt-dlp")
	viper.SetDefault("command.metadata.youtube.args", []string{"--dump-json", "<url>"})
	viper.SetDefault("command.metadata.json.path", "jaq")
	viper.SetDefault("command.metadata.json.args", []string{"--slurp", "."})
	viper.SetDefault("command.download.youtube.path", "yt-dlp")
	viper.SetDefault(
		"command.download.youtube.args",
		[]string{
			// "--cookies",
			// "<cookies>",
			"--format",
			"<format>",
			"--extract-audio",
			"--audio-format",
			"mp3",
			"--audio-quality",
			"0",
			"<url>",
		},
	)

	rootCmd.AddCommand(singleCmd)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
