package dialog_builder

import (
	"context"
	"fmt"
	"github.com/adaptiveteam/aws-utils-go"
	"github.com/adaptiveteam/core-utils-go"
	"github.com/google/go-github/github"
	"github.com/adaptiveteam/dialog-fetcher"
	"golang.org/x/oauth2"
	"os"
	"sort"
	"strings"
	"time"
)

// The following  global variables are to reduce cold-start times in Lambda execution
var (
	ctx                      = context.Background()
	ts                       = oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token},)
	tc                       = oauth2.NewClient(ctx, ts)
	client                   = github.NewClient(tc)
	dynamo                   = awsutils.NewDynamo(os.Getenv("AWS_REGION"), "", "dialog")
	token                    = os.Getenv("GITHUB_API_KEY")
)

const (
	DIALOG_ID_PREFIX = "// DIALOG_ID: "
	LEARN_MORE_PREFIX = "// LEARN_MORE: "
)

func updateFile(
	org string,
	repo string,
	branch string,
	fileName string,
	newContents string,
	commitMessage string,
) (modified bool, err error){
	getOptions := &github.RepositoryContentGetOptions{
		Ref: "heads/"+branch,
	}
	oldFile,_,_,err := client.Repositories.GetContents(
		ctx,
		org,
		repo,
		fileName,
		getOptions,
	)
	if err == nil && oldFile.GetType() == "file" {
		oldContents, err := oldFile.GetContent()
		if err == nil && oldContents != newContents {
			newContentBytes := []byte(newContents)
			sha := oldFile.GetSHA()
			repositoryContentsOptions := &github.RepositoryContentFileOptions{
				Message:   &commitMessage,
				Content:   newContentBytes,
				SHA:       &sha,
				Branch:    &branch,
				Committer: nil,
			}
			_, _, err = client.Repositories.UpdateFile(
				context.Background(),
				org,
				repo,
				fileName,
				repositoryContentsOptions,
			)
			modified = true
		}
	}else {
		err = fmt.Errorf("%s is not a file", fileName)
	}
	return modified, err
}

func storeDialog(
	dc *DialogData,
	dialogCoordinates string,
	dialogSubject string,
	dialog []string,
	comments []string,
	dialogID string,
	learnMoreLink string,
	learnMoreContent string,
) (err error) {
	updated := time.Now().Format("2006-01-02")
	item := fetch_dialog.DialogEntry{
		Context:dialogCoordinates,
		Subject:dialogSubject,
		Updated:updated,
		Dialog:dialog,
		Comments:comments,
		DialogID: dialogID,
		LearnMoreLink: learnMoreLink,
		LearnMoreContent:learnMoreContent,
		BuildBranch:dc.BuildBranch,
		CultivationBranch:dc.CultivationBranch,
		MasterBranch:dc.MasterBranch,
	}
	err = dynamo.PutTableEntry(item, dc.DialogTable)

	return err
}

func getLearnMoreContent(dc *DialogData, dialogID string) (content string, link string, err error){
	getOptions := &github.RepositoryContentGetOptions{
		Ref: "heads/"+dc.CultivationBranch,
	}
	fileName := dc.LearnMoreFolder+"/"+dialogID+".md"
	var response *github.Response
	var learnMoreFile *github.RepositoryContent
	learnMoreFile,_,response,err = client.Repositories.GetContents(
		ctx,
		dc.Organization,
		dc.LearnMoreRepo,
		fileName,
		getOptions,
	)

	if err == nil {
		content,err = learnMoreFile.GetContent()
		if err == nil {
			link = "https://github.com/"
			link = link + dc.Organization+"/"
			link = link + dc.LearnMoreRepo+"/blob/"
			link = link + dc.CultivationBranch+"/"
			link = link + dc.LearnMoreFolder+"/"+dialogID+".md"
		}
	} else if response.StatusCode == 404 {
		content = ""
		link = ""
		err = nil
	}
	return content, link, err
}

func updateDialogFile(dc *DialogData, newContent string, path string, commitMessage string) (err error) {
	dc.Modified, err = updateFile(
		dc.Organization,
		dc.DialogRepo,
		dc.BuildBranch,
		path,
		newContent,
		commitMessage,
	)
	return  err
}

func compileDialogFile(
	dialogComments []string,
	dialogID string,
	learnMoreLink string,
	dialogLines []string,
) (file string){
	for _,dc := range dialogComments {
		file = file+"# "+dc+"\n"
	}
	if dialogID != "" {
		file = file + DIALOG_ID_PREFIX+dialogID+"\n"
	}

	if learnMoreLink != "" {
		file = file + LEARN_MORE_PREFIX+learnMoreLink+"\n"
	}

	for _,dl := range dialogLines {
		file = file + dl+"\n"
	}
	return file
}

