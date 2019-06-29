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
			log.Printf( "added key named %s: ", name)
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