package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tenfyzhong/rime-dict-manager/dict"
)

var (
	addCode   string
	addWeight int
	addGroup  string
)

var addCmd = &cobra.Command{
	Use:   "add [word]",
	Short: "Add or update a word in the user dictionary",
	Long: `Adds a new word to the dictionary or updates it if it already exists.
If the Wubi code is not provided via --code, it will be automatically generated.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		wordToAdd := args[0]

		d := dict.NewDictionary(userDictFile)
		if err := d.Load(); err != nil {
			return err
		}

		finalCode := addCode
		if finalCode == "" {
			fmt.Println("Attempting to auto-generate Wubi code...")
			encoder, err := dict.NewWubiEncoder(mainDictFile)
			if err != nil {
				return fmt.Errorf("could not create wubi encoder: %w", err)
			}
			generated, err := encoder.GenerateCode(wordToAdd)
			if err != nil {
				return fmt.Errorf("failed to generate code: %w. Please provide it manually with --code", err)
			}
			finalCode = generated
			fmt.Printf("Auto-generated code for '%s': %s\n", wordToAdd, finalCode)
		}

		d.AddOrUpdate(wordToAdd, finalCode, addWeight, addGroup)

		fmt.Printf("Saving changes to %s...\n", userDictFile)
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
	addCmd.Flags().StringVarP(&addCode, "code", "c", "", "Manually specify the Wubi code")
	addCmd.Flags().IntVarP(&addWeight, "weight", "w", 100, "Specify the weight for the word")
	addCmd.Flags().StringVarP(&addGroup, "group", "g", "个人", "Specify the group for the word")
	rootCmd.AddCommand(addCmd)
}
