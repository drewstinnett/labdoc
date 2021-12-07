/*
Copyright © 2021 Drew Stinnett <drew@drewlink.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/apex/log"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
)

// ciGitlabCmd represents the ciGitlab command.
var ciGitlabCmd = &cobra.Command{
	Use:   "gitlab",
	Short: "Run CI process inside GitLab CI",
	Run: func(cmd *cobra.Command, args []string) {
		gitDirectory, err := cmd.Flags().GetString("git-directory")
		cobra.CheckErr(err)
		requiredEnvs := []string{
			"GITLAB_USER_EMAIL", "GITLAB_USER_NAME", "GITLAB_TOKEN", "GITLAB_USER_LOGIN", "CI_SERVER_HOST",
			"CI_PROJECT_ROOT_NAMESPACE", "CI_PROJECT_NAME", "GITLAB_TOKEN",
		}
		var missingEnvs []string
		for _, requiredEnv := range requiredEnvs {
			if _, ok := os.LookupEnv(requiredEnv); !ok {
				missingEnvs = append(missingEnvs, requiredEnv)
			}
		}
		if len(missingEnvs) > 0 {
			log.Fatalf("Missing required environment variables: %v", missingEnvs)
		}

		log.Info("Opening Git directory")
		r, err := git.PlainOpen(gitDirectory)
		cobra.CheckErr(err)

		log.Info("Setting Working Directory")
		w, err := r.Worktree()
		cobra.CheckErr(err)

		log.Info("Adding README.md")
		_, err = w.Add("README.md")
		cobra.CheckErr(err)

		log.Info("Getting status")
		status, err := w.Status()
		cobra.CheckErr(err)

		fmt.Println(status)
		_, err = w.Commit("chore: labdoc update", &git.CommitOptions{
			Author: &object.Signature{
				Name:  os.Getenv("GITLAB_USER_NAME"),
				Email: os.Getenv("GITLAB_USER_EMAIL"),
				When:  time.Now(),
			},
		})
		cobra.CheckErr(err)

		remotes, err := r.Remotes()
		cobra.CheckErr(err)
		var hasGitlabRemote bool
		for _, remote := range remotes {
			if remote.Config().Name == "gitlab" {
				hasGitlabRemote = true
			}
		}
		if !hasGitlabRemote {
			log.Info("Creating remote")
			_, err = r.CreateRemote(&config.RemoteConfig{
				Name: "gitlab",
				URLs: []string{
					fmt.Sprintf("https://%v:%v@%v/%v/%v.git",
						os.Getenv("GITLAB_USER_LOGIN"),
						os.Getenv("GITLAB_TOKEN"),
						os.Getenv("CI_SERVER_HOST"),
						os.Getenv("CI_PROJECT_ROOT_NAMESPACE"),
						os.Getenv("CI_PROJECT_NAME"),
					),
				},
			})
			cobra.CheckErr(err)
		}

		log.Info("Doing push")
		err = r.Push(&git.PushOptions{
			RemoteName: "origin",
			RefSpecs: []config.RefSpec{
				config.RefSpec("refs/heads/main:refs/heads/main"),
			},
		})
		cobra.CheckErr(err)
	},
}

func init() {
	ciCmd.AddCommand(ciGitlabCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ciGitlabCmd.PersistentFlags().String("foo", "", "A help for foo")
	ciGitlabCmd.PersistentFlags().StringP("git-directory", "g", ".", "Path to the git directory")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ciGitlabCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
