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
	"students/students"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
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
	if appEnv == "" {

		appEnv = "dev"
		//log.Fatal("APP_ENV is not set")
		fmt.Println("app_env not set")
		fp = filepath
	} else {
		fp = prodfilepath
	}

	fmt.Println(fp)
	creds, err := os.ReadFile(fp)
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error reading file: ", err)
	}

	// yamlContent, err := yaml.JSONToYAML(creds)
	// if err != nil {
	// 	log.Fatal("Error converting JSON to YAML: ", err)
	// }

	fmt.Println("why are u dying and where?")
	config := Conf{}
	err = json.Unmarshal(creds, &config)
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error unmarshalling file: ", err)
	}
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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
			return
		}

		userID, err := ValidateToken(parts[1])
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
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
		// Log the error and return, do not continue executing the function
		fmt.Printf("error parsing token: %v", err)
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
