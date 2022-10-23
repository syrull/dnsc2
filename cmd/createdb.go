package cmd

import (
	"database/sql"

	"github.com/erdaltsksn/cui"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

const dbName = "c2.db"

var createDbCmd = &cobra.Command{
	Use:   "createdb",
	Short: "Initialize a database",
	Run: func(cmd *cobra.Command, args []string) {
		sqlScheme := `
		CREATE TABLE client(
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			machine_id VARCHAR NOT NULL,
			remote_ip VARCHAR NOT NULL,
			last_updated INTEGER NOT NULL,
			created_at INTEGER NOT NULL
		);
		`
		db, err := sql.Open("sqlite3", dbName)
		if err != nil {
			cui.Error("Failed to create a database %s\n", err)
		}
		defer db.Close()
		_, err = db.Exec(sqlScheme)
		if err != nil {
			cui.Error("Failed to exec the scheme %s\n", err)
		}

		cui.Success("Successfully created a database!")
	},
}

func init() {
	rootCmd.AddCommand(createDbCmd)
}
