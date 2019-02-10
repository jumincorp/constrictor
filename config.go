package constrictor

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	app         = &cobra.Command{}
	programName string
)

func App(name string, shortDesc string, longDesc string, run func([]string)) *cobra.Command {
	app.Use = name
	app.Short = shortDesc
	app.Long = longDesc
	app.Run = func(cmd *cobra.Command, args []string) {
		fmt.Printf("args %v\n", args)
		run(args)
		// Do nothing except evaluate variables
	}
	//app.RunE = func(cmd *cobra.Command, args []string) error {
	//fmt.Printf("E args %v\n", args)
	//return nil
	//// Do nothing except evaluate variables
	//}

	programName = name

	cobra.OnInitialize(readConfig)
	return app
}

func readConfig() {
	viper.SetConfigName(programName)
	viper.AddConfigPath(fmt.Sprintf("/etc/"))
	viper.AddConfigPath(fmt.Sprintf("."))

	viper.SetEnvPrefix(strings.Replace(programName, "-", "_", -1))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		fmt.Printf("error reading config file: %s", err)
	}
}

func StringVar(name string, shortName string, defaultVal string, desc string) func() string {
	app.PersistentFlags().StringP(name, shortName, defaultVal, desc)
	viper.BindPFlag(name, app.PersistentFlags().Lookup(name))

	return func() string {
		return viper.GetString(name)
	}
}

func AddressPortVar(name string, shortName string, defaultVal string, desc string) func() string {
	app.PersistentFlags().StringP(name, shortName, defaultVal, desc)
	viper.BindPFlag(name, app.PersistentFlags().Lookup(name))

	return func() string {
		val := viper.GetString(name)
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			return fmt.Sprintf(":%d", i)
		}
		return val
	}
}

func TimeDurationVar(name string, shortName string, defaultVal string, desc string) func() time.Duration {
	app.PersistentFlags().StringP(name, shortName, defaultVal, desc)
	viper.BindPFlag(name, app.PersistentFlags().Lookup(name))

	return func() time.Duration {
		if delay, ok := viper.Get(name).(int); ok {
			return time.Duration(time.Duration(delay) * time.Second)
		}
		return viper.GetDuration(name)
	}
}
