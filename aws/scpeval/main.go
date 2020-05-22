package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

type Policy struct {
	Version   string      `json:"Version"`
	Statement []Statement `json:"Statement"`
}
type StringNotEquals struct {
	AwsRequestedRegion []string `json:"aws:RequestedRegion"`
}
type ArnNotLike struct {
	AwsPrincipalARN string `json:"aws:PrincipalARN"`
}
type Condition struct {
	StringNotEquals StringNotEquals `json:"StringNotEquals"`
	ArnNotLike      ArnNotLike      `json:"ArnNotLike"`
}
type Statement struct {
	NotAction []*string `json:"NotAction,omitempty"`
	Condition Condition `json:"Condition,omitempty"`
	Action    []*string `json:"Action,omitempty"`
}

func main() {
	sess := session.Must(session.NewSession())
	svc := iam.New(sess)

	policy := Policy{}

	// The first argument should be the IAM policy .json file
	policyBytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("File does not exist %s\n", err)
	}

	// The second argument should be the IAM Principal you want to test
	principal := os.Args[2]

	err = json.Unmarshal(policyBytes, &policy)
	if err != nil {
		log.Fatalf("Boom %s\n", err)
	}

	// Evaluate Actions
	for _, statement := range policy.Statement {
		if len(statement.Action) == 0 {
			continue
		}

		fullPolicy := string(policyBytes)
		sim := iam.SimulateCustomPolicyInput{
			ActionNames:     statement.Action,
			PolicyInputList: []*string{&fullPolicy},
			CallerArn:       &principal,
		}

		simRes, err := svc.SimulateCustomPolicy(&sim)
		if err != nil {
			fmt.Println("Error gettting response", err)
		}

		for _, res := range simRes.EvaluationResults {
			// fmt.Println(res)
			if !strings.Contains(*res.EvalDecision, "Deny") {
				fmt.Printf("Action %s was allowed!\n", *res.EvalActionName)
			}
		}
	}

	// Evaluate NotActions
	for _, statement := range policy.Statement {
		if len(statement.NotAction) == 0 {
			continue
		}

		fullPolicy := string(policyBytes)
		sim := iam.SimulateCustomPolicyInput{
			ActionNames:     statement.NotAction,
			PolicyInputList: []*string{&fullPolicy},
			CallerArn:       &principal,
		}

		simRes, err := svc.SimulateCustomPolicy(&sim)
		if err != nil {
			fmt.Println("Error gettting response", err)
		}

		for _, res := range simRes.EvaluationResults {
			if *res.EvalDecision == "allowed" {
				fmt.Printf("NotAction %s works!\n", *res.EvalActionName)
			}
		}
	}

}
