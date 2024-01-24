package main

import (
	"fmt"

	gitlab "github.com/xanzy/go-gitlab"
)

var (
	gitClient  *gitlab.Client
	token      string = "WsgTfAskGjsnns1t1HjM"
	url        string = "https://gitlab.cee.redhat.com"
	projectID  int    = 84399
	mrID       int    = 5
	sha        string = "c7e9ab73603131cc9b91a729d9ab2c27be204b24"
	targetUrl  string = "https://console.redhat.com/preview/application-pipeline/workspaces/hongweiliu/applications/go-test/pipelineruns/go-test-mk45c-bvprc"
	desciption string = "Failed"
	name       string = "RHTAP Integration test - scenario 1 - snapshot 1"
	context    string = "mr 5"
)

func main() {
	var err error
	gitClient, err = gitlab.NewClient(token, gitlab.WithBaseURL(url))
	if err != nil {
		fmt.Printf("initclienterr:%v\n", err)
		panic(err)
	} else {
		fmt.Println("initialize...")
	}

	opt := &gitlab.SetCommitStatusOptions{
		State:       gitlab.BuildStateValue("failed"),
		Name:        gitlab.Ptr(name),
		TargetURL:   gitlab.Ptr(targetUrl),
		Description: gitlab.Ptr(desciption),
	}

	//Create CommitStatus for commit
	commitStatus, _, err := gitClient.Commits.SetCommitStatus(projectID, sha, opt)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("succeed to create commitStatus with %d...\n", commitStatus.PipelineId)
	}

	//Update CommitStatus for Pipeline created above
	opt = &gitlab.SetCommitStatusOptions{
		State:       gitlab.BuildStateValue("failure"),
		Name:        gitlab.Ptr(name),
		TargetURL:   gitlab.Ptr(targetUrl),
		Description: gitlab.Ptr(desciption),
		PipelineID:  gitlab.Ptr(commitStatus.PipelineId),
	}
	commitStatus, _, err = gitClient.Commits.SetCommitStatus(projectID, sha, opt)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("succeed to update commitStatus %d\n", commitStatus.PipelineId)
	}

	//Create a note for merge request
	body := `Integration test for snapshot snapshot 1 and scenario 1 has <b>Failed</b>

| Task | Duration | Test Suite | Status | Details 
|----------|----------|----------|----------|----------|
| task 1 |28 seconds | Default | :x: Failed | :white_check_mark: 5 success(es) <br> :warning: 1 warning(s) <br> :x: 1 failure(s) |
| task 2 |3 minutes  | Default | :white_check_mark: Success | :white_check_mark: 80 success(es) <br> :warning: 3 warning(s)|
| task 3 |2 seconds  | Default | :warning: Warning | :white_check_mark: 10 success(es) <br> :warning: 1 warning(s)`

	mopt_create := &gitlab.CreateMergeRequestNoteOptions{Body: gitlab.Ptr(body)}
	comment, _, err := gitClient.Notes.CreateMergeRequestNote(projectID, mrID, mopt_create)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("succeed to create comment %d\n", comment.ID)
	}

	mopt_update := &gitlab.UpdateMergeRequestNoteOptions{Body: gitlab.Ptr(body)}
	_, _, err = gitClient.Notes.UpdateMergeRequestNote(projectID, mrID, comment.ID, mopt_update)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("succeed to update comment %d\n", comment.ID)
	}
}
