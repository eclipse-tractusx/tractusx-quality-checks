# Eclipse Tractus-X quality checks

The [Eclipse Tractus-X](https://projects.eclipse.org/projects/automotive.tractusx) quality checks is an automation effort
to ensure basic quality alignment across [Tractus-X OSS products](https://github.com/eclipse-tractusx/).

The checks will be aligned with the [Tractus-X release guidelines](https://eclipse-tractusx.github.io/docs/release) and
are available as library to use in multiple usecases. These could be:

- PR checks to test, if contributions follow release guidelines
- Dashboards to show overall alignment of products
- Binaries to check code locally (manual run, pre-commit hooks, ...)
- ...

## Local build

If you want to build the command locally, you just need Golang version 1.20 and run `go build` on the root level of
this repository.

## Custom GitHub action

The `tractusx-quality-checks` can be run as a GitHub action. This action is build as a
[composite action](https://docs.github.com/en/actions/creating-actions/about-custom-actions#types-of-actions).
This approach was chosen, because it allows us to run our Golang code without the need to maintain a container image
for it.

```yaml
...
jobs:
  quality-check:
    steps:
      # Checkout your repo, since the quality action is working on a local copy
      - uses: actions/checkout@v3

      # Setup go as a prerequisite for the quality-check action
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: 1.20

      # Use the quality check aciton
      - name: Run quality check command
        uses: eclipse-tractusx/tractusx-quality-checks@v1
```