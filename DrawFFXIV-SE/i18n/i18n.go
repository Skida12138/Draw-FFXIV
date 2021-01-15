package i18n

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/skida12138/drawffxiv-se/utils"
)

var msgMap map[string]string

// SetLang : set language by language code
func SetLang(lang string) error {
	return utils.NewResult(
		os.Open(fmt.Sprintf("./data/i18n/%s.json", lang)),
	).AndThen(func(msgFile interface{}) (interface{}, error) {
		return ioutil.ReadAll(msgFile.(*os.File))
	}).AndThen(func(msgBytes interface{}) (interface{}, error) {
		msgMap = make(map[string]string)
		return nil, json.Unmarshal(msgBytes.([]byte), &msgMap)
	}).Error()
}

// Msg : return message by key and corresponding language
func Msg(key string) string {
	return msgMap[key]
}
