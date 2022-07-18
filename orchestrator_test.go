package main

import (
	"fmt"
	"reflect"
	"testing"
)

type Test struct {
	cloudfixReccos []byte
	attrMapping    []byte
	expected       map[string]map[string]string
}

var addTests = []Test{
	{
		// Test 1 : All known oppurtunity types with both dynamic and static ideal attribute values
		[]byte(`[
			{
				"id": "0048f1ca-c6c9-4f18-a2a1-822061436f03",
				"region": "us-east-1",
				"primaryImpactedNodeId": "vol-07339113fe87d6a42",
				"otherImpactedNodeIds": [],
				"resourceId": "vol-07339113fe87d6a42",
				"resourceName": "PSPIV-W16-8",
				"difficulty": 1,
				"risk": 1,
				"applicationEnvironment": "staging",
				"annualSavings": 33.15,
				"annualCost": 419.51,
				"status": "Manual Approval",
				"parameters": {},
				"templateApproved": false,
				"customerId": 5,
				"accountId": "631108317415",
				"accountNickname": "dev-avolin-pivotal",
				"opportunityType": "Gp2Gp3",
				"opportunityDescription": "EBS GP2 to Gp3",
				"generatedDate": "2022-06-20T07:01:17.274Z",
				"lastUpdatedDate": "2022-06-20T07:01:17.274Z"
			},
			{
				"id": "008c2866-f34c-4f89-948c-95d437896736",
				"region": "us-east-1",
				"primaryImpactedNodeId": "i-06436a993d6e05590",
				"otherImpactedNodeIds": [],
				"resourceId": "i-06436a993d6e05590",
				"resourceName": "SaaS Portal Ext",
				"difficulty": 1,
				"risk": 1,
				"applicationEnvironment": "staging",
				"annualSavings": 57.31,
				"annualCost": 596.00,
				"status": "Manual Approval",
				"parameters": {
					"Migrating to instance type": "t3a.xmicro",
					"Has recent snapshot": true
				},
				"templateApproved": true,
				"customerId": 5,
				"accountId": "442000216972",
				"accountNickname": "prod-centralfunctions-infragraph",
				"opportunityType": "Ec2IntelToAmd",
				"opportunityDescription": "EC2 Intel to AMD",
				"generatedDate": "2022-05-06T12:10:11.568Z",
				"lastUpdatedDate": "2022-06-20T07:01:21.847Z"
			}
		]`),
		[]byte(`{
			"Gp2Gp3": {
				"Attribute Type": "type",
				"Attribute Value": "gp3"
			},
			"Ec2IntelToAmd": {
				"Attribute Type": "instance_type",
				"Attribute Value": "parameters.Migrating to instance type"
			}}`),
		map[string]map[string]string{
			"vol-07339113fe87d6a42": {
				"type": "gp3",
			},
			"i-06436a993d6e05590": {
				"instance_type": "t3a.xmicro",
			},
		},
	},
	{
		// Test 2 : Unknown oppurtunity types
		[]byte(`[
			{
				"id": "0048f1ca-c6c9-4f18-a2a1-822061436f03",
				"region": "us-east-1",
				"primaryImpactedNodeId": "vol-07339113fe87d6a42",
				"otherImpactedNodeIds": [],
				"resourceId": "vol-07339113fe87d6a42",
				"resourceName": "PSPIV-W16-8",
				"difficulty": 1,
				"risk": 1,
				"applicationEnvironment": "staging",
				"annualSavings": 33.15,
				"annualCost": 419.51,
				"status": "Manual Approval",
				"parameters": {},
				"templateApproved": false,
				"customerId": 5,
				"accountId": "631108317415",
				"accountNickname": "dev-avolin-pivotal",
				"opportunityType": "Gp2Gp4",
				"opportunityDescription": "EBS GP2 to Gp3",
				"generatedDate": "2022-06-20T07:01:17.274Z",
				"lastUpdatedDate": "2022-06-20T07:01:17.274Z"
			},
			{
				"id": "008c2866-f34c-4f89-948c-95d437896736",
				"region": "us-east-1",
				"primaryImpactedNodeId": "i-06436a993d6e05590",
				"otherImpactedNodeIds": [],
				"resourceId": "i-06436a993d6e05590",
				"resourceName": "SaaS Portal Ext",
				"difficulty": 1,
				"risk": 1,
				"applicationEnvironment": "staging",
				"annualSavings": 57.31,
				"annualCost": 596.00,
				"status": "Manual Approval",
				"parameters": {
					"Migrating to instance type": "t3a.xmicro",
					"Has recent snapshot": true
				},
				"templateApproved": true,
				"customerId": 5,
				"accountId": "442000216972",
				"accountNickname": "prod-centralfunctions-infragraph",
				"opportunityType": "NewOpprType",
				"opportunityDescription": "EC2 Intel to AMD",
				"generatedDate": "2022-05-06T12:10:11.568Z",
				"lastUpdatedDate": "2022-06-20T07:01:21.847Z"
			}
		]`),
		[]byte(`{
			"Gp2Gp3": {
				"Attribute Type": "type",
				"Attribute Value": "gp3"
			},
			"Ec2IntelToAmd": {
				"Attribute Type": "instance_type",
				"Attribute Value": "parameters.Migrating to instance type"
			}}`),
		map[string]map[string]string{
			"vol-07339113fe87d6a42": {
				"NoAttributeMarker": "EBS GP2 to Gp3",
			},
			"i-06436a993d6e05590": {
				"NoAttributeMarker": "EC2 Intel to AMD",
			},
		},
	},
	{
		// Test 3 : No reccomendations received
		[]byte(``),
		[]byte(`{
			"Gp2Gp3": {
				"Attribute Type": "type",
				"Attribute Value": "gp3"
			},
			"Ec2IntelToAmd": {
				"Attribute Type": "instance_type",
				"Attribute Value": "parameters.Migrating to instance type"
			}}`),
		map[string]map[string]string{},
	},
	{
		// Test 4 : The strategy json asks to index into a parameter key that is not present
		[]byte(`[
			{
				"id": "0048f1ca-c6c9-4f18-a2a1-822061436f03",
				"region": "us-east-1",
				"primaryImpactedNodeId": "vol-07339113fe87d6a42",
				"otherImpactedNodeIds": [],
				"resourceId": "vol-07339113fe87d6a42",
				"resourceName": "PSPIV-W16-8",
				"difficulty": 1,
				"risk": 1,
				"applicationEnvironment": "staging",
				"annualSavings": 33.15,
				"annualCost": 419.51,
				"status": "Manual Approval",
				"parameters": {},
				"templateApproved": false,
				"customerId": 5,
				"accountId": "631108317415",
				"accountNickname": "dev-avolin-pivotal",
				"opportunityType": "Gp2Gp3",
				"opportunityDescription": "EBS GP2 to Gp3",
				"generatedDate": "2022-06-20T07:01:17.274Z",
				"lastUpdatedDate": "2022-06-20T07:01:17.274Z"
			},
			{
				"id": "008c2866-f34c-4f89-948c-95d437896736",
				"region": "us-east-1",
				"primaryImpactedNodeId": "i-06436a993d6e05590",
				"otherImpactedNodeIds": [],
				"resourceId": "i-06436a993d6e05590",
				"resourceName": "SaaS Portal Ext",
				"difficulty": 1,
				"risk": 1,
				"applicationEnvironment": "staging",
				"annualSavings": 57.31,
				"annualCost": 596.00,
				"status": "Manual Approval",
				"parameters": {
					"Migrating to instance type": "t3a.xmicro",
					"Has recent snapshot": true
				},
				"templateApproved": true,
				"customerId": 5,
				"accountId": "442000216972",
				"accountNickname": "prod-centralfunctions-infragraph",
				"opportunityType": "Ec2IntelToAmd",
				"opportunityDescription": "EC2 Intel to AMD",
				"generatedDate": "2022-05-06T12:10:11.568Z",
				"lastUpdatedDate": "2022-06-20T07:01:21.847Z"
			}
		]`),
		[]byte(`{
			"Gp2Gp3": {
				"Attribute Type": "type",
				"Attribute Value": "gp3"
			},
			"Ec2IntelToAmd": {
				"Attribute Type": "instance_type",
				"Attribute Value": "parameters.Migrating to instance type DOES NOT EXIST"
			}}`),
		map[string]map[string]string{
			"vol-07339113fe87d6a42": {
				"type": "gp3",
			},
			"i-06436a993d6e05590": {
				"NoAttributeMarker": "EC2 Intel to AMD",
			},
		},
	},
	{
		// Test 5 : Error unmarshalling attrMap
		[]byte(`[
			{
				"id": "0048f1ca-c6c9-4f18-a2a1-822061436f03",
				"region": "us-east-1",
				"primaryImpactedNodeId": "vol-07339113fe87d6a42",
				"otherImpactedNodeIds": [],
				"resourceId": "vol-07339113fe87d6a42",
				"resourceName": "PSPIV-W16-8",
				"difficulty": 1,
				"risk": 1,
				"applicationEnvironment": "staging",
				"annualSavings": 33.15,
				"annualCost": 419.51,
				"status": "Manual Approval",
				"parameters": {},
				"templateApproved": false,
				"customerId": 5,
				"accountId": "631108317415",
				"accountNickname": "dev-avolin-pivotal",
				"opportunityType": "Gp2Gp3",
				"opportunityDescription": "EBS GP2 to Gp3",
				"generatedDate": "2022-06-20T07:01:17.274Z",
				"lastUpdatedDate": "2022-06-20T07:01:17.274Z"
			},
			{
				"id": "008c2866-f34c-4f89-948c-95d437896736",
				"region": "us-east-1",
				"primaryImpactedNodeId": "i-06436a993d6e05590",
				"otherImpactedNodeIds": [],
				"resourceId": "i-06436a993d6e05590",
				"resourceName": "SaaS Portal Ext",
				"difficulty": 1,
				"risk": 1,
				"applicationEnvironment": "staging",
				"annualSavings": 57.31,
				"annualCost": 596.00,
				"status": "Manual Approval",
				"parameters": {
					"Migrating to instance type": "t3a.xmicro",
					"Has recent snapshot": true
				},
				"templateApproved": true,
				"customerId": 5,
				"accountId": "442000216972",
				"accountNickname": "prod-centralfunctions-infragraph",
				"opportunityType": "Ec2IntelToAmd",
				"opportunityDescription": "EC2 Intel to AMD",
				"generatedDate": "2022-05-06T12:10:11.568Z",
				"lastUpdatedDate": "2022-06-20T07:01:21.847Z"
			}
		]`),
		[]byte(`{
			"Gp2Gp3": {
				"Attribute Type": "type",
				"Attribute Value": "gp3" WRONG SYNTAX
			},
			"Ec2IntelToAmd": {
				"Attribute Type": "instance_type",
				"Attribute Value": "parameters.Migrating to instance type DOES NOT EXIST"
			}}`),
		map[string]map[string]string{},
	},
	{
		// Test 6 : Error unmarshalling cloudfix reccos
		[]byte(`[
			{
				"id": "0048f1ca-c6c9-4f18-a2a1-822061436f03",
				"region": "us-east-1",
				"primaryImpactedNodeId": "vol-07339113fe87d6a42" WRONG SYNTAX
				"otherImpactedNodeIds": [],
				"resourceId": "vol-07339113fe87d6a42",
				"resourceName": "PSPIV-W16-8",
				"difficulty": 1,
				"risk": 1,
				"applicationEnvironment": "staging",
				"annualSavings": 33.15,
				"annualCost": 419.51,
				"status": "Manual Approval",
				"parameters": {},
				"templateApproved": false,
				"customerId": 5,
				"accountId": "631108317415",
				"accountNickname": "dev-avolin-pivotal",
				"opportunityType": "Gp2Gp3",
				"opportunityDescription": "EBS GP2 to Gp3",
				"generatedDate": "2022-06-20T07:01:17.274Z",
				"lastUpdatedDate": "2022-06-20T07:01:17.274Z"
			},
			{
				"id": "008c2866-f34c-4f89-948c-95d437896736",
				"region": "us-east-1",
				"primaryImpactedNodeId": "i-06436a993d6e05590",
				"otherImpactedNodeIds": [],
				"resourceId": "i-06436a993d6e05590",
				"resourceName": "SaaS Portal Ext",
				"difficulty": 1,
				"risk": 1,
				"applicationEnvironment": "staging",
				"annualSavings": 57.31,
				"annualCost": 596.00,
				"status": "Manual Approval",
				"parameters": {
					"Migrating to instance type": "t3a.xmicro",
					"Has recent snapshot": true
				},
				"templateApproved": true,
				"customerId": 5,
				"accountId": "442000216972",
				"accountNickname": "prod-centralfunctions-infragraph",
				"opportunityType": "Ec2IntelToAmd",
				"opportunityDescription": "EC2 Intel to AMD",
				"generatedDate": "2022-05-06T12:10:11.568Z",
				"lastUpdatedDate": "2022-06-20T07:01:21.847Z"
			}
		]`),
		[]byte(`{
			"Gp2Gp3": {
				"Attribute Type": "type",
				"Attribute Value": "gp3"
			},
			"Ec2IntelToAmd": {
				"Attribute Type": "instance_type",
				"Attribute Value": "parameters.Migrating to instance type DOES NOT EXIST"
			}}`),
		map[string]map[string]string{},
	},
}

func TestParseReccos(t *testing.T) {

	var orches Orchestrator
	for _, test := range addTests {
		expected := test.expected
		got := orches.parseReccos(test.cloudfixReccos, test.attrMapping)
		eq := reflect.DeepEqual(got, expected)
		if !eq {
			fmt.Print(expected)
			fmt.Print(got)
			t.Errorf("Test failed!")
		}
	}
}
