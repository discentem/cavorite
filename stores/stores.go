package stores

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/user"
)

type Options struct {
	CreatePantri   *bool   `json:"create_pantri"`
	RemoveFromRepo *bool   `json:"remove_from_repo"`
	Type           *string `json:"type"`
}

func StorageType() (*string, error) {
	usr, _ := user.Current()
	dir := usr.HomeDir
	b, err := ioutil.ReadFile(fmt.Sprintf("%s/.pantri/config", dir))
	if err != nil {
		return nil, err
	}
	var o struct {
		Options `json:"options"`
	}
	if err := json.Unmarshal(b, &o); err != nil {
		return nil, err
	}
	return o.Type, nil

}
