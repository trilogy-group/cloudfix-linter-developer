# cloudfix-linter

A linting tool for HashiCorp's Terraform to flag recommendations from Cloudfix

Link to the User guide to use Cloudfix linter: https://docs.google.com/document/d/1cLR7cGQ8RrIgAb7R9qwgabU_VjAXO3Ria_9RGZ0USWY


## Usage guide
1) Run command 
```bash
wget https://github.com/trilogy-group/cloudfix-linter/releases/latest/download/install.sh \
 && chmod +x install.sh \
 && ./install.sh
 ```

2). Ensure that terraform can access your AWS account. You can user one of the following
    a) Devconnect with [saml2aws](https://github.com/Versent/saml2aws)
    b) Set the access key and the secret key inside of the provider "aws" block eg: in the main.tf file provider "aws" { region = "us-east-1" access_key = "my-access-key" secret_key = "my-secret-key" } 
    c) Set and export AWS_ACCESS_KEY_ID , AWS_SECRET_ACCESS_KEY , AWS_SESSION_TOKEN as enviroment variables. More information on how to give access can be found [here](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)

3) This version works with CloudFix v3 so make sure you have credentials to https://preview.app.cloudfix.com/

4). From your terraform code working directory do "cloudfix-linter init".
```bash
cd my-terraform-project
cloudfix-linter init
cloudfix-linter --help
```

5). Run "terraform apply" to deploy the resources from your terraform code working directory.
```bash
terraform apply
```

6). To get recommendations from cloudfix and see them through CLI run command "cloudfix-linter flagRecco" 

Note :- If you make any changes to your terraform code, You first have to deploy them using “terraform apply” and then run “cloudfix-linter” command again through working directory of your terraform code to see reccomendations being flagged according to recent changes. 

Note:- If you do not have terraform code template to test this tool. You can use [this](https://github.com/trilogy-group/cloudfixLinter-demo) demo


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
