package main

import "os"

type appSettings struct {
	Port string
}

func CreateAppSettings() *appSettings {
	var (
		port string
		ok   bool
	)
	if port, ok = os.LookupEnv("PLATO_PORT"); !ok {
		port = "3009"
	}
	appSet := &appSettings{
		Port: port,
	}
	return appSet
}
