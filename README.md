# Cloudfix-linter

[![Maintained by Bridgecrew.io](https://img.shields.io/badge/maintained%20by-bridgecrew.io-blueviolet)](https://bridgecrew.io/?utm_source=github&utm_medium=organic_oss&utm_campaign=yor)
![golangci-lint](https://github.com/bridgecrewio/yor/workflows/tests/badge.svg)
[![security](https://github.com/bridgecrewio/yor/actions/workflows/security.yml/badge.svg)](https://github.com/bridgecrewio/yor/actions/workflows/security.yml)
<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-81%25-brightgreen.svg?longCache=true&style=flat)</a>
[![slack-community](https://img.shields.io/badge/Slack-4A154B?style=plastic&logo=slack&logoColor=white)](https://slack.bridgecrew.io/)
[![Go Report Card](https://goreportcard.com/badge/github.com/bridgecrewio/yor)](https://goreportcard.com/report/github.com/bridgecrewio/yor)
[![Go Reference](https://pkg.go.dev/badge/github.com/bridgecrewio/yor.svg)](https://pkg.go.dev/github.com/bridgecrewio/yor)
[![Docker pulls](https://img.shields.io/docker/pulls/bridgecrew/yor.svg)](https://hub.docker.com/r/bridgecrew/yor)
[![Chocolatey downloads](https://img.shields.io/chocolatey/dt/yor?label=chocolatey_downloads)](https://community.chocolatey.org/packages/yor)
[![GitHub All Releases](https://img.shields.io/github/downloads/bridgecrewio/yor/total)](https://github.com/bridgecrewio/yor/releases)

Yor is an open-source tool that helps add informative and consistent tags across infrastructure as code (IaC) frameworks. Today, Yor can automatically add tags to Terraform, CloudFormation, and Serverless Frameworks.

Yor is built to run as a [GitHub Action](https://github.com/bridgecrewio/yor-action) automatically adding consistent tagging logics to your IaC. Yor can also run as a pre-commit hook and a standalone CLI.

It is a command line tool that flags optimisation oppurtunities detected by Cloudfix for the resources that have been deployed using terraform. It'll either flag the specific attribute within the resource that needs to be changed (along with what it needs to be changed to), or in the case that such an attribute does not exist, describe the oppurtunity against the name of the resource about which the oppurtunity is present. It will identify the resources deployed by remote modules and provide recommendations for them.

## Demo
[![](docs/yor_tag_and_trace_recording.gif)](https://raw.githubusercontent.com/bridgecrewio/yor/main/docs/yor_tag_and_trace_recording.gif)

1. An active cloudfix account at https://app.cloudfix.com/
2. Resources deployed on AWS using terraform for which you would like to see reccomendations.

## Usage guide
#### 1. Run command 
- Windows
```
Invoke-WebRequest -URI https://github.com/trilogy-group/cloudfix-linter/releases/latest/download/install.ps1 -OutFile install.ps1; ./install.ps1 (pwd).path
```

- Linux and Devspaces
```bash
read -sp "Enter sudo password: " pass  &&  wget -O - https://github.com/trilogy-group/cloudfix-linter/releases/latest/download/install.sh | bash /dev/stdin $pass
 ```
 - Devspaces
```
wget -O - https://github.com/trilogy-group/cloudfix-linter/releases/latest/download/install.sh | bash
```

#### 2. Ensure that terraform can access your AWS account. You can user one of the following:

- Devconnect with [saml2aws](https://github.com/Versent/saml2aws)
- Set the access key and the secret key inside of the provider "aws" block eg: in the main.tf file provider "aws" { region = "us-east-1" access_key = "my-access-key" secret_key = "my-secret-key" } 
- Set and export AWS_ACCESS_KEY_ID , AWS_SECRET_ACCESS_KEY , AWS_SESSION_TOKEN as enviroment variables. More information on how to give access can be found [here](https://registry.terraform.io/providers/hashicorp/aws/latest/docs).

#### 3. Cloudfix Credentials
This version works with CloudFix v3 so make sure you have credentials to https://app.cloudfix.com/

#### 4. From your terraform code working directory do "cloudfix-linter init".
```bash
cd my-terraform-project
cloudfix-linter init
cloudfix-linter --help
```
__OR__

#### 5. Run "terraform apply" to deploy the resources from your terraform code working directory.
```bash
terraform apply
```

#### 6. To get recommendations from cloudfix and see them through CLI run command 
```
cloudfix-linter recco
```

Docker
```sh
docker pull bridgecrew/yor

docker run --tty --volume /local/path/to/tf:/tf bridgecrew/yor tag --directory /tf
```


GitHub Action
```yaml
name: IaC trace

on:
  # Triggers the workflow on push or pull request events but only for the main branch
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

  # Allows you to run this workflow manually from the Actions tab
  workflow_dispatch:

jobs:
  yor:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
        name: Checkout repo
        with:
          fetch-depth: 0
          ref: ${{ github.head_ref }}
      - name: Run yor action and commit
        uses: bridgecrewio/yor-action@main
```



Pre-commit
```yaml
  - repo: https://github.com/bridgecrewio/yor
    rev: 0.1.143
    hooks:
      - id: yor
        name: yor
        entry: yor tag -d
        args: ["."]
        language: golang
        types: [terraform]
        pass_filenames: false
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
Detailed mapping can be viewed [here](https://github.com/trilogy-group/cloudfix-linter/blob/6ed0a514dc3dd8c865f81e2dcddda456d3012fca/cloudfixIntegration/cloudfixManager.go#L221).

`-o` : Modify output formats.

```sh
yor tag -d . -o cli
# default cli output

yor tag -d . -o json
# json output

yor tag -d . --output cli --output-json-file result.json
# print cli output and additional output to file on json file -- enables programmatic analysis alongside printing human readable result
```

`--skip-tags`: Specify only named tags (allow list) or run all tags except those listed (deny list).

```sh
yor tag -d . --skip-tags yor_trace
## Run all but yor_trace

yor tag -d . --skip-tags yor_trace,git_modifiers
## Run all but yor_trace and git_modifiers

yor tag -d . --skip-tags git*
## Run all tags except tags with specified patterns
```

`skip-dirs` : Skip directory paths you can define paths that will not be tagged.

```sh
yor tag -d path/to/files
## Run on the directory path/to/files

yor tag -d path/to/files --skip-dirs path/to/files/skip,path/to/files/another/skip2
## Run yor on the directory path/to/files, skipping path/to/files/skip/ and path/to/files/another/skip2/
```
Detailed mapping can be viewed [here](https://github.com/trilogy-group/cloudfix-linter/blob/6ed0a514dc3dd8c865f81e2dcddda456d3012fca/cloudfixIntegration/cloudfixManager.go#L221).

`list-tag`

```sh
yor list-tag-groups
 # List tag classes that are built into yor.

yor list-tags
 # List all the tags built into yor
yor list-tags --tag-groups git

 # List all the tags built into yor under the tag group git
```

`dry-run`
```sh
yor tag -d . --dry-run
# Perform a dry run to get a preview in the CLI output of all of the tags that will be added using Yor without applying any changes to your IaC files.
```
### What is Yor trace?
yor_trace is a magical tag creating a unique identifier for an IaC resource code block.

Having a yor_trace in place can help with tracing code block to its cloud provisioned resources without access to sensitive data such as plan or state files.

See demo [here](https://yor.io/4.Use%20Cases/useCases.html)
## Contributing

The project uses a custom ruleset written for [TfLint](https://github.com/terraform-linters/tflint/blob/master/docs/developer-guide/architecture.md) to flag reccomendations from cloudfix. The github repo for the ruleset can be accessed [here](https://github.com/trilogy-group/tflint-ruleset-template)

We are working on extending Yor and adding more parsers (to support additional IaC frameworks) and more taggers (to tag using other contextual data).

To maintain our conventions, please run lint on your branch before opening a PR. To run lint:
```sh
golangci-lint run --fix --skip-dirs tests/yor_plugins
```

## Support

For more support contact us at https://slack.bridgecrew.io/.
