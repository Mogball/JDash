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
	configData, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Unable to read client config file: %v", err)
	}
	config, err := google.ConfigFromJSON(configData, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client conig: %v", err)
	}
	return getClient(ctx, config)
}

func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Printf("Cannot get credential file path: %v", err)
		return config.Client(ctx, getTokenFromWeb(config))
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authUrl := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Manual authorization required:\n%v\n", authUrl)
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}
	tok, err := config.Exchange(context.TODO(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token: %v", err)
	}
	return tok
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
