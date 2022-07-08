package main

import (
	"encoding/json"
	"fmt"

	"github.com/trilogy-group/cloudfix-linter/logger"
)

func main() {
	//modulePaths := extractModulePaths(jsonString)
}

func extractModulePaths(jsonString string) []string {
	appLogger := logger.New()
	byteValue := []byte(jsonString)
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
	appLogger.Info().Println("Successfully extracted module paths from json string containing module paths")
	return modulePaths
}
