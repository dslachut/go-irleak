// Copyright © 2017 David Lachut <dslachut@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"

	"github.com/bgentry/speakeasy"
	"github.com/elithrar/simple-scrypt"
	"github.com/spf13/cobra"
)

// useraddCmd represents the useradd command
var useraddCmd = &cobra.Command{
	Use:   "useradd",
	Short: "Add a user to the database",
	Long: `Add a user to the database of the IRLeak Server

Usage: irleak useradd username password`,
	Run: func(cmd *cobra.Command, args []string) {
		conf()
		var password string
		var err error
		if len(args) == 2 {
			password = args[1]
		} else {
			password, err = speakeasy.Ask("Password for new IRLeak user: ")
			if err != nil {
				log.Fatal(err)
			}
			pass2, err := speakeasy.Ask("Please re-enter the password:")
			if err != nil {
				log.Fatal(err)
			}
			if password != pass2 {
				log.Fatal("New password does not match")
			}
		}
		hashBytes, err := scrypt.GenerateFromPassword([]byte(password), scrypt.DefaultParams)
		if err != nil {
			log.Fatal(err)
		}
		hash := string(hashBytes)
		k := getKB()
		ok := k.AddUser(args[0], hash)
		fmt.Printf("useradd called:\n\tuser: %s\n\tsuccess: %v\n", args[0], ok)
	},
}

func init() {
	RootCmd.AddCommand(useraddCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// useraddCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// useraddCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
