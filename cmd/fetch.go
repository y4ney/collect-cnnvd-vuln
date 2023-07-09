package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch CNNVD vuln info",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("fetch called")
	},
}

func init() {
}
