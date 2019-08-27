package dialog_builder

import (
	"fmt"
	awsutils "github.com/adaptiveteam/aws-utils-go"
	"os"
	"testing"
)

func Test_Dialog_Build(t *testing.T) {
	environments := map[string]string{
		"lexcorp":"us-east-2",
		"hoger":"us-east-1",
		"ivan":"us-east-1",
		"staging":"us-east-1",
	}

	for e,r := range environments {
		fmt.Println("Environment -",e)
		dynamo := awsutils.NewDynamo(r, "", "dialog")
		dc := NewDialogData(
			dynamo,
			e,
			os.Getenv("DIALOG_ORGANIZATION"),
			os.Getenv("DIALOG_REPO"),
			os.Getenv("DIALOG_DIRECTORY"),
			os.Getenv("DIALOG_CATALOG"),
			os.Getenv("DIALOG_TABLE"),
			os.Getenv("ALIAS_DIRECTORY"),
			os.Getenv("LEARN_MORE_REPO"),
			os.Getenv("LEARN_MORE_DIRECTORY"),
			os.Getenv("BUILD_BRANCH"),
			os.Getenv("CULTIVATION_BRANCH"),
			os.Getenv("MASTER_BRANCH"),
		)
		buildErrors := Build(&dc)
		for k,v := range buildErrors {
			fmt.Printf("Ran into an error of type %v with error %v\n", k,v)
		}
	}

}
