package keys

import (
	"errors"
	"log"
	"path/filepath"
	"strings"

	"github.com/magiconair/properties"
)

var (
	keysInstance *switchKeys
	UseKey       bool = true
)

type switchKeys struct {
	keys map[string]string
}

func (k *switchKeys) GetKey(keyName string) string {
	return k.keys[keyName]
}

func SwitchKeys() (*switchKeys, error) {
	return keysInstance, nil
}

func InitSwitchKeys(baseFolder string) (*switchKeys, error) {
	var (
		path string
		p    *properties.Properties
		err  error
	)

	// first, try to read the prod keys from the settings value
	if baseFolder != "" {
		path = baseFolder
		if !strings.HasSuffix(path, ".keys") {
			path = filepath.Join(path, "prod.keys")
		}

		log.Printf("Trying to load prod.keys based on settings.json: %v \n", path)
		p, err = properties.LoadFile(path, properties.UTF8)
	} else {
		err = errors.New("prod.keys not defined in settings.json")
	}

	// second, if not found by settings look into the current folder
	if err != nil {
		path = filepath.Join(baseFolder, "prod.keys")

		log.Printf("Trying to load prod.keys based on current folder: %v \n", path)
		p, err = properties.LoadFile(path, properties.UTF8)
	}

	// third, if not found in current, look in home directory
	if err != nil {
		path = "${HOME}/.switch/prod.keys"

		log.Printf("Trying to load prod.keys based on home directory: %v \n", path)
		p, err = properties.LoadFile(path, properties.UTF8)
	}

	if err != nil {
		log.Println("Unable to find prod.keys")
		return nil, errors.New("Error trying to read prod.keys [reason:" + err.Error() + "]")
	}

	keysInstance = &switchKeys{keys: map[string]string{}}
	for _, key := range p.Keys() {
		value, _ := p.Get(key)
		keysInstance.keys[key] = value
	}

	log.Printf("Loaded prod.keys from: %v \n", path)
	return keysInstance, nil
}
