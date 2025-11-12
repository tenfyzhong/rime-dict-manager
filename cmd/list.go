package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/spf13/cobra"
	"github.com/tenfyzhong/rime-dict-manager/dict"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "漂亮地列出用户词典中的所有词条",
	Long:  `读取并以易于阅读的格式显示用户词典文件中的所有词条, 按分组呈现.`,
	Run: func(cmd *cobra.Command, args []string) {
		d := dict.NewDictionary(userDictFile)
		if err := d.Load(); err != nil {
			log.Fatalf("加载词典文件失败: %v", err)
		}

		fmt.Printf("词典文件: %s\n\n", userDictFile)

		// Define column widths
		wordWidth := 25
		codeWidth := 20

		// Helper function for printing padded strings
		printWithPaddedWidth := func(s string, width int) {
			fmt.Print(s)
			// Calculate padding needed based on visual width
			pad := width - runewidth.StringWidth(s)
			if pad > 0 {
				fmt.Print(strings.Repeat(" ", pad))
			}
		}

		// Print header
		printWithPaddedWidth("词语 (Word)", wordWidth)
		printWithPaddedWidth("编码 (Code)", codeWidth)
		fmt.Println("权重 (Weight)")
		fmt.Println(strings.Repeat("-", wordWidth+codeWidth+10))

		for _, entry := range d.Entries {
			if entry.IsGroup {
				totalWidth := wordWidth + codeWidth + 10
				groupString := fmt.Sprintf(" %s ", entry.Group)
				// Calculate asterisks for centering
				asteriskCount := (totalWidth - runewidth.StringWidth(groupString)) / 2
				if asteriskCount < 0 { // Ensure at least some asterisks if group name is very long
					asteriskCount = 0
				}
				remainingAsterisks := totalWidth - asteriskCount - runewidth.StringWidth(groupString)
				if remainingAsterisks < 0 {
					remainingAsterisks = 0
				}
				fmt.Printf("\n%s%s%s\n", strings.Repeat("*", asteriskCount), groupString, strings.Repeat("*", remainingAsterisks))
			} else if !entry.IsComment && entry.Word != "" {
				printWithPaddedWidth(entry.Word, wordWidth)
				printWithPaddedWidth(entry.Code, codeWidth)
				fmt.Printf("%d\n", entry.Weight)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
