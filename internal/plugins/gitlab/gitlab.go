package gitlabdoc

import (
	"fmt"
	"html/template"
	"os"

	"github.com/drewstinnett/labdoc/pkg/labdoc"
	"github.com/xanzy/go-gitlab"
)

type plug struct{}

func (p *plug) TemplateFunctions() (template.FuncMap, error) {
	templ := template.FuncMap{
		"recentlyCreatedProjects":       recentlyCreatedProjects,
		"recentEvents":                  recentEvents,
		"recentlyAcceptedMergeRequests": recentlyAcceptedMergeRequests,
	}
	return templ, nil
}

func recentlyCreatedProjects(limit int) ([]*gitlab.Project, error) {
	git, me, err := newClient()
	if err != nil {
		return nil, err
	}

	owned := true
	orderBy := "created_at"
	retProjects := []*gitlab.Project{}
	opts := &gitlab.ListProjectsOptions{
		Owned:   &owned,
		OrderBy: &orderBy,
		ListOptions: gitlab.ListOptions{
			PerPage: 50,
			Page:    1,
		},
	}
	for {
		projects, resp, err := git.Projects.ListProjects(opts)
		if err != nil {
			return nil, err
		}
		if resp.NextPage == 0 {
			break
		}
		for _, project := range projects {
			if project.CreatorID == me.ID {
				retProjects = append(retProjects, project)
				if len(retProjects) == limit {
					return retProjects, nil
				}
			}
		}
		opts.Page = resp.NextPage

	}

	return retProjects, nil
}

func newClient() (*gitlab.Client, *gitlab.User, error) {
	token := os.Getenv("GITLAB_TOKEN")
	url := os.Getenv("GITLAB_URL")

	if token == "" {
		return nil, nil, fmt.Errorf("GITLAB_TOKEN not set")
	}

	var git *gitlab.Client
	var err error
	if url != "" {
		git, err = gitlab.NewClient(token, gitlab.WithBaseURL(fmt.Sprintf("%v/api/v4", url)))
		if err != nil {
			return nil, nil, err
		}
	} else {
		git, err = gitlab.NewClient(token)
		if err != nil {
			return nil, nil, err
		}
	}
	me, _, err := git.Users.CurrentUser()
	if err != nil {
		return nil, nil, err
	}
	return git, me, nil
}

func recentlyAcceptedMergeRequests(limit int) ([]EnhancedEvent, error) {
	// return recentEvents(gitlab.MergeRequestEventTargetType, limit)
	git, _, err := newClient()
	if err != nil {
		return nil, err
	}

	var retEnhancedEvents []EnhancedEvent
	t := gitlab.MergeRequestEventTargetType
	opts := &gitlab.ListContributionEventsOptions{
		TargetType: &t,
		ListOptions: gitlab.ListOptions{
			PerPage: 50,
			Page:    1,
		},
	}
	for {
		events, resp, err := git.Events.ListCurrentUserContributionEvents(opts)
		if err != nil {
			return nil, err
		}
		if resp.NextPage == 0 {
			break
		}
		for _, event := range events {
			// Only look at accepted
			if event.ActionName != "accepted" {
				continue
			}
			p, _, err := git.Projects.GetProject(event.ProjectID, nil)
			if err != nil {
				return nil, err
			}
			e := EnhancedEvent{
				Project: p,
				Event:   event,
			}
			retEnhancedEvents = append(retEnhancedEvents, e)
			if len(retEnhancedEvents) == limit {
				return retEnhancedEvents, nil
			}
		}
		opts.Page = resp.NextPage

	}

	return retEnhancedEvents, nil
}

func recentEvents(eventType gitlab.EventTargetTypeValue, limit int) ([]EnhancedEvent, error) {
	git, _, err := newClient()
	if err != nil {
		return nil, err
	}

	var retEnhancedEvents []EnhancedEvent
	// t := gitlab.ProjectEventTargetType
	opts := &gitlab.ListContributionEventsOptions{
		TargetType: &eventType,
		ListOptions: gitlab.ListOptions{
			PerPage: 50,
			Page:    1,
		},
	}
	for {
		events, resp, err := git.Events.ListCurrentUserContributionEvents(opts)
		if err != nil {
			return nil, err
		}
		if resp.NextPage == 0 {
			break
		}
		for _, event := range events {
			p, _, err := git.Projects.GetProject(event.ProjectID, nil)
			if err != nil {
				return nil, err
			}
			e := EnhancedEvent{
				Project: p,
				Event:   event,
			}
			retEnhancedEvents = append(retEnhancedEvents, e)
			if len(retEnhancedEvents) == limit {
				return retEnhancedEvents, nil
			}
		}
		opts.Page = resp.NextPage

	}

	return retEnhancedEvents, nil
}

func init() {
	labdoc.Add("gitlab", func() labdoc.Plugin {
		return &plug{}
	})
}

type EnhancedEvent struct {
	Event   *gitlab.ContributionEvent
	Project *gitlab.Project
}
