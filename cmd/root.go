/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"goisolator/internal/config"
	"goisolator/internal/docker"
	"os"
	"os/signal"
	"syscall"

	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goisolator",
	Short: "goisolator is a docker utility to isolate containers from each others",
	Long:  `goisolator is a docker utility to isolate containers from each others`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		config.GetConfig()

		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			logrus.Fatal(err)
		}

		dockerService := docker.NewDockerService(cli)

		ctx := context.Background()
		ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
		defer cancel()

		dockerService.Start(ctx)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goisolator.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
