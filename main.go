package main

import (
	"fmt"
	"flag"
	"net/http"
	"encoding/json"
	"io/ioutil"
	//"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter"
	"os"
)

type repositories struct {
	name string
	repo string
}

type Issues struct {
	TotalCount int `json:"total_count"`
	IncompleteResults bool `json:"incomplete_results"`
	Items Items `json:"items"`
}

type Items []struct{
	Url string `json:"url"`
	Repository_url string `json:"repository_url"`
	Labels_url	string `json:"labels_url"`
	Comments_url string `json:"comments_url"`
	Events_url string `json:"events_url"`
	Html_url string `json:"html_url"`
	Id int `json:"id"`
	Number int `json:"number"`
	Title string `json:"title"`
	User User `json:"user"`
	Labels Labels `json:"labels"`
	State string `json:"state"`
	Locked bool `json:"locked"`
	Assignee string `json:"assignee"`
	Assignees Assignees `json:"assignees"`
	Milestone string `json:"milestone"`
	Comments int `json:"comments"`
	Created_at string `json:"create_at"`
	Updated_at string `json:"update_at"`
	Closed_at string `json:"closed_at"`
	Author_association string `json:"author_association"`
	Pull_requests Pull_requests `json:"pull_requests"`
	Body string `json:"body"`
	Score float32 `json:"score"`
}

type User struct {
	Login string `json:"login"`
	Id int `json:"id"`
	Avatar_url string `json:"avatar_url"`
	Gravatar_id string `json:"gravatar_id"`
	Url string `json:"url"`
	Html_url string `json:"html_url"`
	Followers_url string `json:"followers_url"`
	Following_url string `json:"following_url"`
	Gists_url string `json:"gists_url"`
	Starred_url string `json:"starred_url"`
	Subscribtions_url string `json:"subscribtions_url"`
	Organizations_url string `json:"organizations_url"`
	Repos_url string `json:"repos_url"`
	Events_url string `json:"events_url"`
	Received_events_url string `json:"received_events_url"`
	Type string `json:"type"`
	Site_admin bool `json:"site_admin"`
}

type Labels []struct {
	Id int `json:"id"`
	Url string `json:"url"`
	Name string `json:"name"`
	Color string `json:"color"`
	Default bool `json:"default"`
}

type Assignees []struct {}
type Pull_requests []struct {
	Url string `json:"url"`
	Html_url string `json:"html_url"`
	Diff_url string `json:"diff_url"`
	Patch_url string `json:"patch_url"`
}

var repo = flag.String("repo", "", "Name of project to search.")
var issue = flag.String("issue", "", "Specify issue that you are looking for.")

func _set_repositories() []repositories {
	v := []repositories{
		{
			name: "terraform",
			repo: "hashicorp/terraform",
		},
		{
			name: "ansible",
			repo: "ansible/ansible",
		},
	}
	return v
}

func _search() {
	r := *repo
	i := *issue
	repos := _set_repositories()

	for key, value := range repos {
		if repos[key].name == r {
			client := &http.Client{}
			repo_name := value
			req, err := http.NewRequest("GET", "https://api.github.com/search/issues", nil)
			if err != nil {
				fmt.Println("[ERROR] Some issue with request.")
			}
			q := req.URL.Query()
			q.Add("q", fmt.Sprintf("%s+repo:%s", i, repo_name.repo))
			req.URL.RawQuery = q.Encode()
			resp, _ := client.Do(req)
			body, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			data := print_results(body)
			render_table(data)
		}
	}
}

func print_results(body []byte) [][]string {
	var i Issues
	err := json.Unmarshal(body, &i)
	if err != nil {
		fmt.Println("[ERROR] Cannot unmarshal data. %s", err)
	}
	for _, value := range i.Items {
		data := [][]string{
			[]string{value.Title, value.State, value.Html_url},
		}
		return data
	}
}

func main() {
	flag.Parse()
	_search()
}