package main

import (
	"encoding/json"
	"fmt"

	"github.com/trilogy-group/cloudfix-linter/logger"
)

func extractModulePaths(jsonString string) ([]string, error) {
	appLogger := logger.New()
	var modulePaths []string
	byteValue := []byte(jsonString)
	/*
		Initialising a variable result which stores the data in the format of defined structure.
		Structure: "map(key->string,value->(array of map(key->string,value->interface))"
	*/
	var result map[string][]map[string]interface{}
	err := json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		appLogger.Error().Println("Failed to unmarshall module paths from json string")
		return modulePaths, err
	}
	appLogger.Info().Println("Unmarshalled module paths succesfully!")
	noOfModules := len(result["issues"])
	modulePaths = make([]string, noOfModules)
	for key, element := range result["issues"] {
		modulePaths[key] = fmt.Sprint(element["message"])
	}
	appLogger.Info().Println("Extracted module paths succesfully!")
	return modulePaths, nil
}
