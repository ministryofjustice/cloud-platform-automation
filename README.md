# cloud-platform-automation

This repository contains library of scripts user to automate the cloud platform tasks

## Contents
- [utils](utils/): is a collection of utility scripts that can be used to automate common tasks.
- [commit-check](commit-check/): is a tool that checks list of files in a commit and validates them against a set of rules. The file must be the only file in the commit or it like fail the check.
- [annotation-check](annotation-check/): is a tool that checks the annotations namespace.yaml files and validates them against a set of rules.
- [iam-policy-usage](iam-policy-usage/): is a tool that checks the IAM policies in the AWS account and checking for last used date, policy associations, and tagged and untagged resources.


## GO Modules
- The repository uses Go modules for dependency management.
- The go.mod file contains the list of dependencies and their versions.
- The go.sum file contains the checksums of the dependencies.
- There is one go.mod file at the root of the repository. All the scripts use the same go.mod file to manage dependencies.

## Release Process
- Binaries are built and released using the [GoReleaser](https://goreleaser.com/) tool. 
- The release process is automated using GitHub Actions. 
- The release process is triggered by creating a new tag in the repository. The tag should be in the format `X.Y.Z` where `X.Y.Z` is the version number.
- The release process will build and release the binaries for the tag and create a GitHub release with the binaries attached.
- The release process will also update the `latest` tag to point to the new release.
- All the scripts use the same Dockerfile to build the binaries. The Dockerfile is located in the root of the repository. It will be populated with the necessary tools and dependencies to build the binary, this is done by the `.goreleaser.yaml` file. that sits in the each of the script directories.

