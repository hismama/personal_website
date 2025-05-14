package config

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/firestore"
	secret "cloud.google.com/go/secretmanager/apiv1"
	secretpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

var (
	Ctx          = context.Background()
	SMTPServer   string
	SMTPPort     string
	SMTPUser     string
	SMTPPass     string
	ToEmail      string
	CookieSecret []byte
	ProjectID    string
	FS           *firestore.Client
)

func init() {
	_ = godotenv.Load()
	cfg := LoadEnv(Ctx)
	SMTPServer = cfg.SMTPServer
	SMTPPort = cfg.SMTPPort
	SMTPUser = cfg.SMTPUser
	SMTPPass = cfg.SMTPPass
	ToEmail = cfg.ContactTo
	CookieSecret = cfg.CookieSecret
	ProjectID = cfg.ProjectId
	FS = FirestoreClient(Ctx, ProjectID)

}

type Config struct {
	SMTPServer   string
	SMTPPort     string
	SMTPUser     string
	SMTPPass     string
	ContactTo    string
	CookieSecret []byte
	ProjectId    string
}

func LoadEnv(ctx context.Context) *Config {
	// GCS pull from Secret Manager
	proj := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if proj == "" {
		if v, err := metadata.ProjectIDWithContext(ctx); err == nil {
			proj = v
		} else {
			log.Fatalf("project ID not found: set GOOGLE_CLOUD_PROJECT or run on GCP")
		}
	}
	get := func(key string) string {
		// Local pull from .env
		if v := os.Getenv(key); v != "" {
			return v
		}

		cl, err := secret.NewClient(ctx)
		if err != nil {
			log.Fatalf("secret client: %v", err)
		}
		name := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", proj, key)
		res, err := cl.AccessSecretVersion(ctx, &secretpb.AccessSecretVersionRequest{
			Name: name,
		})
		if err != nil {
			log.Fatalf("secret %s: %v", key, err)
		}
		return string(res.Payload.Data)
	}

	return &Config{
		SMTPServer:   strings.TrimSpace(get("SMTP_Server")),
		SMTPPort:     strings.TrimSpace(get("SMTP_Port")),
		SMTPUser:     strings.TrimSpace(get("SMTP_USER")),
		SMTPPass:     strings.TrimSpace(get("SMTP_PASS")),
		ContactTo:    strings.TrimSpace(get("CONTACT_TO")),
		CookieSecret: []byte(strings.TrimSpace(get("COOKIE_SECRET"))),
		ProjectId:    proj,
	}
}

func FirestoreClient(ctx context.Context, projID string) *firestore.Client {

	client, err := firestore.NewClientWithDatabase(ctx, projID, os.Getenv("Database_ID"))
	if err != nil {
		log.Fatalf("firestore init: %v", err)
	}
	return client
}
