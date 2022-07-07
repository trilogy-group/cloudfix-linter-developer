package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/trilogy-group/cloudfix-linter/logger"
)

func main() {
	_, err := extractModulePaths("sample.json")
	if err != nil {
		fmt.Println("Could not Extract module names:- ", err)
		return
	}
}

func extractModulePaths(fileName string) ([]string, error) {
	var emptyArray []string
	appLogger := logger.New()
	jsonFile, err := os.Open(fileName)
	if err != nil {
		appLogger.Error().Println("Error occurred while opening json file containing module names in orchestrator")
		return emptyArray, err
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		appLogger.Error().Println("Error occurred while Reading from json file containing module names in orchestrator")
		return emptyArray, err
	}
	/*
		Initialising a variable result which stores the data in the format of defined structure.
		Structure: "map(key->string,value->(array of map(key->string,value->interface))"
	*/
	var result map[string][]map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)
	noOfModules := len(result["issues"])
	modulePaths := make([]string, noOfModules)
	for key, element := range result["issues"] {
		modulePaths[key] = fmt.Sprint(element["message"])
	}
	return modulePaths, nil
}
