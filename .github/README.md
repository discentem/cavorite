# GitHub Actions

## Linting

`linter.yaml` runs [Super-Linter][sl] against the code base on pull request, push, and when manually triggered.

Most relevantly, the [golangci-lint][gll] linter is used for Go.

[sl]: <https://github.com/marketplace/actions/super-linter>
[gll]: <https://github.com/golangci/golangci-lint>
