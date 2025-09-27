// Package main implements the GophKeeper CLI client.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gophkeeper/internal/client"
	"gophkeeper/internal/config"
)

var (
	version   = "1.0.0"
	buildDate = "unknown"
)

func main() {
	var (
		serverURL   = flag.String("server", "", "Server URL (override env)")
		configDir   = flag.String("config", "", "Configuration directory (override env)")
		showVersion = flag.Bool("version", false, "Show version information")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("GophKeeper Client v%s\n", version)
		fmt.Printf("Build Date: %s\n", buildDate)
		os.Exit(0)
	}

	cfg := config.LoadClientConfig()
	if *serverURL != "" {
		cfg.ServerURL = *serverURL
	}
	if *configDir != "" {
		cfg.ConfigDir = *configDir
	}

	// Create client
	cli, err := client.NewClient(cfg.ServerURL, cfg.ConfigDir)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Parse command
	if len(os.Args) < 2 {
		showHelp()
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "register":
		if len(args) < 3 {
			fmt.Println("Usage: gophkeeper register <username> <email> <password>")
			os.Exit(1)
		}
		if err := cli.Register(args[0], args[1], args[2]); err != nil {
			log.Fatalf("Registration failed: %v", err)
		}
		fmt.Println("Registration successful!")

	case "login":
		if len(args) < 2 {
			fmt.Println("Usage: gophkeeper login <username> <password>")
			os.Exit(1)
		}
		if err := cli.Login(args[0], args[1]); err != nil {
			log.Fatalf("Login failed: %v", err)
		}
		fmt.Println("Login successful!")

	case "add":
		if len(args) < 2 {
			fmt.Println("Usage: gophkeeper add <type> <title> [data]")
			fmt.Println("Types: login_password, text, binary, bank_card")
			os.Exit(1)
		}
		if err := cli.AddData(args[0], args[1], args[2:]); err != nil {
			log.Fatalf("Failed to add data: %v", err)
		}
		fmt.Println("Data added successfully!")

	case "list":
		if err := cli.ListData(); err != nil {
			log.Fatalf("Failed to list data: %v", err)
		}

	case "get":
		if len(args) < 1 {
			fmt.Println("Usage: gophkeeper get <id>")
			os.Exit(1)
		}
		if err := cli.GetData(args[0]); err != nil {
			log.Fatalf("Failed to get data: %v", err)
		}

	case "delete":
		if len(args) < 1 {
			fmt.Println("Usage: gophkeeper delete <id>")
			os.Exit(1)
		}
		if err := cli.DeleteData(args[0]); err != nil {
			log.Fatalf("Failed to delete data: %v", err)
		}
		fmt.Println("Data deleted successfully!")

	case "sync":
		if err := cli.SyncData(); err != nil {
			log.Fatalf("Failed to sync data: %v", err)
		}
		fmt.Println("Data synchronized successfully!")

	case "history":
		if len(args) < 1 {
			fmt.Println("Usage: gophkeeper history <id>")
			os.Exit(1)
		}
		if err := cli.ShowHistory(args[0]); err != nil {
			log.Fatalf("Failed to show history: %v", err)
		}

	case "help":
		showHelp()

	default:
		fmt.Printf("Unknown command: %s\n", command)
		showHelp()
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Println("GophKeeper - Secure Password Manager")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  register <username> <email> <password>  Register a new user")
	fmt.Println("  login <username> <password>             Login to your account")
	fmt.Println("  add <type> <title> [data]               Add new data")
	fmt.Println("  list                                    List all data")
	fmt.Println("  get <id>                                Get specific data")
	fmt.Println("  delete <id>                             Delete data")
	fmt.Println("  sync                                    Synchronize with server")
	fmt.Println("  history <id>                            Show data history")
	fmt.Println("  help                                    Show this help")
	fmt.Println("  version                                 Show version information")
	fmt.Println("")
	fmt.Println("Data types:")
	fmt.Println("  login_password  Login and password data")
	fmt.Println("  text           Text data")
	fmt.Println("  binary         Binary data")
	fmt.Println("  bank_card      Bank card data")
}
