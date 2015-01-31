// milkpasswd is a password manager written in Go
package main

func main() {
	err := setupCli()
	if err != nil {
		panic(err)
	}
}
