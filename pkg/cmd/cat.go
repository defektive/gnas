package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/defektive/gnas/pkg/files"
	"github.com/spf13/cobra"
)

// CatCmd represents the cat command
var CatCmd = &cobra.Command{
	Use:   "cat",
	Short: "cat files",
	Long:  `cat files`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			return
		}

		for _, fileToCat := range args {
			if files.Exists(fileToCat) {

				fh, err := os.Open(args[0])
				if err != nil {
					fmt.Println(err)
					continue
				}
				defer fh.Close()
				io.Copy(os.Stdout, fh)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(CatCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// CatCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// CatCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
