package main

import (
	"fmt"
	"flag"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
)

type repositories struct {
	name string
	repo string
}

type Issues struct {
	TotalCount        int   `json:"total_count"`
	IncompleteResults bool  `json:"incomplete_results"`
	Items             Items `json:"items"`
}

type Items []struct {
	Url               string       `json:"url"`
	RepositoryUrl     string       `json:"repository_url"`
	LabelsUrl         string       `json:"labels_url"`
	CommentsUrl       string       `json:"comments_url"`
	EventsUrl         string       `json:"events_url"`
	HtmlUrl           string       `json:"html_url"`
	Id                int          `json:"id"`
	Number            int          `json:"number"`
	Title             string       `json:"title"`
	User              User         `json:"user"`
	Labels            Labels       `json:"labels"`
	State             string       `json:"state"`
	Locked            bool         `json:"locked"`
	Assignee          string       `json:"assignee"`
	Assignees         Assignees    `json:"assignees"`
	Milestone         string       `json:"milestone"`
	Comments          int          `json:"comments"`
	CreatedAt         string       `json:"create_at"`
	UpdatedAt         string       `json:"update_at"`
	ClosedAt          string       `json:"closed_at"`
	AuthorAssociation string       `json:"author_association"`
	PullRequests      PullRequests `json:"pull_requests"`
	Body              string       `json:"body"`
	Score             float64      `json:"score"`
}

type User struct {
	Login             string `json:"login"`
	Id                int    `json:"id"`
	AvatarUrl         string `json:"avatar_url"`
	GravatarId        string `json:"gravatar_id"`
	Url               string `json:"url"`
	HtmlUrl           string `json:"html_url"`
	FollowersUrl      string `json:"followers_url"`
	FollowingUrl      string `json:"following_url"`
	GistsUrl          string `json:"gists_url"`
	StarredUrl        string `json:"starred_url"`
	SubscribtionsUrl  string `json:"subscribtions_url"`
	OrganizationsUrl  string `json:"organizations_url"`
	ReposUrl          string `json:"repos_url"`
	EventsUrl         string `json:"events_url"`
	ReceivedEventsUrl string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type Labels []struct {
	Id      int    `json:"id"`
	Url     string `json:"url"`
	Name    string `json:"name"`
	Color   string `json:"color"`
	Default bool   `json:"default"`
}

type Assignees []struct{}
type PullRequests []struct {
	Url      string `json:"url"`
	HtmlUrl  string `json:"html_url"`
	DiffUrl  string `json:"diff_url"`
	PatchUrl string `json:"patch_url"`
}

var repo = flag.String("repo", "", "Name of project to search.")
var issue = flag.String("issue", "", "Specify issue that you are looking for.")
//
// Creates default repository array where repository names are stored.
// TODO: Create a method that allows to read a repository list from configuration file.
//
func setRepositories() []repositories {
	repos := []repositories{
		{
			name: "terraform",
			repo: "hashicorp/terraform",
		},
		{
			name: "ansible",
			repo: "ansible/ansible",
		},
	}
	return repos
}

//
// Methods that implements search issues in repository
//
func search() {
	r := *repo
	i := *issue
	repos := setRepositories()

	for key, value := range repos {
		if repos[key].name == r {
			client := &http.Client{}
			repoName := value
			req, err := http.NewRequest("GET", "https://api.github.com/search/issues", nil)
			if err != nil {
				fmt.Println("[ERROR] Some issue with request.")
			}
			q := req.URL.Query()
			q.Add("q", fmt.Sprintf("%s+repo:%s", i, repoName.repo))
			req.URL.RawQuery = q.Encode()
			resp, _ := client.Do(req)
			body, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			printResults(body)
		}
	}
}

//
// Method prints table results to stdout.
//
func printResults(body []byte) {
	var i Issues
	err := json.Unmarshal(body, &i)
	if err != nil {
		fmt.Println("[ERROR] Cannot unmarshal data. %s", err)
	}
	data := [][]string{}
	for _, value := range i.Items {
		if value.Score >= 1.000 {
			score := strconv.FormatFloat(value.Score, 'f', 4, 32)
			issueValue := []string{value.Title, value.State, value.HtmlUrl, score}
			data = append(data, issueValue)
		}
	}
	// Creating table for stdout.
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Title", "State", "URL", "Search score"})

	for _, v := range data {
		table.Append(v)
	}
	table.SetRowLine(true)
	total := strconv.Itoa(i.TotalCount)
	table.SetFooter([]string{"", "", "Found total:", total})
	table.Render()
}

func main() {
	flag.Parse()
	search()
}