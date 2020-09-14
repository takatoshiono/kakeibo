package mf

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/takatoshiono/kakeibo/backend/internal/cmd/mf/db"
)

var cfgFile string

func init() {
	cobra.OnInitialize(initConfig)
}

// NewCmd creates the `mf` command.
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mf",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		//	Run: func(cmd *cobra.Command, args []string) { },
	}

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mf.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	cmd.AddCommand(driveCmd)
	cmd.AddCommand(db.NewCmdDB(
		&db.Options{
			ImportOption: &db.ImportOption{
				DriverName: os.Getenv("DB_DRIVER_NAME"),
				DSN:        os.Getenv("DB_DSN"),
			},
		},
	))
	cmd.AddCommand(csvCmd)

	return cmd
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

		// Search config in home directory with name ".mf" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".mf")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
