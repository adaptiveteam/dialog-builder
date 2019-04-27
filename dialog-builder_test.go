package dialog_builder

import (
	"os"
	"testing"
)

func Dialog_Build(t *testing.T) {
	dc := NewDialogData(
		os.Getenv("DIALOG_ORGANIZATION"),
		os.Getenv("DIALOG_REPO"),
		os.Getenv("DIALOG_DIRECTORY"),
		os.Getenv("DIALOG_CATALOG"),
		os.Getenv("DIALOG_TABLE"),
		os.Getenv("LEARN_MORE_REPO"),
		os.Getenv("LEARN_MORE_DIRECTORY"),
		os.Getenv("BUILD_BRANCH"),
		os.Getenv("CULTIVATION_BRANCH"),
		os.Getenv("MASTER_BRANCH"),
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
