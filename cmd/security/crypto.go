/*
Copyright © 2020 @mas2020 andrea.genovesi@gmail.com

*/
package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"github.com/mas2020-golang/ion/packages/utils"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var (
	rows                     int = 10
	file                     string
	encrypt, decrypt, remove bool
)

// NewCryptoCmd represents the crypto command
func NewCryptoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "crypto ...",
		Example: ``,
		Short:   "Crypto",
		Long:    `Long Crypto`,
		Run: func(cmd *cobra.Command, args []string) {
			if encrypt && decrypt {
				fmt.Printf("%s you cannot specify both --encrypt and --decrypt flags\n", utils.RedS("Error:"))
				return
			}
			// TODO: ask the key reading from the standard input
			key := "test1234test1234"

			if encrypt {
				err := encryptFile(file, key)
				utils.Check(err)
				fmt.Printf("%s the file %q has been encrypted\n", utils.GreenS("Success:"), file)
			}
			if decrypt {
				err := decryptFile(file, key)
				utils.Check(err)
				fmt.Printf("%s the file %q has been decrypted\n", utils.GreenS("Success:"), file)
			}
			removeFile(file)
		},
	}

	// flags
	cmd.Flags().StringVarP(&file, "file", "f", "", "file to encrypt/decrypt")
	cmd.Flags().BoolVarP(&encrypt, "encrypt", "e", false, "specify if the to encrypt")
	cmd.Flags().BoolVarP(&decrypt, "decrypt", "d", false, "specify if the to decrypt")
	cmd.Flags().BoolVarP(&remove, "remove", "r", false, "remove the original file")
	return cmd
}

func removeFile(path string) {
	if remove {
		utils.Check(os.Remove(file))
	}
}
func encryptFile(path, key string) (err error) {
	inFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		utils.Check(f.Close())
	}(inFile)

	// The key should be 16 bytes (AES-128), 24 bytes (AES-192) or
	// 32 bytes (AES-256)
	k := []byte(key)
	if len(k) != 16 && len(k) != 24 && len(k) != 32 {
		return fmt.Errorf("the key must be 16,24,32 bytes length")
	}
	block, err := aes.NewCipher(k)
	if err != nil {
		return fmt.Errorf("the Cyther block has not been created: %v", err)
	}
	iv := make([]byte, block.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return fmt.Errorf("error creating the iv factor: %v", err)
	}

	// destination file
	outFile, err := os.OpenFile(path+".crypto", os.O_RDWR|os.O_CREATE, 0777)
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
			stream.XORKeyStream(b, b[:n])
			// Write into file
			outFile.Write(b[:n])
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("read %d bytes: %v", n, err)
		}
	}
	// Append the IV
	if _, err = outFile.Write(iv); err != nil {
		return fmt.Errorf("error writing the iv factor to the file: %v", err)
	}

	return nil
}

func decryptFile(path, key string) (err error) {
	inFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		utils.Check(f.Close())
	}(inFile)

	// The key should be 16 bytes (AES-128), 24 bytes (AES-192) or
	// 32 bytes (AES-256)
	k := []byte(key)
	if len(k) != 16 && len(k) != 24 && len(k) != 32 {
		return fmt.Errorf("the key must be 16,24,32 bytes length")
	}
	block, err := aes.NewCipher(k)
	if err != nil {
		return fmt.Errorf("the Cyther block has not been created: %v", err)
	}

	fi, err := inFile.Stat()
	if err != nil {
		return fmt.Errorf("error reading the size of the input file %q: %v", path, err)
	}
	// read the iv factor at the end of the file
	iv := make([]byte, block.BlockSize())
	msgLen := fi.Size() - int64(len(iv))
	_, err = inFile.ReadAt(iv, msgLen)
	if err != nil {
		return fmt.Errorf("error reading the iv factor from the file %q: %v", path, err)
	}

	// destination file
	outFile, err := os.OpenFile(path[:len(path)-7], os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("error creating the output file %q: %v", path[:len(path)-7], err)
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
			if n > int(msgLen) {
				n = int(msgLen)
			}
			msgLen -= int64(n)
			stream.XORKeyStream(b, b[:n])
			// Write into file
			outFile.Write(b[:n])
			if msgLen == 0 {
				break
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("read %d bytes: %v", n, err)
		}
	}

	return nil
}

func encryptText(text string) (buf []byte, err error) {
	return nil, nil
}

func decryptText(text string) (buf []byte, err error) {
	return nil, nil
}
