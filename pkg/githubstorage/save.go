package githubstorage

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

var repoName string = "openai-prompts-save"
var fileName string = "prompts.md"

func SaveInput(content string, response string) error {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	owner, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return err
	}
	repo, _, err := client.Repositories.Get(ctx, owner.GetLogin(), repoName)
	if err != nil && !strings.Contains(err.Error(), "404 Not Found") {
		return err
	}

	if repo == nil {
		repo, _, err = client.Repositories.Create(ctx, "", &github.Repository{
			Name: github.String(repoName),
		})
		if err != nil {
			return err
		}
		err = initRepo(repo)
		if err != nil {
			return err
		}

	}

	// clone the repo
	tmpDir, err := os.MkdirTemp("", "openai-prompts-save")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)
	r, err := git.PlainClone(tmpDir, false, &git.CloneOptions{
		URL: repo.GetCloneURL(),
	})

	if err != nil {
		return err
	}

	return addCommitPush(tmpDir, repo, fileName, content, response, r)

}

func initRepo(repo *github.Repository) error {

	// Initialize a local repository in a temporary directory
	localRepoPath, _ := os.MkdirTemp("", repoName)
	_, _ = git.PlainInit(localRepoPath, false)

	// Create a README.md file
	readmePath := localRepoPath + "/README.md"
	_ = os.WriteFile(readmePath, []byte("#Init"), 0644)

	// Open the local repository
	r, _ := git.PlainOpen(localRepoPath)

	// Stage the changes
	return addCommitPush(localRepoPath, repo, "README.md", "Inital Commit", "", r)
}

func addCommitPush(tmpDir string, repo *github.Repository, fileName string, content string, body string, r *git.Repository) error {
	// Create or update the file
	filePath := fmt.Sprintf("%s/%s", tmpDir, fileName)

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString("### " + content + "\n")
	if err != nil {
		return err
	}

	arrOfStrings := strings.Split(body, "\n")
	for _, str := range arrOfStrings {
		_, err = file.WriteString("> " + str + "\n")
		if err != nil {
			return err
		}
	}

	// Add the changes to the index
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	_, err = w.Add(fileName)
	if err != nil {
		return err
	}

	// Commit the changes
	_, err = w.Commit(content, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "bot",
			Email: "dathanvp+gptgitsave@gmail.com",
			When:  time.Now(),
		},
	})

	if err != nil {
		return err
	}

	// Set the remote URL
	_, err = r.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{*repo.CloneURL},
	})
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return err
	}

	err = r.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: "dathan",
			Password: string(os.Getenv("GITHUB_TOKEN")),
		},
	})

	if err != nil {
		return err
	}

	return nil

}
