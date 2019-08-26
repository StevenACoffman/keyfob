// Package cmd is the entry points for all commands
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os/exec"
	"strings"
)

const (
	execPathKeychain = "/usr/bin/security"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		List("keyfob")
	},
}

func parseValue(line string) string {
	words := strings.FieldsFunc(line, func(r rune) bool {
		if r == '"' {
			return true
		}
		return false
	})
	if len(words) > 3 {
		return words[3]
	}
	return ""
}

// List shows secret key names, identified by service, from the keyring.
func List(service string) (string, error) {
	out, err := exec.Command(
		execPathKeychain,
		"dump-keychain").CombinedOutput()

	if err != nil {
		log.Fatal(err)
	}
	outString := string(out)

	parseDump(outString)

	return "", nil
}

func parseDump(keychainDump string) {
	lines := strings.FieldsFunc(keychainDump, func(r rune) bool {
		if r == '\n' {
			return true
		}
		return false
	})
	classMatches := false
	account := ""
	serviceMatches := false
	for _, line := range lines {

		if strings.HasPrefix(line, "keychain:") {
			if classMatches && serviceMatches {
				fmt.Println(account)
			}
			classMatches = false
			account = ""
			serviceMatches = false
		}
		if strings.HasPrefix(line, "class:") {
			classMatches = line == "class: \"genp\""
		}
		if strings.HasPrefix(line, "    \"acct\"<blob>=\"") {
			account = parseValue(line)
		}
		if strings.HasPrefix(line, "    \"svce\"<blob>=\"") {
			serviceMatches = parseValue(line) == "keyfob"
		}
	}
	// if the very last one was a match, this catches it
	if classMatches && serviceMatches {
		fmt.Println(account)
	}
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
