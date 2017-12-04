package app

import (
	"firebase.google.com/go"
	"jdash/config"
	"cloud.google.com/go/firestore"
	"golang.org/x/net/context"
	"log"
	"google.golang.org/api/option"
)

var IsInitialized = false

var firestoreClient *firestore.Client
var appConfig *config.Config

func Init() {
	if !IsInitialized {
		log.Println("Initializing App and global parameters")
		if appConfig == nil {
			InitConfig()
		}
		if firestoreClient == nil {
			InitFirebaseApp()
		}
		IsInitialized = true
	}
}

func Config() *config.Config {
	if appConfig == nil {
		InitConfig()
	}
	return appConfig
}

func FirestoreClient() *firestore.Client {
	if firestoreClient == nil {
		InitFirebaseApp()
	}
	return firestoreClient
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
	firestoreClient = client
}

func InitConfig() {
	log.Println("Initializing app configuration")
	appConfig = config.Make()
}
