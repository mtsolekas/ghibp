// Copyright (C) 2022 Marios Tsolekas <marios.tsolekas@gmail.com>
// SPDX-License-Identifier: GPL-3.0-or-later
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

var breachCmd = &cobra.Command{
	Use:   "breach DOMAIN...",
	Short: "Query if the provided domains have been breached",
	Long: `Search the HaveIBeenPwned database for breaches in the provided
domains and if a match is found display all information about the given breach.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for i, domain := range args {
			items, err := queryBreaches(domain)
			if err != nil {
				logger.Fatal(err)
			}

			for j, item := range items {
				if i < len(args)-1 || j < len(items)-1 {
					item += "\n"
				}

				fmt.Println(item)
			}

			if i < len(args)-1 {
				r, err := time.ParseDuration(rateLimit)
				if err != nil {
					logger.Fatal("failed to parse rate limit")
				}

				time.Sleep(r)
			}
		}
	},
}

func queryBreaches(domain string) ([]string, error) {
	match, err := findBreach(domain)
	if err != nil {
		return nil, err
	}

	var items []string

	for _, m := range match {
		msg := fmt.Sprintf(
			"Title: %s\nDomain: %s\nDate: %s\nCount: %d\n",
			m.Title,
			m.Domain,
			m.BreachDate,
			m.PwnCount,
		) + "Data Leaked:\n"
		for _, d := range m.DataClasses {
			msg += fmt.Sprintf("  %s\n", d)
		}

		items = append(items, msg+m.Description)
	}

	return items, nil
}

func findBreach(domain string) ([]Breach, error) {
	resp, err := http.Get(breachLink + "?domain=" + domain)
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

	var breaches []Breach
	if err = json.Unmarshal(body, &breaches); err != nil {
		return nil, err
	}

	if len(breaches) == 0 {
		return nil, errors.New("no matching breaches found")
	}

	return breaches, nil
}
