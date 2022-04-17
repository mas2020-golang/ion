/*
Copyright © 2020 @mas2020 andrea.genovesi@gmail.com

*/
package security

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/mas2020-golang/ion/packages/utils"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	rows            int = 10
	remove, decrypt bool
)

// NewCryptoCmd represents the crypto command
func NewCryptoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "crypto ...",
		Args:    cobra.MinimumNArgs(1),
		Example: ``,
		Short:   "Crypto",
		Long:    `Long Crypto`,
		Run: func(cmd *cobra.Command, args []string) {
			crypto(args)
		},
	}

	// flags
	cmd.Flags().BoolVarP(&remove, "remove", "r", false, "remove the original file")
	cmd.Flags().BoolVarP(&decrypt, "decrypt", "d", false, "if --decrypt is passed, it will decrypt the content, otherwise it will encrypt")

	// help
	//cmd.SetHelpFunc(func(command *cobra.Command, strings []string) {
	//	fmt.Fprintf(cmd.OutOrStdout(), "my custom help")
	//})
	return cmd
}

func removeFile(path string) {
	if remove {
		utils.Check(os.Remove(path))
		//fmt.Printf("%s the file %q has been successfully removed\n", utils.GreenS("Success:"), file)
	}
}

// crypto is the starting point of the command
func crypto(args []string) {
	// check the arg to see if encrypt/decrypt a file or a folder
	_, err := os.Stat(args[0])
	utils.Check(err)

	// ask for password
	key, err := utils.ReadPassword("Password: ")
	utils.Check(err)
	fmt.Println("")
	if !decrypt {
		key2, err := utils.ReadPassword("Repeat the password:")
		fmt.Println("")
		utils.Check(err)
		if key != key2 {
			utils.Error("the passwords need to be the same")
			return
		}
	}
	if len(key) < 6 {
		utils.Warning("The password is too short, use at least a 6 chars length key")
		return
	}

	err = execute(args[0], key)
	utils.Check(err)
}

// cryptoFolder takes care to encrypt or decrypt the folder and all the nested content
func execute(path, key string) error {
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
			if err := execute(filepath.Join(path, fi.Name()), key); err != nil {
				return err
			}
		}
	} else {
		// case is a file to encrypt or decrypt, take a look at the extension
		if decrypt && !strings.HasSuffix(path, ".crypto") {
			utils.Warning(fmt.Sprintf("the file '%s' is already decrypted, skipped", utils.BoldS(path)))
			return nil
		}
		if !decrypt && strings.HasSuffix(path, ".crypto") {
			utils.Warning(fmt.Sprintf("the file '%s' is already encrypted, skipped", utils.BoldS(path)))
			return nil
		}

		// encrypt/decrypt the file
		if !decrypt {
			fmt.Printf("Encrypting the file '%s'...", utils.BoldS(path))
		} else {
			fmt.Printf("Decrypting the file '%s'...", utils.BoldS(path))
		}
		if err := cryptoFile(path, key, !decrypt); err != nil {
			fmt.Printf("%s\n", utils.RedS("KO"))
			return err
		}
		fmt.Printf("%s\n", utils.GreenS("DONE"))
		removeFile(path)
	}

	return nil
}