func parseDialogFile(blob string) (
	dialog []string,
	comments []string,
	dialogID string,
	learnMoreLink string,
) {

	dialog = make([]string,0)
	comments = make([]string,0)
	var fileLines = strings.Split(blob,"\n")
	for _,d := range fileLines {
		if len(d) > 0 {
			var trimmed string
			if strings.HasPrefix(d,"#") {
				trimmed = strings.Trim(strings.Trim(d,"#")," ")
				comments = append(comments,trimmed)
			} else if strings.HasPrefix(d,DIALOG_ID_PREFIX) {
				trimmed = strings.Trim(strings.Trim(d,DIALOG_ID_PREFIX)," ")
				dialogID = trimmed
			} else if strings.HasPrefix(d,LEARN_MORE_PREFIX) {
				trimmed = strings.Trim(strings.Trim(d,LEARN_MORE_PREFIX)," ")
				learnMoreLink = trimmed
				if learnMoreLink == "" {
					learnMoreLink = "ERROR!"
				}
			} else {
				dialog = append(dialog,strings.Trim(d," "))
			}
		}
	}
	return dialog,comments,dialogID, learnMoreLink
}

func postDialog(
	dc *DialogData,
	dialog *github.RepositoryContent,
)(err error) {
	getOptions := &github.RepositoryContentGetOptions{
		Ref:"heads/"+dc.BuildBranch,
	}
	contents,_,_,err := client.Repositories.GetContents(
		ctx,
		dc.Organization,
		dc.DialogRepo,
		dialog.GetPath(),
		getOptions,
	)

	if err == nil {
		dialogBlob,err := contents.GetContent()
		if err == nil {
			var learnMoreContent string
			var commitMessage string
			dialogLines,dialogComments,dialogID,oldLearnMoreLink := parseDialogFile(dialogBlob)

			// If there is no dialog ID then add one
			if dialogID == "" {
				dialogID = core_utils_go.Uuid()
				commitMessage = commitMessage + "Adding dialog UUID. "
			}

			learnMoreContent, newLearnMoreLink, err := getLearnMoreContent(dc,dialogID)

			if err == nil {
				// Now construct the learn more web page link
				// the newLearnMoreLink is meant for cultivators
				// the learnMoreWebLink is meant for end users
				var learnMoreWebLink string
				if newLearnMoreLink != "" {
					learnMoreWebLink = "https://"+dc.LearnMoreRepo+"/"+dc.LearnMoreFolder+"/"+dialogID
				}

				// Make sure the learn more links are the same.
				// The repo could have moved or there might have been an error!
				if newLearnMoreLink != oldLearnMoreLink {
					if newLearnMoreLink != "" {
						commitMessage = commitMessage + "Adding Learn More link.\n"
					} else {
						commitMessage = commitMessage + "Removing Learn More link.\n"
					}
				}

				newContent := compileDialogFile(
					dialogComments,
					dialogID,
					newLearnMoreLink,
					dialogLines,
				)
				if newContent != dialogBlob {
					commitMessage = commitMessage + "Fixing up content"
				}

				// If there were any changes then we need to update the file
				if commitMessage != "" {
					err = updateDialogFile(
						dc,
						newContent,
						contents.GetPath(),
						commitMessage,
					)
				}

				if err == nil {
					// Strip off .txt from the subject
					subject := strings.Split(contents.GetName(),".")[0]
					path := dialog.GetPath()
					// remove the suject file and then replace the slashes in the path with #
					parts := strings.Split(path, "/")
					contextCoordinates := strings.Join(parts[:len(parts)-1],"#")+"#"
					err = storeDialog(
						dc,
						contextCoordinates,
						subject,
						dialogLines,
						dialogComments,
						dialogID,
						learnMoreWebLink,
						learnMoreContent,
					)
				}
			}
		}
	}
	return err
}

func crawlContext(
	dc *DialogData,
	path string,
) (err error) {
	getOptions := &github.RepositoryContentGetOptions{
		Ref:"heads/"+dc.BuildBranch,
	}
	_,directories,_,err := client.Repositories.GetContents(
		ctx,
		dc.Organization,
		dc.DialogRepo,
		path,
		getOptions,
	)

	for i:=0; i < len(directories) && err == nil;i++ {
		if directories[i].GetType() == "file" && strings.HasSuffix(directories[i].GetName(),".txt"){
			err = postDialog(
				dc,
				directories[i],
			)
		} else if directories[i].GetType() == "dir" {
			err = crawlContext(
				dc,
				directories[i].GetPath(),
			)
		}
	}
	return err
}

func pullRequestExists(
	dc *DialogData,
	repo string,
	head string,
	base string,
) (found bool){
	prs,_,err := client.PullRequests.List(
		ctx,
		dc.Organization,
		repo,
		nil,
		)
	found = false
	for i := 0; err == nil && i < len(prs) && !found; i++ {
		if *prs[i].Head.Ref==head && *prs[i].Base.Ref==base {
			found = true
		}
	}

	return found
}

func createPullRequest(
	dc *DialogData,
	prSummmary string,
	prDetails,
	repo string,
	sourceBranch string,
	targetBranch string,
) (pr *github.PullRequest, err error) {
	newPR := &github.NewPullRequest{
		Title:               github.String(prSummmary),
		Head:                github.String(sourceBranch),
		Base:                github.String(targetBranch),
		Body:                github.String(prDetails),
		MaintainerCanModify: github.Bool(true),
	}

	pr, _, err = client.PullRequests.Create(context.Background(), dc.Organization, repo, newPR)
	if err != nil {
		fmt.Println(err)
		return
	}

	return pr, err
}

