package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/tidwall/buntdb"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete host",
	Long:  `Delete host`,
	Run: func(cmd *cobra.Command, args []string) {

		home, _ := homedir.Dir()
		db, _ := buntdb.Open(home + "/ssh-manager.db")

		defer db.Close()

		var labels []string

		db.View(func(tx *buntdb.Tx) error {
			tx.Ascend("", func(key, value string) bool {
				labels = append(labels, fmt.Sprintf("%s: %s", key, value))
				return true
			})
			return nil
		})

		prompt := promptui.Select{
			Label: "Select host to delete",
			Items: labels,
		}

		_, result, err := prompt.Run()

		if err != nil {
			os.Exit(1)
		}

		index := 0
		for i, host := range labels {
			if host == strings.Split(result, ": ")[1] {
				index = i
			}
		}

		db.Update(func(tx *buntdb.Tx) error {
			key := strings.Split(labels[index], ": ")[0]
			tx.Delete(key)
			return nil
		})

		fmt.Println("Host deleted")

	},
}
