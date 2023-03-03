package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/cli/browser"
	"github.com/cli/oauth/device"
)

type exitCode int

const (
	exitOK     exitCode = 0
	exitError  exitCode = 1
	exitCancel exitCode = 2
	exitAuth   exitCode = 4
)

func main() {
	code := mainRun()
	os.Exit(int(code))
}

func mainRun() exitCode {
	clientID := "0dd9bfc230aae2a3aec6"
	//clientID := "178c6fc778ccc68e1d6a"
	scopes := []string{"repo", "read:org", "gist"}
	httpClient := http.DefaultClient

	code, err := device.RequestCode(httpClient, "https://github.com/login/device/code", clientID, scopes)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Copy code: %s\n", code.UserCode)
	fmt.Printf("then open: %s\n", code.VerificationURI)

	_ = waitForEnter(os.Stdin)

	if err := browser.OpenURL(code.VerificationURI); err != nil {
		fmt.Printf("Failed opening a web browser at %s\n", code.VerificationURI)
		fmt.Printf("  %s\n", err)
		fmt.Print("  Please try entering the URL in your browser manually\n")

		//fmt.Fprintf(w, "%s Failed opening a web browser at %s\n", cs.Red("!"), authURL)
		//fmt.Fprintf(w, "  %s\n", err)
		//fmt.Fprint(w, "  Please try entering the URL in your browser manually\n")
	}

	accessToken, err := device.Wait(context.TODO(), httpClient, "https://github.com/login/oauth/access_token", device.WaitOptions{
		ClientID:   clientID,
		DeviceCode: code,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Access token: %s\n", accessToken.Token)

	return exitOK
}

func waitForEnter(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	return scanner.Err()
}
