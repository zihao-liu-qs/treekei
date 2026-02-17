package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/zihao-liu-qs/treekei/internal/printer"
	"github.com/zihao-liu-qs/treekei/internal/scanner"
)

func init() {
	rootCmd.Flags().BoolVarP(
		&showAll,
		"all",
		"a",
		false, // default value
		"Include hidden files",
	)

	rootCmd.Flags().IntVarP(
		&maxDepth,
		"level",
		"L",
		0, // 0 for infinite depth
		"Max display depth (0 means unlimited)",
	)

	rootCmd.Flags().StringVarP(
		&sortBy,
		"sort",
		"s",
		"lines",                 // default sort by lines
		"Sort by: lines | name", // lines = scanner.SortByLines name = scanner.SortByName
	)

	rootCmd.Flags().BoolVar(
		&noColor,
		"no-color",
		false,
		"Disable colored output",
	)

	rootCmd.Flags().BoolVarP(
		&dirOnly,
		"dir-only",
		"d",
		false,
		"List directories only",
	)

	rootCmd.Flags().StringVarP(
		&langs,
		"lang",
		"l",
		"",
		"Filter by language(s), comma-separated (e.g. go,ts)",
	)
}

var (
	showAll  bool
	maxDepth int
	sortBy   string
	noColor  bool
	dirOnly  bool
	langs    string
)

var rootCmd = &cobra.Command{
	Use:   "treekei [dirPath]",
	Short: "A CLI tool to show file line counts in a tree shape",
	Args:  cobra.MaximumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		targetPath := "."
		if len(args) == 1 {
			targetPath = args[0]
		}

		absPath, err := filepath.Abs(targetPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: cannot resolve path: %v\n", err)
			os.Exit(1)
		}

		var s *scanner.Scanner
		switch sortBy {
		case string(scanner.SortByLines):
			s = scanner.NewScanner(showAll, scanner.SortByLines, langs)
		case string(scanner.SortByName):
			s = scanner.NewScanner(showAll, scanner.SortByName, langs)
		}

		p := printer.NewPrinter(maxDepth, noColor, dirOnly)

		root, errors := s.ScanDirectory(absPath)
		if root == nil {
			fmt.Fprintf(os.Stderr, "Error: cannot scan directory\n")
			os.Exit(1)
		}

		p.PrintTree(root, "", true, true, 0)

		if len(errors) > 0 {
			fmt.Fprintf(os.Stderr, "\n%d error(s) encountered:\n", len(errors))
			for _, err := range errors {
				fmt.Fprintf(os.Stderr, "  - %v\n", err)
			}
		}
	},
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("Failed to execute treekei command")
	}
}
