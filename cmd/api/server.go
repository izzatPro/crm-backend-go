package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"

	mw "restapi/internal/api/middlewares"
	"restapi/internal/api/router"
	"restapi/pkg/utils"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env from current working directory (project root)
	// If .env not found — continue and rely on OS env vars
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found. Using OS environment variables.")
	}

	port := os.Getenv("API_PORT")
	if port == "" {
		// sensible default
		port = ":3000"
	}

	cert := os.Getenv("CERT_FILE")
	key := os.Getenv("KEY_FILE")

	// If you run HTTPS/TLS, you must provide cert & key
	if cert == "" || key == "" {
		log.Println("Warning: CERT_FILE or KEY_FILE is empty. Server will start WITHOUT TLS.")
	}

	// Safer minimum TLS version
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// HPP options (example)
	hppOptions := mw.HPPOptions{
		CheckQuery:                  true,
		CheckBody:                   true,
		CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
		Whitelist:                   []string{"sortBy", "sortOrder", "name", "age", "class"},
	}

	// Router
	r := router.MainRouter()

	// JWT middleware excluded paths (public routes)
	jwtMiddleware := mw.MiddlewaresExcludePaths(
		mw.JWTMiddleware,
		"/execs/login",
		"/execs/forgotpassword",
		"/execs/resetpassword/reset",
	)

	// Apply middlewares
	secureMux := utils.ApplyMiddlewares(
		r,
		mw.SecurityHeaders,
		mw.Compression,
		mw.Hpp(hppOptions),
		mw.XSSMiddleware,
		jwtMiddleware,
		mw.ResponseTimeMiddleware,
		mw.Cors,
	)

	server := &http.Server{
		Addr:      port,
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server is running on:", port)

	// Start server: prefer TLS if cert & key are provided
	var err error
	if cert != "" && key != "" {
		err = server.ListenAndServeTLS(cert, key)
	} else {
		err = server.ListenAndServe()
	}

	if err != nil && err != http.ErrServerClosed {
		log.Fatalln("Error starting the server:", err)
	}
}
