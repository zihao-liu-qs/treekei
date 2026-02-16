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

var (
	showAll bool
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

		s := scanner.NewScanner(showAll)
		p := &printer.Printer{}

		root, errors := s.ScanDirectory(absPath)
		if root == nil {
			fmt.Fprintf(os.Stderr, "Error: cannot scan directory\n")
			os.Exit(1)
		}

		p.PrintTree(root, "", true, true)

		if len(errors) > 0 {
			fmt.Fprintf(os.Stderr, "\n%d error(s) encountered:\n", len(errors))
			for _, err := range errors {
				fmt.Fprintf(os.Stderr, "  - %v\n", err)
			}
		}
	},
}

func init() {
	// default is false
	rootCmd.Flags().BoolVarP(
		&showAll,
		"all",
		"a",
		false,
		"Include hidden files",
	)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("Failed to execute treekei command")
	}
}
