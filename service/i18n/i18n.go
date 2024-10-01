package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

//go:embed locales
var LocaleFS embed.FS

type Lang map[string]map[string]string

type LangMap map[string]Lang

func LoadLocaleFiles() LangMap {
	files, err := LocaleFS.ReadDir("locales")
	if err != nil {
		log.Fatal(err)
	}
	langMap := LangMap{}
	for _, file := range files {
		var lang Lang
		f, err := LocaleFS.ReadFile("locales/" + file.Name())
		if err != nil {
			fmt.Println(err)
		}

		fileName := file.Name()
		if pos := strings.LastIndexByte(fileName, '.'); pos != -1 {
			fileName = fileName[:pos]
		}

		if err := json.Unmarshal(f, &lang); err != nil {
			fmt.Println(err)
		}
		langMap[fileName] = lang
	}
	return langMap
}
