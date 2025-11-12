package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/tenfyzhong/rime-dict-manager/dict"
)

var setWeightCmd = &cobra.Command{
	Use:   "set-weight [word] [weight]",
	Short: "Set the weight for a word in the dictionary",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		wordToUpdate := args[0]
		newWeight, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid weight value: %s. Must be an integer", args[1])
		}

		d := dict.NewDictionary(userDictFile)
		if err := d.Load(); err != nil {
			return err
		}

		found := false
		for i := range d.Entries {
			if d.Entries[i].Word == wordToUpdate {
				d.Entries[i].Weight = newWeight
				found = true
				// break // Uncomment if you only want to update the first occurrence
			}
		}

		if !found {
			return fmt.Errorf("word '%s' not found in the dictionary", wordToUpdate)
		}

		fmt.Printf("Updating weight for '%s' to %d...\n", wordToUpdate, newWeight)
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
	rootCmd.AddCommand(setWeightCmd)
}
