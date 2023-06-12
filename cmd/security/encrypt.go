/*
Copyright © 2020 @mas2020 andrea.genovesi@gmail.com
*/
package security

import (
	"os"

	"github.com/mas2020-golang/ion/packages/utils"
	"github.com/spf13/cobra"
)

var (
	remove bool
)

// NewEncryptCmd represents the crypto command
func NewEncryptCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "encrypt <PATH>",
		Args:    cobra.MinimumNArgs(1),
		Example: `$ ion encrypt /tmp --remove
$ ion encrypt /tmp/myfile.txt`,
		Short:   "An easy way to encrypt file and folders",
		Long: `An easy way to encrypt file and folders using the AES algo with a 256 bits key.
`,
		Run: func(cmd *cobra.Command, args []string) {
			_, err := os.Stat(args[0])
			utils.Check(err)

			// ask for password
			key, err := askForPassword(false)
			utils.Check(err)

			err = cryptographyExec(args[0], key, true)
			utils.Check(err)
		},
	}

	// flags
	cmd.Flags().BoolVarP(&remove, "remove", "r", false, "remove the original file")

	// help
	//cmd.SetHelpFunc(func(command *cobra.Command, strings []string) {
	//	fmt.Fprintf(cmd.OutOrStdout(), "my custom help")
	//})
	return cmd
}
