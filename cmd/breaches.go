// Copyright (C) 2022 Marios Tsolekas <marios.tsolekas@gmail.com>
// SPDX-License-Identifier: GPL-3.0-or-later
package cmd

import (
	"fmt"

	"github.com/microcosm-cc/bluemonday"
	"github.com/mmcdole/gofeed"
	"github.com/spf13/cobra"
)

var breachesCmd = &cobra.Command{
	Use:   "breaches",
	Short: "Get the latest breaches",
	Long: `Display information on the most recent breaches as they appear on
https://feeds.feedburner.com/HaveIBeenPwnedLatestBreaches. For each breach
display its title and total number of breached accounts, the date of the breach,
a short description and a link to the listing on https://haveibeenpwned.com.`,
	Args: cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		items, err := queryNewest()
		if err != nil {
			logger.Fatal(err)
		}

		for i := len(items) - 1; i >= 0; i-- {
			fmt.Println(items[i])
		}
	},
}

func queryNewest() ([]string, error) {
	fp := gofeed.NewParser()

	feed, err := fp.ParseURL(breachFeed)
	if err != nil {
		return nil, err
	}

	p := bluemonday.StrictPolicy()

	items := make([]string, feed.Len())

	for i, item := range feed.Items {
		items[i] = item.Title + "\n" +
			item.Published + "\n" +
			item.Link + "\n" +
			item.Description

		if i > 0 {
			items[i] += "\n"
		}

		items[i] = p.Sanitize(items[i])
	}

	return items, nil
}
