package isearch

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"testing"
)

func TestSetRepositories(t *testing.T) {
	repos := setRepositories()

	if reflect.TypeOf(repos) != reflect.TypeOf(ConfigurationFile{}) {
		t.Error()
	}
}

func TestPassedArguments(t *testing.T) {
	repo := ""
	issue := ""
	oldArgs := os.Args
	os.Args = []string{"main.go", "-repo=ansible", "-issue=ec2"}

	tmpRepo, tmpIssue := parseArgs()
	actualRepo := *tmpRepo
	actualIssue := *tmpIssue
	if actualRepo != repo {
		t.Errorf("Test failed. Expected %s got %s", repo, actualRepo)
	}
	if actualIssue != issue {
		t.Errorf("Test failed. Expected %s got %s", issue, actualIssue)
	}

	os.Args = oldArgs
}

func TestSearch(t *testing.T) {
	search()
}

func TestPrintResult(t *testing.T) {
	client := &http.Client{}
	repoName := "Rocklviv/isearch"
	issue := "test issue"
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		t.Error()
	}
	q := req.URL.Query()
	q.Add("q", fmt.Sprintf("%s+repo:%s", issue, repoName))
	req.URL.RawQuery = q.Encode()
	resp, _ := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	printResults(body)
}
