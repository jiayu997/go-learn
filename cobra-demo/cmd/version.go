package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version short help",
	Long:  "version long help",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			log.Fatal("version can't use args")
		} else {
			fmt.Println("cobra-demo version alhpa1")
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
