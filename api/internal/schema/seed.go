package schema

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jmoiron/sqlx"
)

const folder = "/internal/schema/seeds/"
const ext = ".sql"

func Seed(db *sqlx.DB, filename string) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	path, err := os.Getwd()
	if err != nil {
		return err
	}

	src := fmt.Sprintf("%s%s%s%s", path, folder, filename, ext)
	dat, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	if _, err := tx.Exec(string(dat)); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	return tx.Commit()
}
