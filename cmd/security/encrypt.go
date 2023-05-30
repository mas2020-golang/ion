/*
Copyright © 2020 @mas2020 andrea.genovesi@gmail.com
*/
package security

import (
	"fmt"
	"github.com/mas2020-golang/ion/packages/utils"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	remove bool
)

// NewCryptoCmd represents the crypto command
func NewEncryptCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "encrypt ...",
		Args:    cobra.MinimumNArgs(1),
		Example: ``,
		Short:   "An easy way to encrypt file and folders",
		Long: `An easy way to encrypt file and folders using the AES algo with a 256 bits key.

Examples:
$ ion encrypt /tmp --remove
$ ion encrypt /tmp/myfile.txt`,
		Run: func(cmd *cobra.Command, args []string) {
			// ask for password
			key, err := askForPassword()
			utils.Check(err)

			err = encrypt(args[0], key)
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

// cryptoFolder takes care to encrypt or decrypt the folder and all the nested content
func encrypt(path, key string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat error for the path %s, error: %v", path, err)
	}

	if fi.IsDir() {
		fis, err := ioutil.ReadDir(path)
		if err != nil {
			return fmt.Errorf("error reading for the path %s, error: %v", path, err)
		}
		for _, fi := range fis {
			if err := encrypt(filepath.Join(path, fi.Name()), key); err != nil {
				return err
			}
		}
	} else {
		// take a look at the extension first
		if strings.HasSuffix(path, ".crypto") {
			utils.Warning(fmt.Sprintf("-- the file '%s' is already encrypted, skipped", utils.BoldS(path)))
			return nil
		}
		fmt.Printf(">> encrypting the file '%s'...", utils.BoldS(path))
		if err := cryptoFile(path, key, true); err != nil {
			fmt.Printf("%s\n", utils.RedS("KO"))
			return err
		}
		fmt.Printf("%s\n", utils.GreenS("DONE"))
		removeFile(path)
	}

	return nil
}

func askForPassword() (string, error) {
	key, err := utils.ReadPassword("Password: ")
	utils.Check(err)
	fmt.Println("")
	key2, err := utils.ReadPassword("Repeat the password:")
	fmt.Println("")
	utils.Check(err)
	if key != key2 {
		return "", fmt.Errorf("the passwords need to be the same")
	}
	if len(key) < 6 {
		return "", fmt.Errorf("The password is too short, use at least a 6 chars length key")
	}
	return key, nil
}
