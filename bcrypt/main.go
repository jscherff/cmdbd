package main

import (
	`fmt`
	`os`
	`path/filepath`
	`golang.org/x/crypto/bcrypt`
)

func main() {

	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <word> [<word> ...]\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	for _, word := range os.Args[1:] {

		if hash, err := bcrypt.GenerateFromPassword([]byte(word), bcrypt.DefaultCost); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("%s: %s\n", word, string(hash))
		}
	}
}
