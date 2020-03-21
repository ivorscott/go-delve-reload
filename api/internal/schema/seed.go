package schema

import (
	"fmt"
	"io/ioutil"

	"github.com/jmoiron/sqlx"
)

const folder = "/seeds/"
const ext = ".sql"

func Seed(db *sqlx.DB, filename string) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	src := fmt.Sprintf("%s%s%s%s", RootDir(), folder, filename, ext)
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
