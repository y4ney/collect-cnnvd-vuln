package cmd

import (
	"fmt"
	"github.com/y4ney/collect-cnnvd-vuln/internal/config"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "打印 collect-cnnvd-vuln 的版本",
	Long:  `所有的软件导游版本，collect-cnnvd-vuln也不例外`,
	RunE:  runPrintVersion,
}

func runPrintVersion(_ *cobra.Command, _ []string) error {
	fmt.Fprintf(out, "%s version %s\n", config.AppName, config.AppVersion)
	fmt.Fprintf(out, "build date: %s\n", config.BuildTime)
	fmt.Fprintf(out, "commit: %s\n\n", config.LastCommitHash)
	fmt.Fprintln(out, "https://github.com/y4ney/collect-cnnvd-vuln")

	return nil
}
