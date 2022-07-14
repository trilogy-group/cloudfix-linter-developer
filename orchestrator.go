package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
)

type Parameter struct {
	IdealType      string `json:"Migrating to instance type"`
	RecentSnapshot bool   `json:"Has recent snapshot"`
}

//structure for unmarshalling the reccomendation json response from cloudfix
type ResponseReccos struct {
	Id                     string
	Region                 string
	PrimaryImpactedNodeId  string
	OtherImpactedNodeIds   []string
	ResourceId             string
	ResourceName           string
	Difficulty             int
	Risk                   int
	ApplicationEnvironment string
	AnnualSavings          float32
	AnnualCost             float32
	Status                 string
	Parameters             Parameter
	TemplateApproved       bool
	CustomerId             int
	AccountId              string
	AccountNickname        string
	OpportunityType        string
	OpportunityDescription string
	GeneratedDate          string
	LastUpdatedDate        string
}

type Orchestrator struct {
	// No Data Fields are required for this class
}

// Memeber functions for the Orchestrator class follow:

func (o *Orchestrator) parseReccos(reccos []byte) map[string]map[string]string {
	//appLogger := logger.New()
	mapping := map[string]map[string]string{}
	var responses []ResponseReccos
	err := json.Unmarshal(reccos, &responses)
	if err != nil {
		//appLogger.Error().Println("Failed to unmarshall reccomendations")
		panic(err)
	}
	//appLogger.Info().Println("Reccomendations unamrshalled succesfully!")
	//fmt.Println(oppurMap["Ec2IntelToAmd"])
	for _, recco := range responses {
		awsID := recco.ResourceId
		oppurType := recco.OpportunityType
		temp := map[string]string{}
		switch oppurType {
		case "Gp2Gp3":
			temp["type"] = "gp3"
			mapping[awsID] = temp
		case "Ec2IntelToAmd":
			var idealType = recco.Parameters.IdealType
			temp["instance_type"] = idealType
			mapping[awsID] = temp
		default:
			//appLogger.Warning().Printf("Unknown Oppurtunity Type for resource ID: \"%s\"", awsID)
			temp["NoAttributeMarker"] = recco.OpportunityDescription
			mapping[awsID] = temp
		}
	}
	//appLogger.Info().Println("Reccomendation mapping made!")
	return mapping
}

func (o *Orchestrator) extractModulePaths(jsonString []byte) ([]string, error) {
	//appLogger := logger.New()
	var modulePaths []string
	//byteValue := []byte(jsonString)
	/*
		Initialising a variable result which stores the data in the format of defined structure.
		Structure: "map(key->string,value->(array of map(key->string,value->interface))"
	*/
	var result map[string][]map[string]interface{}
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
	//for root module resources
	for _, rootResource := range tfState.Values.RootModule.Resources {
		o.addPairToTagMap(rootResource, tagToID)

	}
	// for all the resources present in child modules under the root module
	for _, childModule := range tfState.Values.RootModule.ChildModules {
		for _, childResource := range childModule.Resources {
			o.addPairToTagMap(childResource, tagToID)
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

func main() {

	var orches Orchestrator
	var persist Persistance
	reccosFileName := "recos.txt"
	currPWD, _ := exec.Command("pwd").Output()
	currPWDStr := string(currPWD[:])
	currPWDStrip := strings.Trim(currPWDStr, "\n")
	currPWDStrip += "/reccos.json"
	fileR, errR := ioutil.ReadFile(currPWDStrip)
	if errR != nil {
		//Add Error Log
		panic(errR)
	}
	reccosMapping := orches.parseReccos(fileR)
	errP := persist.store_reccos(reccosMapping, reccosFileName)
	if errP != nil {
		panic(errP)
	}
	os.Setenv("ReccosMapFile", reccosFileName)
	tagFileName := "tagsID.txt"
	tagToIDMap, errG := orches.getTagToID()
	if errG != nil {
		panic(errG)
	}
	errT := persist.store_tagMap(tagToIDMap, tagFileName)
	if errT != nil {
		panic(errT)
	}
	os.Setenv("TagsMapFile", tagFileName)
	output, _ := exec.Command("tflint", "--module", "--disable-rule=module_source").Output()
	fmt.Print(string(output))
}
