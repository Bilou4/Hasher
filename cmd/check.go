package cmd

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:           "check",
	Short:         "",
	Long:          ``,
	SilenceErrors: true,
	SilenceUsage:  true,
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

func checkHash(h hash.Hash) error {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		input = strings.ReplaceAll(input, "\n", "")
		inputs := strings.Split(input, " ")
		if len(inputs) != 2 {
			return fmt.Errorf("malformed entry: HASH FilePath")
		}
		hashToCheck := strings.Trim(inputs[0], " ")
		path := strings.Trim(inputs[1], " ")
		if hex.DecodedLen(len(hashToCheck)) != h.Size() {
			return fmt.Errorf("malformed hash. Expected %d bytes, got %d", h.Size(), hex.DecodedLen(len(hashToCheck)))
		}
		hash, err := computeHash(path, h)
		if err != nil {
			return err
		}
		if hashToCheck != fmt.Sprintf("%x", hash) {
			return fmt.Errorf("%s FAILED", path)
		}
		fmt.Printf("%s OK\n", path)
		h.Reset()
	}
	return nil
}
