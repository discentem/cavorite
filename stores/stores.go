package stores

type PantriConfig struct {
	Type string `json:"type"`
}

type Options struct {
	RemoveFromSourceRepo *bool
}
