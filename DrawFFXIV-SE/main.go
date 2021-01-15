package main

import (
	"log"

	"github.com/skida12138/drawffxiv-se/i18n"
	_ "github.com/skida12138/drawffxiv-se/routes"
)

func main() {
	if err := i18n.SetLang("zh-cmn-Hans"); err != nil {
		log.Fatal("fatal error occurs while setting languages")
		log.Fatal(err)
	}
}
