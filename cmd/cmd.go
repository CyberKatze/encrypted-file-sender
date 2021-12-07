package cmd

import (
	"flag"
	"fmt"
)

var app string

func init() {
	flag.StringVar(&app, "app", "client", "app can be client or server by passing 'server' or 'client' ")
	flag.Parse()
}

func GetApp() (string, error) {
	if app == "server" || app == "client" {
		return app, nil
	} else {
		return "", fmt.Errorf("app flag should be 'client' or 'server' not : %q", app)
	}

}
