package register

import (
	"database/sql"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/xuanmingyi/wingo/errno"
)

type SQliteDriver struct {
	DB    *sql.DB
	Valid bool
}

// type:
// 0  int
// 1  float
// 2  string
func (d *SQliteDriver) init() {
	create_table_sqls := []string{
		`CREATE TABLE IF NOT EXISTS "register" (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"key" VARCHAR(128) UNIQUE,
		"type" int(4),
		"value" VARCHAR(128));`,
	}

	for _, sql := range create_table_sqls {
		_, err := d.DB.Exec(sql)
		if err != nil {
			return
		}
	}

	d.Valid = true
}

func (d *SQliteDriver) IsValid() bool {
	return d.Valid
}

func (d *SQliteDriver) Open(path string) (err error) {
	// path: sqlite://my.db
	d.DB, err = sql.Open("sqlite3", path[9:])
	if err != nil {
		return err
	}
	d.init()
	return nil
}

func (d *SQliteDriver) Close() {
	if d.IsValid() {
		d.DB.Close()
	}
}

func (d *SQliteDriver) Create(key string, value Record) (err error) {
	if !d.IsValid() {
		return errno.RegisterDriverNotValid
	}
	var type_flag int
	stmt, err := d.DB.Prepare("INSERT INTO `register`(key, type, value) values (?, ?, ?)")
	if err != nil {
		return err
	}
	switch value.(type) {
	case int, int64:
		type_flag = 0
	case float32, float64:
		type_flag = 1
	case string:
		type_flag = 2
	}
	_, err = stmt.Exec(key, type_flag, value)
	if err != nil {
		return err
	}
	return nil
}

func (d *SQliteDriver) Delete(key string) (err error) {
	if !d.IsValid() {
		return errno.RegisterDriverNotValid
	}

	stmt, err := d.DB.Prepare("DELETE FROM `register_string` WHERE `key` = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(key)
	if err != nil {
		return err
	}
	return nil
}

func (d *SQliteDriver) Update(key string, value Record) (err error) {
	if !d.IsValid() {
		return errno.RegisterDriverNotValid
	}
	stmt, err := d.DB.Prepare("UPDATE `register` SET `value` = ? WHERE `key` = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(value, key)
	if err != nil {
		return err
	}
	return nil
}

func (d *SQliteDriver) Search(key string) (records map[string]Record, err error) {
	if !d.IsValid() {
		return nil, errno.RegisterDriverNotValid
	}

	records = make(map[string]Record)
	rows, err := d.DB.Query("SELECT `key`, `type`, `value` FROM `register`")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var key string
		var type_flag int
		var record Record
		err = rows.Scan(&key, &type_flag, &record)
		if err != nil {
			return nil, err
		}

		switch type_flag {
		case 0:
			record, _ = strconv.ParseInt(record.(string), 10, 64)
		case 1:
			record, _ = strconv.ParseFloat(record.(string), 64)
		case 2:
			record = record.(string)
		}

		records[key] = record
	}
	return records, nil
}
