package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"zabbix-maint/internal/cli"
	"zabbix-maint/internal/config"
	"zabbix-maint/internal/log"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "\033[31mError: %v\033[0m\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	// 初始化日志系统
	logDir, err := os.UserHomeDir()
	if err == nil {
		logFile := filepath.Join(logDir, ".zbx-cli", "logs", "zbx-cli.log")
		if err := log.Init(logFile, "info"); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: log init failed: %v\n", err)
		}
	}

	var (
		instanceFlag = flag.String("i", "", "specify instance name")
		helpFlag     = flag.Bool("h", false, "show help")
		helpLongFlag = flag.Bool("help", false, "show help")
		versionFlag  = flag.Bool("v", false, "show version")
	)
	flag.Parse()

	log.Debugf("Command args: %v", os.Args)

	if *helpFlag || *helpLongFlag {
		printHelp()
		return nil
	}

	if *versionFlag {
		printVersion()
		return nil
	}

	args := flag.Args()
	log.Infof("Starting command: %v", args)

	if len(args) == 0 {
		instanceName := *instanceFlag
		if instanceName == "" {
			instanceName = getDefaultInstance()
			if instanceName == "" {
				fmt.Println("Welcome to Zabbix CLI!")
				fmt.Println()
				fmt.Println("No default instance configured.")
				fmt.Println()
				fmt.Println("Usage:")
				fmt.Println("  zbx-cli -i <instance>              # interactive mode")
				fmt.Println("  zbx-cli config add --name ...      # add instance")
				fmt.Println("  zbx-cli -h                         # show help")
				fmt.Println()
				fmt.Println("Or use: zbx-cli config list          # list configured instances")
				return nil
			}
		}
		return cli.RunInteractiveMode(ctx, instanceName)
	}

	cmd := args[0]
	cmdArgs := args[1:]

	switch cmd {
	case "config":
		return handleConfig(ctx, cmdArgs)
	case "version":
		printVersion()
		return nil
	default:
		instanceName := *instanceFlag
		if instanceName == "" {
			return fmt.Errorf("one-shot mode requires -i <instance>, use -h for help")
		}
		handler, err := cli.NewOneShotHandlerWithInstance(ctx, instanceName)
		if err != nil {
			return err
		}
		return handler.Execute(ctx, append([]string{cmd}, cmdArgs...))
	}
}

func handleConfig(ctx context.Context, args []string) error {
	if len(args) == 0 {
		fmt.Println("Config commands:")
		fmt.Println("  zbx-cli config add --name <name> --url <url> --user <user> --pass <password>")
		fmt.Println("  zbx-cli config list")
		fmt.Println("  zbx-cli config remove <name>")
		fmt.Println("  zbx-cli config test <name>")
		return nil
	}

	subCmd := args[0]
	cmdArgs := args[1:]

	handler := cli.NewConfigHandler()
	var name, url, user, pass string

	for i := 0; i < len(cmdArgs); i++ {
		switch cmdArgs[i] {
		case "--name":
			if i+1 < len(cmdArgs) {
				i++
				name = cmdArgs[i]
			}
		case "--url":
			if i+1 < len(cmdArgs) {
				i++
				url = cmdArgs[i]
			}
		case "--user":
			if i+1 < len(cmdArgs) {
				i++
				user = cmdArgs[i]
			}
		case "--pass":
			if i+1 < len(cmdArgs) {
				i++
				pass = cmdArgs[i]
			}
		}
	}

	switch subCmd {
	case "add":
		if name == "" || url == "" || user == "" || pass == "" {
			return fmt.Errorf("config add requires --name, --url, --user, --pass")
		}
		return handler.Add(ctx, name, url, user, pass)
	case "list":
		return handler.List(ctx)
	case "remove":
		if len(cmdArgs) < 1 {
			return fmt.Errorf("config remove requires <name>")
		}
		return handler.Remove(ctx, cmdArgs[0])
	case "test":
		if len(cmdArgs) < 1 {
			return fmt.Errorf("config test requires <name>")
		}
		return handler.Test(ctx, cmdArgs[0])
	default:
		return fmt.Errorf("unknown config command: %s", subCmd)
	}
}

func printHelp() {
	fmt.Println("Zabbix CLI Tool v1.0")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  zbx-cli -i <instance>                          # interactive mode")
	fmt.Println("  zbx-cli -i <instance> <command> [args...]      # one-shot mode")
	fmt.Println("  zbx-cli config <subcommand> [args...]          # config management")
	fmt.Println()
	fmt.Println("Global Options:")
	fmt.Println("  -i <instance>   Specify instance name")
	fmt.Println("  -h              Show this help")
	fmt.Println("  -v              Show version")
	fmt.Println()
	fmt.Println("Config Commands:")
	fmt.Println("  config add --name <name> --url <url> --user <user> --pass <password>")
	fmt.Println("  config list")
	fmt.Println("  config remove <name>")
	fmt.Println("  config test <name>")
	fmt.Println()
	fmt.Println("One-Shot Commands (require -i <instance>):")
	fmt.Println("  user create --alias <name> --group <id> --password <pwd>")
	fmt.Println("  user disable --id <user_id>")
	fmt.Println("  user enable --id <user_id>")
	fmt.Println("  user passwd --id <user_id> --password <new_pwd>")
	fmt.Println("  user delete --id <user_id>")
	fmt.Println("  user list")
	fmt.Println("  user detail --id <user_id>")
	fmt.Println("  usergroup create --name <name>")
	fmt.Println("  usergroup disable --id <group_id>")
	fmt.Println("  usergroup clear --id <group_id>")
	fmt.Println("  usergroup delete --id <group_id>")
	fmt.Println("  usergroup list")
	fmt.Println("  host clone --src <host_id> --name <new_name>")
	fmt.Println("  host list")
	fmt.Println("  host detail --id <host_id>")
	fmt.Println("  hostgroup create --name <name>")
	fmt.Println("  hostgroup clear --id <group_id>")
	fmt.Println("  hostgroup delete --id <group_id>")
	fmt.Println("  hostgroup list")
	fmt.Println("  dashboard create --name <name>")
	fmt.Println("  dashboard list")
	fmt.Println("  system version")
	fmt.Println("  system status")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  zbx-cli -i prod-zbx5")
	fmt.Println("  zbx-cli -i prod-zbx5 user list")
	fmt.Println("  zbx-cli -i prod-zbx5 user create --alias zhangsan --group 7 --password Temp123!")
	fmt.Println("  zbx-cli config add --name prod-zbx5 --url http://zabbix5/api_jsonrpc.php --user Admin --pass zabbix")
	fmt.Println()
}

func printVersion() {
	fmt.Println("zbx-cli version 1.0.0")
}

func getDefaultInstance() string {
	cfg, err := config.LoadAll()
	if err != nil {
		return ""
	}
	return cfg.Global.DefaultInstance
}
