package gmail

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
	"jdash/app"
	"io/ioutil"
	"google.golang.org/api/gmail/v1"
	"golang.org/x/oauth2/google"
)

func CreateClient() (*gmail.Service, error) {
	ctx := app.Context
	conf, err := ioutil.ReadFile("client_config.json")
	if err != nil {
		log.Fatalf("Failed to read gmail config file: %v", err)
	}
	config, err := google.ConfigFromJSON(conf, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client config file: %v", err)
	}
	client := getClient(ctx, config)
	return gmail.New(client)
}

func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	var tok *oauth2.Token
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Printf("Unable to get path to cached file: %v", err)
		tok = getTokenFromWeb(config)
	} else {
		tok, err := tokenFromFile(cacheFile)
		if err != nil {
			tok = getTokenFromWeb(config)
			saveToken(cacheFile, tok)
		}
	}
	return config.Client(ctx, tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authUrl := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	log.Printf("Manual authorization required:\n%v", authUrl)
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}
	tok, err := config.Exchange(context.TODO(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
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
	return filepath.Join(tokenCacheDir, url.QueryEscape("jdash-gmail-token.json")), nil
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
	log.Printf("Saving credential file to: %s", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Printf("Unable to cache oauth token: %v", err)
		return
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
