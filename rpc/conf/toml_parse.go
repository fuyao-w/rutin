package conf

import (
	"github.com/BurntSushi/toml"
	"log"
)

func Decode(path string, obj interface{}) {
	_, err := toml.DecodeFile(path, obj)
	log.Fatalf("decode err %s", err.Error())
}
