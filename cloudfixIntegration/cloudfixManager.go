package cloudfixIntegration

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/trilogy-group/cloudfix-linter-developer/logger"
)

// Structure for unmarshalling the oppurtunityType to Attributes mapping (the mapping is present in "mappingAttributes.json")
type IdealAttributes struct {
	AttributeType  string `json:"Attribute Type"`
	AttributeValue string `json:"Attribute Value"`
	EnableQuickFix bool   `json:"EnableQuickFix"`
}

// structure for unmarshalling the reccomendation json response from cloudfix
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

type Recommendation struct {
	Recommendation map[string][]IdealAttributes
}

// type RecommendationDetails struct {
// 	AttributeValue string
// }

type ErrorCodes int

const (
	GENERIC_ERROR     ErrorCodes = iota //all other cases
	CRED_ERROR                          //Could not find Creds
	STORAGE_ERROR                       //Could not store the token
	UNAUTHCREDS_ERROR                   //Creds found, but server said Incorrect Creds
)

type CloudfixManager struct {
	//no data fields required
}

type customError struct {
	StatusCode ErrorCodes
	Message    string
}

func (e *customError) Error() string {
	return e.Message
}

//Member functions follow:

func (c *CloudfixManager) getReccosFromCloudfix(token string) ([]byte, *customError) {
	dlog := logger.DevLogger()
	var reccos []byte
	req, err := http.NewRequest("GET", RECOMMENDATIONS_ENDPOINT, nil)
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
		dlog.WithField("statusCode", statusCode).Error("Failed to fetch reccomendations")
		body, errI := ioutil.ReadAll(response.Body)
		if errI == nil {
			dlog.WithField("statusCode", statusCode).Error(body)
		}
		return reccos, &customError{GENERIC_ERROR, "Failed to fetch reccomendations"}
	}
	reccos, errI := ioutil.ReadAll(response.Body)
	if errI != nil {
		return []byte{}, &customError{GENERIC_ERROR, "Internal Error"}
	}
	return reccos, nil
}

func (c *CloudfixManager) createMap(reccos []byte, attrMapping []byte) map[string]Recommendation {
	mapping := map[string]Recommendation{} //this is the map that has to be returned in the end
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
		attributeTypeToValue := map[string][]IdealAttributes{}
		attributes, ok := attrMap[oppurType]
		var recommendationDetails IdealAttributes
		if ok {
			recommendationDetails = attributes
			//known oppurtunity type has been encountered
			atrValueByPeriod := strings.Split(attributes.AttributeValue, ".")
			if atrValueByPeriod[0] == "parameters" {
				//the ideal value needs to be picked up from cloudfix reccomendations
				valueFromReccos, ok := recco.Parameters[atrValueByPeriod[1]]
				if !ok {
					//log that attribute is not present
					//if the code reaches here, then this means that the strategy for parsing has not been made correctly.
					// So we are resorting to showing the reccomendation against the resource name with the description for the oppurtunity
					recommendationDetails.AttributeValue = recco.OpportunityDescription
					recommendationDetails.AttributeType = "NoAttributeMarker"
				} else {
					idealAtrValue := valueFromReccos.(string) //extracting the ideal value as a string from cloudfix reccomendations
					recommendationDetails.AttributeValue = idealAtrValue
				}
			}
		} else {
			//unknown oppurtunity type has been encountered
			//So we are resorting to showing the reccomendation against the resource name with the description for the oppurtunity
			recommendationDetails = IdealAttributes{AttributeType: "NoAttributeMarker", AttributeValue: recco.OpportunityDescription}
		}
		attributeTypeToValue[recommendationDetails.AttributeType] = append(attributeTypeToValue[recommendationDetails.AttributeType], recommendationDetails)
		_, exist := mapping[awsID]
		if exist == true {
			// awsID has multiple recommendations associated with it
			// merge all the recommendations
			for key, value := range attributeTypeToValue {
				_, exits := mapping[awsID].Recommendation[key]
				if exits {
					for _, val := range value {
						mapping[awsID].Recommendation[key] = append(mapping[awsID].Recommendation[key], val)
					}
				} else {
					mapping[awsID].Recommendation[key] = value
				}
			}
		} else {
			mapping[awsID] = Recommendation{Recommendation: attributeTypeToValue}
		}
	}
	return mapping
}

