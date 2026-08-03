package cmd

import (
	"fmt"

	"gcrypto/secbox"

	"github.com/spf13/cobra"
)

var secboxCmd = &cobra.Command{
	Use:   "secbox",
	Short: "Encrypt/Decrypt files using secbox",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 3 {
			fmt.Println("Usage: xbs secbox <encrypt|decrypt> <srcPath> <dstPath>")
			return
		}

		action := args[0]
		srcPath := args[1]
		dstPath := args[2]

		switch action {
		case "encrypt":
			err := secbox.EncryptFile(srcPath, dstPath)
			if err != nil {
				fmt.Printf("Error encrypting file: %v\n", err)
				return
			}
			fmt.Println("File encrypted successfully.")
		case "decrypt":
			err := secbox.DecryptFile(srcPath, dstPath)
			if err != nil {
				fmt.Printf("Error decrypting file: %v\n", err)
				return
			}
			fmt.Println("File decrypted successfully.")
		default:
			fmt.Println("Invalid action. Use 'encrypt' or 'decrypt'.")
		}
	},
}

func init() {
	rootCmd.AddCommand(secboxCmd)
}
