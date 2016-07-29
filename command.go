// operations related to command line interface input & parsing
package main

import (
	"flag"
	"fmt"
	"os"
)

func printAppUsage() {
	fmt.Fprintln(os.Stderr, "Web middleware app to log, proxy requests etc.\n")
	fmt.Fprintf(os.Stderr, "Usage: %s [options] [arg]\n\n", os.Args[0])
	fmt.Fprintln(os.Stderr, `Options:`)
	fmt.Fprintln(os.Stderr, "  -help | -h show Print this help")

	flag.PrintDefaults()

	fmt.Fprintf(os.Stderr, "\nArguments:\n")
	fmt.Fprintf(os.Stderr, "  %s - initialize example config for running application. If file exists, then it will be reset to defaults\n", argCfgInit)
	fmt.Fprintf(os.Stderr, "\nExample:\n  %s %s \n\n", os.Args[0], argCfgInit)
}

// parseAppInput parses input options and args from command line. Returns false when app should stop running
// after function execution
func parseAppInput(cfg string) bool {
	if flag.NArg() > 0 {
		switch flag.Arg(0) {
		case argCfgInit:
			if ok, err := saveToFile(appConfigExample, cfg); !ok {
				fmt.Fprintf(os.Stderr, "Failed to init file \"%s\". Reason: %s\n", cfg, err.Error())
			}
			fmt.Printf("File \"%s\" was initialized successfully. Customize it and run application\n", cfg)
		default:
			fmt.Fprintf(os.Stderr, "Unknown command: %s\n", flag.Arg(0))
			printAppUsage()
		}

		return false
	}

	if _, err := os.Stat(cfg); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "File \"%s\" was not found. To init config run: %s %s\n", cfg, os.Args[0], argCfgInit)
		return false
	}

	return true
}
