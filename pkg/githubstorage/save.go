package githubstorage

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

var repoName string = "openai-prompts-save"
var fileName string = "prompts.md"

func SaveInput(content string) error {

	fmt.Printf("Saving input to GitHub: %s\n", content)

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
	if err != nil {
		return err
	}

	if repo == nil {
		fmt.Println("Creating repo")
		repo, _, err = client.Repositories.Create(ctx, "", &github.Repository{
			Name: github.String(repoName),
		})
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

	// Create or update the file
	filePath := fmt.Sprintf("%s/%s", tmpDir, fileName)
	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return err
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
	commit, err := w.Commit(content, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "bot",
			Email: "dathanvp+gptgitsave@gmail.com",
			When:  time.Now(),
		},
	})

	if err != nil {
		return err
	}

	fmt.Printf("Commit %s created\n", commit.String())
	err = r.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth: &http.TokenAuth{
			Token: string(os.Getenv("GITHUB_TOKEN")),
		},
	})

	if err != nil {
		return err
	}
	return nil
}
