package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
			CREATE TABLE IF NOT EXISTS scheduler (
				id INTEGER PRIMARY KEY AUTOINCREMENT, 
				date CHAR(8) NOT NULL DEFAULT "",
				title VARCHAR(128) NOT NULL,
				comment TEXT,
				repeat VARCHAR(128)
			);
			
			CREATE INDEX idx_scheduler_date ON scheduler (date);
		`

// Инициализация базы данных
// Принимает на вход указатель на логгер, имя БД и имя драйвера БД
// Возвращает:
//   - указатель на БД (*sql.DB)
//   - ошибку (error)
func Init(dbFile string, dbDriver string) (*sql.DB, error) {

	_, err := os.Stat(dbFile)
	var install bool
	if err != nil {
		install = true
	}

	db, err := sql.Open(dbDriver, dbFile)
	if err != nil {
		return nil, err
	}

	if install {
		if _, err = db.Exec(schema); err != nil {
			return nil, err
		}
	}

	return db, nil
}
