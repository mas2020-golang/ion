/*
Copyright © 2020 @mas2020 andrea.genovesi@gmail.com

*/
package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/mas2020-golang/ion/packages/utils"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
)

var (
	rows   int = 10
	file   string
	remove bool
)

// NewCryptoCmd represents the crypto command
func NewCryptoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "crypto ...",
		Example: ``,
		Short:   "Crypto",
		Long:    `Long Crypto`,
		Run: func(cmd *cobra.Command, args []string) {
			key, err := utils.ReadPassword()
			utils.Check(err)
			fmt.Println("")

			// understand of encrypt or decrypt by the presence of the suffix .crypto
			encrypt := !strings.HasSuffix(file, ".crypto")
			err = crypto(file, key, encrypt)
			utils.Check(err)

			if encrypt {
				fmt.Printf("%s the file %q has been encrypted\n", utils.GreenS("Success:"), file)
			} else {
				fmt.Printf("%s the file %q has been decrypted\n", utils.GreenS("Success:"), file)
			}
			removeFile()
		},
	}

	// flags
	cmd.Flags().StringVarP(&file, "file", "f", "", "file to encrypt/decrypt")
	cmd.Flags().BoolVarP(&remove, "remove", "r", false, "remove the original file")

	// help
	//cmd.SetHelpFunc(func(command *cobra.Command, strings []string) {
	//	fmt.Fprintf(cmd.OutOrStdout(), "my custom help")
	//})
	return cmd
}

func removeFile() {
	if remove {
		utils.Check(os.Remove(file))
		//fmt.Printf("%s the file %q has been successfully removed\n", utils.GreenS("Success:"), file)
	}
}

// cryptoFolder takes care to encrypt or decrypt the folder and all the nested content
//func cryptoFolder() {
//
//}

// crypto is the function to encrypt or decrypt a filee
func crypto(path, key string, encrypt bool) (err error) {
	var (
		msgLen  int64
		outFile *os.File
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
	} else {
		// in case of decryption get the iv factor at the of the input file
		fi, err := inFile.Stat()
		if err != nil {
			return fmt.Errorf("error reading the size of the input file %q: %v", path, err)
		}
		// read the iv factor at the end of the file
		msgLen = fi.Size() - int64(len(iv))
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
	}

	return nil
}

//func encryptFile(path, key string) (err error) {
//	// input file
//	inFile, err := os.Open(path)
//	if err != nil {
//		return err
//	}
//	defer func(f *os.File) {
//		utils.Check(f.Close())
//	}(inFile)
//
//	// The key is transformed to be 32 bytes long (AES-256)
//	block, err := getCypher(key)
//	if err != nil {
//		return fmt.Errorf("the Cyther block has not been created: %v", err)
//	}
//	iv := make([]byte, block.BlockSize())
//	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
//		return fmt.Errorf("error creating the iv factor: %v", err)
//	}
//
//	// destination file
//	outFile, err := os.OpenFile(path+".crypto", os.O_RDWR|os.O_CREATE, 0777)
//	if err != nil {
//		return fmt.Errorf("error creating the output file %q: %v", path+".crypto", err)
//	}
//	defer func(f *os.File) {
//		if err := f.Close(); err != nil {
//			utils.Check(f.Close())
//		}
//	}(outFile)
//
//	// The buffer size must be multiple of 16 bytes
//	b := make([]byte, 1024)
//	stream := cipher.NewCTR(block, iv)
//
//	for {
//		n, err := inFile.Read(b)
//		if n > 0 {
//			stream.XORKeyStream(b, b[:n])
//			// Write into file
//			outFile.Write(b[:n])
//		}
//
//		if err == io.EOF {
//			break
//		}
//
//		if err != nil {
//			return fmt.Errorf("read %d bytes: %v", n, err)
//		}
//	}
//	// Append the IV
//	if _, err = outFile.Write(iv); err != nil {
//		return fmt.Errorf("error writing the iv factor to the file: %v", err)
//	}
//
//	return nil
//}
//
//func decryptFile(path, key string) (err error) {
//	// input file
//	inFile, err := os.Open(path)
//	if err != nil {
//		return err
//	}
//	defer func(f *os.File) {
//		utils.Check(f.Close())
//	}(inFile)
//
//	// The key is transformed to be 32 bytes long (AES-256)
//	block, err := getCypher(key)
//	if err != nil {
//		return fmt.Errorf("the Cyther block has not been created: %v", err)
//	}
//
//	fi, err := inFile.Stat()
//	if err != nil {
//		return fmt.Errorf("error reading the size of the input file %q: %v", path, err)
//	}
//	// read the iv factor at the end of the file
//	iv := make([]byte, block.BlockSize())
//	msgLen := fi.Size() - int64(len(iv))
//	_, err = inFile.ReadAt(iv, msgLen)
//	if err != nil {
//		return fmt.Errorf("error reading the iv factor from the file %q: %v", path, err)
//	}
//
//	// destination file
//	outFile, err := os.OpenFile(path[:len(path)-7], os.O_RDWR|os.O_CREATE, 0777)
//	if err != nil {
//		return fmt.Errorf("error creating the output file %q: %v", path[:len(path)-7], err)
//	}
//	defer func(f *os.File) {
//		if err := f.Close(); err != nil {
//			utils.Check(f.Close())
//		}
//	}(outFile)
//
//	// The buffer size must be multiple of 16 bytes
//	b := make([]byte, 1024)
//	stream := cipher.NewCTR(block, iv)
//
//	for {
//		n, err := inFile.Read(b)
//		if n > 0 {
//			if n > int(msgLen) {
//				n = int(msgLen)
//			}
//			msgLen -= int64(n)
//			stream.XORKeyStream(b, b[:n])
//			// Write into file
//			outFile.Write(b[:n])
//			if msgLen == 0 {
//				break
//			}
//		}
//
//		if err == io.EOF {
//			break
//		}
//
//		if err != nil {
//			return fmt.Errorf("read %d bytes: %v", n, err)
//		}
//	}
//
//	return nil
//}

//func encryptText(text string) (buf []byte, err error) {
//	return nil, nil
//}
//
//func decryptText(text string) (buf []byte, err error) {
//	return nil, nil
//}

// getCypher return the Cipher
func getCypher(key string) (cipher.Block, error) {
	k := sha256.Sum256([]byte(key))
	return aes.NewCipher(k[:])
}
