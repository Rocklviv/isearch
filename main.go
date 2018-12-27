package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
)

type ConfigurationFile []struct {
	Name string `json:"name"`
	Repo string `json:"repo"`
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

var (
	configFile = ".isearch.json"
	apiUrl = "https://api.github.com/search/issues"
	repo string
	issue string
)

//
// Sets repository struct from configuration file or use default.
//
func setRepositories() ConfigurationFile {
	usr, _ := user.Current()
	dir := usr.HomeDir
	path := filepath.Join(dir, configFile)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("[WARNING] File %s does not exists. Going to use default repo configuration. \r\n", path)
		repos := ConfigurationFile{
			{
				Name: "terraform",
				Repo: "hashicorp/terraform",
			},
			{
				Name: "ansible",
				Repo: "ansible/ansible",
			},
		}
		return repos
	} else {
		jsonFile, err := os.Open(path)
		if err != nil {
			fmt.Printf("[ERROR] Cannot open %s. Check is file exists and have proper permissions. \r\n", path)
			fmt.Printf("[ERROR] %s", err.Error())
			os.Exit(1)
		}
		defer jsonFile.Close()
		file, _ := ioutil.ReadAll(jsonFile)

		var repos ConfigurationFile
		errRead := json.Unmarshal(file, &repos)
		if errRead != nil {
			fmt.Printf("[ERROR] %s \r\n", errRead.Error())
			os.Exit(2)
		}
		return repos
	}
}

//
// Method that implements search issues in repository
//
func search() {
	repos := setRepositories()

	for key, value := range repos {
		if repos[key].Name == repo {
			client := &http.Client{}
			repoName := value
			req, err := http.NewRequest("GET", apiUrl, nil)
			if err != nil {
				fmt.Println("[ERROR] Some issue with request.")
			}

			q := req.URL.Query()
			q.Add("q", fmt.Sprintf("%s+repo:%s", issue, repoName.Repo))
			req.URL.RawQuery = q.Encode()
			resp, _ := client.Do(req)
			body, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()

			printResults(body)
		}
	}
}

//
// Method prints results into table
//
func printResults(body []byte) {
	var i Issues
	err := json.Unmarshal(body, &i)
	if err != nil {
		fmt.Println("[ERROR] Cannot unmarshal data", err)
	}
	data := [][]string{}
	for _, value := range i.Items {
		if value.Score >= 0.500 {
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
	table.SetFooter([]string{"", "", "Found issues:", total})
	table.Render()
}

func parseArgs() (*string, *string) {
	flagRepo := flag.String("repo", "", "Name of project to search.")
	flagIssue := flag.String("issue", "", "Specify issue that you are looking for.")
	return flagRepo, flagIssue
}

func main() {
	tmpRepo, tmpIssue := parseArgs()
	repo = *tmpRepo
	issue = *tmpIssue
	search()
}
