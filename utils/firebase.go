package utils

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type FirebaseServices struct {
	FCM       *messaging.Client
	Firestore *firestore.Client
}

func InitFirebase() *FirebaseServices {
	ctx := context.Background()
	opt := option.WithAuthCredentialsFile(option.ServiceAccount, "serviceAccountKey.json")

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing firebase app: %v", err)
	}

	fcmClient, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting FCM client: %v", err)
	}

	// Tambahkan inisialisasi Firestore
	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error getting Firestore client: %v", err)
	}

	return &FirebaseServices{
		FCM:       fcmClient,
		Firestore: firestoreClient,
	}
}
