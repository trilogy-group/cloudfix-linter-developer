package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/trilogy-group/cloudfix-linter-developer/cloudfixIntegration"
	"github.com/trilogy-group/cloudfix-linter-developer/logger"
)

// Structures for Marshalling JSON outputs

// JSONIssue is a temporary structure for converting TFLint issues to JSON.
type JSONIssue struct {
	Rule    JSONRule    `json:"rule"`
	Message string      `json:"message"`
	Range   JSONRange   `json:"range"`
	Callers []JSONRange `json:"callers"`
}

// JSONRule is a temporary structure for converting TFLint rules to JSON.
type JSONRule struct {
	Name     string `json:"name"`
	Severity string `json:"severity"`
	Link     string `json:"link"`
}

// JSONRange is a temporary structure for converting ranges to JSON.
type JSONRange struct {
	Filename string  `json:"filename"`
	Start    JSONPos `json:"start"`
	End      JSONPos `json:"end"`
}

// JSONPos is a temporary structure for converting positions to JSON.
type JSONPos struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// JSONError is a temporary structure for converting errors to JSON.
type JSONError struct {
	Summary  string     `json:"summary,omitempty"`
	Message  string     `json:"message"`
	Severity string     `json:"severity"`
	Range    *JSONRange `json:"range,omitempty"` // pointer so omitempty works
}

// JSONOutput is a temporary structure for converting to JSON.
type JSONOutput struct {
	Issues []JSONIssue `json:"issues"`
	Errors []JSONError `json:"errors"`
}

type Orchestrator struct {
	// No Data Fields are required for this class
}

// Giving reference to tflint.exe file if present in windows
func tflint() string {
	if runtime.GOOS == "windows" {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		basePath := filepath.Dir(ex)
		return basePath + "\\tflint.exe"
	}
	return "tflint"
}

// Memeber functions for the Orchestrator class follow:

func (o *Orchestrator) extractModulePaths(jsonString []byte) ([]string, error) {
	dlog := logger.DevLogger()
	var modulePaths []string
	//byteValue := []byte(jsonString)
	/*
		Initialising a variable result which stores the data in the format of defined structure.
		Structure: "map(key->string,value->(array of map(key->string,value->interface))"
	*/
	var result map[string][]map[string]interface{}
	if len(jsonString) == 0 {
		//Empty string has been sent. No modules present
		return modulePaths, nil
	}
	err := json.Unmarshal(jsonString, &result)
	if err != nil {
		dlog.Error("Failed to unmarshall module paths from json string", err)
		return modulePaths, err
	}
	dlog.Debug("Unmarshalled module paths succesfully!")
	noOfModules := len(result["issues"])
	modulePaths = make([]string, noOfModules)
	for key, element := range result["issues"] {
		modulePaths[key] = fmt.Sprint(element["message"])
	}
	dlog.Error("Extracted module paths succesfully!")
	return modulePaths, nil
}

func (o *Orchestrator) runReccos(jsonFlag bool) {
	dlog := logger.DevLogger()
	var persist Persistance
	var cloudfixMan cloudfixIntegration.CloudfixManager
	var terraMan TerraformManager
	reccosFileName := "recos.txt"
	reccosMapping, errC := cloudfixMan.GetReccos()
	if errC != nil {
		logger.Info("Something went wrong. More logs in the log directory. ", errC)
		dlog.Error("Failed to get Reccomendations from CloudFix: ", errC)
		fmt.Printf(`{"error": "%s"}`, errC)
		return
	}
	if len(reccosMapping) == 0 {
		//log that no reccomendations could be received
		logger.Info("No oppurtunities found")
		fmt.Println(`[]`)
		//exit gracefully
		return
	}
	errP := persist.store_reccos(reccosMapping, reccosFileName)
	if errP != nil {
		panic(errP)
	}
	os.Setenv("ReccosMapFile", reccosFileName)
	tagFileName := "tagsID.txt"
	tagToIDMap, errG := terraMan.getTagToIDMapping()
	if errG != nil {
		panic(errG)
	}
	errT := persist.store_tagMap(tagToIDMap, tagFileName)
	if errT != nil {
		panic(errT)
	}
	os.Setenv("TagsMapFile", tagFileName)
	modulesJson, _ := exec.Command(tflint(), "--only=module_source", "-f=json").Output()
	modulePaths, errM := o.extractModulePaths(modulesJson)
	if errM != nil {
		//log failure in extracting module paths
		dlog.Error("Failed to extract module paths", errM)
		return
	}
	if !jsonFlag {
		output, _ := exec.Command(tflint(), "--module", "--disable-rule=module_source").Output()
		fmt.Print(string(output))
		for _, module := range modulePaths {
			outputM, _ := exec.Command(tflint(), module, "--module", "--disable-rule=module_source").Output()
			fmt.Print(string(outputM))
		}
	} else {
		var flaggedIssues []JSONIssue
		output, _ := exec.Command(tflint(), "--module", "--disable-rule=module_source", "-f=json").Output()
		var jsonOutRoot JSONOutput
		errMR := json.Unmarshal(output, &jsonOutRoot)
		if errMR != nil {
			fmt.Println("Error getting JSON output")
			return

		}
		if len(jsonOutRoot.Issues) != 0 {
			flaggedIssues = append(flaggedIssues, jsonOutRoot.Issues...)
		}
		for _, module := range modulePaths {
			outputM, _ := exec.Command(tflint(), module, "--module", "--disable-rule=module_source", "-f=json").Output()
			var jsonOutModules JSONOutput
			errMM := json.Unmarshal(outputM, &jsonOutModules)
			if errMM != nil {
				fmt.Println("Error getting JSON output")
				return
			}
			if len(jsonOutModules.Issues) != 0 {
				flaggedIssues = append(flaggedIssues, jsonOutModules.Issues...)
			}
		}
		if len(flaggedIssues) == 0 {
			fmt.Println(`[]`)
			return
		}
		out, err := json.Marshal(flaggedIssues)
		if err != nil {
			logger.Info("Failed to parse and display opportunities. More info in the log directory")
			dlog.Error("Failed to parse flagged issues: ", err)
			return
		}
		fmt.Println(string(out[:]))
	}
}
