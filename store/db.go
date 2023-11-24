package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type DB struct {
	Snippets []Snippet
	DbFiles  []string
	rootPath string
}

const DB_FILE_NAME = "db.json"

func Load(path string) (*DB, error) {
	snippets, dbFiles, err := loadSnippets(path)
	if err != nil {
		return nil, err
	}
	m := map[string]bool{}
	for _, v := range snippets {
		m[v.configPath] = true
	}
	sort.Slice(dbFiles, func(i, j int) bool {
		if len(dbFiles[i]) != len(dbFiles[j]) {
			return len(dbFiles[i]) < len(dbFiles[j])
		} else {
			return dbFiles[i] < dbFiles[j]
		}
	})
	return &DB{Snippets: snippets,
		DbFiles:  dbFiles,
		rootPath: path}, nil
}

func (db *DB) CreateDbFile(target string) error {
	path := filepath.Join(db.rootPath, target)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}
	file := filepath.Join(path, DB_FILE_NAME)
	if _, err := os.Stat(file); os.IsExist(err) {
		return nil
	}
	return os.WriteFile(file, []byte("[]"), 0644)
}

func (db *DB) Save() error {
	db.distinct()
	for _, dbFile := range db.DbFiles {
		snippets := []Snippet{}
		for _, snippet := range db.Snippets {
			if snippet.configPath == dbFile {
				var err error
				snippet.Path, err = snippetRelPath(db.rootPath, dbFile, snippet.Path)
				if err != nil {
					return err
				}
				snippets = append(snippets, snippet)
			}
		}
		bs, err := json.MarshalIndent(snippets, "", "  ")
		if err != nil {
			return err
		}
		err = os.WriteFile(dbFile, bs, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) Add(snippet Snippet, script []byte) error {
	dbFile, err := findClosestDbFile(db.rootPath, db.DbFiles, snippet.Path)
	if err != nil {
		return err
	}
	filePath := filepath.Dir(snippet.Path)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return err
		}
	}
	if err := os.WriteFile(snippet.Path, script, 0644); err != nil {
		return err
	}
	snippet.configPath = dbFile
	db.Snippets = append(db.Snippets, snippet)
	ShowSnippet(snippet)

	return db.Save()
}

func (db *DB) distinct() error {
	m := map[string]bool{}
	snippets := []Snippet{}
	for _, v := range db.Snippets {
		if m[v.Path] {
			continue
		}
		m[v.Path] = true
		snippets = append(snippets, v)
	}
	db.Snippets = snippets
	return nil
}

func (db *DB) Remove(snippet Snippet) error {
	for i, s := range db.Snippets {
		if s.Path == snippet.Path {
			db.Snippets = append(db.Snippets[:i], db.Snippets[i+1:]...)
			break
		}
	}
	return db.Save()
}

func (db *DB) Find(search Search) []Snippet {
	return findSnippet(search, db.Snippets)
}

func loadSnippets(root string) ([]Snippet, []string, error) {
	dbFiles := []string{}
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return err
		}
		if info.Name() == DB_FILE_NAME {
			dbFiles = append(dbFiles, path)
		}
		return nil
	})
	snippets := []Snippet{}
	for _, dbFile := range dbFiles {
		pets := []Snippet{}
		if bs, err := os.ReadFile(dbFile); err == nil {
			err = json.Unmarshal(bs, &pets)
			if err != nil {
				return nil, nil, err
			}
		} else {
			return nil, nil, err
		}
		f := filepath.Dir(dbFile)
		for i := 0; i < len(pets); i++ {
			pets[i].Path = filepath.Join(f, pets[i].Path)
			pets[i].configPath = dbFile
		}
		snippets = append(snippets, pets...)
	}
	return snippets, dbFiles, nil
}

func (db *DB) FuzzySnippet(keys []string) []Snippet {
	r := []Snippet{}
	dup := map[string]bool{}
	for _, s := range keys {
		for _, snippet := range db.Snippets {
			if dup[snippet.Path] {
				continue
			}
			if !strings.Contains(snippet.ShortName, s) {
				continue
			}
			if !strings.Contains(snippet.Name, s) {
				continue
			}
			if !strings.Contains(snippet.Description, s) {
				continue
			}
			if !strings.Contains(strings.Join(snippet.Aliases, " "), s) {
				continue
			}
			dup[snippet.Path] = true
			r = append(r, snippet)
		}
	}
	return r
}

// 寻找最近的db文件
// 多个 db 根据长短进行排序后, 从短到长进行匹配
// 通过计算和 db 文件的相对路径, 相对路径最短的就是目标路径
func snippetRelPath(root, dbfile, snipfile string) (string, error) {
	dbRel, err := filepath.Rel(root, dbfile)
	if err != nil {
		return "", err
	}
	snipRel, err := filepath.Rel(root, snipfile)
	if err != nil {
		return "", err
	}
	dbDir := filepath.Dir(dbRel)
	n, err := filepath.Rel(dbDir, snipRel)
	if err != nil {
		return "", err
	}
	return n, nil
}

// 寻找最近的db文件
// 多个 db 根据长短进行排序后, 从短到长进行匹配
// 通过计算和 db 文件的相对路径, 相对路径最短的就是目标路径
// 无法返回正确路径或者出现 `../` 情况都是有问题的
func findClosestDbFile(root string, dbFiles []string, target string) (string, error) {
	closestDb := ""
	latest := len(dbFiles[len(dbFiles)-1])
	errCount := 0
	for _, db := range dbFiles {
		nn, err := snippetRelPath(root, db, target)
		if err != nil {
			errCount++
			continue
		}
		if strings.HasPrefix(nn, "../") {
			errCount++
			continue
		}
		if len(nn) < latest {
			latest = len(nn)
			closestDb = db
		}
	}
	if errCount == len(dbFiles) {
		return "", errors.New("no db file found")
	}
	if closestDb == "" {
		return dbFiles[0], nil
	}
	return closestDb, nil
}
