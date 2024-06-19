package tagger

import (
	"database/sql"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type MetadataDB struct {
	db      *sql.DB
	tagFile string
	root    string
}

func NewMetadataDB(root string, tagFile string) *MetadataDB {
	return &MetadataDB{
		tagFile: tagFile,
		root:    root,
	}
}

// var DatabaseErr error = fmt.Errorf("db err")

func (md *MetadataDB) GetTagsForFile(name string, fresh bool) ([]string, error) {
	db, err := md.getDb(fresh)

	if err != nil {
		return []string{}, err
	}

	rows, err := db.Query("SELECT tag from tags where name = ?", name)

	if err != nil {
		return nil, err
	}

	tags := []string{}

	for rows.Next() {
		tag := ""

		if err := rows.Scan(&tag); err != nil {
			return tags, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func (md *MetadataDB) GetDescriptionForFile(name string, fresh bool) (string, error) {
	db, err := md.getDb(fresh)

	if err != nil {
		return "", err
	}

	row := db.QueryRow("SELECT desc from desc where name = ?", name)
	desc := ""

	if row.Err() == nil {
		if err := row.Scan(&desc); err != nil {
			return "", err
		}
	}

	return desc, nil
}

func (md *MetadataDB) TagFile(name string, tags []string, fresh bool) error {
	db, err := md.getDb(fresh)

	if err != nil {
		return err
	}

	for _, t := range tags {
		_, err := db.Exec("INSERT INTO tags(name, tag) VALUES (?,?)", name, t)

		if err != nil {
			return err
		}
	}

	return nil
}

func (md *MetadataDB) SetDescriptionForFile(name string, desc string, fresh bool) error {
	db, err := md.getDb(fresh)

	if err != nil {
		return err
	}

	if desc != "" {
		_, err := db.Exec("INSERT INTO desc(name, desc) VALUES (?, ?) ON CONFLICT(name) DO UPDATE SET desc=excluded.desc", name, desc)

		if err != nil {
			return err
		}
	}

	return nil
}

func (md *MetadataDB) Close() error {
	if md.db != nil {
		return md.db.Close()
	}

	return nil
}

func (md *MetadataDB) getDb(fresh bool) (*sql.DB, error) {
	if md.db != nil && !fresh {
		return md.db, nil
	}

	if md.db != nil {
		md.db.Close()
	}

	dbPath := filepath.Join(md.root, md.tagFile)

	var err error

	db, err := sql.Open("sqlite3", dbPath)

	if err != nil {
		return nil, err
	}

	_, err = db.Exec("VACUUM")

	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS tags (name text, tag varchar(50))")

	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS desc (name text PRIMARY KEY, desc text)")

	if err != nil {
		return nil, err
	}

	md.db = db

	return db, err
}
