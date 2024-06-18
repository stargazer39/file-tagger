package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type BrowseFlags struct {
	path    *string
	tagFile string
}

type TagFlags struct {
	file    *string
	tags    *arrayFlags
	tagFile string
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
		path:    browseCMD.String("f", ".", "set the root filepath for browsing"),
		tagFile: ".tag",
	}

	var tags arrayFlags

	tagFlags := TagFlags{
		file:    tagCMD.String("f", ".", "tag file"),
		tags:    &tags,
		tagFile: ".tag",
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

	_, err2 := db.Exec("CREATE TABLE IF NOT EXISTS tags (name text, tag varchar(50))")

	check(err2)

	return db, err
}

func Browse(options BrowseFlags) {
	info, tf, listError, tfError := GetListOfFiles(*options.path, options.tagFile)

	for _, i := range *info {
		log.Print(i.Name())
		if tf {
			db, err := GetDB(filepath.Join(*options.path, options.tagFile))

			check(err)
			// Read the tag base
			rows, err := db.Query("SELECT tag from tags where name = ?", i.Name())

			check(err)

			for rows.Next() {
				tag := ""
				check(rows.Scan(&tag))

				fmt.Printf(" %s", tag)
			}
		}

		println()
	}

	if listError != nil || tfError != nil {
		log.Println(listError, tfError)
	}

}

func Tag(options TagFlags) {
	dir, file := filepath.Split(*options.file)

	dbPath := filepath.Join(dir, options.tagFile)

	log.Println(*options.file, *options.tags)

	db, err := GetDB(dbPath)

	check(err)

	for _, t := range *options.tags {
		_, err := db.Exec("INSERT INTO tags(name, tag) VALUES (?,?)", file, t)
		check(err)
	}

	db.Close()
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