// cryptoFile is the function to encrypt or decrypt a file. When encrypt the file
// the checksum of the original file is created and appended to the end of the file, after the iv
// factor. When decrypting the checksum is taken and compared with the checksum of the decrypted
// file, if the checksums are not equal means that the decryption password was wrong.
func cryptoFile(path, key string, encrypt bool) (err error) {
	var (
		msgLen   int64
		outFile  *os.File
		checksum = make([]byte, 32)
	)

	// input file
	inFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		utils.Check(f.Close())
	}(inFile)

	// The key is transformed to be 32 bytes long (AES-256)
	block, err := getCypher(key)
	if err != nil {
		return fmt.Errorf("the Cyther block has not been created: %v", err)
	}
	iv := make([]byte, block.BlockSize())

	if encrypt {
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return fmt.Errorf("error creating the iv factor: %v", err)
		}
		// create the checksum from the original file
		hash := sha256.New()
		if _, err := io.Copy(hash, inFile); err != nil {
			return fmt.Errorf("error calculating the checksum of the input file %q: %v", path, err)
		}
		checksum = hash.Sum(nil)
		if _, err = inFile.Seek(0, 0); err != nil {
			return fmt.Errorf("error seeking the file %q: %v", path, err)
		}
	} else {
		// in case of decryption get the iv factor and checksum at the of the input file
		fi, err := inFile.Stat()
		if err != nil {
			return fmt.Errorf("error reading the size of the input file %q: %v", path, err)
		}
		// read the checksum at the end of the file
		msgLen = fi.Size() - int64(len(checksum))
		_, err = inFile.ReadAt(checksum, msgLen)
		if err != nil {
			return fmt.Errorf("error reading the checksum from the file %q: %v", path, err)
		}
		// read the iv factor at the end of the file
		msgLen -= int64(len(iv))
		_, err = inFile.ReadAt(iv, msgLen)
		if err != nil {
			return fmt.Errorf("error reading the iv factor from the file %q: %v", path, err)
		}
	}

	// destination file
	if encrypt {
		outFile, err = os.OpenFile(path+".crypto", os.O_RDWR|os.O_CREATE, 0777)
	} else {
		outFile, err = os.OpenFile(path[:len(path)-7], os.O_RDWR|os.O_CREATE, 0777)
	}
	if err != nil {
		return fmt.Errorf("error creating the output file %q: %v", path+".crypto", err)
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			utils.Check(f.Close())
		}
	}(outFile)

	// The buffer size must be multiple of 16 bytes
	b := make([]byte, 1024)
	stream := cipher.NewCTR(block, iv)

	for {
		n, err := inFile.Read(b)
		if n > 0 {
			if !encrypt {
				// for decryption only
				if n > int(msgLen) {
					n = int(msgLen)
				}
				msgLen -= int64(n)
			}
			stream.XORKeyStream(b, b[:n])
			// Write into file
			_, err = outFile.Write(b[:n])
			if err != nil {
				return fmt.Errorf("error writing the output file: %v", err)
			}
			if !encrypt {
				if msgLen == 0 {
					break
				}
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("read %d bytes: %v", n, err)
		}
	}
	if encrypt {
		// Append the IV
		if _, err = outFile.Write(iv); err != nil {
			return fmt.Errorf("error writing the iv factor to the file: %v", err)
		}
		// Append the checksum
		if _, err = outFile.Write(checksum); err != nil {
			return fmt.Errorf("error writing the checksum to the file: %v", err)
		}
	} else {
		// in case of decryption the 2 checksums have to be compared
		// calculate the checksum of the decrypted file and checks with the checksum
		// of the encrypted file
		if _, err = outFile.Seek(0, 0); err != nil {
			return fmt.Errorf("error seeking the file %q: %v", path[:len(path)-7], err)
		}
		hash := sha256.New()
		if _, err := io.Copy(hash, outFile); err != nil {
			return fmt.Errorf("error calculating the checksum of the output file %q: %v", path[:len(path)-7], err)
		}
		checkOutFile := hash.Sum(nil)
		if !bytes.Equal(checksum, checkOutFile) {
			if err = os.Remove(path[:len(path)-7]); err != nil {
				return fmt.Errorf("something was wrong, the pwd is incorrect or the content has been altered"+
					"\nfor the file %q. Error trying to delete the file", path[:len(path)-7])
			}
			return fmt.Errorf("the password is wrong or the content has been altered for the file %q", path[:len(path)-7])
		}
	}

	return nil
}

// getCypher return the Cipher
func getCypher(key string) (cipher.Block, error) {
	k := sha256.Sum256([]byte(key))
	return aes.NewCipher(k[:])
}
