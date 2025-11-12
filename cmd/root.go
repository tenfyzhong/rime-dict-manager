package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/tenfyzhong/rime-dict-manager/config"
)

var (
	// These variables are now accessible to all commands in the cmd package.
	userDictFile  string
	mainDictFile  string
	deployCommand string
	noDeploy      bool
)

var rootCmd = &cobra.Command{
	Use:   "rime-dict-manager",
	Short: "A CLI tool to manage Rime user dictionaries.",
	Long: `A command-line tool to query, add, modify, and delete entries
in a Rime user dictionary file, with automatic Wubi code generation
and Rime redeployment capabilities.`,
	Version: config.Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}

func init() {
	defaultUserDictFile := os.ExpandEnv("$HOME/Library/Rime/wubi86_jidian_user.dict.yaml")
	defaultMainDictFile := os.ExpandEnv("$HOME/Library/Rime/wubi86_jidian.dict.yaml")

	rootCmd.PersistentFlags().StringVarP(&userDictFile, "file", "f", defaultUserDictFile, "Path to the Rime user dictionary file.")
	rootCmd.PersistentFlags().StringVar(&mainDictFile, "main-dict", defaultMainDictFile, "Path to the main dictionary for Wubi code generation.")
	rootCmd.PersistentFlags().StringVar(&deployCommand, "deploy-cmd", `/Library/Input\ Methods/Squirrel.app/Contents/MacOS/Squirrel --reload`, "The command to execute for Rime redeployment.")
	rootCmd.PersistentFlags().BoolVar(&noDeploy, "no-deploy", false, "Disable automatic Rime redeployment after an operation.")
}

func runDeployCommand() error {
	fmt.Printf("Executing deployment command: %s\n", deployCommand)

	var cmd *exec.Cmd
	cmd = exec.Command("bash", "-c", deployCommand)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command failed: %w\nOutput: %s", err, string(output))
	}

	fmt.Println(string(output))
	return nil
}
