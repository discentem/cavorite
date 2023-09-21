package objects

import "fmt"

type KeyModifier interface {
	Modify(string) string
}

type AddPrefixToKey struct {
	Prefix string
}

func (p AddPrefixToKey) Modify(original string) string {
	if p.Prefix == "" {
		return original
	}
	return fmt.Sprintf("%s/%s", p.Prefix, original)
}

func ModifyMultipleKeys(modder KeyModifier, originalKeys ...string) []string {
	newKeys := []string{}
	for _, key := range originalKeys {
		newKeys = append(newKeys, modder.Modify(key))
	}
	return newKeys
}
