package main

import (
	"log"
	"os"

	"github.com/UndeadDemidov/ya-pr-diplomb/config"
	"github.com/UndeadDemidov/ya-pr-diplomb/internal/repo/postgres"
	"github.com/UndeadDemidov/ya-pr-diplomb/internal/servers"
	"github.com/UndeadDemidov/ya-pr-diplomb/pkg/telemetry"
)

func main() {
	log.Println("Starting gophkeeper")

	configPath := config.GetPath(os.Getenv("config"))
	cfg, err := config.Get(configPath)
	if err != nil {
		log.Fatalf("Loading config: %v", err)
	}

	appLogger := telemetry.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()
	appLogger.Infof(
		"AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %v",
		cfg.Server.AppVersion,
		cfg.Logger.Level,
		cfg.Server.Mode,
		cfg.Server.SSL,
	)
	appLogger.Infof("Success parsed config: %#v", cfg.Server.AppVersion)

	db, err := postgres.NewPostgresDB(cfg.Postgres)
	if err != nil {
		appLogger.Fatalf("Postgresql init: %s", err)
	}
	defer db.Close()

	grpcSrv := servers.NewGRPC(appLogger, cfg, db)
	appLogger.Fatal(grpcSrv.Run())

	// go startHTTP()
	//
	// // Block forever
	// var wg sync.WaitGroup
	//
	// wg.Add(1)
	// wg.Wait()
}

// func serveSwagger(w http.ResponseWriter, r *http.Request) {
// 	http.ServeFile(w, r, "www/swagger.json")
// }
//
// func startHTTP() {
// 	// Serve the swagger,
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/swagger.json", serveSwagger)
//
// 	fs := http.FileServer(http.Dir("www/swagger-ui"))
// 	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui", fs))
//
// 	log.Println("REST server ready...")
// 	log.Println("Serving Swagger at: http://localhost:8080/swagger-ui/")
//
// 	const timeout = 5
//
// 	srv := &http.Server{ //nolint:exhaustruct
// 		Addr:         "localhost:8000",
// 		Handler:      mux,
// 		ReadTimeout:  timeout * time.Second,
// 		WriteTimeout: timeout * time.Second,
// 	}
//
// 	err := srv.ListenAndServe()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
