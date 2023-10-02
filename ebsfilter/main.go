package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func get_instances(ec2_svc *ec2.EC2, instanceName string) (*ec2.DescribeInstancesOutput, error){
	// Retrieve EC2 Instances
	var input *ec2.DescribeInstancesInput

	if instanceName == "" {
		input = &ec2.DescribeInstancesInput{}
	} else {
		input = &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				{
					Name: aws.String("tag:Name"),
					Values: []*string{
						aws.String(instanceName),
					},
				},
			},
		}
	}
	
	result, err := ec2_svc.DescribeInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
	}

	return result, err
}

func getNameTag( instance *ec2.Instance ) string {
	// Extracts Name Tag from Instance. Sets it to a default if a name tag is not found (e.g EKS Nodes)
	for _, tag := range instance.Tags {
		if *tag.Key == "Name" {
			return *tag.Value
		}
	}
	return "NameNotSet"
}

func getPublicIP(instance *ec2.Instance) string{
	// Public IP addresses are not mandatory
	// so we need to check if it exists before attemting to collect it
	
	if instance.PublicIpAddress != nil {
		return *instance.PublicIpAddress
	}

	return "Public IP not set"
}

func getVolumeSize(ec2_svc *ec2.EC2, instance *ec2.Instance ) (error, float64) {
	//Get the EBS volumes attached to 
	var volumeTotal int64
	var volumeTotalFloat float64
	var err error

	for _, volume := range instance.BlockDeviceMappings {
		input := &ec2.DescribeVolumesInput{
			VolumeIds: []*string{aws.String(*volume.Ebs.VolumeId)},
		}
		
		volume, err := ec2_svc.DescribeVolumes(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return  err, volumeTotalFloat
		}
		volumeTotal = volumeTotal + *volume.Volumes[0].Size
		volumeTotalFloat = float64(volumeTotal) * 1.074
	}
	
	return  err, volumeTotalFloat
}

func main() {

	var instanceName = flag.String("instancename", "", "EC2 Instance Name Tag")

	flag.Parse()
	instanceNameValue := *instanceName
	sess := session.Must(session.NewSession())
	svc := ec2.New(sess)

	var outputArray []InstanceOutput 
	instances, err := get_instances(svc, instanceNameValue)
	if err != nil {
		log.Fatalln(err)
	}

	for _, instances := range instances.Reservations {
			for _, instance := range instances.Instances{
				instanceOut := InstanceOutput{
					InstanceId: *instance.InstanceId,
					Name: getNameTag(instance),
					PublicIP: getPublicIP(instance),
					PrivateIP: *instance.PrivateIpAddress,
					InstanceType: *instance.InstanceType,
					InstanceState: *instance.State.Name,
				}

				_, TotalVolumeSpace := getVolumeSize(svc,instance) 
				instanceOut.TotalVolumeSpace = TotalVolumeSpace
				outputArray = append(outputArray, instanceOut)
			}
	}

	for _, out := range outputArray {
		fmt.Printf("%s, %s, %s, %s, %s, %s, %.1f\n", out.InstanceId, out.Name, out.PublicIP, out.PrivateIP, out.InstanceType, out.InstanceState, out.TotalVolumeSpace)
	}

}
