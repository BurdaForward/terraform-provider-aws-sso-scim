# AWS SSO SCIM Terraform provider

This repository holds a terraform provider enabling you to use the SCIM provisioning interface of the AWS SSO service. Using this
provider you are able to provision users and groups within AWS SSO when using a external identity source.

This code is currently not in a production ready state. Use on your own risk.

This code provides the following:

- [ ] Data
  - [ ] `aws_sso_scim_user`
  - [ ] `aws_sso_scim_group`
- [ ] Resources
  - [ ] `aws_sso_scim_user`
  - [ ] `aws_sso_scim_group`
  - [ ] `aws_sso_scim_group_member`


## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
-	[Go](https://golang.org/doc/install) >= 1.15

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:
```sh
$ go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

Fill this in for each provider

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```
