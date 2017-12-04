package app

import (
	"firebase.google.com/go"
	"jdash/config"
	"cloud.google.com/go/firestore"
	"golang.org/x/net/context"
	"log"
	"google.golang.org/api/option"
)

func Init() {
	if !IsInitialized {
		log.Println("Initializing App and global parameters")
		InitConfig()
		InitFirebaseApp()
		IsInitialized = true
	}
}

var IsInitialized = false
var FirebaseApp *firebase.App
var FirestoreClient *firestore.Client

var appConfig *config.Config

func Config() *config.Config {
	if appConfig == nil {
		InitConfig()
	}
	return appConfig
}

func InitFirebaseApp() {
	log.Println("Initializing Firebase and Firestore")
	configFile := Config().Word[config.FIREBASE_CONFIG_FILE]
	opt := option.WithCredentialsFile(configFile)
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalln(err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	FirebaseApp = app
	FirestoreClient = client
}

func InitConfig() {
	log.Println("Initializing app configuration")
	appConfig = config.Make()
}
