package client

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
)

// Create creates a new MariaDB chart release
func Create(id string) error {
	db_host := os.Getenv("MARIADB_HOST")
	db_port := os.Getenv("MARIADB_PORT")
	db_user := os.Getenv("MARIADB_USER")
	db_pass := os.Getenv("MARIADB_PASS")

	db, err := sql.Open("mysql", db_user+":"+db_pass+"@tcp("+db_host+":"+db_port+")/")
	if err != nil {
		return err
	}
	defer db.Close()

	database := HashedValue("db" + id)
	username := HashedValue("user" + id)
	password := HashedValue("pass" + id)
	glog.Infof("Create database: %s\n", database)
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS `" + database + "`")
	if err != nil {
		panic(err)
	}

	glog.Infof("Create user: %s\n", username)
	_, err = db.Exec("GRANT ALL ON `" + database + "`.* TO `" + username + "`@'%' IDENTIFIED BY '" + password + "'")
	if err != nil {
		panic(err)
	}

	return nil
}

// Delete deletes a MariaDB chart release
func Delete(id string) error {
	db_host := os.Getenv("MARIADB_HOST")
	db_port := os.Getenv("MARIADB_PORT")
	db_user := os.Getenv("MARIADB_USER")
	db_pass := os.Getenv("MARIADB_PASS")

	db, err := sql.Open("mysql", db_user+":"+db_pass+"@tcp("+db_host+":"+db_port+")/")
	if err != nil {
		return err
	}
	defer db.Close()

	database := HashedValue("db" + id)
	username := HashedValue("user" + id)
	glog.Infof("Drop database: %s\n", database)
	_, err = db.Exec("DROP DATABASE IF EXISTS `" + database + "`")
	if err != nil {
		return err
	}

	glog.Infof("Drop user: %s\n", username)
	_, err = db.Exec("DROP USER IF EXISTS `" + username + "`")
	if err != nil {
		return err
	}

	return nil
}

func HashedValue(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))[0:8]
}
