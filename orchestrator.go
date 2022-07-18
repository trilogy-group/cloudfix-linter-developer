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

//Structure for unmarshalling the oppurtunityType to Attributes mapping (the mapping is present in "mappingAttributes.json")
type IdealAttributes struct {
	AttributeType  string `json:"Attribute Type"`
	AttributeValue string `json:"Attribute Value"`
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
	Parameters             map[string]interface{}
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

func (o *Orchestrator) parseReccos(reccos []byte, attrMapping []byte) map[string]map[string]string {
	//function to process the reccomendations from cloudfix and turn that into a map
	//the structure of the map is resourceID -> Attribute type that needs to be targetted -> Ideal Attribute Value
	// If there is no attribute that has to be targetted, attribute type would be filled with "NoAttributeMarker" and
	//Attribute Value would be filled with any message that in the end has to be displayed to the user
	mapping := map[string]map[string]string{} //this is the map that has to be returned in the end
	var responses []ResponseReccos
	if len(reccos) == 0 {
		//log that no reccomendations have been received
		return mapping
	}
	errR := json.Unmarshal(reccos, &responses) //the reccomendations from cloudfix are being unmarshalled
	if errR != nil {
		// add log
		return mapping
	}
	var attrMap map[string]IdealAttributes
	errM := json.Unmarshal(attrMapping, &attrMap) //the mapping that defines how to parse an oppurtunity type is being unmarshalled here
	if errM != nil {
		//add log
		return mapping
	}
	for _, recco := range responses { //iterating through the recommendations one by one
		awsID := recco.ResourceId
		oppurType := recco.OpportunityType
		attributeTypeToValue := map[string]string{}
		attributes, ok := attrMap[oppurType]
		if ok {
			//known oppurtunity type has been encountered
			atrValueByPeriod := strings.Split(attributes.AttributeValue, ".")
			if atrValueByPeriod[0] == "parameters" {
				//the ideal value needs to be picked up from cloudfix reccomendations
				valueFromReccos, ok := recco.Parameters[atrValueByPeriod[1]]
				if !ok {
					//log that attribute is not present
					//if the code reaches here, then this means that the strategy for parsing has not been made correctly.
					// So we are resorting to showing the reccomendation against the resource name with the description for the oppurtunity
					attributeTypeToValue["NoAttributeMarker"] = recco.OpportunityDescription
				} else {
					idealAtrValue := valueFromReccos.(string) //extracting the ideal value as a string from cloudfix reccomendations
					attributeTypeToValue[attributes.AttributeType] = idealAtrValue
				}
			} else {
				//the ideal value is static and can be directly added
				attributeTypeToValue[attributes.AttributeType] = attributes.AttributeValue
			}
		} else {
			//unknown oppurtunity type has been encountered
			//So we are resorting to showing the reccomendation against the resource name with the description for the oppurtunity
			attributeTypeToValue["NoAttributeMarker"] = recco.OpportunityDescription
		}
		mapping[awsID] = attributeTypeToValue
	}
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
	attrMapping := []byte(`{
		"Gp2Gp3": {
			"Attribute Type": "type",
			"Attribute Value": "gp3"
		},
		"Ec2IntelToAmd": {
			"Attribute Type": "instance_type",
			"Attribute Value": "parameters.Migrating to instance type"
		},
		"StandardToSIT": {
			"Attribute Type": "NoAttributeMarker",
			"Attribute Value": "Enable Intelligent Tiering for this S3 Block by writing a aws_s3_bucket_intelligent_tiering_configuration resource block"
		},
		"EfsInfrequentAccess": {
			"Attribute Type": "NoAttributeMarker",
			"Attribute Value": "Enable Intelligent Tiering for EFS File by declaring a sub-block called lifecycle_policy within this resource block"
		}
	}`)
	reccosMapping := orches.parseReccos(fileR, attrMapping)
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
	tagToIDMap, errG := orches.getTagToID()
	if errG != nil {
		panic(errG)
	}
	errT := persist.store_tagMap(tagToIDMap, tagFileName)
	if errT != nil {
		panic(errT)
	}
	os.Setenv("TagsMapFile", tagFileName)
	modulesJson, _ := exec.Command("tflint", "--only=module_source", "-f=json").Output()
	modulePaths, errM := orches.extractModulePaths(modulesJson)
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
