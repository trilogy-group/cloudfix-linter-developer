package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
)

type TerraformManager struct {
	//No data types required
}

// Giving reference to terraform.exe file if present in windows
func terraform() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	basePath := filepath.Dir(ex)
	if runtime.GOOS == "windows" {
		return basePath + "\\terraform.exe"
	}
	return basePath + "/terraform"
}
func (t *TerraformManager) getTagToID(TfLintOutData []byte) (map[string]map[string][]string, error) {
	tagToID := make(map[string]map[string][]string)
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
			t.addPairToTagMap(rootResource, tagToID)
		}
	}
	// for all the resources present in child modules under the root module
	for _, childModule := range tfState.Values.RootModule.ChildModules {
		for _, childResource := range childModule.Resources {
			if childResource != nil {
				t.addPairToTagMap(childResource, tagToID)
			}
		}
	}
	return tagToID, nil
}

func (t *TerraformManager) addPairToTagMap(resource *tfjson.StateResource, tagToID map[string]map[string][]string) {
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

	var resourceType string = resource.Type
	resourceType = strings.Trim(resourceType,"\n")
	resourceType = strings.Trim(resourceType,`"`)

	var resourceName string = resource.Name
	resourceName = strings.Trim(resourceName,"\n")
	resourceName = strings.Trim(resourceName,`"`)

	innerMapKey := resourceType+"&"+resourceName

	_, exists := tagToID[yorTagTrim]

	if (!exists) {
		tagToID[yorTagTrim]= map[string][]string{}
	}

	tagToID[yorTagTrim][innerMapKey] = append(tagToID[yorTagTrim][innerMapKey], AWSResourceIDTrim)
}

func (t *TerraformManager) getTagToIDMapping() (map[string]map[string][]string, error) {
	tagToID := make(map[string]map[string][]string)
	var TfLintOutData []byte
	var errT error
	var modeBoolval bool
	val, present := os.LookupEnv("CLOUDFIX_TERRAFORM_LOCAL")
	if present {
		modeBoolval, _ = strconv.ParseBool(val)
	}
	if present && modeBoolval {
		TfLintOutData, errT = os.ReadFile("tf.show")
	} else {
		TfLintOutData, errT = exec.Command(terraform(), "show", "-json").Output()
	}
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
