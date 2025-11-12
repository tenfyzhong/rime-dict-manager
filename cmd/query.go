package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tenfyzhong/rime-dict-manager/dict"
)

var queryCmd = &cobra.Command{
	Use:   "query [word]",
	Short: "Query a word in the user dictionary",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		wordToQuery := args[0]

		d := dict.NewDictionary(userDictFile)
		if err := d.Load(); err != nil {
			return err
		}

		found := false
		currentGroup := "Default" // Default group if no '##' is specified
		for _, entry := range d.Entries {
			if entry.IsGroup {
				currentGroup = entry.Group
				continue
			}

			if !entry.IsComment && entry.Word == wordToQuery {
				if !found {
					fmt.Printf("Found entries for '%s':\n", wordToQuery)
					found = true
				}
				fmt.Printf("- Word:   %s\n", entry.Word)
				fmt.Printf("  Code:   %s\n", entry.Code)
				fmt.Printf("  Weight: %d\n", entry.Weight)
				fmt.Printf("  Group:  %s\n", currentGroup)
				fmt.Println("---")
			}
		}

		if !found {
			fmt.Printf("Word '%s' not found in %s\n", wordToQuery, userDictFile)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)
}
