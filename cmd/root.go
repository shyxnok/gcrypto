package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gcrypto",
	Short: "A small tool for file encryption and decryption",
}

// Execute 运行根命令。错误由 cobra 打印到 stderr，并以退出码 1 结束。
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
