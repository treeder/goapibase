package main

import (
	"context"
	"log"
	"net/http"

	"github.com/treeder/firetils"
	"github.com/treeder/gcputils"
	"github.com/treeder/goapibase"
	"github.com/treeder/gotils"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	// Setup logging
	var err error
	var l *zap.Logger
	env := gcputils.GetEnvVar("ENV", "development")
	if env == "production" {
		l, err = zap.NewProduction()
	} else {
		l, err = zap.NewDevelopment()
	}
	if err != nil {
		log.Fatal(err)
	}
	ctx = gotils.WithLogger(ctx, l)

	// gProjectID := gcputils.GetEnvVar("G_PROJECT_ID", "")
	// gProjectID2, err := metadata.ProjectID()
	// if err != nil {
	// 	fmt.Println("gprojectID2 error:", err)
	// }
	// fmt.Println("PROJECT_ID FROM METADATA: ", gProjectID2)
	// CAN ALSO GET FROM THE JSON
	opts, gProjectID, err := gcputils.CredentialsAndProjectIDFromEnv("G_KEY", "G_PROJECT_ID")
	if err != nil {
		log.Fatal(err)
	}

	firebaseApp, err := firetils.New(ctx, gProjectID, opts)
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
