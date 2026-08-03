package cmd

import (
	"fmt"

	"gcrypto/secbox"

	"github.com/spf13/cobra"
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt a file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("Usage:  secbox encrypt <srcPath> <dstPath>")
			return
		}
		err := secbox.EncryptFile(args[0], args[1])
		if err != nil {
			fmt.Printf("Error encrypting file: %v\n", err)
			return
		}
		fmt.Println("File encrypted successfully.")
	},
}

var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt a file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("Usage:  secbox decrypt <srcPath> <dstPath>")
			return
		}
		err := secbox.DecryptFile(args[0], args[1])
		if err != nil {
			fmt.Printf("Error decrypting file: %v\n", err)
			return
		}
		fmt.Println("File decrypted successfully.")
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)
	rootCmd.AddCommand(decryptCmd)
}
