// Package cmd is the entry points for all commands
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
	Args: cobra.ExactArgs(1),
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
