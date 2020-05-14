package main

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/organizations"
)

func main() {
	sess := session.Must(session.NewSession())
	svc := organizations.New(sess)

	var err error
	result := &organizations.ListAccountsOutput{}
	input := &organizations.ListAccountsInput{
		NextToken: new(string), // New this up for the first loop iteration to not be nil
	}

	for input.NextToken != nil {
		input.NextToken = result.NextToken
		// List all the accounts in the organization
		result, err = svc.ListAccounts(input)
		if err != nil {
			log.Fatalf("Failed to list accounts: %v", err)
		}

		// Sleep here to avoid rate limits when calling with NextToken
		time.Sleep(100 * time.Millisecond)

		// Get the tags for all accounts in the org
		for _, acc := range result.Accounts {
			nr, err := svc.ListTagsForResource(
				&organizations.ListTagsForResourceInput{ResourceId: acc.Id},
			)

			if err != nil {
				fmt.Printf("Failed to list tags for resource %s %s\n", *acc.Id, err)
				break
			}

			// Filter the tags on the "status" key and print their value
			for _, tags := range nr.Tags {
				if *tags.Key == "status" {
					fmt.Printf("%s: %s \n", *acc.Id, *tags.Value)
				}
			}
		}
	}
}
