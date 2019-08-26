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
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
	"log"
	"strings"
	"time"
)

// otpCmd represents the otp command
var otpCmd = &cobra.Command{
	Use:   "otp [key name]",
	Short: "Generate a One Time Password for the named key",
	Long: `otp [key name] prints a two-factor authentication code from the key with the given name. 
If -clip is specified, otp also copies to the code to the system clipboard.
With no arguments, otp prints two-factor authentication codes from all known time-based keys.

The default time-based authentication codes are derived from a hash of the key and the current time,
so it is important that the system clock have at least one-minute accuracy.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		service := "keyfob"
		keyName := args[0]

		codeText, err := generateTOTP(service, keyName)
		if err != nil {
			log.Fatal(err)
			return
		}

		if clip {
			clipboard.WriteAll(codeText)
		}
		//fmt has no prefix, log does
		fmt.Printf("%s\n", codeText)

	},
}

func generateTOTP(service, keyName string) (string, error) {
	secret, err := keyring.Get(service, keyName)
	if err != nil {
		return "", err
	}
	raw, err := decodeKey(secret)
	if err != nil {
		return "", fmt.Errorf("%s: malformed key", secret)
	}
	code := totp(raw, time.Now(), 6)
	codeText := fmt.Sprintf("%0*d", 6, code)

	return codeText, nil
}

var clip bool

func init() {
	rootCmd.AddCommand(otpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// otpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	otpCmd.Flags().BoolVarP(&clip, "clip", "c", false, "If -clip is specified, also copies the code to the system clipboard.")
}

func decodeKey(key string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(strings.ToUpper(key))
}

func hotp(key []byte, counter uint64, digits int) int {
	h := hmac.New(sha1.New, key)
	binary.Write(h, binary.BigEndian, counter)
	sum := h.Sum(nil)
	v := binary.BigEndian.Uint32(sum[sum[len(sum)-1]&0x0F:]) & 0x7FFFFFFF
	d := uint32(1)
	for i := 0; i < digits && i < 8; i++ {
		d *= 10
	}
	return int(v % d)
}

func totp(key []byte, t time.Time, digits int) int {
	return hotp(key, uint64(t.UnixNano())/30e9, digits)
}
