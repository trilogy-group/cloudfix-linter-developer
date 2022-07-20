package main

import (
	"os/exec"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var (
	rootCmd = &cobra.Command{
		Use:   "cloudfix-linter",
		Short: "This tool helps flag reccomendations from Cloudfix in your terraform code",
		Long:  "This tool helps flag reccomendations from Cloudfix in your terraform code",
	}
	recccoCmd = &cobra.Command{
		Use:   "flagRecco",
		Short: "To flag reccomendations",
		Long:  "Running this command will parse through your terraform code and flag reccomendations from Cloudfix for resources that it finds",
		Run: func(cmd *cobra.Command, args []string) {
			var orches Orchestrator
			orches.runReccos()
		},
	}
	currptFlag = &cobra.Command{
		Use:   "addTags",
		Short: "Add tags to your terraform code to trace them back to the cloud",
		Long:  "Add tags to your terraform code to trace them back to the cloud. You will need to run this command if the tool detects that there are no tags for a resource in your terraform code. You will be asked to run this command in that instance",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := exec.Command("yor", "tag", "-d", ".", "--tag-groups", "code2cloud").Output()
			if err != nil {
				return err
			}
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(recccoCmd)
	rootCmd.AddCommand(currptFlag)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
