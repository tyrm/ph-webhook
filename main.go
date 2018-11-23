package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"./models"
	"./uniwebhook"
	"github.com/gorilla/mux"
	"github.com/juju/loggo"
)

var logger *loggo.Logger

func main() {
	loggo.ConfigureLoggers("<root>=TRACE")

	newLogger := loggo.GetLogger("webhook")
	logger = &newLogger

	config := CollectConfig()

	// Connect DB
	models.InitDB(config.DBEngine)
	defer models.CloseDB()

	// Create Top Router
	r := mux.NewRouter()
	r.Use(uniwebhook.LoggingMiddleware)
	r.PathPrefix("/").HandlerFunc(uniwebhook.HandleCat) // Top Router Catch All

	go http.ListenAndServe(":8080", r)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	nch := make(chan os.Signal)
	signal.Notify(nch, syscall.SIGINT, syscall.SIGTERM)
	logger.Infof("%s", <-nch)

	logger.Infof("Done!")
}
