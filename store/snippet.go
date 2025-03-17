package store

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"strings"
)

type Language = string

const (
	Other  = Language("other")
	ANKO   = Language("anko")
	PY2    = Language("py")
	PY3    = Language("py")
	GO     = Language("golang")
	SH     = Language("sh")
	Prompt = Language("prompt")
)

type Snippet struct {
	ShortName   string   `json:"short_name"`
	Path        string   `json:"path"`
	Name        string   `json:"name"`
	Aliases     []string `json:"aliases"`
	Tags        []string `json:"tags"`
	Language    Language `json:"language"`
	Description string   `json:"description"`
	CanExec     bool     `json:"can_exec"`
	configPath  string
}

type Search struct {
	Snippet
	Aliases string
	Fuzzy   string
}

func (s *Search) isEmpty() bool {
	return s.ShortName == "" && s.Name == "" && s.Aliases == "" && s.Fuzzy == "" && len(s.Tags) == 0
}

func SnippetPrintTable(snippets []Snippet) {
	header := []string{fmt.Sprintf("short_name(%d)", len(snippets)), "name", "aliases", "can_exec", "description", "tags", "path"}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetRowLine(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_CENTER)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	var TableHeaderColor = make([]tablewriter.Colors, len(header))
	for i := range TableHeaderColor {
		TableHeaderColor[i] = tablewriter.Colors{tablewriter.Bold, tablewriter.FgGreenColor}
	}
	table.SetHeaderColor(TableHeaderColor...)
	for _, s := range snippets {
		table.Append([]string{
			s.ShortName,
			s.Name,
			strings.Join(s.Aliases, "\n"),
			fmt.Sprintf("%t", s.CanExec),
			s.Description,
			strings.Join(s.Tags, "\n"),
			s.Path,
		})
	}
	table.Render()
}

func ShowSnippet(snippet Snippet) (string, error) {
	sb := strings.Builder{}
	write := func(key, value string) {
		sb.WriteString(key)
		sb.WriteString(": ")
		sb.WriteString(value)
		sb.WriteString("\n")
	}
	write("Name", snippet.Name)
	write("ShortName", snippet.ShortName)
	write("Aliases", strings.Join(snippet.Aliases, ","))
	write("Path", snippet.Path)
	write("Language", snippet.Language)
	write("CanExec", fmt.Sprintf("%v", snippet.CanExec))

	sb.WriteString("\nscript:\n\n")
	script, err := ReadScript(snippet)
	if err != nil {
		return sb.String(), err
	}
	sb.WriteString(script)
	return sb.String(), nil
}

type findHandler func(string) bool

func findSnippet(search Search, snippets []Snippet, keywords []string) []Snippet {
	if search.isEmpty() && len(keywords) == 0 {
		return snippets
	}
	short := match(search.ShortName)
	name := match(search.Name)
	alias := match(search.Aliases)
	tags := anyMatch(search.Tags)
	fuzzy := contain(search.Fuzzy)
	repeating := map[string]bool{}
	result := []Snippet{}
	addResult := func(snippet Snippet) {
		if repeating[snippet.Path] {
			return
		}
		repeating[snippet.Path] = true
		result = append(result, snippet)
	}
Outer:
	for _, s := range snippets {
		if short(s.ShortName) || name(s.Name) {
			addResult(s)
		}
		for _, a := range s.Aliases {
			if alias(a) {
				addResult(s)
				continue Outer
			}
		}
		for _, t := range s.Tags {
			if tags(t) {
				addResult(s)
				continue Outer
			}
		}
		if fuzzy(s.Description) {
			addResult(s)
		}
	}
	for _, s := range snippets {
		flag := true
		if len(keywords) == 0 {
			continue
		}
		for _, t := range keywords {
			f := strings.Contains(s.ShortName, t)
			f = f || strings.Contains(s.Name, t)
			for _, a := range s.Aliases {
				f = f || strings.Contains(a, t)
			}
			for _, tg := range s.Tags {
				f = f || strings.Contains(tg, t)
			}
			flag = flag && f
		}
		if flag {
			addResult(s)
		}
	}
	return result
}

func contain(key string) findHandler {
	return func(s string) bool {
		if key == "" {
			return false
		}
		if strings.Contains(s, key) {
			return true
		}
		if strings.Contains(strings.ToLower(s), strings.ToLower(key)) {
			return true
		}
		return false
	}
}

func match(key string) findHandler {
	return func(s string) bool {
		return key != "" && s == key
	}
}

func anyMatch(keys []string) findHandler {
	ms := []findHandler{}
	for _, k := range keys {
		ms = append(ms, match(k))
	}
	return func(s string) bool {
		for _, v := range ms {
			if v(s) {
				return true
			}
		}
		return false
	}
}

func ReadScript(snip Snippet) (string, error) {
	script, err := os.ReadFile(snip.Path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(script)), nil
}
