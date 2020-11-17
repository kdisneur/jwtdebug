package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"

	"github.com/kdisneur/jwtdebug/internal"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	var showVersion bool
	var jwkurl string
	var hs256 string

	fs := flag.NewFlagSet("jwtdebug", flag.ExitOnError)
	fs.StringVar(&jwkurl, "jwk", "", "url to the jwk set")
	fs.StringVar(&hs256, "hs256", "", "hmac S256 secret key or JWT_DEBUG_HS256 env variable")
	fs.BoolVar(&showVersion, "v", false, "displays the current version")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}

	if showVersion {
		fmt.Println(internal.GetVersionInfo().String())
		return nil
	}

	if hs256 == "" {
		hs256 = os.Getenv("JWT_DEBUG_HS256")
	}

	if hs256 == "" && jwkurl == "" {
		return fmt.Errorf("one of the following needs to be set: -jwk or -hs256")
	}

	if hs256 != "" && jwkurl != "" {
		return fmt.Errorf("only one of the following can be set: -jwk or -hs256")
	}

	input, err := getInput(os.Stdin, strings.TrimSpace(strings.Join(fs.Args(), " ")))
	if err != nil {
		return fmt.Errorf("can't get input data: %v", err)
	}

	var options []jwt.Option

	if jwkurl != "" {
		jwkset, err := jwk.FetchHTTP(jwkurl)
		if err != nil {
			return fmt.Errorf("can't fetch jwk set from '%s': %v", jwkurl, err)
		}

		options = append(options, jwt.WithKeySet(jwkset))
	}

	if hs256 != "" {
		options = append(options, jwt.WithVerify(jwa.HS256, hs256))
	}

	token, err := jwt.ParseString(input, options...)
	if err != nil {
		return fmt.Errorf("can't parse jwt token: %v", err)
	}

	values, err := token.AsMap(context.Background())
	if err != nil {
		return fmt.Errorf("can't get jwt claims: %v", err)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(values); err != nil {
		return fmt.Errorf("can't JSON-ify jwt claims: %v", err)
	}

	return nil
}

func getInput(f *os.File, arg string) (string, error) {
	fi, err := f.Stat()
	if err != nil {
		if arg == "" {
			return "", fmt.Errorf("args are empty and STDIN is not readabale: %v", err)
		}
		return arg, nil
	}

	if fi.Size() == 0 && fi.Mode()&os.ModeCharDevice != 0 {
		if arg == "" {
			return "", fmt.Errorf("args and STDIN are empty")
		}

		return arg, nil
	}

	input, err := ioutil.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("can't read from STDIN: %v", err)
	}

	if arg != "" {
		return "", fmt.Errorf("args and STDIN are both set: choose only one")
	}

	return string(input), nil
}
