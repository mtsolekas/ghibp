// Copyright (C) 2022, 2024 Marios Tsolekas <marios.tsolekas@gmail.com>
// SPDX-License-Identifier: GPL-3.0-or-later
package cmd

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/tobischo/gokeepasslib/v3"
	"golang.org/x/term"
)

var (
	ErrPasswordNotFound = errors.New("password not found")

	passwordCmd = &cobra.Command{
		Use:   "password PASSWORD...",
		Short: "Lookup the given passwords for breaches in the HaveIBeenPwned database",
		Long: `Query the HaveIBeenPwned password database for each of the given
passwords and return the number of hits. Passwords are not transmitted in cleartext
but only the first 5 digits of its sha1 hash are sent to the servers and the rest of
the lookup is done locally.`,
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if passKdbx != "" {
				sourceKdbx(passKdbx)
			} else if passStdin {
				for {
					sourceStdin()
				}
			} else {
				sourceArgs(args)
			}
		},
	}
)

func findPass(pass string) ([]string, error) {
	hash := fmt.Sprintf("%x", sha1.Sum([]byte(pass)))

	client := &http.Client{}

	req, err := http.NewRequest("GET", passLink+hash[0:5], nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Add-Padding", "true")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(
			"request returned HTTP status %d",
			resp.StatusCode,
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile("(?i)" + hash[5:] + ":[0-9]+")
	match := string(re.Find(body))
	if match == "" {
		return nil, ErrPasswordNotFound
	}

	return append(strings.Split(match, ":"), hash[0:5]), nil
}

func readPassword() ([]byte, error) {
	fmt.Print("Enter Password (Ctrl-c to exit): ")
	bytepw, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}

	fmt.Println()

	return bytepw, nil
}

func traverseKdbxGrps(groups []gokeepasslib.Group) {
	if len(groups) == 0 {
		return
	}

	for _, grp := range groups {
		if grp.Name == "Recycle Bin" {
			continue
		}

		for _, entry := range grp.Entries {
			match, err := findPass(entry.GetPassword())
			if err != nil {
				if errors.Is(err, ErrPasswordNotFound) {
					logger.Printf(
						"%s for %s/%s",
						err,
						grp.Name,
						entry.GetTitle(),
					)
				} else {
					logger.Fatal(err)
				}

				continue
			}
			fmt.Printf(
				"Password found (%s/%s):\n Hash: %s\n Hits: %s\n",
				grp.Name,
				entry.GetTitle(),
				strings.ToLower(match[2]+match[0]),
				match[1],
			)
		}

		traverseKdbxGrps(grp.Groups)
	}
}

func sourceArgs(args []string) {
	for _, pass := range args {
		match, err := findPass(pass)
		if err != nil {
			if errors.Is(err, ErrPasswordNotFound) {
				logger.Printf("%s (%s)", err, pass)
			} else {
				logger.Fatal(err)
			}
		} else {
			fmt.Printf(
				"Password found (%s):\n  Hash: %s\n  Hits: %s\n",
				pass,
				strings.ToLower(match[2]+match[0]),
				match[1],
			)
		}
	}
}

func sourceStdin() {
	bytepw, err := readPassword()
	if err != nil {
		logger.Fatal(err)
	}

	pass := string(bytepw)
	match, err := findPass(pass)
	if err != nil {
		if errors.Is(err, ErrPasswordNotFound) {
			logger.Printf("%s (%s)", err, pass)
		} else {
			logger.Fatal(err)
		}

		return
	}

	fmt.Printf(
		"Password found (%s):\n  Hash: %s\n  Hits: %s\n",
		pass,
		strings.ToLower(match[2]+match[0]),
		match[1],
	)
}

func sourceKdbx(path string) {
	f, err := os.Open(path)
	if err != nil {
		logger.Fatal(err)
	}

	defer f.Close()

	dbPassword, err := readPassword()
	if err != nil {
		logger.Fatal(err)
	}

	db := gokeepasslib.NewDatabase()
	db.Credentials = gokeepasslib.NewPasswordCredentials(string(dbPassword))
	err = gokeepasslib.NewDecoder(f).Decode(db)
	if err != nil {
		logger.Fatal(err)
	}

	db.UnlockProtectedEntries()
	traverseKdbxGrps(db.Content.Root.Groups)
}
