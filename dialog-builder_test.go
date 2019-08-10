package dialog_builder

import (
	"os"
	"testing"
	utils "github.com/adaptiveteam/adaptive-utils-go"
)

func Dialog_Build_Test(t *testing.T) {
	dc := NewDialogData(
		utils.NonEmptyEnv("DIALOG_ORGANIZATION"),
		utils.NonEmptyEnv("DIALOG_REPO"),
		utils.NonEmptyEnv("DIALOG_DIRECTORY"),
		utils.NonEmptyEnv("DIALOG_CATALOG"),
		utils.NonEmptyEnv("DIALOG_TABLE"),
		utils.NonEmptyEnv("LEARN_MORE_REPO"),
		utils.NonEmptyEnv("LEARN_MORE_DIRECTORY"),
		utils.NonEmptyEnv("BUILD_BRANCH"),
		utils.NonEmptyEnv("CULTIVATION_BRANCH"),
		utils.NonEmptyEnv("MASTER_BRANCH"),
	)

	errorBuild, errorCultivatePR, errorMasterPR, errorLearnMorePR := Build(&dc)
	if errorBuild != nil ||
		errorCultivatePR != nil ||
		errorMasterPR != nil ||
		errorLearnMorePR != nil {
		t.Errorf("Build errorS %v,%v,%v,%v",
			errorBuild,
			errorCultivatePR,
			errorMasterPR,
			errorLearnMorePR,
		)
	}
}
