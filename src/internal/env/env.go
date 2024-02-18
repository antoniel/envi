package envs

import (
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
)

type EnvVars struct {
	GOOGLE_SECRET       string `validate:"required"`
	GOOGLE_CLIENT_ID    string `validate:"required"`
	GOOGLE_CALLBACK_URL string `validate:"required"`
}

var Env = EnvVars{}

func init() {
	Env = EnvVars{
		GOOGLE_CLIENT_ID:    os.Getenv("GOOGLE_CLIENT_ID"),
		GOOGLE_SECRET:       os.Getenv("GOOGLE_SECRET"),
		GOOGLE_CALLBACK_URL: os.Getenv("GOOGLE_CALLBACK_URL"),
	}

	err := validator.New().Struct(Env)

	if err != nil {
		log.Fatalf("❌ env vars validation error: %v", err)
	}
	fmt.Println("✅ Env loaded")
}
