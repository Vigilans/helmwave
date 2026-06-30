package helper

import (
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	helm "helm.sh/helm/v3/pkg/cli"
)

func Dotenv(path string) {
	if path == "" {
		if _, err := os.Stat(".env"); err == nil {
			path = ".env"
		}
	}

	if path != "" {
		if err := godotenv.Load(path); err != nil {
			log.Fatalf("Error loading env file %q: %s", path, err)
		}
	}

	Helm = helm.New() // Recreate helm instance to respect helm variables from .env file
}
