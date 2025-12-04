package cmd

import (
	"fmt"
	"log"
	"path"
	"path/filepath"

	"github.com/defektive/gnas/pkg/files"
	"github.com/spf13/cobra"
)

// UploadCmd represents the fetch command
var UploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		remoteDir, _ := cmd.Flags().GetString("remote-dir")
		for _, fileToPut := range args {

			remoteFile := filepath.Clean(fileToPut)
			remoteDir := path.Clean(remoteDir)

			remotePath := path.Join(remoteDir, remoteFile)

			serverURL := fmt.Sprintf("%s/%s", UploadEndpoint, remotePath)
			err := files.UploadFile(serverURL, UploadToken, fileToPut)
			if err != nil {
				log.Println(err)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(UploadCmd)

	UploadCmd.Flags().StringVarP(&UploadEndpoint, "server", "s", UploadEndpoint, "Server URL")
	UploadCmd.Flags().StringVarP(&UploadToken, "token", "t", UploadToken, "Auth Token")
	UploadCmd.Flags().StringP("remote-dir", "d", "", "Remote directory")
}
