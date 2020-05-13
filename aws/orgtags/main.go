package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/organizations"
)

var (
	nextToken *string
)

func main() {
	input := &organizations.ListAccountsInput{}
	sess := session.Must(session.NewSession())
	svc := organizations.New(sess)

	for {
		if nextToken != nil {
			input.NextToken = nextToken
		} else {
			nextToken = nil
		}

		// List all the accounts in the organization
		result, err := svc.ListAccounts(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case organizations.ErrCodeAccessDeniedException:
					fmt.Println(organizations.ErrCodeAccessDeniedException, aerr.Error())
				case organizations.ErrCodeAWSOrganizationsNotInUseException:
					fmt.Println(organizations.ErrCodeAWSOrganizationsNotInUseException, aerr.Error())
				case organizations.ErrCodeInvalidInputException:
					fmt.Println(organizations.ErrCodeInvalidInputException, aerr.Error())
				case organizations.ErrCodeServiceException:
					fmt.Println(organizations.ErrCodeServiceException, aerr.Error())
				case organizations.ErrCodeTooManyRequestsException:
					fmt.Println(organizations.ErrCodeTooManyRequestsException, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
		}
		// Sleep here to avoid rate limits when calling with NextToken
		time.Sleep(100 * time.Millisecond)

		// Get the tags for all accounts in the org
		for _, acc := range result.Accounts {
			nr, _ := svc.ListTagsForResource(
				&organizations.ListTagsForResourceInput{ResourceId: acc.Id},
			)
			// Filter the tags on the "status" key and print their value
			for _, tags := range nr.Tags {
				if *tags.Key == "status" {
					fmt.Printf("%s: %s \n", *acc.Id, *tags.Value)
				}
			}
		}

		// Set the NextToken from the result, and call it again
		nextToken = result.NextToken

		if nextToken == nil {
			break
		}

	}
}
