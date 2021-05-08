package main

import (
	"net/http"
	"os"

	"github.com/wadells/drone-approval/plugin"

	"github.com/drone/drone-go/plugin/validator"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// config contains all configuration as environment variables
type config struct {
	Bind   string `envconfig:"DRONE_BIND"`
	Secret string `envconfig:"DRONE_SECRET"`
}

func main() {
	log := logrus.New()
	log.Out = os.Stdout

	var cfg config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Secret == "" {
		log.Fatalln("missing drone plugin secret")
	}
	if cfg.Bind == "" {
		cfg.Bind = ":5433"
	}

	handler := validator.Handler(
		cfg.Secret,
		plugin.New(log),
		log,
	)

	log.Infof("server listening on address %s", cfg.Bind)

	http.Handle("/", handler)
	log.Fatal(http.ListenAndServe(cfg.Bind, nil))
}
