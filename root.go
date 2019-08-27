/*
Copyright © 2019 Azim Sonawalla <azim.sonawalla@gmail.com>

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

package main

import (
	"fmt"
	"log"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "unbake",
	Short: "Convert docker buildx bake files into sequential docker commands",
	Long: `docker/buildx is a plugin to the most recent docker distribution
that allows users to define high level build information (bake files)
and transmit that build info to the buildkit-powered docker daemon
in a single invocation.

Since this is a relatively new piece of technolgoy, many environments (including
primarily CI and CD systems) don't yet have the most recent version of docker or
the buildx plugin, and it introduces significant overhead to install these.

unbake solves for this discrepancy by taking a bake file as an input and generating
plain docker commands to build those targets. You'd want to run the unbake container
in your CI against a bake file and pipe the output of it to shell to continue with a
multi-invocation docker build process.`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, file := range files {
			var commands, err = unbake(file)
			if err != nil {
				log.Panic(err)
			}
			for _, command := range commands {
				fmt.Println(command)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var files []string

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default is $HOME/.unbake.yaml)")

	rootCmd.Flags().StringArrayVarP(&files, "file", "f", []string{},
		"bake file to process")
	_ = rootCmd.MarkFlagRequired("file")

	rootCmd.Flags().BoolVarP(&buildKit, "buildkit", "b", false,
		"prepend DOCKER_BUILDKIT=1 to commands")

	rootCmd.Flags().StringVar(&dockerCfg, "docker-config", "",
		"custom docker client config directory")

	rootCmd.Flags().BoolVarP(&quiet, "quiet", "q", false,
		"docker build only prints results")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".unbake" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".unbake")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
