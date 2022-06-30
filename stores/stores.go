package stores

type PantriConfig struct {
	Type string `json:"type"`
}

type Options struct {
	// TODO(discentem) remove this option. See #15
	RemoveFromSourceRepo *bool `json:"remove_from_sourcerepo"`
}
