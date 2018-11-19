package main

import (
	"reflect"
	"testing"
)

func TestSetRepositories(t *testing.T) {
	repos := setRepositories()

	if reflect.TypeOf(repos) != reflect.TypeOf([]repositories{}) {
		t.Error()
	}
}

func TestSearch(t *testing.T) {

}
