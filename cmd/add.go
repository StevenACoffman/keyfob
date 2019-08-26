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
	"bufio"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
	"log"
	"os"
	"strings"
	"unicode"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [key name] [optional key value]",
	Short: "adds or overwrites a new key to the keychain with the given name",
	Long: `adds or overwrites a new key to the keychain with the given name.
It prints a prompt to standard error and reads a two-factor key from standard input.
Two-factor keys are short case-insensitive strings of letters A-Z and digits 2-7.`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {

		service := "keyfob"
		name := args[0]

		var text string

		if len(args) == 1 {
			log.Printf("enter key secret to add for %s: ", name)
			text, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				log.Fatalf("error reading key: %v", err)
			}
			text = strings.Map(noSpace, text)
			text += strings.Repeat("=", -len(text)&7) // pad to 8 bytes

		} else {

			text = args[1]
		}

		if _, err := decodeKey(text); err != nil {
			log.Fatalf("invalid key: %v", err)
		}

		err := keyring.Set(service, name, text)
		if err != nil {
			log.Fatalf("Unable to write application password for keyfob: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func noSpace(r rune) rune {
	if unicode.IsSpace(r) {
		return -1
	}
	return r
}
