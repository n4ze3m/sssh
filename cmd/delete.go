package cmd

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/tidwall/buntdb"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete host",
	Long:  `Delete host`,
	Run: func(cmd *cobra.Command, args []string) {
		db, _ := buntdb.Open("database.db")

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
			fmt.Printf("Prompt failed %v \n", err)

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