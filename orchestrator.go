package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

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

func (o *Orchestrator) runReccos() {

	var persist Persistance
	var cloudfixMan CloudfixManager
	var terraMan TerraformManager
	reccosFileName := "recos.txt"
	reccosMapping := cloudfixMan.parseReccos()
	if len(reccosMapping) == 0 {
		//log that no reccomendations could be received
		//exit gracefully
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
	output, _ := exec.Command("tflint", "--module", "--disable-rule=module_source").Output()
	fmt.Print(string(output))
	for _, module := range modulePaths {
		outputM, _ := exec.Command("tflint", module, "--module", "--disable-rule=module_source").Output()
		fmt.Print(string(outputM))
	}
}
