/*
Copyright © 2020 @mas2020 andrea.genovesi@gmail.com
*/
package security

import (
	"os"

	"github.com/mas2020-golang/goutils/output"
	"github.com/spf13/cobra"
)

// NewCryptoCmd represents the crypto command
func NewDecryptCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "decrypt <PATH>",
		Args: cobra.MinimumNArgs(1),
		Example: `$ ion decrypt /tmp --remove
$ ion decrypt /tmp/myfile.txt.crypto`,
		Short: "An easy way to decrypt file and folders",
		Long:  `An easy way to decrypt file and folders using the AES algo with a 256 bits key.`,
		Run: func(cmd *cobra.Command, args []string) {
			_, err := os.Stat(args[0])
			output.CheckErrorAndExit("", "", err)

			// ask for password
			key, err := askForPassword(true)
			output.CheckErrorAndExit("", "", err)

			err = cryptographyExec(args[0], key, false)
			output.CheckErrorAndExit("", "", err)
		},
	}
	cmd.GroupID = "sec"
	// flags
	cmd.Flags().BoolVarP(&remove, "remove", "r", false, "remove the .crypto file")
	return cmd
}
