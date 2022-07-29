# cloudfix-linter
A linting tool for HashiCorp's Terraform to flag recommendations from Cloudfix

Link to the User guide to use Cloudfix linter: https://docs.google.com/document/d/1cLR7cGQ8RrIgAb7R9qwgabU_VjAXO3Ria_9RGZ0USWY

## Guide on how to add support for new Cloudfix Oppurtunity Types:

A pure json mapping has been made so that support for new insights can be added easily.
Sample mapping json:

```
{
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
	}
```

For each new oppurtunity type, create a new block in the json by its name. If the opportunity type targets an attribute in specific, put in the name of the attribute for the Attribute Type. If it does not target any attribute, put in "NoAttributeMarker" instead. For the Attribute Value, if that needs to picked up from the parameters field of the cloudfix recommendation, set that as parameters.{Name of field within parameters block} (for reference take a look at the block for Ec2IntelToAmd). In case the value for the attribute is static and need not be picked up from the parameters field, it can be hardcoded directly in the json (for reference take a look at the block for Gp2Gp3). If the oppurtunity type does not target any attribute in specific, for the attribute value, put in the message that you want displayed to the user (for reference see the block for EfsInfrequentAccess)

This mapping is currently part of the code itself, but can be easily hosted online. 
