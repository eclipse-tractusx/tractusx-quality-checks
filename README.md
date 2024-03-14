# Eclipse Tractus-X quality checks

[![status: archive](https://github.com/GIScience/badges/raw/master/status/archive.svg)](https://github.com/GIScience/badges#archive)

This is content that is moved to another repository. The new location is [sig-release/release-automation](https://github.com/eclipse-tractusx/sig-release/tree/main/release-automation)

>❗**Note**❗
> This repository is currently not under active maintenance and therefore archived. If you plan on picking up this product again, feel free to reach out ot the Eclipse Tractus-X projects leads via [mailing list](https://accounts.eclipse.org/mailing-list/tractusx-dev).

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

## Git pre-commit hook

It is possible to add quality checks as git pre-commit hook in the repository which will run checks triggered by git commit command. To achieve that please follow below steps:

1. Download latest version of the tractusx-quality-checks tool for your platform under Asset section at https://github.com/eclipse-tractusx/tractusx-quality-checks/releases/latest
2. Add executable permission if necessary and copy the binary (changing file name to tractusx-quality-checks) into one of the $PATH location depending on your OS.

#### macOS example:

```
$ chmod +x tractusx-quality-checks-0.8.0-darwin-arm64
$ sudo cp tractusx-quality-checks-0.8.0-darwin-arm64 /usr/local/bin/tractusx-quality-checks
```

3. Inside your repo copy below to .git/hooks/pre-commit:

```
    #!/usr/bin/env bash
    # ^ Note the above "shebang" line. This says "This is an executable shell script"
    # Name this script "pre-commit" and place it in the ".git/hooks/" directory

    # Exit immediately with that command's exit status if the command fails
    set -eo pipefail

    # Run Tractus-X quality checks
    tractusx-quality-checks checkLocal
    echo -e "\n\n--- TRG checks passed! ---\n\n"
```

5. Make .git/hooks/pre-commit executable.
6. Pre-commit setup is complete and initiates checks each time git commit is ran. 

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