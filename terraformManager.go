package main

import (
	"os/exec"
	"strings"
	"os"
	"runtime"
	"path/filepath"
	tfjson "github.com/hashicorp/terraform-json"
)

type TerraformManager struct {
	//No data types required
}
// Giving reference to terraform.exe file if present in windows
func terraform() string{
	if(runtime.GOOS=="windows"){
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		basePath := filepath.Dir(ex)
		return basePath+"\\terraform.exe"
	}
	return "terraform"
}
func (t *TerraformManager) getTagToID(TfLintOutData []byte) (map[string]string, error) {
	tagToID := make(map[string]string)
	tagCountMap := make(map[string]int)
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
			t.addPairToTagMap(rootResource, tagToID, tagCountMap)
		}
	}
	// for all the resources present in child modules under the root module
	for _, childModule := range tfState.Values.RootModule.ChildModules {
		for _, childResource := range childModule.Resources {
			if childResource != nil {
				t.addPairToTagMap(childResource, tagToID, tagCountMap)
			}
		}
	}
	return tagToID, nil
}

func (t *TerraformManager) addPairToTagMap(resource *tfjson.StateResource, tagToID map[string]string, tagCountMap map[string]int) {
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

	tags, ok := tagsRaw.(map[string]interface{})
	if !ok {
		//log that tags are not present
		return
	}
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
	var tagCount int = tagCountMap[yorTagTrim]
	tagCountMap[yorTagTrim] += 1
	if tagCount!=0 {
		yorTagTrim += "$"+strconv.Itoa(tagCount)
	}
	tagToID[yorTagTrim] = AWSResourceIDTrim
}

func (t *TerraformManager) getTagToIDMapping() (map[string]string, error) {
	tagToID := make(map[string]string)
	TfLintOutData, errT := exec.Command(terraform(), "show", "-json").Output()
	if errT != nil {
		//Add Log
		return tagToID, errT
	}
	tagToID, err := t.getTagToID(TfLintOutData)
	if err != nil {
		//Add Log
		return tagToID, err
	}
	return tagToID, nil
}
