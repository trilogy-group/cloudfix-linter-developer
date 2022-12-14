package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/trilogy-group/cloudfix-linter-developer/logger"
)

func yor() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	basePath := filepath.Dir(ex)
	if runtime.GOOS == "windows" {
		return basePath + "\\yor.exe"
	}
	return basePath + "/yor"
}

// rootCmd represents the base command when called without any subcommands
var (
	os_type = runtime.GOOS
	rootCmd = &cobra.Command{
		Use:   "cloudfix-linter",
		Short: "This tool helps flag reccomendations from Cloudfix in your terraform code",
		Long:  "This tool helps flag reccomendations from Cloudfix in your terraform code",
	}
	jsonFlag  bool
	recccoCmd = &cobra.Command{
		Use:   "recco",
		Short: "To flag reccomendations",
		Long:  "Running this command will parse through your terraform code and flag reccomendations from Cloudfix for resources that it finds",
		Run: func(cmd *cobra.Command, args []string) {
			dirname, err := os.UserHomeDir()
			if err != nil {
				log.Fatal(err)
			}
			logger.InitLogger(dirname, jsonFlag)
			logger.Info("Cloudfix-linter starting")
			homeDir, err := os.Getwd()
			if err != nil {
				fmt.Println(`{ "error": "Failed. Unable to find working directory"}`)
			}
			hclFilePath := homeDir + "/.tflint.hcl"
			if _, err := os.Stat(hclFilePath); errors.Is(err, os.ErrNotExist) {
				fmt.Println(`{ "error" : "The current directory needs to be initialised. Run "cloudfix-linter init" to initialise" }`)
				return
			}
			var orches Orchestrator
			orches.runReccos(jsonFlag)
		},
	}
	currptFlag = &cobra.Command{
		Use:   "addTags",
		Short: "Add tags to your terraform code to trace them back to the cloud",
		Long:  "Add tags to your terraform code to trace them back to the cloud. You will need to run this command if the tool detects that there are no tags for a resource in your terraform code. You will be asked to run this command in that instance",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := exec.Command(yor(), "tag", "-d", ".", "--tag-groups", "code2cloud").Output()
			if err != nil {
				return errors.New(fmt.Sprintf(`{ "error" : "%s", "message" : "yor tag error" }`, err))
			}
			return nil
		},
	}
	initCmd = &cobra.Command{
		Use:   "init",
		Short: "To initialise the directory. Run this before asking for recommendations",
		Long:  "Running this command will initialise the directory and add tags to your terraform resources",
		Run: func(cmd *cobra.Command, args []string) {
			tflintRecco, e := cmd.Flags().GetString("enableNonCloudfixRecco")
			if e != nil {
				fmt.Println(`{ "error": "Failed to initialise", "message": "error in finding enableNonCloudfixRecco"}`)
			}
			var default_recco bool = true
			if tflintRecco != "" {
				if tflintRecco == "true" {
					default_recco = false
				}
			}
			err := initDir(default_recco)
			if err != nil {
				fmt.Println(fmt.Sprintf(`{ "error" : "%s", "message" : "Failed to initialise" }`, err))
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(recccoCmd)
	rootCmd.AddCommand(currptFlag)
	rootCmd.AddCommand(initCmd)
	initCmd.PersistentFlags().String("enableNonCloudfixRecco", "", "Enables Extra recommendations coming from tflint")
	recccoCmd.Flags().BoolVarP(&jsonFlag, "json", "j", false, "to get output in json format")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error occurred while execution")
	}
}
