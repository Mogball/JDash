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
		InitFirebaseApp()
		InitConfig()
		IsInitialized = true
	}
}

var IsInitialized = false
var FirebaseApp *firebase.App
var FirestoreClient *firestore.Client
var Context context.Context
var Config *config.Config

func InitFirebaseApp() {
	log.Println("Initializing Firebase and Firestore")
	opt := option.WithCredentialsFile("firebase_config.json")
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
	Context = ctx
}

func InitConfig() {
	log.Println("Initializing app configuration")
	Config = config.Make()
}
