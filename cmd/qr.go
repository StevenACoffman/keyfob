/*
The MIT License (MIT)

Copyright © 2019 StevenACoffman

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"encoding/base32"
	"fmt"
	"github.com/mdp/qrterminal"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
	"log"
	"os"
	osUser "os/user"
)

// qrCmd represents the qr command
var qrCmd = &cobra.Command{
	Use:   "qr [key name]",
	Short: "Generate a QR Code for the named key",
	Long: `qr [key name] prints a QR Code for the key with the given name.
This can be useful for backing up QR Codes to Google Authenticator or Authy or whatever.`,

	Run: func(cmd *cobra.Command, args []string) {

		service := "keyfob"
		keyName := args[0]

		err := generateQRCode(service, keyName)
		if err != nil {
			log.Fatal(err)
			return
		}
	},
}

func generateQRCode(service, keyName string) error {
	secret, err := keyring.Get(service, keyName)
	if err != nil {
		return err
	}
	raw, err := decodeKey(secret)
	if err != nil {
		return fmt.Errorf("%s: malformed key", secret)
	}

	currentUser, err := osUser.Current()
	if err != nil {
		return err
	}
	uri := fmt.Sprintf("otpauth://totp/%s@%s?secret=%s&issuer=%s",
		keyName+":"+currentUser.Username,
		keyName,
		base32.StdEncoding.EncodeToString(raw),
		keyName,
	)

	qrterminal.Generate(uri, qrterminal.L, os.Stderr)
	return nil
}
func init() {
	rootCmd.AddCommand(qrCmd)
}
