package initialize

import (
	"fmt"

	"github.com/discentem/pantri_but_go/stores"
	localstore "github.com/discentem/pantri_but_go/stores/local"
)

func Initalize(sourceRepo, backend, address string, opts stores.Options) error {
	if backend != "local" {
		return fmt.Errorf("%s is not supported. only %q is currently supported", backend, "local")
	}
	_, err := localstore.NewWithOptions(sourceRepo, address, opts)
	if err != nil {
		return err
	}
	return nil

}
