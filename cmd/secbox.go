package cmd

import (
	"gcrypto/secbox"

	"github.com/spf13/cobra"
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt <srcPath> <dstPath>",
	Short: "Encrypt a file",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := secbox.EncryptFile(args[0], args[1]); err != nil {
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
		if err := secbox.DecryptFile(args[0], args[1]); err != nil {
			return err
		}
		cmd.Println("File decrypted successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd, decryptCmd)
}
