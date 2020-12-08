package mf

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/takatoshiono/kakeibo/backend/internal/cmd/mf/csv"
	"github.com/takatoshiono/kakeibo/backend/internal/cmd/mf/db"
	"github.com/takatoshiono/kakeibo/backend/internal/cmd/mf/drive"
)

// Option is the option for the `mf` command.
type Option struct {
	DriverName           string
	DSN                  string
	MoneyForwardEmail    string
	MoneyForwardPassword string
}

var cfgFile string

func init() {
	cobra.OnInitialize(initConfig)
}

// NewCmd creates the `mf` command.
func NewCmd(o *Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mf",
		Short: "Manage data of Money Forward ME",
		Long:  `Manage data of Money Forward ME.`,
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

	cmd.AddCommand(drive.NewCmdDrive(
		&drive.Options{
			UploadOption:   &drive.UploadOption{},
			DownloadOption: &drive.DownloadOption{},
			DeleteOption:   &drive.DeleteOption{},
		},
	))
	cmd.AddCommand(db.NewCmdDB(
		&db.Options{
			ImportOption: &db.ImportOption{
				DriverName: o.DriverName,
				DSN:        o.DSN,
			},
		},
	))
	cmd.AddCommand(csv.NewCmdCSV(
		&csv.Options{
			DownloadOption: &csv.DownloadOption{
				MoneyForwardEmail:    o.MoneyForwardEmail,
				MoneyForwardPassword: o.MoneyForwardPassword,
			},
		},
	))

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
