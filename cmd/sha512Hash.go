package cmd

import (
	"crypto/sha512"

	"github.com/spf13/cobra"
)

// sha512HashCmd represents the sha512Hash command
var sha512HashCmd = &cobra.Command{
	Use:   "sha512 [FILE]",
	Short: "Display SHA-512 checksums (512 bits).",
	Long: `Display SHA-512 checksums (512 bits).

Without FILE or when FILE is '-', read the standard input.
If the list of FILE contains a directory, it will be proceed recursively.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		filesToCheck, err := getFilesToCompute(args)
		if err != nil {
			return err
		}

		numJobs := len(filesToCheck)
		jobs := make(chan JobsParam, numJobs)
		results := make(chan HashResult, numJobs)

		initWorkers(jobs, results)

		for _, filePath := range filesToCheck {
			h := sha512.New()
			jobs <- JobsParam{filePath, h}
		}
		close(jobs)

		return waitForResult(numJobs, results)
	},
}

func init() {
	hashCmd.AddCommand(sha512HashCmd)
}
