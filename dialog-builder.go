package dialog_builder


// NewDialogData is helper function that enables users to correctly
// create instantiate a DialogData structure
func NewDialogData(
	organization string,
	dialogRepo string,
	dialogFolder string,
	dialogCatalog string,
	dialogTable string,
	learnMoreRepo string,
	learnMoreFolder string,
	buildBranch string,
	cultivationBranch string,
	masterBranch string,
) (rv DialogData) {
	if organization       == "" ||
		dialogRepo        == "" ||
		dialogFolder      == "" ||
		dialogCatalog     == "" ||
		dialogTable       == "" ||
		learnMoreRepo     == "" ||
		learnMoreFolder   == "" ||
		buildBranch       == "" ||
		cultivationBranch == "" ||
		masterBranch      == ""	{
		panic("cannot have empty initialization values")
	}

	rv.Organization = organization
	rv.DialogRepo = dialogRepo
	rv.DialogFolder = dialogFolder
	rv.DialogCatalog = dialogCatalog
	rv.DialogTable = dialogTable
	rv.LearnMoreRepo = learnMoreRepo
	rv.LearnMoreFolder = learnMoreFolder
	rv.BuildBranch = buildBranch
	rv.CultivationBranch = cultivationBranch
	rv.MasterBranch = masterBranch
	rv.Modified = false
	return rv
}

// DialogData stores all of the infrastructure data necessary to work with GitHub & Dynamo
// organization is the organization that owns the dialog repo
// dialogRepo is the name of the repo where the dialog can be found
// dialogFolder is the folder where the dialog can be found
// dialogCatalog is the file for the dialog catalog file
// dialogTable is the DynamoDB table use to store the dialog
// learnMoreRepo is the repo where the Learn More content can be found
// learnMoreFolder is the directory where the Learn More content can be found
// buildBranch is the branch to build against
// cultivationBranch is the branch to update with cultivation work
// masterBranch is the master branch
type DialogData struct {
	Organization string
	DialogRepo string
	DialogFolder string
	DialogCatalog string
	DialogTable string
	LearnMoreRepo string
	LearnMoreFolder string
	BuildBranch string
	CultivationBranch string
	MasterBranch string
	Modified bool
}

func Build(dc *DialogData) (
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

	if errorBuild == nil && dc.Modified {
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
