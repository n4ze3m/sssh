package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/tidwall/buntdb"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all hosts",
	Long:  `List all hosts`,
	Run: func(cmd *cobra.Command, args []string) {
		db, _ := buntdb.Open("database.db")

		defer db.Close()

		var hosts []string
		var labels []string

		db.View(func(tx *buntdb.Tx) error {
			tx.Ascend("", func(key, value string) bool {
				hosts = append(hosts, value)
				labels = append(labels, fmt.Sprintf("%s: %s", key, value))
				return true
			})
			return nil
		})

		prompt := promptui.Select{
			Label: "Select host",
			Items: labels,
		}

		_, result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v \n", err)

		}

		index := 0
		for i, host := range hosts {
			if host == strings.Split(result, ": ")[1] {
				index = i
			}
		}

		c := exec.Command("ssh", hosts[index])

		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Run()
	},
}
