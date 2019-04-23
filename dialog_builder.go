package dialog_builder

import (
	. "github.com/adaptiveteam/adaptive-utils-go/models"
)


func Build(dc DialogData) (
	errorBuild error,
	errorCultivatePR error,
	errorMasterPR error,
	errorLearnMorePR error,
){
	errorBuild = loadDialog(dc)
	if errorBuild == nil {
		errorBuild = updateCatalog(
			dc,
			dc.DialogCatalog,
		)
	}

	if errorBuild == nil {
		if !pullRequestExists(
			dc,
			dc.DialogRepo,
			dc.BuildBranch,
			dc.CultivationBranch,
		) {
			_, errorCultivatePR = createPullRequest(
				dc,
				"Moving from build to cultivation",
				"Moving changes from the build branch created by the build process back down to cultivation",
				dc.DialogRepo,
				dc.BuildBranch,
				dc.CultivationBranch,
			)
		}
		if !pullRequestExists(
			dc,
			dc.DialogRepo,
			dc.BuildBranch,
			dc.MasterBranch,
		) {
			_, errorMasterPR = createPullRequest(
				dc,
				"Moving from build to master",
				"Moving changes from the build branch created by the build process up to the master branch",
				dc.DialogRepo,
				dc.BuildBranch,
				dc.MasterBranch,
			)
		}
		if !pullRequestExists(
			dc,
			dc.DialogRepo,
			dc.BuildBranch,
			dc.MasterBranch,
		) {
			_,errorLearnMorePR = createPullRequest(
				dc,
				"Moving from cultivation",
				"Moving changes from the cultivation branch up to the master branch",
				dc.LearnMoreRepo,
				dc.CultivationBranch,
				dc.MasterBranch,
			)
		}
	}

	return errorBuild, errorCultivatePR, errorMasterPR, errorLearnMorePR
}
