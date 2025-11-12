package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tenfyzhong/rime-dict-manager/dict"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [word]",
	Short: "Delete a word from the user dictionary",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		wordToDelete := args[0]

		d := dict.NewDictionary(userDictFile)
		if err := d.Load(); err != nil {
			return err
		}

		var newEntries []dict.Entry
		found := false
		for _, entry := range d.Entries {
			if entry.Word == wordToDelete {
				found = true
				// Skip this entry to delete it
				continue
			}
			newEntries = append(newEntries, entry)
		}

		if !found {
			return fmt.Errorf("word '%s' not found in the dictionary", wordToDelete)
		}

		d.Entries = newEntries

		fmt.Printf("Deleting word '%s'...\n", wordToDelete)
		if err := d.Save(); err != nil {
			return fmt.Errorf("failed to save dictionary: %w", err)
		}
		fmt.Println("Successfully saved.")

		if !noDeploy {
			fmt.Println("Triggering Rime redeployment...")
			if err := runDeployCommand(); err != nil {
				return fmt.Errorf("deployment failed: %w", err)
			}
			fmt.Println("Deployment command executed.")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
