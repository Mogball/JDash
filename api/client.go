package api

import (
	"golang.org/x/oauth2"
	"golang.org/x/net/context"
	"net/http"
	"log"
	"os"
	"encoding/json"
	"os/user"
	"path/filepath"
	"net/url"
	"fmt"
	"jdash/config"
	"jdash/app"
	"io/ioutil"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

func CreateGmailClient() *http.Client {
	configFile := app.Config().Word[config.CLIENT_CONFIG_FILE]
	ctx := context.Background()
	return createAuthorizedClient(ctx, configFile)
}

func createAuthorizedClient(ctx context.Context, configFile string) *http.Client {
	clientConfig, err := getConfig(configFile)
	if err != nil {
		log.Fatalf("Unable to parse client conig: %v", err)
	}
	return getClient(ctx, clientConfig)
}

func GetConfig() (*oauth2.Config, error) {
	configFile := app.Config().Word[config.CLIENT_CONFIG_FILE]
	return getConfig(configFile)
}

func getConfig(configFile string) (*oauth2.Config, error) {
	configData, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	conf, err := google.ConfigFromJSON(configData, gmail.GmailReadonlyScope)
	if err != nil {
		return nil, err
	}
	conf.RedirectURL = app.Config().Word[config.OAUTH_REDIRECT]
	return conf, err
}

func GetCacheToken(config *oauth2.Config) *oauth2.Token {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Printf("Cannot get credential file path: %v", err)
		tok, err := GetTokenFromWeb(config)
		if err != nil {
			log.Fatal(err)
		}
		return tok
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil || !tok.Valid() {
		tok, err = GetTokenFromWeb(config)
		if err != nil {
			log.Fatal(err)
		}
		saveToken(cacheFile, tok)
	}
	return tok
}

func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	return config.Client(ctx, GetCacheToken(config))
}

func GetTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authUrl := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Manual authorization required:\n%v\n", authUrl)
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		return nil, err
	}
	tok, err := config.Exchange(context.TODO(), code)
	if err != nil {
		return nil, err
	}
	return tok, nil
}

func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape(app.Config().Word[config.GMAIL_TOKEN_FILE])), err
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

func saveToken(file string, token *oauth2.Token) {
	log.Printf("Saving token to file: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Printf("Unable to save token: %v", err)
		return
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
