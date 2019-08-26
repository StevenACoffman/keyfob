// Package cmd is the entry points for all commands
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os/exec"
)

// vaultCmd represents the vault command
var vaultCmd = &cobra.Command{
	Use:   "vault [key name] [aws profile]",
	Short: "AWS credential helper using AWS Vault and Time-based One Time Password",
	Long: `"vault [key name] [aws profile] will act as an AWS credential helper using
AWS Vault and Time-based One Time Password
Ref: https://docs.aws.amazon.com/cli/latest/topic/config-vars.html#sourcing-credentials-from-external-processes`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		service := "keyfob"
		user := args[0]
		profile := args[1]
		codeText, err := generateTOTP(service, user)
		if err != nil {
			log.Fatal(err)
			return
		}
		out, err := exec.Command(
			"aws-vault", "exec", "--mfa-token="+codeText, "-j", profile).CombinedOutput()
		fmt.Println(string(out))
		if err != nil {
			log.Fatalf("aws-vault returned %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(vaultCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// vaultCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// vaultCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
