package cmd

import (
	"encoding/base64"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/lyricat/goutils/social/line"
)

var (
	tokenType string
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "generate tokens",
	Run: func(cmd *cobra.Command, args []string) {
		switch tokenType {
		case "line-jwt":
			pub, key, err := line.GenerateJWKPair()
			if err != nil {
				cmd.Printf("failed to generate token: %v\n", err)
				return
			}

			decoded, err := base64.StdEncoding.DecodeString(pub)
			if err != nil {
				cmd.Printf("failed to decode public key: %v\n", err)
				return
			}

			// payload := map[string]string{
			// 	"line_jwk_pub":         pub,
			// 	"line_jwk_pub_decoded": string(decoded),
			// }
			// jsonPayload, err := json.Marshal(payload)
			// if err != nil {
			// 	cmd.Printf("failed to marshal payload: %v\n", err)
			// 	return
			// }

			cmd.Printf("jwt_public_key:\n%s\n\n", pub)
			cmd.Printf("jwt_private_key:\n%v\n\n", key)
			cmd.Printf("jwt_public_key_decoded (copy and paste at line developer console):\n")
			// pretty print json
			fmt.Println(string(decoded))
		}
	},
}

func init() {
	rootCmd.AddCommand(tokenCmd)
	tokenCmd.Flags().StringVarP(&tokenType, "type", "t", "", "line-jwt")
}
