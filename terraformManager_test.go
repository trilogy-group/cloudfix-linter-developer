package main

import (
	"fmt"
	"reflect"
	"testing"
)

type TestTerraManager struct {
	tflintOut []byte
	expected  map[string]string
}

var addTestsTerra = []TestTerraManager{
	{
		// No resources deployed
		[]byte(``),
		map[string]string{},
	},
	{
		//Child resources
		[]byte(`{
			"format_version": "1.0",
			"terraform_version": "1.2.5",
			"values": {
				"root_module": {
					"child_modules": [
						{
							"resources": [
								{
									"address": "module.server-1.aws_s3_bucket.b",
									"mode": "managed",
									"type": "aws_s3_bucket",
									"name": "b",
									"provider_name": "registry.terraform.io/hashicorp/aws",
									"schema_version": 0,
									"values": {
										"acceleration_status": "",
										"acl": null,
										"arn": "arn:aws:s3:::my-tf-test-bucket-sabhya",
										"bucket": "my-tf-test-bucket-sabhya",
										"bucket_domain_name": "my-tf-test-bucket-sabhya.s3.amazonaws.com",
										"bucket_prefix": null,
										"bucket_regional_domain_name": "my-tf-test-bucket-sabhya.s3.amazonaws.com",
										"cors_rule": [],
										"force_destroy": false,
										"grant": [
											{
												"id": "8164d6483be5b82aab905d79b9ea8743c5ec8f09e9df4877e9cabd5d79e4a645",
												"permissions": [
													"FULL_CONTROL"
												],
												"type": "CanonicalUser",
												"uri": ""
											}
										],
										"hosted_zone_id": "Z3AQBSTGFYJSTF",
										"id": "my-tf-test-bucket-sabhya",
										"lifecycle_rule": [],
										"logging": [],
										"object_lock_configuration": [],
										"object_lock_enabled": false,
										"policy": "",
										"region": "us-east-1",
										"replication_configuration": [],
										"request_payer": "BucketOwner",
										"server_side_encryption_configuration": [],
										"tags": {
											"yor_trace": "1259a3dd-a0e7-4fb7-b30b-06bfff7b0d8e"
										},
										"tags_all": {
											"yor_trace": "1259a3dd-a0e7-4fb7-b30b-06bfff7b0d8e"
										},
										"versioning": [
											{
												"enabled": false,
												"mfa_delete": false
											}
										],
										"website": [],
										"website_domain": null,
										"website_endpoint": null
									},
									"sensitive_values": {
										"cors_rule": [],
										"grant": [
											{
												"permissions": [
													false
												]
											}
										],
										"lifecycle_rule": [],
										"logging": [],
										"object_lock_configuration": [],
										"replication_configuration": [],
										"server_side_encryption_configuration": [],
										"tags": {},
										"tags_all": {},
										"versioning": [
											{}
										],
										"website": []
									}
								}
							],
							"address": "module.server-1"
						},
						{
							"resources": [
								{
									"address": "module.server-2.aws_efs_file_system.foo",
									"mode": "managed",
									"type": "aws_efs_file_system",
									"name": "foo",
									"provider_name": "registry.terraform.io/hashicorp/aws",
									"schema_version": 0,
									"values": {
										"arn": "arn:aws:elasticfilesystem:us-east-1:575616755290:file-system/fs-0d77e387d74077730",
										"availability_zone_id": "",
										"availability_zone_name": "",
										"creation_token": "my-product",
										"dns_name": "fs-0d77e387d74077730.efs.us-east-1.amazonaws.com",
										"encrypted": false,
										"id": "fs-0d77e387d74077730",
										"kms_key_id": "",
										"lifecycle_policy": [],
										"number_of_mount_targets": 0,
										"owner_id": "575616755290",
										"performance_mode": "generalPurpose",
										"provisioned_throughput_in_mibps": 0,
										"size_in_bytes": [
											{
												"value": 6144,
												"value_in_ia": 0,
												"value_in_standard": 6144
											}
										],
										"tags": {
											"yor_trace": "1bd42041-0b31-4821-a88f-e5517f88473e"
										},
										"tags_all": {
											"yor_trace": "1bd42041-0b31-4821-a88f-e5517f88473e"
										},
										"throughput_mode": "bursting"
									},
									"sensitive_values": {
										"lifecycle_policy": [],
										"size_in_bytes": [
											{}
										],
										"tags": {},
										"tags_all": {}
									}
								}
							],
							"address": "module.server-2"
						}
					]
				}
			}
		}`),
		map[string]string{
			"1bd42041-0b31-4821-a88f-e5517f88473e": "fs-0d77e387d74077730",
			"1259a3dd-a0e7-4fb7-b30b-06bfff7b0d8e": "my-tf-test-bucket-sabhya",
		},
	},
	{
		//Root resources
		[]byte(`{
			"format_version": "1.0",
			"terraform_version": "1.2.5",
			"values": {
				"root_module": {
					"resources": [
						{
							"address": "aws_ebs_volume.example",
							"mode": "managed",
							"type": "aws_ebs_volume",
							"name": "example",
							"provider_name": "registry.terraform.io/hashicorp/aws",
							"schema_version": 0,
							"values": {
								"arn": "arn:aws:ec2:us-east-1:575616755290:volume/vol-0f9cc016d72167664",
								"availability_zone": "us-east-1a",
								"encrypted": false,
								"final_snapshot": false,
								"id": "vol-0f9cc016d72167664",
								"iops": 100,
								"kms_key_id": "",
								"multi_attach_enabled": false,
								"outpost_arn": "",
								"size": 1,
								"snapshot_id": "",
								"tags": {
									"yor_trace": "36b527db-6108-4062-8810-167b3753086b"
								},
								"tags_all": {
									"yor_trace": "36b527db-6108-4062-8810-167b3753086b"
								},
								"throughput": 0,
								"timeouts": null,
								"type": "gp2"
							},
							"sensitive_values": {
								"tags": {},
								"tags_all": {}
							}
						},
						{
							"address": "aws_instance.showcase-1",
							"mode": "managed",
							"type": "aws_instance",
							"name": "showcase-1",
							"provider_name": "registry.terraform.io/hashicorp/aws",
							"schema_version": 1,
							"values": {
								"ami": "ami-09d56f8956ab235b3",
								"arn": "arn:aws:ec2:us-east-1:575616755290:instance/i-069e995a00cb9f8a9",
								"associate_public_ip_address": true,
								"availability_zone": "us-east-1b",
								"capacity_reservation_specification": [
									{
										"capacity_reservation_preference": "open",
										"capacity_reservation_target": []
									}
								],
								"cpu_core_count": 1,
								"cpu_threads_per_core": 1,
								"credit_specification": [
									{
										"cpu_credits": "standard"
									}
								],
								"disable_api_stop": false,
								"disable_api_termination": false,
								"ebs_block_device": [],
								"ebs_optimized": false,
								"enclave_options": [
									{
										"enabled": false
									}
								],
								"ephemeral_block_device": [],
								"get_password_data": false,
								"hibernation": false,
								"host_id": null,
								"iam_instance_profile": "",
								"id": "i-069e995a00cb9f8a9",
								"instance_initiated_shutdown_behavior": "stop",
								"instance_state": "running",
								"instance_type": "t2.micro",
								"ipv6_address_count": 0,
								"ipv6_addresses": [],
								"key_name": "",
								"launch_template": [],
								"maintenance_options": [
									{
										"auto_recovery": "default"
									}
								],
								"metadata_options": [
									{
										"http_endpoint": "enabled",
										"http_put_response_hop_limit": 1,
										"http_tokens": "optional",
										"instance_metadata_tags": "disabled"
									}
								],
								"monitoring": false,
								"network_interface": [],
								"outpost_arn": "",
								"password_data": "",
								"placement_group": "",
								"placement_partition_number": null,
								"primary_network_interface_id": "eni-0b8fa0aa513353ddf",
								"private_dns": "ip-172-31-87-203.ec2.internal",
								"private_dns_name_options": [
									{
										"enable_resource_name_dns_a_record": false,
										"enable_resource_name_dns_aaaa_record": false,
										"hostname_type": "ip-name"
									}
								],
								"private_ip": "172.31.87.203",
								"public_dns": "ec2-184-72-72-101.compute-1.amazonaws.com",
								"public_ip": "184.72.72.101",
								"root_block_device": [
									{
										"delete_on_termination": true,
										"device_name": "/dev/sda1",
										"encrypted": false,
										"iops": 100,
										"kms_key_id": "",
										"tags": {},
										"throughput": 0,
										"volume_id": "vol-04a145b09011683db",
										"volume_size": 8,
										"volume_type": "gp2"
									}
								],
								"secondary_private_ips": [],
								"security_groups": [
									"default"
								],
								"source_dest_check": true,
								"subnet_id": "subnet-04661104769c18a80",
								"tags": {
									"yor_trace": "99e674ba-7be0-4969-87a7-0478ad792dcd"
								},
								"tags_all": {
									"yor_trace": "99e674ba-7be0-4969-87a7-0478ad792dcd"
								},
								"tenancy": "default",
								"timeouts": null,
								"user_data": null,
								"user_data_base64": null,
								"user_data_replace_on_change": false,
								"volume_tags": null,
								"vpc_security_group_ids": [
									"sg-006cdb061b9f623c9"
								]
							},
							"sensitive_values": {
								"capacity_reservation_specification": [
									{
										"capacity_reservation_target": []
									}
								],
								"credit_specification": [
									{}
								],
								"ebs_block_device": [],
								"enclave_options": [
									{}
								],
								"ephemeral_block_device": [],
								"ipv6_addresses": [],
								"launch_template": [],
								"maintenance_options": [
									{}
								],
								"metadata_options": [
									{}
								],
								"network_interface": [],
								"private_dns_name_options": [
									{}
								],
								"root_block_device": [
									{
										"tags": {}
									}
								],
								"secondary_private_ips": [],
								"security_groups": [
									false
								],
								"tags": {},
								"tags_all": {},
								"vpc_security_group_ids": [
									false
								]
							}
						}
					]
				}
			}
		}`),
		map[string]string{
			"99e674ba-7be0-4969-87a7-0478ad792dcd": "i-069e995a00cb9f8a9",
			"36b527db-6108-4062-8810-167b3753086b": "vol-0f9cc016d72167664",
		},
	},
	{
		//Tags not present
		[]byte(`{
			"format_version": "1.0",
			"terraform_version": "1.2.5",
			"values": {
				"root_module": {
					"resources": [
						{
							"address": "aws_ebs_volume.example",
							"mode": "managed",
							"type": "aws_ebs_volume",
							"name": "example",
							"provider_name": "registry.terraform.io/hashicorp/aws",
							"schema_version": 0,
							"values": {
								"arn": "arn:aws:ec2:us-east-1:575616755290:volume/vol-0f9cc016d72167664",
								"availability_zone": "us-east-1a",
								"encrypted": false,
								"final_snapshot": false,
								"id": "vol-0f9cc016d72167664",
								"iops": 100,
								"kms_key_id": "",
								"multi_attach_enabled": false,
								"outpost_arn": "",
								"size": 1,
								"snapshot_id": "",
								"tags": {
									"yor_trace": "36b527db-6108-4062-8810-167b3753086b"
								},
								"tags_all": {
									"yor_trace": "36b527db-6108-4062-8810-167b3753086b"
								},
								"throughput": 0,
								"timeouts": null,
								"type": "gp2"
							},
							"sensitive_values": {
								"tags": {},
								"tags_all": {}
							}
						},
						{
							"address": "aws_instance.showcase-1",
							"mode": "managed",
							"type": "aws_instance",
							"name": "showcase-1",
							"provider_name": "registry.terraform.io/hashicorp/aws",
							"schema_version": 1,
							"values": {
								"ami": "ami-09d56f8956ab235b3",
								"arn": "arn:aws:ec2:us-east-1:575616755290:instance/i-069e995a00cb9f8a9",
								"associate_public_ip_address": true,
								"availability_zone": "us-east-1b",
								"capacity_reservation_specification": [
									{
										"capacity_reservation_preference": "open",
										"capacity_reservation_target": []
									}
								],
								"cpu_core_count": 1,
								"cpu_threads_per_core": 1,
								"credit_specification": [
									{
										"cpu_credits": "standard"
									}
								],
								"disable_api_stop": false,
								"disable_api_termination": false,
								"ebs_block_device": [],
								"ebs_optimized": false,
								"enclave_options": [
									{
										"enabled": false
									}
								],
								"ephemeral_block_device": [],
								"get_password_data": false,
								"hibernation": false,
								"host_id": null,
								"iam_instance_profile": "",
								"id": "i-069e995a00cb9f8a9",
								"instance_initiated_shutdown_behavior": "stop",
								"instance_state": "running",
								"instance_type": "t2.micro",
								"ipv6_address_count": 0,
								"ipv6_addresses": [],
								"key_name": "",
								"launch_template": [],
								"maintenance_options": [
									{
										"auto_recovery": "default"
									}
								],
								"metadata_options": [
									{
										"http_endpoint": "enabled",
										"http_put_response_hop_limit": 1,
										"http_tokens": "optional",
										"instance_metadata_tags": "disabled"
									}
								],
								"monitoring": false,
								"network_interface": [],
								"outpost_arn": "",
								"password_data": "",
								"placement_group": "",
								"placement_partition_number": null,
								"primary_network_interface_id": "eni-0b8fa0aa513353ddf",
								"private_dns": "ip-172-31-87-203.ec2.internal",
								"private_dns_name_options": [
									{
										"enable_resource_name_dns_a_record": false,
										"enable_resource_name_dns_aaaa_record": false,
										"hostname_type": "ip-name"
									}
								],
								"private_ip": "172.31.87.203",
								"public_dns": "ec2-184-72-72-101.compute-1.amazonaws.com",
								"public_ip": "184.72.72.101",
								"root_block_device": [
									{
										"delete_on_termination": true,
										"device_name": "/dev/sda1",
										"encrypted": false,
										"iops": 100,
										"kms_key_id": "",
										"tags": {},
										"throughput": 0,
										"volume_id": "vol-04a145b09011683db",
										"volume_size": 8,
										"volume_type": "gp2"
									}
								],
								"secondary_private_ips": [],
								"security_groups": [
									"default"
								],
								"source_dest_check": true,
								"subnet_id": "subnet-04661104769c18a80",
								"tags": {},
								"tags_all": {},
								"tenancy": "default",
								"timeouts": null,
								"user_data": null,
								"user_data_base64": null,
								"user_data_replace_on_change": false,
								"volume_tags": null,
								"vpc_security_group_ids": [
									"sg-006cdb061b9f623c9"
								]
							},
							"sensitive_values": {
								"capacity_reservation_specification": [
									{
										"capacity_reservation_target": []
									}
								],
								"credit_specification": [
									{}
								],
								"ebs_block_device": [],
								"enclave_options": [
									{}
								],
								"ephemeral_block_device": [],
								"ipv6_addresses": [],
								"launch_template": [],
								"maintenance_options": [
									{}
								],
								"metadata_options": [
									{}
								],
								"network_interface": [],
								"private_dns_name_options": [
									{}
								],
								"root_block_device": [
									{
										"tags": {}
									}
								],
								"secondary_private_ips": [],
								"security_groups": [
									false
								],
								"tags": {},
								"tags_all": {},
								"vpc_security_group_ids": [
									false
								]
							}
						}
					]
				}
			}
		}`),
		map[string]string{
			"36b527db-6108-4062-8810-167b3753086b": "vol-0f9cc016d72167664",
		},
	},
}

func TestTerraMan(t *testing.T) {

	var terraMan TerraformManager
	for _, test := range addTestsTerra {
		expected := test.expected
		got, _ := terraMan.getTagToID(test.tflintOut)
		eq := reflect.DeepEqual(got, expected)
		if !eq {
			fmt.Print(expected)
			fmt.Print(got)
			t.Errorf("Test failed!")
		}
	}
}
