package cmd

import (
	"github.com/defektive/gnas/pkg/files"
	"github.com/spf13/cobra"
)

// ServeCmd represents the serve command
var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		dir, _ := cmd.Flags().GetString("dir")
		allowPut, _ := cmd.Flags().GetBool("put")

		files.HTTPServer(ServerListener, dir, allowPut, UploadToken)
	},
}

func init() {
	RootCmd.AddCommand(ServeCmd)

	ServeCmd.Flags().StringVarP(&ServerListener, "listener", "l", ServerListener, "server listen address")
	ServeCmd.Flags().StringP("dir", "d", ".", "directory to serve")
	ServeCmd.Flags().BoolP("put", "p", false, "allow put for file uploads")
	ServeCmd.Flags().StringVarP(&UploadToken, "token", "t", UploadToken, "Auth token")
}
