// Package constrictor lightly wraps around cobra and viper.
// It is meant to help you write a simple one-command app which
// could use either command-line parameters or a configuration
// file or a combination of both.
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

// App returns a simple one-command application that combines cobra and viper configuration
func App(name string, shortDesc string, longDesc string, run func([]string) error) *cobra.Command {
	app.Use = name
	app.Short = shortDesc
	app.Long = longDesc

	app.RunE = func(cmd *cobra.Command, args []string) error {
		if err := readConfig(); err != nil {
			switch err.(type) {
			case viper.ConfigFileNotFoundError:
				// Silently ignore this.
			default:
				return err
			}
		}
		return run(args)
	}

	programName = name

	return app
}

func readConfig() error {
	viper.SetConfigName(programName)
	viper.AddConfigPath(fmt.Sprintf("/etc/"))
	viper.AddConfigPath(fmt.Sprintf("."))

	viper.SetEnvPrefix(strings.Replace(programName, "-", "_", -1))
	viper.AutomaticEnv()

	return viper.ReadInConfig()
}

// StringVar returns a function that can be called to evaluate that string configuration parameter.
func StringVar(name string, shortName string, defaultVal string, desc string) func() string {
	app.PersistentFlags().StringP(name, shortName, defaultVal, desc)
	viper.BindPFlag(name, app.PersistentFlags().Lookup(name))

	return func() string {
		return viper.GetString(name)
	}
}

// AddressPortVar returns a function that can be called to evaluate that Address/Port parameter.
// An Address/Port can be stated as a string containing both values, such as "localhost:80"
// The same address could also be configured using only a port string (":80") or even just a port number (80).
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

// TimeDurationVar returns a function that can be called to evaluate an time duration parameter.
// A Time duration can be configured using a valid Go time duration such as "1m11s" but whereas Go
// would assume a string without suffix such as "42" to mean 42 nanoseconds, constrictor defaults to
// seconds and would this interpret it as 42 seconds.
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
