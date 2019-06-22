/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
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
	Short: "Generate a One Time Password",
	Long: `otp name prints a two-factor authentication code from the key with the given name. 
If -clip is specified, otp also copies to the code to the system clipboard.
With no arguments, otp prints two-factor authentication codes from all known time-based keys.

The default time-based authentication codes are derived from a hash of the key and the current time,
so it is important that the system clock have at least one-minute accuracy.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		service := "keyfob"
		user := args[0]
		secret, err := keyring.Get(service, user)
		if err != nil {
			log.Fatal(err)
		}
		raw, err := decodeKey(secret)
		if err == nil {
			code := totp(raw, time.Now(), 6)
			codeText := fmt.Sprintf("%0*d", 6, code)

			if clip {
				clipboard.WriteAll(codeText)
			}

			fmt.Printf("%s\n", codeText)
			return
		}
		log.Printf("%s: malformed key", secret)
	},
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