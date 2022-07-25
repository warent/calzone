package cmd

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/warent/calzone/service/structures"
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
		if len(cobraArgs) != 1 {
			fmt.Println("Unknown arguments. Please supply the calzone package name.")
			return
		}
		calzone := cobraArgs[0]
		fmt.Println("\n# Calzone\n")
		fmt.Println("Starting up!")
		rpcConn, err := rpc.Dial("tcp", fmt.Sprintf("%v:%v", ADDR, PORT))
		if err != nil {
			fmt.Println("Error connecting to Calzone")
			return
		}
		defer rpcConn.Close()
		rpcArgs := &structures.BeginInstallArgs{
			Calzone: calzone,
		}
		resp := structures.BeginInstallResponse{}
		fmt.Printf("Fetching %s... ", calzone)
		err = rpcConn.Call("Calzone.BeginInstall", rpcArgs, &resp)
		if err != nil {
			log.Fatal("arith error:", err)
		}
		completeArgs := &structures.CompleteInstallArgs{
			Calzone:    calzone,
			Parameters: map[string]string{},
		}
		fmt.Println("Got it!")
		if len(resp.Parameters) > 0 {
			fmt.Println("This calzone needs some more information; please supply the parameters below.")
			fmt.Printf("\n## %s\n\n", calzone)
		}
		for key, param := range resp.Parameters {
			fmt.Printf("%s (default=%s): ", param.Description, param.Default)
			reader := bufio.NewReader(os.Stdin)
			userParam, _ := reader.ReadString('\n')
			if userParam == "" {
				userParam = param.Default
			}
			completeArgs.Parameters[key] = strings.Trim(userParam, "\n ")
		}
		fmt.Println("\n## Wrapping up\n")
		fmt.Printf("Calzone is cooking up %s...", calzone)
		completeResp := structures.CompleteInstallResponse{}
		err = rpcConn.Call("Calzone.CompleteInstall", completeArgs, &completeResp)
		if err != nil {
			log.Fatal("arith error:", err)
		}
		fmt.Println("\n\n## Done!\n")
		fmt.Printf("%s is ready on 127.0.0.1:%v", calzone, completeResp.Port)
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
