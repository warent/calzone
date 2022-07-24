package cmd

import (
	"fmt"
	"log"
	"net/rpc"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/warent/calzone/service/structures/args"
)

const ADDR = "127.0.0.1"
const PORT = 61895

var rootCmd = &cobra.Command{
	Use:   "calzone",
	Short: "Calzone is an easy app management tool",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
		// do actual work
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// func init() {
// cobra.OnInitialize(initConfig)
// }

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(installCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, cobraArgs []string) {
		fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
	},
}

var installCmd = &cobra.Command{
	Use: "install",
	Run: func(cmd *cobra.Command, cobraArgs []string) {
		rpcConn, err := rpc.Dial("tcp", fmt.Sprintf("%v:%v", ADDR, PORT))
		if err != nil {
			fmt.Println("Error connecting to Calzone")
			return
		}
		defer rpcConn.Close()
		rpcArgs := &args.Install{
			Calzone: cobraArgs[0],
		}
		err = rpcConn.Call("Calzone.Install", rpcArgs, nil)
		if err != nil {
			log.Fatal("arith error:", err)
		}
	},
}

func initConfig() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	viper.AddConfigPath(home + "/.calzone")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}

}
