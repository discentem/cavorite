package stores

type Store interface {
	Upload(sourceRepo string, objects ...string) error
	Retrieve(sourceRepo string, objects ...string) error
}
type Options struct {
	// TODO(discentem) remove this option. See #15
	RemoveFromSourceRepo *bool `json:"remove_from_sourcerepo" mapstructure:"remove_from_sourcerepo"`
}
