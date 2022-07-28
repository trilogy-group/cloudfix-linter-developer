package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/trilogy-group/cloudfix-linter/cloudfixIntegration"
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

// Memeber functions for the Orchestrator class follow:

func (o *Orchestrator) extractModulePaths(jsonString []byte) ([]string, error) {
	//appLogger := logger.New()
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
		//appLogger.Error().Println("Failed to unmarshall module paths from json string")
		return modulePaths, err
	}
	//appLogger.Info().Println("Unmarshalled module paths succesfully!")
	noOfModules := len(result["issues"])
	modulePaths = make([]string, noOfModules)
	for key, element := range result["issues"] {
		modulePaths[key] = fmt.Sprint(element["message"])
	}
	//appLogger.Info().Println("Extracted module paths succesfully!")
	return modulePaths, nil
}

func (o *Orchestrator) runReccos(jsonFlag bool) {
	var persist Persistance
	var cloudfixMan cloudfixIntegration.CloudfixManager
	var terraMan TerraformManager
	reccosFileName := "recos.txt"
	reccosMapping, errC := cloudfixMan.GetReccos()
	if errC != nil {
		fmt.Println(errC.Message)
		return
	}
	if len(reccosMapping) == 0 {
		//log that no reccomendations could be received
		fmt.Println("No oppurtunities exist for your system")
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
	modulesJson, _ := exec.Command("tflint", "--only=module_source", "-f=json").Output()
	modulePaths, errM := o.extractModulePaths(modulesJson)
	if errM != nil {
		//log failure in extracting module paths
		return
	}
	if !jsonFlag {
		output, _ := exec.Command("tflint", "--module", "--disable-rule=module_source").Output()
		fmt.Print(string(output))
		for _, module := range modulePaths {
			outputM, _ := exec.Command("tflint", module, "--module", "--disable-rule=module_source").Output()
			fmt.Print(string(outputM))
		}
	} else {
		var flaggedIssues []JSONIssue
		output, _ := exec.Command("tflint", "--module", "--disable-rule=module_source", "-f=json").Output()
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
			outputM, _ := exec.Command("tflint", module, "--module", "--disable-rule=module_source", "-f=json").Output()
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
			fmt.Println(`{}`)
			return
		}
		out, err := json.Marshal(flaggedIssues)
		if err != nil {
			fmt.Println("Error getting JSON output")
			return
		}
		fmt.Println(string(out[:]))
	}
}
