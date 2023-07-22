package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cobra-demo",
	Short: "cobra-demo short help",
	Long:  "cobra-demo long help",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("test")
	},
}
var cfgFile string

// 持久标志,父级定义，子可用
// 本地标志：可以在本地分配一个标志，该标志仅适用于该特定命令
// 父命令上的本地标志: 默认情况下，Cobra 仅解析目标命令上的本地标志，而忽略父命令上的任何本地标志。通过启用 Command.TraverseChildren，Cobra 将在执行目标命令之前解析每个命令上的本地标志

func init() {
	//	rootCmd.PersistentFlags().Bool("viper", true, "是否使用viper来解析配置文件")
	//	rootCmd.PersistentFlags().StringP("author", "0", "your name", "作者")
	//	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "配置文件")
	//	rootCmd.Flags().StringP("source", "s", "", "source desc")
}

func Execute() {
	rootCmd.Execute()
}
