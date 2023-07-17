package utils

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/xerrors"
	"os"
	"strconv"
	"time"
)

// BindFlags binds the viper config values to the flags
func BindFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := f.Name

		if !f.Changed && viper.IsSet(configName) {
			val := viper.Get(configName)
			err := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	})
}
func FormatMonth(month time.Month) (int, error) {
	m := fmt.Sprintf("%d", month)
	mm, err := strconv.Atoi(m)
	if err != nil {
		return 0, xerrors.Errorf("failed to convert month %v:%v", month, err)
	}
	return mm, nil
}
