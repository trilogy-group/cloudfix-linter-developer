package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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

type ErrorCodes int

const (
	GENERIC_ERROR     ErrorCodes = iota //all other cases
	CRED_ERROR                          //Could not find Creds
	STORAGE_ERROR                       //Could not store the token
	UNAUTHCREDS_ERROR                   //Creds found, but server returned 401
)

type CloudfixManager struct {
	//no data fields required
}

type customError struct {
	statusCode ErrorCodes
	message    string
}

func (e *customError) Error() string {
	return e.message
}

//Member functions follow:

func (c *CloudfixManager) getReccosFromCloudfix(token string) ([]byte, *customError) {
	var reccos []byte
	req, err := http.NewRequest("GET", "https://w9lnd111rl.execute-api.us-east-1.amazonaws.com/default/api/v1/ui/recommendations", nil)
	if err != nil {
		return reccos, &customError{GENERIC_ERROR, "Internal Error"}
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", token)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return reccos, &customError{GENERIC_ERROR, "Internal Error"}
	}
	defer response.Body.Close()
	statusCode := response.StatusCode
	if statusCode != http.StatusOK {
		return reccos, &customError{GENERIC_ERROR, "Failed to fetch reccomendations"}
	}
	reccos, errI := ioutil.ReadAll(response.Body)
	if errI != nil {
		return []byte{}, &customError{GENERIC_ERROR, "Internal Error"}
	}
	//fmt.Println(string(reccos))
	return reccos, nil
}

func (c *CloudfixManager) createMap(reccos []byte, attrMapping []byte) map[string]map[string]string {
	mapping := map[string]map[string]string{} //this is the map that has to be returned in the end
	var responses []ResponseReccos
	if len(reccos) == 0 {
		//log that no reccomendations have been received
		return mapping
	}
	errU := json.Unmarshal(reccos, &responses) //the reccomendations from cloudfix are being unmarshalled
	if errU != nil {
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

func (c *CloudfixManager) getReccos() (map[string]map[string]string, *customError) {
	//function to process the reccomendations from cloudfix and turn that into a map
	//the structure of the map is resourceID -> Attribute type that needs to be targetted -> Ideal Attribute Value
	// If there is no attribute that has to be targetted, attribute type would be filled with "NoAttributeMarker" and
	//Attribute Value would be filled with any message that in the end has to be displayed to the user

	var cloudAuth CloudfixAuth
	mapping := make(map[string]map[string]string)
	token, errA := cloudAuth.getToken()
	if errA != nil && errA.statusCode != STORAGE_ERROR {
		return mapping, errA
	}
	// currPWD, _ := exec.Command("pwd").Output()
	// currPWDStr := string(currPWD[:])
	// currPWDStrip := strings.Trim(currPWDStr, "\n")
	// currPWDStrip += "/reccos.json"
	// reccos, errR := ioutil.ReadFile(currPWDStrip)
	// if errR != nil {
	// 	//Add Error Log
	// 	panic(errR)
	// }
	reccos, errT := c.getReccosFromCloudfix(token)
	if errT != nil {
		fmt.Println(errT.message)
		return mapping, errT
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
	mapping = c.createMap(reccos, attrMapping)
	return mapping, nil
}
