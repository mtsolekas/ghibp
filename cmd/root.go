// Copyright (C) 2022 Marios Tsolekas <marios.tsolekas@gmail.com>
// SPDX-License-Identifier: GPL-3.0-or-later
package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type Breach struct {
	Name         string
	Title        string
	Domain       string
	BreachDate   string
	AddedDate    string
	ModifiedDate string
	PwnCount     int
	Description  string
	LogoPath     string
	DataClasses  []string
	IsVerified   bool
	IsFabricated bool
	IsSensitive  bool
	IsRetired    bool
	IsSpamList   bool
	IsMalware    bool
}

const (
	rateLimit  = "1500ms"
	passLink   = "https://api.pwnedpasswords.com/range/"
	breachLink = "https://haveibeenpwned.com/api/v3/breaches"
	breachFeed = "https://feeds.feedburner.com/HaveIBeenPwnedLatestBreaches"
)

var (
	disableManpages bool
	disableMarkdown bool
	manpageDst      string
	markdownDst     string

	logger  *log.Logger
	rootCmd = &cobra.Command{
		Use:   "ghibp command",
		Short: "HaveIBeenPwned public API util",
		Long: `Query the HaveIBeenPwned database for information on breaches and
passwords. See each command's documentation for further details.
Powered by https://haveibeenpwned.com/`,
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	logger = log.New(os.Stderr, filepath.Base(os.Args[0])+": ", 0)

	rootCmd.AddCommand(passwordCmd)
	rootCmd.AddCommand(breachCmd)
	rootCmd.AddCommand(breachesCmd)
	rootCmd.AddCommand(docsCmd)

	docsCmd.Flags().
		BoolVarP(&disableManpages, "disable-manpages", "", false, "Disable manpage generation")
	docsCmd.Flags().
		BoolVarP(&disableMarkdown, "disable-markdown", "", false, "Disable markdown generation")

	docsCmd.Flags().
		StringVarP(&markdownDst, "markdown-destination", "", "./docs/",
			"Top level destination directory for markdown documents")
	docsCmd.Flags().
		StringVarP(&manpageDst, "manpage-destination", "", "./dist/man/",
			"Top level destination directory for manpage documents")

	rootCmd.DisableAutoGenTag = true
}
