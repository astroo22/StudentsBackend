package sqlgeneric

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

const prodfilepath = "/var/www/backend/config/postgrescreds.prod.json"
const filepath = "config/postgrescreds.dev.yml"

// INFO: Looks like aws saves these yml values as json which is in conflict with setup
// to avoid a full setup change right now which is not really a bad thing or hard to do
// I'm going to just convert it for now to avoid conflicts with other things.
// TODO: Update this so that it just switches based on env or create the json files here
// either one doesn't matter just do it later.
func getYMLcreds() Config {
	fp := ""
	config := Config{}
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "prod" {

		//log.Fatal("APP_ENV is not set")
		fp = prodfilepath
	} else {
		fp = filepath
	}
	creds, err := os.ReadFile(fp)
	if err != nil {
		fmt.Println("error reading yml secret")
		log.Fatal("Error reading file: ", err)
	}

	if appEnv == "prod" {
		err = json.Unmarshal(creds, &config)
		if err != nil {
			log.Println("in yml unmarshal file might not exist maybe?")
			log.Fatal("Error unmarshalling file: ", err)
		}
	} else {
		err = yaml.Unmarshal(creds, &config)
		if err != nil {
			log.Println("in yml unmarshal file might not exist maybe?")
			log.Fatal("Error unmarshalling file: ", err)
		}
	}
	return config
}
func fmtPsqlConn(data Config) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", data.Host, data.Port, data.Username, data.Password, data.Database)

}
func Init() (*sql.DB, error) {
	return sql.Open("postgres", fmtPsqlConn(getYMLcreds()))
}
