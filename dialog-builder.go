package dialog_builder

import "github.com/adaptiveteam/core-utils-go"

// NewDialogData is helper function that enables users to correctly
// create instantiate a DialogData structure
func NewDialogData(
	organization,
	dialogRepo ,
	dialogFolder ,
	dialogCatalog ,
	dialogTable ,
	aliasFolder ,
	learnMoreRepo ,
	learnMoreFolder ,
	buildBranch ,
	cultivationBranch ,
	masterBranch string,
) (rv DialogData) {
	if organization       == "" ||
		dialogRepo        == "" ||
		dialogFolder      == "" ||
		dialogCatalog     == "" ||
		dialogTable       == "" ||
		aliasFolder      == "" ||
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
	rv.AliasFolder = aliasFolder
	rv.LearnMoreRepo = learnMoreRepo
	rv.LearnMoreFolder = learnMoreFolder
	rv.BuildBranch = buildBranch
	rv.CultivationBranch = cultivationBranch
	rv.MasterBranch = masterBranch
	rv.Modified = false
	rv.BuildID = core_utils_go.Uuid()
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
	AliasFolder string
	LearnMoreRepo string
	LearnMoreFolder string
	BuildBranch string
	CultivationBranch string
	MasterBranch string
	Modified bool
	BuildID string
}

func Build(dc *DialogData) (
	errors map[string]error,
){
	var err error
	errors = make(map[string]error,0)
	err = loadFile(
		dc,
		dc.DialogFolder,
		loadDialog,
	)
	if err != nil {
		errors["load"] = err
	}

	if err == nil {
		err = loadFile(
			dc,
			dc.AliasFolder,
			loadAliases,
		)
		if err != nil {
			errors["aliases"] = err
		}
	}

	if  err == nil {
		err = cleanUp(dc)
		if err != nil {
			errors["cleanup"] = err
		}
	}

	if err == nil {

		err = updateCatalog(
			dc,
			dc.DialogCatalog,
		)
		if err != nil {
			errors["update-catalog"] = err
		}
	}

	if err == nil && dc.Modified {
		if !pullRequestExists(
			dc,
			dc.DialogRepo,
			dc.BuildBranch,
			dc.CultivationBranch,
		) {
			_, err = createPullRequest(
				dc,
				"Moving dialog from build to cultivation",
				"Moving changes from the build branch created by the build process back down to cultivation",
				dc.DialogRepo,
				dc.BuildBranch,
				dc.CultivationBranch,
			)
			if err != nil {
				errors["build-to-cultivation-build-pr"] = err
			}
		}
		if err == nil && !pullRequestExists(
			dc,
			dc.DialogRepo,
			dc.BuildBranch,
			dc.MasterBranch,
		) {
			_, err = createPullRequest(
				dc,
				"Moving dialog from build to master",
				"Moving changes from the build branch created by the build process up to the master branch",
				dc.DialogRepo,
				dc.BuildBranch,
				dc.MasterBranch,
			)
			if err != nil {
				errors["build-to-master-pr"] = err
			}
		}
		if err == nil && !pullRequestExists(
			dc,
			dc.DialogRepo,
			dc.BuildBranch,
			dc.MasterBranch,
		) {
			_,err = createPullRequest(
				dc,
				"Moving from cultivation",
				"Moving changes from the cultivation branch up to the master branch",
				dc.LearnMoreRepo,
				dc.CultivationBranch,
				dc.MasterBranch,
			)
			if err != nil {
				errors["learn-more-pr"] = err
			}
		}
	}

	return errors
}
