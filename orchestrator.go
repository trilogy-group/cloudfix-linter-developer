package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
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

func (o *Orchestrator) getTagToID() (map[string]string, error) {
	tagToID := make(map[string]string)
	TfLintOutData, errT := exec.Command("terraform", "show", "-json").Output()
	if errT != nil {
		return tagToID, errT
	}
	var tfState tfjson.State
	errU := tfState.UnmarshalJSON(TfLintOutData)
	if errU != nil {
		return tagToID, errU
	}
	if tfState.Values == nil {
		//log that no resources have been deployed
		return tagToID, nil
	}
	//for root module resources
	for _, rootResource := range tfState.Values.RootModule.Resources {
		if rootResource != nil {
			o.addPairToTagMap(rootResource, tagToID)
		}
	}
	// for all the resources present in child modules under the root module
	for _, childModule := range tfState.Values.RootModule.ChildModules {
		for _, childResource := range childModule.Resources {
			if childResource != nil {
				o.addPairToTagMap(childResource, tagToID)
			}
		}
	}
	return tagToID, nil
}

func (o *Orchestrator) addPairToTagMap(resource *tfjson.StateResource, tagToID map[string]string) {
	AWSResourceIDRaw, ok := resource.AttributeValues["id"]
	if !ok {
		//log that id is not present
		return
	}
	AWSResourceID := AWSResourceIDRaw.(string)
	tagsRaw, ok := resource.AttributeValues["tags"]
	if !ok {
		//log that tags are not present
		return
	}
	tags := tagsRaw.(map[string]interface{})
	yorTagRaw, ok := tags["yor_trace"]
	if !ok {
		//log that yor_trace is not present
		return
	}
	yorTag := yorTagRaw.(string)
	AWSResourceIDStrip := strings.Trim(AWSResourceID, "\n")
	AWSResourceIDTrim := strings.Trim(AWSResourceIDStrip, `"`)
	yorTagStrip := strings.Trim(yorTag, "\n")
	yorTagTrim := strings.Trim(yorTagStrip, `"`)
	if yorTagTrim == "" || AWSResourceIDTrim == "" {
		return
	}
	tagToID[yorTagTrim] = AWSResourceIDTrim
}

func (o *Orchestrator) runReccos() {

	var persist Persistance
	var cloudfixMan CloudfixManager
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
	tagToIDMap, errG := o.getTagToID()
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
