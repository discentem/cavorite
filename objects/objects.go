package objects

import (
	"errors"
	"fmt"
	"strings"
)

type KeyModifier interface {
	Modify(string) string
	Original(string) string
}

type AddPrefixToKey struct {
	Prefix string
}

var (
	_ = KeyModifier(AddPrefixToKey{})
)

func (p AddPrefixToKey) Modify(original string) string {
	if p.Prefix == "" {
		return original
	}
	return fmt.Sprintf("%s/%s", p.Prefix, original)
}

func (p AddPrefixToKey) Original(modified string) string {
	if p.Prefix == "" {
		return modified
	}
	prefix := fmt.Sprintf("%s/", p.Prefix)
	return strings.Replace(modified, prefix, "", 1)
}

func ModifyMultipleKeys(modder KeyModifier, originalKeys ...string) ([]string, error) {
	if modder == nil {
		return []string{}, errors.New("modder cannot be nil")
	}
	newKeys := []string{}
	for _, key := range originalKeys {
		newKeys = append(newKeys, modder.Modify(key))
	}
	return newKeys, nil
}