func (c *CloudfixManager) GetReccos() (map[string]Recommendation, *customError) {
	//function to process the reccomendations from cloudfix and turn that into a map
	//the structure of the map is resourceID -> Attribute type that needs to be targetted -> Ideal Attribute Value
	// If there is no attribute that has to be targetted, attribute type would be filled with "NoAttributeMarker" and
	//Attribute Value would be filled with any message that in the end has to be displayed to the user
	dlog := logger.DevLogger()
	var cloudAuth CloudfixAuth
	mapping := make(map[string]Recommendation)
	var reccos []byte
	val, present := os.LookupEnv("CLOUDFIX_FILE")
	var modeBoolval bool
	if present {
		modeBoolval, _ = strconv.ParseBool(val)
	}
	if present && modeBoolval {
		var errR error
		dlog.Info("CLOUDFIX_FILE mode on. Reading from reccos.json")
		currPWDStrip := ""
		currPWDStr := ""
		currPWDStrip1 := ""
		if runtime.GOOS == "windows" {
			currPWD, _ := exec.Command("powershell", "-NoProfile", "(pwd).path").Output()
			currPWDStr = string(currPWD[:])
			currPWDStrip = strings.Trim(currPWDStr, "\n")
			currPWDStrip = strings.TrimSuffix(currPWDStrip, "\r")
			currPWDStrip = strings.TrimSuffix(currPWDStrip, "cloudfix-linter")
			currPWDStrip1 = currPWDStrip + "\\reccos.json"
		} else {
			currPWD, _ := exec.Command("pwd").Output()
			currPWDStr = string(currPWD[:])
			currPWDStrip = strings.Trim(currPWDStr, "\n")
			currPWDStrip = strings.TrimSuffix(currPWDStrip, "cloudfix-linter")
			currPWDStrip1 = currPWDStrip + "/reccos.json"
		}

		reccos, errR = ioutil.ReadFile(currPWDStrip1)
		if errR != nil {
			//Add Error Log
			return mapping, &customError{GENERIC_ERROR, "Could not read reccos from file: " + currPWDStrip1}
		}
	} else {
		dlog.Info("CLOUDFIX_FILE mode off. Calling CLoudFix")
		token, errA := cloudAuth.getToken()
		if errA != nil && errA.StatusCode != STORAGE_ERROR {
			dlog.Error("Failed to retrieve and store the CloudFix token. ", errA)
			return mapping, errA
		}
		var errT *customError
		reccos, errT = c.getReccosFromCloudfix(token)

		if errT != nil {
			return mapping, errT
		}
	}
	attrMapping := []byte(`{
						"Gp2Gp3": {
							"Attribute Type": "type",
							"Attribute Value": "gp3",
							"EnableQuickFix" : true
						},
						"Ec2IntelToAmd": {
							"Attribute Type": "instance_type",
							"Attribute Value": "parameters.Migrating to instance type",
							"EnableQuickFix" : false
						},
						"StandardToSIT": {
							"Attribute Type": "NoAttributeMarker",
							"Attribute Value": "Enable Intelligent Tiering for this S3 Block by writing a aws_s3_bucket_intelligent_tiering_configuration resource block",
							"EnableQuickFix" : false
						},
						"EfsInfrequentAccess": {
							"Attribute Type": "NoAttributeMarker",
							"Attribute Value": "Enable Intelligent Tiering for EFS File by declaring a sub-block called lifecycle_policy within this resource block",
							"EnableQuickFix" : false
						},
						"IoToGp3": {
							"Attribute Type": "type",
							"Attribute Value": "gp3",
							"EnableQuickFix" : false
						},
						"DuplicateCloudTrail": {
							"Attribute Type": "enabled",
							"Attribute Value": "false",
							"EnableQuickFix" : false
						},
						"UnusedEBSVolumes": {
							"Attribute Type": "NoAttributeMarker",
							"Attribute Value": "Unattached EBS Volumes, Remove this to save the cost",
							"EnableQuickFix" : false
						},
						"VpcIdleEndpoint": {
							"Attribute Type": "NoAttributeMarker",
							"Attribute Value": "Idle VPC Endpoint, Remove this to save the cost",
							"EnableQuickFix" : false
						},
						"EfsIntelligentTiering": {
							"Attribute Type": "NoAttributeMarker",
							"Attribute Value": "Enable Intelligent Tiering for EFS File by declaring a sub-block called lifecycle_policy within this resource block",
							"EnableQuickFix" : false
						},
						"NeptuneCleanupIdleClusters": {
							"Attribute Type": "NoAttributeMarker",
							"Attribute Value": "Idle Neptune Cluster, Remove this to save the cost",
							"EnableQuickFix" : false
						},
						"InstallSSMAgentWindows": {
							"Attribute Type": "NoAttributeMarker",
							"Attribute Value": "Install SSM agent for Windows",
							"EnableQuickFix" : false
						},
						"InstallSSMAgentLinuxMacSSH": {
							"Attribute Type": "NoAttributeMarker",
							"Attribute Value": "Install SSM agent for Mac and Linux via SSH",
							"EnableQuickFix" : false
						},
						"VpcIdleNatGateway": {
							"Attribute Type": "NoAttributeMarker",
							"Attribute Value": "Idle VPC NAT Gateway, Remove this to save the cost",
							"EnableQuickFix" : false
						},
						"FixVPCDNSForAgents": {
							"Attribute Type": "NoAttributeMarker",
							"Attribute Value": "FixVPCDNSForAgents",
							"EnableQuickFix" : false
						},
						"EsOptimizeStorage": {
							"Attribute Type": "NoAttributeMarker",
							"Attribute Value": "Shrink AWS OpenSearch volumes",
							"EnableQuickFix" : false
						},
						"S3DDBTrafficToGWEndpoint": {
							"Attribute Type": "GlobalAttributeMarker",
							"Attribute Value": "S3/DynamoDB Traffic to Gateway Endpoint",
							"EnableQuickFix" : false
						},
						"DynamoDbProvisioning": {
							"Attribute Type": "billing_mode",
							"Attribute Value": "PROVISIONED",
							"EnableQuickFix" : false
						},
						"ArchiveOldEbsVolumeSnapshots": {
							"Attribute Type": "GlobalAttributeMarker",
							"Attribute Value": "Archive old EBS volume snapshots",
							"EnableQuickFix" : false
						},
						"DynamoDbInfrequentAccess": {
							"Attribute Type": "billing_mode",
							"Attribute Value": "PAY_PER_REQUEST",
							"EnableQuickFix" : false
						},
						"FixInstanceProfileForAgents": {
							"Attribute Type": "NoAttributeMarker",
							"Attribute Value": "FixInstanceProfileForAgents",
							"EnableQuickFix" : false
						},
						"CloudFrontCompression": {
							"Attribute Type": "ordered_cache_behavior.compress",
							"Attribute Value": "true",
							"EnableQuickFix" : false
						},
						"ElbCleanUpIdle": {
							"Attribute Type": "NoAttributeMarker",
							"Attribute Value": "Idle Elb, cleanup to save cost.",
							"EnableQuickFix" : false
						},
						"EC2CleanupUnusedAMIs": {
							"Attribute Type": "NoAttributeMarker",
							"Attribute Value": "Cleanup unused AMIs",
							"EnableQuickFix" : false
						}
						}`)
	mapping = c.createMap(reccos, attrMapping)
	return mapping, nil
}
