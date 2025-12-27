package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var KeyCmd = &cobra.Command{
	Use:   "key:generate",
	Short: "Generate a new application crypto key",
	Long:  "Generate a new 32-byte base64 encoded key and optionally update .env file",
	Run:   runKeyGenerate,
}

var showOnly bool

// For testing
var (
	randReader io.Reader = rand.Reader
	exitFunc             = os.Exit
)

func init() {
	KeyCmd.Flags().BoolVar(&showOnly, "show", false, "Only display the key, don't update .env")
}

func runKeyGenerate(cmd *cobra.Command, args []string) {
	key, err := generateKey()
	if err != nil {
		fmt.Printf("Error generating key: %v\n", err)
		exitFunc(1)
		return
	}

	fullKey := "base64:" + key

	if showOnly {
		fmt.Println(fullKey)
		return
	}

	if err := updateEnvFile(fullKey); err != nil {
		fmt.Printf("Error updating .env: %v\n", err)
		fmt.Printf("\nGenerated key (add manually):\nCRYPTO_KEY=%s\n", fullKey)
		exitFunc(1)
		return
	}

	fmt.Println("Application key set successfully.")
}

func generateKey() (string, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(randReader, key); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

func updateEnvFile(key string) error {
	envPath := ".env"

	content, err := os.ReadFile(envPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf(".env file not found")
		}
		return err
	}

	lines := strings.Split(string(content), "\n")
	found := false
	for i, line := range lines {
		if strings.HasPrefix(line, "CRYPTO_KEY=") {
			lines[i] = "CRYPTO_KEY=" + key
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("CRYPTO_KEY not found in .env")
	}

	return os.WriteFile(envPath, []byte(strings.Join(lines, "\n")), 0644)
}
