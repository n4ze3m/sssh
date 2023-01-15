package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/tidwall/buntdb"
)

var rootCmd = &cobra.Command{
	Use:   "sssh",
	Short: "Depot for your SSH hosts",
	Args:  cobra.MinimumNArgs(1),
	Long:  `sssh is a tool for managing your SSH hosts. It allows you to store your hosts in a central location and then connect to them using a simple command.`,
	Run: func(cmd *cobra.Command, args []string) {
		home, _ := homedir.Dir()
		db, _ := buntdb.Open(home + "/ssh-manager.db")

		defer db.Close()

		name := cmd.Flag("name").Value.String()
		ssh := strings.Join(args, "")

		if name == "list" {
			fmt.Println("Name cannot be 'list'")
			os.Exit(1)
		}

		switch {
		case name != "" && ssh != "":
			var value string
			err := db.View(func(tx *buntdb.Tx) error {
				value, _ = tx.Get(name)
				return nil
			})

			if err != nil {
				fmt.Println(err)
			}

			if value != "" && value != ssh {
				fmt.Println("Host already exists")
				os.Exit(1)
			}

			err = db.Update(func(tx *buntdb.Tx) error {
				_, _, err := tx.Set(name, ssh, nil)
				return err
			})

			if err != nil {
				fmt.Println(err)
			}

		case name != "" && ssh == "":
			var value string
			err := db.View(func(tx *buntdb.Tx) error {
				value, _ = tx.Get(name)
				return nil
			})

			if err != nil {
				fmt.Println(err)
			}

			if value == "" {
				fmt.Println("Host does not exist")
				os.Exit(1)
			}

			ssh = value

		case name == "" && ssh != "":
			var value string
			err := db.View(func(tx *buntdb.Tx) error {
				value, _ = tx.Get(ssh)
				return nil
			})
			if err != nil {
				fmt.Println(err)
			}
			if value == "" {
				db.Update(func(tx *buntdb.Tx) error {
					_, _, err := tx.Set(ssh, ssh, nil)
					return err
				})
			}
			ssh = value

		default:
			fmt.Println("No host specified")
			os.Exit(1)

		}
		c := exec.Command("ssh", ssh)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Run()

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("name", "", "Name of the host")
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(deleteCmd)
}
