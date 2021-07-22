package main

import (
	"context"
	"log"
	"net/http"

	"github.com/treeder/firetils"
	"github.com/treeder/gcputils"
	"github.com/treeder/goapibase"
	"github.com/treeder/gotils/v2"
)

func main() {
	ctx := context.Background()

	opts, projectID, err := gcputils.CredentialsAndProjectIDFromEnv("G_KEY", "G_PROJECT_ID")
	if err != nil {
		log.Fatal(err)
	}

	// Setup logging, optional, typically will work fine without this, but depends on GCP service you're using
	// gcputils.InitLogging()

	firebaseApp, err := firetils.New(ctx, projectID, opts)
	if err != nil {
		gotils.L(ctx).Sugar().Fatalf("couldn't init firebase newapp: %v\n", err)
	}
	firestore, err := firebaseApp.Firestore(ctx)
	if err != nil {
		gotils.L(ctx).Sugar().Fatalf("couldn't init firestore: %v\n", err)
	}
	// if you want auth:
	// fireauth, err := firebaseApp.Auth(ctx)
	// if err != nil {
	// 	gotils.L(ctx).Sugar().Fatalf("error getting firebase auth client: %v\n", err)
	// }

	// add something to firestore just to be sure it's working
	tmp := firestore.Collection("tmp")
	tmp.Add(ctx, TmpType{Name: "wall-e"})

	r := goapibase.InitRouter(ctx)
	// Setup your routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	// Start server
	_ = goapibase.Start(ctx, gotils.Port(8080), r)
}

type TmpType struct {
	Name string `firestore:"name"`
}
