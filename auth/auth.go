package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"students/logger"
	"students/students"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v2"
)

const prodfilepath = "/var/www/backend/config/secrets.json"
const filepath = "config/secrets.yml"

type Conf struct {
	SecretKey string `json:"secretKey"`
}

var secretKey []byte

// TODO: UPDATE THis to work with aws secrets manager json. READ sqlgenerics for more
// info as this is the same issue and solution
func getYMLsecrets() (Conf, error) {
	// switcher for env could probably make this nicer later
	fp := ""
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "prod" {

		//log.Fatal("APP_ENV is not set")
		fp = prodfilepath
	} else {
		fp = filepath
	}

	creds, err := os.ReadFile(fp)
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error reading file: ", err)
		return Conf{}, err
	}

	// yamlContent, err := yaml.JSONToYAML(creds)
	// if err != nil {
	// 	log.Fatal("Error converting JSON to YAML: ", err)
	// }

	config := Conf{}
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
	fmt.Println("got dem secrets")
	return config, err
}
func LoadSecretKey() error {
	conf, err := getYMLsecrets()
	if err != nil {
		fmt.Println("failed secretkey")
		return err
	}
	secretKey = []byte(conf.SecretKey)
	return nil
}

func AuthRequired(next http.Handler) http.Handler {
	// errors make less since the closer it is to done

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.Log.Warning("Authorization header missing")
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Log.Warning("Invalid Authorization header: %v", authHeader)
			http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
			return
		}

		userID, err := ValidateToken(parts[1])
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		logger.Log.Info("Token validated for userID: %v", userID)
		// Store the userID for later use
		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// using bycrypt and jwt here
func AuthenticateUser(userName string, password string) (bool, students.User, error) {
	user, err := students.GetUserByUserName(userName)
	if err != nil {
		return false, students.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return false, students.User{}, err
	}
	user.HashedPassword = ""
	return true, user, nil
}

// ValidateToken validates the JWT token
func ValidateToken(t string) (string, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})
	if err != nil {
		logger.Log.Warning("error parsing token: %v", err)
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["user_id"].(string)
		return userID, nil
	} else {
		return "", err
	}
}

// generate token for future auth
func GenerateToken(ownerID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = ownerID
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	t, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return t, nil

}
