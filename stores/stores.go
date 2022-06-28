package stores

type Options struct {
	RemoveFromSourceRepo *bool `json:"remove_from_sourcerepo"`
}

type Store interface {
	Upload(objects []string) error
	Retrieve(objects []string) error
}
