package main

import (
	"fmt"
	"os"

	"github.com/howeyc/gopass"
	"github.com/jawher/mow.cli"

	"github.com/milk/milkpasswd"
)

func setupCli() error {
	app := cli.App("milkpasswd", "Store your passwords in a safe place")

	app.Command("add", "Store new credentials", func(cmd *cli.Cmd) {
		path := cmd.StringOpt("p path", "/", "Storage location")
		name := cmd.StringArg("NAME", "", "Credential name")
		username := cmd.StringArg("USERNAME", "", "Your username")
		website := cmd.StringArg("WEBSITE", "", "Website URL")
		description := cmd.StringArg("DESC", "", "Description")

		cmd.Spec = "[-p] NAME USERNAME [WEBSITE [DESC]]"

		cmd.Action = func() {
			fmt.Printf("Master password: ")
			master := gopass.GetPasswd()

			fmt.Printf("Cipher text: ")
			cipher := gopass.GetPasswd()

			fmt.Printf("Password to store: ")
			password := gopass.GetPasswd()

			key := milkpasswd.Sha256sum(master)
			ciphertext := milkpasswd.Md5sum(cipher)

			encrypted, err := milkpasswd.Encrypt(key, ciphertext, password)

			entry, err := milkpasswd.CreateEntry(*path,
				*name,
				*username,
				encrypted,
				*website,
				*description)
			if err != nil {
				panic(err)
			}
			entry.Save()
		}
	})

	app.Command("get", "Retrieve stored credentials", func(cmd *cli.Cmd) {
		path := cmd.StringArg("PATH", "", "Credential path")

		cmd.Spec = "PATH"

		cmd.Action = func() {
			entry, err := milkpasswd.GetEntry(*path)
			if err != nil {
				panic(err)
			}

			if entry == nil {
				fmt.Println("Entry not found!")
				return
			}

			fmt.Printf("Master password: ")
			master := gopass.GetPasswd()

			fmt.Printf("Cipher text: ")
			cipher := gopass.GetPasswd()

			key := milkpasswd.Sha256sum(master)
			ciphertext := milkpasswd.Md5sum(cipher)

			decrypted, err := milkpasswd.Decrypt(key, ciphertext, entry.Password)

			fmt.Println(entry.String())
			fmt.Println("Password:", decrypted)
		}
	})

	app.Command("del", "Remove stored credentials", func(cmd *cli.Cmd) {
		path := cmd.StringArg("PATH", "", "Credential path")

		cmd.Spec = "PATH"

		cmd.Action = func() {
			err := milkpasswd.DeleteEntry(*path)
			if err != nil {
				panic(err)
			}
		}
	})

	app.Command("list", "List the stored credentials", func(cmd *cli.Cmd) {
		search := cmd.StringArg("SEARCH", "", "Search string")

		cmd.Spec = "[SEARCH]"

		cmd.Action = func() {
			var data map[string]*milkpasswd.Entry
			var err error
			if *search == "" {
				data, err = milkpasswd.ListEntries()
			} else {
				data, err = milkpasswd.SearchEntries(*search)
			}
			if err != nil {
				panic(err)
			}

			if len(data) == 0 {
				fmt.Println("No matching credentials found.")
				return
			}

			for k, v := range data {
				fmt.Println("Key:", k, "Value:", v.String())
			}
		}
	})

	return app.Run(os.Args)
}
