package analyzer

import (
	"fmt"
	"github.com/google/logger"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"os"
	"path/filepath"
	"strings"
)

// clone a repo and return its path
func CloneRepo(url string) (string, string) {
	var lastCommitHash string

	// 1. check GOPATH
	s := os.Getenv("GOPATH")
	if s == "" {
		logger.Fatal("could not find GOPATH," +
			" please make sure you have set the GOPATH!")
		return "", ""
	}
	// url string has /n /r
	url = strings.TrimSuffix(url, "\n")
	url = strings.TrimSuffix(url, "\r")
	var projectPath = filepath.Join(s, "src", "github.com", url)
	var repo *git.Repository
	if _, err := os.Stat(projectPath); os.IsExist(err) {
		// path/to/whatever exists, delete or git pull ?
		// TODO: delete or git pull
		logger.Infof("The path exists: %v", err)

	} else {
		// create the path
		err := os.MkdirAll(projectPath, 0755)
		if err != nil {
			logger.Fatalf("Failed to create path : %v", err)
			return "", ""
		}

		// Clones the repository into the given dir, just as a normal git clone does
		logger.Infof("Start cloning project: %s", url)
		r, err := git.PlainClone(projectPath, false, &git.CloneOptions{
			URL:      "https://github.com/" + url,
			Progress: os.Stdout,
		})
		repo = r
		if err != nil {
			logger.Info(err)
			// git pull origin master
			r, err := git.PlainOpen(projectPath)
			if err != nil {
				logger.Info(err)
			}
			repo = r
		}
	}

	//// 2. Tempdir to clone the repository
	//dir, err := ioutil.TempDir("", url)
	//if err != nil {
	//	log.Fatal(err)
	//}

	// get last commithash
	head, _ := repo.Head()
	if head != nil {
		cIter, err := repo.Log(&git.LogOptions{From: head.Hash()})
		var commits []*object.Commit
		err = cIter.ForEach(func(c *object.Commit) error {
			commits = append(commits, c)
			return nil
		})

		lastCommitHash = fmt.Sprintf("%s", commits[0].Hash)
		if err != nil {
			logger.Fatal(err)
		}
	}

	return projectPath, lastCommitHash
}
