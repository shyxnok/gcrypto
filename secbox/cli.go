package secbox

import (
	"os"

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
	Use:   "encrypt <srcPath> <dstPath>",
	Short: "Encrypt a file",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := EncryptFile(args[0], args[1]); err != nil {
			return err
		}
		cmd.Println("File encrypted successfully.")
		return nil
	},
}

var decryptCmd = &cobra.Command{
	Use:   "decrypt <srcPath> <dstPath>",
	Short: "Decrypt a file",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := DecryptFile(args[0], args[1]); err != nil {
			return err
		}
		cmd.Println("File decrypted successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd, decryptCmd)
}
