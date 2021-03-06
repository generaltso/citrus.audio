// replace actual implementation with whatever
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func checkErr(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

func fileExists(filename string) bool {
	f, err := os.Open(filename)
	f.Close()
	if os.IsNotExist(err) {
		return false
	}
	checkErr(err)
	return true
}

type datastore struct {
	path string
}

func (db *datastore) init() {
	p, err := filepath.Abs(".db")
	checkErr(err)
	db.path = p
}

func (db *datastore) get(key string) Suprême {
	var val Suprême

	fileName := fmt.Sprintf("%s%c%s.json", db.path, os.PathSeparator, key)
	if fileExists(fileName) {
		contents, err := ioutil.ReadFile(fileName)
		checkErr(err)
		checkErr(json.Unmarshal(contents, &val))
	}
	return val
}

func (db *datastore) set(key string, val Suprême) {
	contents, err := json.Marshal(val)
	checkErr(err)

	fileName := fmt.Sprintf("%s%c%s.json", db.path, os.PathSeparator, key)
	ioutil.WriteFile(fileName, contents, 0666)
}

// temp?

type lmod struct {
	Name    string
	ModTime time.Time
}

type lmodSlice []lmod

func (l lmodSlice) Len() int {
	return len(l)
}
func (l lmodSlice) Less(i, j int) bool {
	return l[i].ModTime.Before(l[j].ModTime)
}
func (l lmodSlice) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (db *datastore) getRange(start, limit int) []string {
	dir, err := os.Open(db.path)
	checkErr(err)
	files, err := dir.Readdir(-1)
	checkErr(err)
	ls := lmodSlice{}
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		ls = append(ls, lmod{strings.TrimRight(file.Name(), ".json"), file.ModTime()})
	}
	dir.Close()
	sort.Sort(sort.Reverse(ls))
	len := min(start+limit, len(ls))
	size := len - start
	if size <= 0 {
		return []string{}
	}
	ret := make([]string, size)
	j := 0
	for i := start; i < len; i++ {
		l := ls[i]
		ret[size-j-1] = l.Name
		j++
	}
	return ret
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
