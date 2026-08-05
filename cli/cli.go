package cli

import (
	"os"

	"gcrypto/dlh"
	"gcrypto/elh"

	"github.com/spf13/cobra"
)

// rootCmd 为根命令，子命令见下方 encryptCmd/decryptCmd。
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

var encryptCmd = &cobra.Command{
	Use:   "encrypt <srcPath> <suffix>",
	Short: "Encrypt a file to any suffix",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		dst, err := elh.EncryptFileExt(args[0], args[1])
		if err != nil {
			return err
		}
		cmd.Printf("File encrypted successfully: %s\n", dst)
		return nil
	},
}

var decryptCmd = &cobra.Command{
	Use:   "decrypt <srcPath>",
	Short: "Decrypt a file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dst, err := dlh.DecryptFile(args[0])
		if err != nil {
			return err
		}
		cmd.Printf("File decrypted successfully: %s\n", dst)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd, decryptCmd)
}
