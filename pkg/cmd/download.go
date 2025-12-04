package cmd

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/defektive/gnas/pkg/files"
	"github.com/spf13/cobra"
)

// DownloadCmd represents the fetch command
var DownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		for _, fileToGet := range args {

			parse, err := url.Parse(fileToGet)
			if err != nil {
				log.Println(err)
				continue
			}

			saveFile := filepath.Base(parse.Path)

			if files.Exists(saveFile) {
				fmt.Printf("File exists. Would you like to delete it? ")
				var input string
				fmt.Scanln(&input)
				if strings.ToLower(input)[0] == 'y' {
					os.Remove(saveFile)
				}

			}
			err = files.DownloadURL(parse, saveFile)
			if err != nil {
				log.Println(err)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(DownloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// DownloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// DownloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