func loadDialog(dc *DialogData) error {
	getOptions := &github.RepositoryContentGetOptions{
		Ref: "heads/"+dc.BuildBranch,
	}
	_,directory,_,err := client.Repositories.GetContents(
		ctx,
		dc.Organization,
		dc.DialogRepo,
		"",
		getOptions,
	)
	found := false

	for i := 0; i < len(directory) && !found && err == nil; i++{
		if directory[i].GetName() == dc.DialogFolder {
			err = crawlContext(
				dc,
				directory[i].GetPath(),
			)
			found  = true
		}
	}

	if err ==  nil && !found {
		err = fmt.Errorf("unable to find dialog directory %s", dc.DialogFolder)
	}
	return err
}

func getAllContent(dialogTable string) (dialogEntries []fetch_dialog.DialogEntry, err error) {
	//scan the table after deletion
	dialogEntries = make([]fetch_dialog.DialogEntry,0)
	err = dynamo.ScanTable(dialogTable, &dialogEntries)
	if err == nil {
		sort.SliceStable(
			dialogEntries,
			func(i, j int) bool {
				return dialogEntries[i].Context <= dialogEntries[j].Context
			},
		)
	}
	return dialogEntries,err
}

func generateCatalog(
	dc *DialogData,
	fileName string,
) (report string, err error) {
	allContent, err := getAllContent(dc.DialogTable)
	if err == nil {
		baseURL  := "https://github.com/"+dc.Organization
		baseDialogEditURL := baseURL+"/"+dc.DialogRepo+"/edit/"+dc.CultivationBranch+"/"
		baseDialogViewURL := baseURL+"/"+dc.DialogRepo+"/blob/"+dc.CultivationBranch+"/"
		baseCreateLearnMoreURL := baseURL+"/"+dc.LearnMoreRepo+"/create/"+dc.CultivationBranch+"/"+dc.LearnMoreFolder+"/"
		baseEditLearnMoreURL := baseURL+"/"+dc.LearnMoreRepo+"/edit/"+dc.CultivationBranch+"/"+dc.LearnMoreFolder+"/"
		baseViewLearnMoreURL := baseURL+"/"+dc.LearnMoreRepo+"/blob/"+dc.CultivationBranch+"/"+dc.LearnMoreFolder+"/"

		var tableOfContents = "# Table of Contents\n"
		var currentContext string
		var newContext string
		for _,de := range allContent {
			newContext = strings.Replace(de.Context,"#","/",-1)
			quickLink := baseDialogViewURL+fileName+"#context-"+strings.Replace(de.Context,"#","", -1)
			if newContext != currentContext {
				report = report + "## *Context: "+newContext+"*\n"
				currentContext = newContext
				tableOfContents = tableOfContents+"\n  * ["+currentContext+"]("+quickLink+")"
			}
			report = report + "### Subject: "+de.Subject+"\n"
			report = report + " [[Edit]]("+baseDialogEditURL+currentContext+de.Subject+".txt)"
			report = report + "[[View]]("+baseDialogViewURL+currentContext+de.Subject+".txt)\n"

			report = report + "#### Return to: [[Context]]("+quickLink+")"
			report = report + "[[Table of Contents]]("+baseDialogViewURL+fileName+"#table-of-contents"+")\n\n"
			report = report + "#### Dialog ID: "+de.DialogID+"\n\n"

			if len(de.Comments) > 0 {
				for _,c := range de.Comments {
					report = report +"*"+c+"*\n\n"
				}
			} else {
				report = report + "*No comments provided for this dialog.*\n\n"
			}
			for _,do := range de.Dialog {
				report = report + "  - "+do+"\n\n"
			}

			if de.LearnMoreContent == "" {
				report = report + "__*Learn More Page*__ for "+de.Subject+": [[Create]]("+baseCreateLearnMoreURL+de.DialogID+".md)\n"
			}  else {
				report = report + "#### Learn More Page"
				report = report + " [[Edit]]("+baseEditLearnMoreURL+de.DialogID+".md)"
				report = report + "[[View]]("+baseViewLearnMoreURL+de.DialogID+".md)\n"
				report = report + de.LearnMoreContent+"\n"
			}
		}
		report = tableOfContents+"\n\n"+report
	}
	return report, err
}

func updateCatalog(
	dc *DialogData,
	fileName string,
) (err error) {
	var newCatalogContents string
	commitMessage := "Catalog updated on "+time.Now().Format("2006-01-02 at 15:04:05")
	newCatalogContents,err = generateCatalog(
		dc,
		fileName,
	)
	var modified bool
	modified,err = updateFile(
		dc.Organization,
		dc.DialogRepo,
		dc.BuildBranch,
		fileName,
		newCatalogContents,
		commitMessage,
	)
	if !modified {
		err = fmt.Errorf("expected to modify dialog library but did not")
	}

	return err
}
