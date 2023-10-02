package main

import (
	"flag"
	"fmt"
	"reflect"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

func get_iam_users(iam_svc *iam.IAM) *iam.ListUsersOutput {
	//retrieve users
	input := &iam.ListUsersInput{}

	result, err := iam_svc.ListUsers(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeServiceFailureException:
				fmt.Println(iam.ErrCodeServiceFailureException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
	}
	return result
}

func get_old_iam_user_access_keys(iam_svc *iam.IAM, iam_user string, monthsolderthan int) []aws_access_key {
	//retrieve keys from a specific user
	var aws_access_keys []aws_access_key

	input := &iam.ListAccessKeysInput{
		UserName: aws.String(iam_user),
	}

	result, err := iam_svc.ListAccessKeys(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				fmt.Println(iam.ErrCodeNoSuchEntityException, aerr.Error())
			case iam.ErrCodeServiceFailureException:
				fmt.Println(iam.ErrCodeServiceFailureException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
	}

	if reflect.ValueOf(result.AccessKeyMetadata).IsZero() {
		fmt.Printf("-\n")
	} else {
		for _, k := range result.AccessKeyMetadata {
			var awsAccessKey aws_access_key
			awsAccessKey.AccessKeyId = *k.AccessKeyId
			awsAccessKey.CreateDate = *k.CreateDate
			awsAccessKey.Status = *k.Status
			awsAccessKey.Username = *k.UserName

			too_old := older_than(awsAccessKey.CreateDate, monthsolderthan)
			if too_old {
				aws_access_keys = append(aws_access_keys, awsAccessKey)
			}
		}
	}
	return aws_access_keys
}

func older_than(key_age time.Time, monthsolderthan int) bool {
	//was the key created
	today := time.Now()
	time_ago := today.AddDate(0, -monthsolderthan, 0) // minus months
	return key_age.Before(time_ago)
}

func main() {
	//TODO improve credential detection: Profiles, MFA etc
	var months = flag.Int("monthsolderthan", 3, "Specify the age of the keys in months you wish to find. )")

	flag.Parse()

	monthsolderthan := *months

	sess := session.Must(session.NewSession())
	iam_svc := iam.New(sess)
	iam_users := get_iam_users(iam_svc)

	var all_aws_access_key []aws_access_key

	for _, user := range iam_users.Users {
		access_keys := get_old_iam_user_access_keys(iam_svc, *user.UserName, monthsolderthan)
		for _, k := range access_keys {
			all_aws_access_key = append(all_aws_access_key, k)
		}
	}
	//TODO maybe use templates
	for _, key := range all_aws_access_key {
		fmt.Printf("Account: %s\n", key.Username)
		fmt.Printf("Key ID: %s\n", key.AccessKeyId)
		fmt.Printf("Key Create Date: %s\n", key.CreateDate)
		fmt.Println()
	}

}
