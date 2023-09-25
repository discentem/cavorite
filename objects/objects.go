package objects

import (
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
	prefix := fmt.Sprintf("%s/", p.Prefix)
	return strings.ReplaceAll(modified, prefix, "")
}

func ModifyMultipleKeys(modder KeyModifier, originalKeys ...string) []string {
	newKeys := []string{}
	for _, key := range originalKeys {
		newKeys = append(newKeys, modder.Modify(key))
	}
	return newKeys
}
