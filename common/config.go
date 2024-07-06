package common

import (
	"os"
)

type Configuration struct {
	AppName            string
	AppDomain          string
	AppEmail           string
	TokenSecret        string
	RefreshTokenSecret string
	//database
	DBHost     string
	DBPort     string
	DBUsername string
	DBPassword string
	//mail server
	SMTPHost     *string
	SMTPPort     *string
	SMTPUsername *string
	SMTPPassword *string
	//mailjet
	MailjetPublicKey  *string
	MailjetPrivateKey *string
}

var AppConfig *Configuration

func InitConfig() {
	// string: key, bool: required
	variables := map[string]bool{
		"APP_NAME":             true,
		"APP_DOMAIN":           true,
		"APP_EMAIL":            true,
		"TOKEN_SECRET":         true,
		"REFRESH_TOKEN_SECRET": true,
		"DB_HOST":              true,
		"DB_PORT":              true,
		"DB_USERNAME":          true,
		"DB_PASSWORD":          true,
		"SMTP_HOST":            false,
		"SMTP_PORT":            false,
		"SMTP_USERNAME":        false,
		"SMTP_PASSWORD":        false,
		"MAILJET_PUBLIC_KEY":   false,
		"MAILJET_PRIVATE_KEY":  false,
	}

	tmp := map[string]*string{}
	for k, v := range variables {
		val, found := os.LookupEnv(k)
		if !found && v {
			panic(k + " not set")
		}
		tmp[k] = &val
	}

	AppConfig = &Configuration{
		AppName:            *tmp["APP_NAME"],
		AppDomain:          *tmp["APP_DOMAIN"],
		AppEmail:           *tmp["APP_EMAIL"],
		TokenSecret:        *tmp["TOKEN_SECRET"],
		RefreshTokenSecret: *tmp["REFRESH_TOKEN_SECRET"],
		DBHost:             *tmp["DB_HOST"],
		DBPort:             *tmp["DB_PORT"],
		DBUsername:         *tmp["DB_USERNAME"],
		DBPassword:         *tmp["DB_PASSWORD"],
		SMTPHost:           tmp["SMTP_HOST"],
		SMTPPort:           tmp["SMTP_PORT"],
		SMTPUsername:       tmp["SMTP_USERNAME"],
		SMTPPassword:       tmp["SMTP_PASSWORD"],
		MailjetPublicKey:   tmp["MAILJET_PUBLIC_KEY"],
		MailjetPrivateKey:  tmp["MAILJET_PRIVATE_KEY"],
	}
}
