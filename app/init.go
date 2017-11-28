package app

import (
    "log"

    "github.com/revel/revel"

    "golang.org/x/net/context"
    "google.golang.org/api/option"
    "cloud.google.com/go/firestore"
    "firebase.google.com/go"

    "jdash/app/config"
)

var (
    // AppVersion revel app version (ldflags)
    AppVersion string

    // BuildTime revel app build-time (ldflags)
    BuildTime string
)

func init() {
    // Filters is the default set of global filters.
    revel.Filters = []revel.Filter{
        revel.PanicFilter,             // Recover from panics and display an error page instead.
        revel.RouterFilter,            // Use the routing table to select the right Action
        revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
        revel.ParamsFilter,            // Parse parameters into Controller.Params.
        revel.SessionFilter,           // Restore and write the session cookie.
        revel.FlashFilter,             // Restore and write the flash cookie.
        revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
        revel.I18nFilter,              // Resolve the requested language
        HeaderFilter,                  // Add some security based headers
        revel.InterceptorFilter,       // Run interceptors around the action.
        revel.CompressFilter,          // Compress the result.
        revel.ActionInvoker,           // Invoke the action.
    }

    // revel.DevMode, revel.RunMode
    revel.OnAppStart(InitFirebaseApp)
    revel.OnAppStart(InitConfig)
    revel.OnAppStart(ScheduleTasks)
}

// HeaderFilter adds common security headers
// There is a full implementation of a CSRF filter in
// https://github.com/revel/modules/tree/master/csrf
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
    c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
    c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
    c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")

    fc[0](c, fc[1:]) // Execute the next filter stage.
}

var FirebaseApp *firebase.App
var FirestoreClient *firestore.Client
var Context context.Context
var Config *config.Config

func InitFirebaseApp() {
    log.Println("Initializing Firebase and Firestore")
    opt := option.WithCredentialsFile("conf/firebase_config.json")
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

func ScheduleTasks() {
    log.Println("Scheduling dashboard tasks")
}
