package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/Julian-Kai/gitlab-search/internal/helpers"
	"github.com/Julian-Kai/gitlab-search/internal/services"
)

const (
	DelayCallSeconds = 1
	MaxSearchResults = 5
)

var SearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Perform a thorough search of your GitLab projects",
	Long:  "Perform a thorough search of your GitLab projects",
	RunE:  SearchCmdRun,
}

func main() {
	InitialCmdFlags()
	rootCmd := &cobra.Command{Use: "gs"}
	rootCmd.AddCommand(SearchCmd)
	rootCmd.Execute()
}

func InitialCmdFlags() {
	var url string
	var token string
	var keyword string

	SearchCmd.Flags().StringVarP(&url, "url", "u", "", "gitlab url")
	SearchCmd.Flags().StringVarP(&token, "token", "t", "", "personal access token")
	SearchCmd.Flags().StringVarP(&keyword, "keyword", "k", "", "search keyword")
}

func SearchCmdRun(cmd *cobra.Command, args []string) error {
	url, err := cmd.Flags().GetString("url")
	println("url: ", url)
	if err != nil {
		return err
	}
	token, err := cmd.Flags().GetString("token")
	println("token: ", token)
	if err != nil {
		return err
	}
	keyword, err := cmd.Flags().GetString("keyword")
	println("keyword: ", keyword)
	if err != nil {
		return err
	}

	svc, err := services.NewGitLabService(url, token)
	if err != nil {
		log.Fatalf("failed to create gitlab client: %v", err)
	}

	// get groups
	groupIDs, err := getGroups(svc)
	if err != nil {
		log.Fatalf("failed to get groups: %v", err)
		return err
	}

	// get projects
	projects, err := getProjects(svc, groupIDs)
	if err != nil {
		log.Fatalf("failed to get projects: %v", err)
		return err
	}

	// do search
	for _, p := range projects {
		blobs, costTimeDuration, err := svc.Search(p.ID, keyword, MaxSearchResults+1)
		if err != nil {
			log.Fatalf("failed to search keyword: %v", err)
			return err
		}

		printResults(p.Name, blobs, costTimeDuration)
		time.Sleep(DelayCallSeconds * time.Second)
	}
	return nil
}

func printResults(projectName string, blobs []*services.Blob, costTimeDuration time.Duration) {
	if len(blobs) > 0 {
		var size string
		var comment string

		if len(blobs) > MaxSearchResults {
			size = fmt.Sprintf("%d+", MaxSearchResults)
			comment = fmt.Sprintf("(only show %d results)", MaxSearchResults)
		} else {
			size = strconv.Itoa(len(blobs))
		}

		fmt.Printf("🔍 Project [%s] has [%s] results %s, cost %d ms \n\n", projectName, size, comment, costTimeDuration.Milliseconds())

		for i := 0; i < helpers.Min(MaxSearchResults, len(blobs)); i++ {
			b := blobs[i]
			fmt.Printf("👉 %s\n\n", b.Path)
			fmt.Printf("```# branch: %s, line: %d\n", b.Ref, b.Line)
			fmt.Printf("%s\n", strings.Trim(strings.Replace(b.Data, "\t", "  ", -1), "\n"))
			fmt.Printf("```\n\n")
		}
	} else {
		fmt.Printf("🔍 Project [%s] has no code results, cost %d ms\n\n", projectName, costTimeDuration.Milliseconds())
	}
}

func getProjects(svc services.GitLabSvc, groupIDs []int) ([]*services.Project, error) {
	res := make([]*services.Project, 0)
	for _, gid := range groupIDs {
		projects, err := svc.GetProjects(gid)
		if err != nil {
			return nil, err
		}
		res = append(res, projects...)
	}
	fmt.Printf("There are [%d] projects\n", len(res))
	return res, nil
}

func getGroups(svc services.GitLabSvc) ([]int, error) {
	groupIDs, err := svc.GetGroups()
	if err != nil {
		return nil, err
	}
	fmt.Printf("There are [%d] groups\n", len(groupIDs))
	return groupIDs, nil
}
