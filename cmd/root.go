/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
)

type CompulsoryFlag struct {
	Port   string
	Origin string
}

type OptionalFlag struct {
	ClearCache bool
}

var Newflag CompulsoryFlag
var NewOptionalFlag OptionalFlag

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "CacheFlow",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		clearCache, _ := cmd.Flags().GetBool("clear-cache")
		if clearCache {
			address := os.Getenv("CACHEFLOW_REDIS_ADDRESS")
			port := os.Getenv("CACHEFLOW_REDIS_PORT")
			password := os.Getenv("CACHEFLOW_REDIS_PASSWORD")
			database := os.Getenv("CACHEFLOW_REDIS_DATABASE")

			num, err := strconv.Atoi(database)
			if err != nil {
				log.Fatalf("database incorrect %v", err)
			}

			fullAddress := fmt.Sprintf("%s:%s", address, port)

			rdb := redis.NewClient(&redis.Options{
				Addr:     fullAddress,
				Password: password,
				DB:       num,
			})
			ctx := context.Background()
			searchPattern := "*" + Newflag.Origin + "*"
			delCount := 0
			iter := rdb.Scan(ctx, 0, searchPattern, 0).Iterator()
			for iter.Next(ctx) {
				rdb.Del(ctx, iter.Val())
				delCount++;
			}
			log.Printf("Cache clear successful, Delete Count: %d", delCount)
		}
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

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.CacheFlow.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.Flags().StringVar(&Newflag.Port, "port", "3000", "Set the port for the application to run (port required)")
	rootCmd.Flags().StringVar(&Newflag.Origin, "origin", "https://dummyjson.com", "Set the origin for the application to run (origin required)")
	rootCmd.Flags().BoolP("clear-cache", "c", false, "Delete the cache associated with the origin")
	rootCmd.MarkFlagRequired("port")
	rootCmd.MarkFlagRequired("origin")
}
