package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func initFlag(namespace, output *string, getCmd *cobra.Command) {
	getCmd.PersistentFlags().StringVarP(namespace, "namespace", "n", "default", "--namespace default")
	getCmd.PersistentFlags().StringVarP(output, "output", "o", "", "yaml|json|raw")
	getCmd.MarkFlagRequired("namespace")
}

func initCommand() {
	var namespace, output string
	var getCmd = &cobra.Command{
		Use:   "get",
		Short: "get short help",
		Long:  "get long help",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 0 {
				cmd.Help()
			}
		},
	}

	var getPodsCmd = &cobra.Command{
		Use:   "pods",
		Short: "pods short help",
		Long:  "pods long help",
		Run: func(cmd *cobra.Command, args []string) {
			if namespace != "" && output != "" {
				cmd.Printf("cobora-demo get pods -n %s -o %s", namespace, output)
			} else {
				cmd.Println("cobra-demo get pods")
			}
		},
		TraverseChildren: true,
	}

	var getNodesCmd = &cobra.Command{
		Use:   "nodes",
		Short: "nodes for short help",
		Long:  "nodes for long help",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("nodes")
		},
		TraverseChildren: true,
	}

	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(getPodsCmd)
	getCmd.AddCommand(getNodesCmd)

	initFlag(&namespace, &output, getCmd)
}

func init() {
	initCommand()
}
