package secbox

import (
	"fmt"
	"os"

	"gcrypto/encrypto"
)

func EncryptFile(srcPath, dstPath string) error {
	raw, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	encrypted, err := encrypto.EncryptText(raw)
	if err != nil {
		return fmt.Errorf("failed to encrypt file: %w", err)
	}

	return os.WriteFile(dstPath, encrypted, 0644)
}
func DecryptFile(srcPath, dstPath string) error {
	encBin, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	raw, err := encrypto.DecryptText(encBin)
	if err != nil {
		return fmt.Errorf("failed to decrypt file: %w", err)
	}
	return os.WriteFile(dstPath, raw, 0644)
}
