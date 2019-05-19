package dialog_builder

import (
	"fmt"
	"os"
	"testing"
)

func Test_Dialog_Build(t *testing.T) {
	dc := NewDialogData(
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
