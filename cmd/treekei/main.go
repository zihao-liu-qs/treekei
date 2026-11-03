package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/zihao-liu-qs/treekei/internal/filter"
	"github.com/zihao-liu-qs/treekei/internal/printer"
	"github.com/zihao-liu-qs/treekei/internal/scanner"
)

var version = "dev"

func init() {
	rootCmd.Version = version

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

	rootCmd.Flags().IntVarP(
		&minLineCount,
		"above",
		"m",
		0,
		"Only show entries that are above a certain line count threshold",
	)

	rootCmd.Flags().StringArrayVar(
		&matchPatterns,
		"match",
		nil,
		"Only show files whose relative path matches the regex (repeatable)",
	)

	rootCmd.Flags().StringArrayVar(
		&noMatchPatterns,
		"no-match",
		nil,
		"Skip files whose relative path matches; prune matching directories (repeatable)",
	)

	rootCmd.Flags().BoolVarP(
		&ignoreCase,
		"ignore-case",
		"i",
		false,
		"Case-insensitive regex matching",
	)

	rootCmd.AddCommand(listLangsCmd)
}

var (
	showAll         bool
	maxDepth        int
	sortBy          string
	noColor         bool
	dirOnly         bool
	langs           string
	minLineCount    int
	matchPatterns   []string
	noMatchPatterns []string
	ignoreCase      bool
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

		matcher, err := filter.Compile(matchPatterns, noMatchPatterns, ignoreCase)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid regex: %v\n", err)
			os.Exit(2)
		}

		sortMode := scanner.SortByLines
		switch sortBy {
		case string(scanner.SortByLines):
			sortMode = scanner.SortByLines
		case string(scanner.SortByName):
			sortMode = scanner.SortByName
		default:
			fmt.Fprintf(os.Stderr, "Error: unknown sort mode %q (use lines|name)\n", sortBy)
			os.Exit(2)
		}

		s := scanner.NewScanner(scanner.Options{
			ShowAll: showAll,
			SortBy:  sortMode,
			Langs:   langs,
			Matcher: matcher,
		})

		p := printer.NewPrinter(maxDepth, noColor, dirOnly, minLineCount)

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

var listLangsCmd = &cobra.Command{
	Use:   "list-langs",
	Short: "List languages supported by the --lang filter",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		for _, lang := range scanner.SupportedLangs() {
			fmt.Println(lang)
		}
	},
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("Failed to execute treekei command")
	}
}
