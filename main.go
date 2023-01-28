package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	yadisk "github.com/MOZGIII/yadisk-api"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/yandex"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)


var (
	yandexOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID: os.Getenv("YANDEX_CLIENT_ID"),
		ClientSecret: os.Getenv("YANDEX_CLIENT_SECRET"),
		Endpoint: yandex.Endpoint,
	}
	randomState = "random-state-0711"
	overwrite = flag.Bool("overwrite", false, "do you want to overwrite the file if it exists")
	verbose = flag.Bool("verbose", false, "print info during execution")
	server *http.Server
)


func main() {
	// Parse command line arguments
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: go run main.go [options] local_file.txt disk:/remote_file.txt\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		fmt.Println("\nError:\n",
			"You must provide a local file path and a remote file path.\n",
			"") // last line is for the linter
		os.Exit(2)
		return
	}

	// Get token from environment variable
	token := os.Getenv("YANDEX_TOKEN")
	if token == "" {
		startOAuthFlow()
		return
	} 

	uploadFile(token)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	url := yandexOauthConfig.AuthCodeURL(randomState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != randomState {
		http.Error(w, "invalid oauth state", http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	token, err := yandexOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Code exchange failed", http.StatusInternalServerError)
		return
	}

	storeTokenInDotEnv(token.AccessToken)
	
	fmt.Fprint(w, "<html><head><script>window.close();</script></head><body>Authorization successful! This page can be closed now.</body></html>")

	go func() {
		uploadFile(token.AccessToken)
		server.Shutdown(context.Background())
	}()
}

func startOAuthFlow() {
	// Print the URL for the user to authorize the application
	url := yandexOauthConfig.AuthCodeURL(randomState)
	fmt.Println("Visit the URL for the auth dialog:\n", url)
			

	// Start a local HTTP server to handle the OAuth callback
	server = &http.Server{Addr: ":8080"}
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/callback", handleCallback)
	server.ListenAndServe()
}

func storeTokenInDotEnv(token string) {
	env,_ := godotenv.Read(".env")
	env["YANDEX_TOKEN"] = token
	godotenv.Write(env, ".env")
}

func uploadFile(token string) {
	src := flag.Arg(0)
	dst := flag.Arg(1)

	if *verbose {
		fmt.Printf("Using token \"%s\"...\n", token)
		fmt.Printf("Uploading \"%s\" to \"%s\"...\n", src, dst)
	}

	log.SetPrefix("error: ")

	file, err := os.Open(src)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer file.Close()

	api := yadisk.NewAPI(token)

	if err := api.Upload(file, dst, *overwrite); err != nil {
		log.Fatalln(err)
		return
	}

	fmt.Println("File uploaded successfully!")
}

