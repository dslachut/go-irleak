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
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"lachut.net/gogs/dslachut/go-irleak/api"
	"lachut.net/gogs/dslachut/go-irleak/kb"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Runs the irleak upload API server",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf()
		//activeKB := kb.NewSQLiteKB("tmp.db", nil)
		activeKB := getKB()

		http.HandleFunc("/api/temp", func(w http.ResponseWriter, r *http.Request) {
			api.TemperatureHandler(w, r, activeKB)
		})

		http.HandleFunc("/api/auth", func(w http.ResponseWriter, r *http.Request) {
			api.AuthHandler(w, r, activeKB)
		})

		port := viper.GetString("port")
		log.Printf("serving IRLeak API on port %s\n", port)
		http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	},
}

func conf() {
	viper.SetDefault("port", "11021")
	viper.SetDefault("dbtype", "sqlite")
	viper.SetDefault("dbparams", map[string]string{"file": "tmp.db"})
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.irleak")
	viper.AddConfigPath("$HOME/.config/irleak/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err)
		log.Println("config file not found, using defaults")
	}
	log.Printf("%v\n", viper.AllSettings())
}

func getKB() kb.KB {
	switch {
	case viper.GetString("dbtype") == "sqlite":
		params := viper.GetStringMapString("dbparams")
		file := params["file"]
		delete(params, "file")
		if len(params) > 0 {
			return kb.NewSQLiteKB(file, params)
		} else {
			return kb.NewSQLiteKB(file, nil)
		}
	}
	return nil
}

func init() {
	RootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
