/*
Copyright © 2020 @mas2020 andrea.genovesi@gmail.com
*/
package security

import (
	"os"

	"github.com/mas2020-golang/ion/packages/utils"
	"github.com/spf13/cobra"
)

// NewCryptoCmd represents the crypto command
func NewDecryptCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "decrypt <PATH>[ PATH]",
		Args:    cobra.MinimumNArgs(1),
		Example: ``,
		Short:   "An easy way to decrypt file and folders",
		Long: `An easy way to decrypt file and folders using the AES algo with a 256 bits key.

Examples:
$ ion decrypt /tmp --remove
$ ion decrypt /tmp/myfile.txt.crypto`,
		Run: func(cmd *cobra.Command, args []string) {
			_, err := os.Stat(args[0])
			utils.Check(err)

			// ask for password
			key, err := askForPassword(true)
			utils.Check(err)

			err = cryptographyExec(args[0], key, false)
			utils.Check(err)
		},
	}

	// flags
	cmd.Flags().BoolVarP(&remove, "remove", "r", false, "remove the .crypto file")
	return cmd
}
