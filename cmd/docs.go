// Copyright (C) 2024 Marios Tsolekas <marios.tsolekas@gmail.com>
// SPDX-License-Identifier: GPL-3.0-or-later
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Generate documentation",
	Long:  "Generate documentation for all commands and sub-commands in markdown and manpage format",
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if !disableManpages {
			generateManpages()
		}

		if !disableMarkdown {
			generateMarkdown()
		}
	},
}

func generateMarkdown() {
	err := doc.GenMarkdownTree(rootCmd, markdownDst)
	if err != nil {
		logger.Fatal(err)
	}
}

func generateManpages() {
	header := &doc.GenManHeader{
		Title:   "ghibp",
		Section: "1",
	}

	err := doc.GenManTree(rootCmd, header, manpageDst)
	if err != nil {
		logger.Fatal(err)
	}
}
