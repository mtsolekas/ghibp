// Copyright (C) 2022 Marios Tsolekas <marios.tsolekas@gmail.com>
// SPDX-License-Identifier: GPL-3.0-or-later
package cmd

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var passwordCmd = &cobra.Command{
	Use:   "password PASSWORD...",
	Short: "Lookup the given passwords for breaches in the HaveIBeenPwned database",
	Long: `Query the HaveIBeenPwned password database for each of the given
passwords and return the number of hits. Passwords are not transmitted in cleartext
but only the first 5 digits of its sha1 hash are sent to the servers and the rest of
the lookup is done locally.`,
	Args: cobra.MinimumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		for _, pass := range args {
			match, err := findPass(pass)
			if err != nil {
				logger.Fatal(err)
			}

			fmt.Printf(
				"Password found (%s):\n  Hash: %s\n  Hits: %s\n",
				pass,
				strings.ToLower(match[2]+match[0]),
				match[1],
			)
		}
	},
}

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
		return nil, errors.New("password not found")
	}

	return append(strings.Split(match, ":"), hash[0:5]), nil
}
