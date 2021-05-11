package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/wadells/drone-approval/plugin"

	"github.com/drone/drone-go/plugin/validator"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

var (
	version = "unknown"
	commit  = "unknown"
)

// config contains all configuration as environment variables
type config struct {
	Bind   string `envconfig:"DRONE_BIND"`
	Secret string `envconfig:"DRONE_SECRET"`
}

func main() {
	log := logrus.New()
	log.Out = os.Stdout

	if len(os.Args) > 1 {
		if strings.Contains(os.Args[1], "version") {
			log.Info("version:\t" + version)
			log.Info("commit: \t" + commit)
			return
		} else {
			log.Fatal("unexpected arguments: " + strings.Join(os.Args[1:], " "))
		}
	}

	var cfg config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Secret == "" {
		log.Fatalln("missing drone plugin secret")
	}
	if cfg.Bind == "" {
		cfg.Bind = ":80"
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
