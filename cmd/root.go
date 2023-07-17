package cmd

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/y4ney/collect-cnnvd-vuln/internal/utils"
	"golang.org/x/term"
	"io"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

const (
	confDir = ".config/y4ney"
)

var (
	verbose bool
	out     io.WriteCloser = os.Stdout
)
var rootCmd = &cobra.Command{
	Use:   "collect-cnnvd-vuln",
	Short: "CNNVD 漏洞收集",
	Long: `

              _  _              _                                    _                   _        
             | || |            | |                                  | |               	| |       
  ___   ___  | || |  ___   ___ | |_      ___  _ __   _ __ __   __ __| |   __   __ _   _ | | _ __  
 / __| / _ \ | || | / _ \ / __|| __|    / __|| '_ \ | '_ \\ \ / // _' |   \ \ / /| | | || || '_ \ 
| (__ | (_) || || ||  __/| (__ | |_    | (__ | | | || | | |\ V /| (_| |    \ V / | |_| || || | | |
 \___| \___/ |_||_| \___| \___| \__|    \___||_| |_||_| |_| \_/  \__,_|     \_/   \__,_||_||_| |_|

collect-cnnvd-vuln 为非官方应用程序，可以自动化的查询、下载和订阅漏洞信息

官网：https://www.cnnvd.org.cn
项目：https://github.com/y4ney/collect-cnnvd-vuln
`,
	//http://www.network-science.de/ascii/
	//doom
	PersistentPreRun: setup,
}

func setup(cmd *cobra.Command, _ []string) {
	initLogger()
	initConfig()
	utils.BindFlags(cmd)
}

func initLogger() {
	defaultLogger := zerolog.New(os.Stderr)

	logLevel := zerolog.InfoLevel
	if verbose {
		logLevel = zerolog.TraceLevel
	}

	zerolog.SetGlobalLevel(logLevel)

	// use color logger when run in terminal
	if isTerminal() {
		defaultLogger = zerolog.New(zerolog.NewConsoleWriter())
	}

	log.Logger = defaultLogger.With().Timestamp().Stack().Logger()
}

func initConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	viper.AddConfigPath(home)
	viper.SetConfigType("json")
	viper.SetConfigName(path.Join(confDir, "collect-cnnvd-vuln.json"))

	if err := viper.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Println("error reading config file")
			os.Exit(1)
		}
	}

	// Environment variables can't have dashes
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}

func isTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(schemaCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(fetchCmd)
	rootCmd.AddCommand(saveCmd)
	// TODO Subscribe 订阅 cmd
	// TODO 数字大屏

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose logging (DEBUG and below)")

	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true
}
