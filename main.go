package main

import (
	"database/sql"
	"file-tagger/tagger"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type BrowseFlags struct {
	Path    *string
	TagFile string
}

type TagFlags struct {
	File    *string
	Tags    *arrayFlags
	TagFile string
	Desc    *string
}

type FileInfo struct {
	Tags []string
	Desc string
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join(*i, ",")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	browseCMD := flag.NewFlagSet("browse", flag.ExitOnError)
	tagCMD := flag.NewFlagSet("tag", flag.ExitOnError)

	browseFlags := BrowseFlags{
		Path:    browseCMD.String("f", ".", "set the root filepath for browsing"),
		TagFile: ".tag",
	}

	var tags arrayFlags

	tagFlags := TagFlags{
		File:    tagCMD.String("f", ".", "tag file"),
		Tags:    &tags,
		TagFile: ".tag",
		Desc:    tagCMD.String("d", "", "Description"),
	}

	tagCMD.Var(&tags, "t", "Set tags")

	switch os.Args[1] {
	case "browse":
		browseCMD.Parse(os.Args[2:])
		Browse(browseFlags)
		fmt.Println("browse mode")
	case "tag":
		tagCMD.Parse(os.Args[2:])
		Tag(tagFlags)
		fmt.Println("tag mode")
	default:
		log.Fatalf("[ERROR] unknown subcommand '%s', see help for more details.", os.Args[1])
	}
}

func GetDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)

	check(err)

	_, err4 := db.Exec("VACUUM")
	check(err4)

	_, err2 := db.Exec("CREATE TABLE IF NOT EXISTS tags (name text, tag varchar(50))")

	check(err2)

	_, err3 := db.Exec("CREATE TABLE IF NOT EXISTS desc (name text PRIMARY KEY, desc text)")

	check(err3)

	return db, err
}

func Browse(options BrowseFlags) {
	t := tagger.NewTagger()
	log.Println(t.ListFiles(*options.Path))
}

func Tag(options TagFlags) {
	t := tagger.NewTagger()
	log.Println(t.TagFile(*options.File, *options.Tags))
}

func GetListOfFiles(root string, tagFile string) (*[]os.FileInfo, bool, error, error) {
	tf := filepath.Join(root, tagFile)

	info, err := os.Stat(tf)

	infoList := []os.FileInfo{}

	err2 := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.Name() == root || info.Name() == tagFile {
			return nil
		}

		if info.IsDir() || strings.HasPrefix(".", info.Name()) {
			return filepath.SkipDir
		}

		infoList = append(infoList, info)

		return nil
	})

	return &infoList, err == nil && !info.IsDir(), err2, err
}

func GetTagData(name string) {

}

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}
