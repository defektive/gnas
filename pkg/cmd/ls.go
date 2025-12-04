package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/defektive/gnas/pkg/files"
	"github.com/spf13/cobra"
)

// ListCmd represents the ls command
var ListCmd = &cobra.Command{
	Use:   "ls",
	Short: "List files",
	Long:  `List files`,
	Run: func(cmd *cobra.Command, args []string) {

		// Default path to list if no arguments are provided
		listPath := "."

		var err error
		// Check for command-line arguments
		if len(args) > 0 {
			listPath = args[0]
		}

		// Normalize the path to handle potential issues with backslashes on Windows
		listPath, err = filepath.Abs(listPath)
		if err != nil {
			log.Println(err)
		}

		// Get the file information
		myFiles, err := files.GetFileInfo(listPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Print the file list
		files.PrintFileList(myFiles)
	},
}

func init() {
	RootCmd.AddCommand(ListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
