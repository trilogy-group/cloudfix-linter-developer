package main

import (
	"encoding/json"

	"github.com/trilogy-group/cloudfix-linter/logger"
)

//structure for unmarhsalling the Parameters field of the reccomendation output from CLoudFix
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

func parseReccos(reccos []byte) map[string]map[string]string {
	appLogger := logger.New()
	mapping := map[string]map[string]string{}
	var responses []ResponseReccos
	err := json.Unmarshal(reccos, &responses)
	if err != nil {
		appLogger.Error().Println("Failed to unmarshall reccomendations")
		panic(err)
	}
	appLogger.Info().Println("Reccomendations unamrshalled succesfully!")
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
			appLogger.Warning().Printf("Unknown Oppurtunity Type for resource ID: \"%s\"", awsID)
			temp["NoAttributeMarker"] = recco.OpportunityDescription
			mapping[awsID] = temp
		}
	}
	appLogger.Info().Println("Reccomendation mapping made!")
	return mapping
}

func getTagToID() map[string]string {

}
