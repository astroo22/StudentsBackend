package sqlgeneric

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

func getYMLcreds() Config {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {

		appEnv = "dev"
		//log.Fatal("APP_ENV is not set")
	}
	filepath := fmt.Sprintf("config/postgrescreds.%s.yml", appEnv)
	fmt.Println(filepath)
	creds, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal("Error reading file: ", err)
	}
	config := Config{}
	err = yaml.Unmarshal(creds, &config)
	if err != nil {
		log.Println("in yml unmarshal file might not exist maybe?")
		log.Fatal("Error unmarshalling file: ", err)
	}
	return config
}
func fmtPsqlConn(data Config) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", data.Host, data.Port, data.Username, data.Password, data.Database)

}
func Init() (*sql.DB, error) {
	return sql.Open("postgres", fmtPsqlConn(getYMLcreds()))
}
