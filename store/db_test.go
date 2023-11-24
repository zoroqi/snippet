package store

import (
	"fmt"
	"sort"
	"testing"
)

func TestLoadSnippets(t *testing.T) {
	//db := Db{Path: "testdata/snippets.json"}
	if pets, _, err := loadSnippets("_data"); err != nil {
		t.Fatal(err)
	} else {
		for _, p := range pets {
			t.Logf("%+v\n", p)
		}
	}
}

func Test_Path(t *testing.T) {
	td := [][]string{
		{"../a/b", "../a/b/db.json", "../a/b/c/d.sh", "c/d.sh"},
		{"../a/b", "../a/b/db.json", "../a/b/c/d/e.sh", "c/d/e.sh"},
		{"./a", "./a/b/db.json", "./a/b/c/d/e.sh", "c/d/e.sh"},
		{"~/a", "~/a/db.json", "~/a/b/c/d.sh", "b/c/d.sh"},
		{"a", "a/db.json", "a/b/c/d.sh", "b/c/d.sh"},
		{"", "./db.json", "a/b/c/d.sh", "a/b/c/d.sh"},
		{"/a", "/a/c/db.json,", "/a/b/d.sh", "../b/d.sh"},
	}
	for _, v := range td {
		n, err := snippetRelPath(v[0], v[1], v[2])
		if err != nil {
			t.Fatal(err)
		}
		if n != v[3] {
			t.Fatal(fmt.Sprintf("want %s, got %s", v[3], n))
		}
	}
}

func Test_DbPath(t *testing.T) {
	r := "../a/b"
	dbFiles := []string{
		"../a/b/db.json",
		"../a/d/db.json",
		"../a/b/c/db.json",
		"../a/b/c/d/db.json",
	}
	sort.Slice(dbFiles, func(i, j int) bool {
		if len(dbFiles[i]) != len(dbFiles[j]) {
			return len(dbFiles[i]) < len(dbFiles[j])
		} else {
			return dbFiles[i] < dbFiles[j]
		}
	})
	td := [][]string{
		{"../a/b/c/d/e/f.sh", "../a/b/c/d/db.json"},
		{"../a/b/c.sh", "../a/b/db.json"},
		{"../a/b/c/d.sh", "../a/b/c/db.json"},
		{"../a/d/c/d/e/f/sh", "../a/d/db.json"},
	}
	for _, v := range td {
		n, err := findClosestDbFile(r, dbFiles, v[0])
		if err != nil {
			t.Fatal(err)
		}
		if n != v[1] {
			t.Fatalf("want %s, got %s", v[1], n)
		}
	}
}
