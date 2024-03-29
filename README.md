# Terraform Provider for Megaport

- [![Build Status](https://travis-ci.org/utilitywarehouse/terraform-provider-megaport.svg?branch=master)](https://travis-ci.org/utilitywarehouse/terraform-provider-megaport)

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.12+
- [Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)

## Credentials

To use the Megaport API, this Provider requires that a Megaport API access token
is provided. To simplify the process of obtaining a new access token, there is a
utility tool you can use. It is recommended that the token is provided via the
`MEGAPORT_TOKEN` environment variable.

To retrieve a new token for the megaport api and export it as a variable:
```sh
$ export $(make reset-token)
```

Alternatively, you can use the helper tool directly:
```sh
$ cd util/megaport_token
$ go run .
```
To revoke the current token (and get a new one) you can pass the `--reset` flag.

By default, the above script will obtain a token to use with megaport staging
api environment. You can set `MEGAPORT_ENDPOINT` environment variable to specify
an alternative endpoint. For example, to get a token for megaport production
api, you can use the above target as:
```sh
MEGAPORT_ENDPOINT=https://api.megaport.com make reset-token
```
## Developing the Provider

If you wish to work on the provider, you'll first need
[Go](http://www.golang.org) installed on your machine (please check the
[requirements](#requirements) before proceeding).

```sh
$ git clone https://github.com/utilitywarehouse/terraform-provider-megaport.git
...
$ cd terraform-provider-megaport
```

To compile the provider, run `make build`. This will build the provider and put
the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-megaport
...
```

## Using the Provider

You can simply use the provider published in the Terraform Registry [here](https://registry.terraform.io/providers/utilitywarehouse/megaport/latest)
to manage your terraform resources.

To use a custom-built provider in your Terraform environment (e.g. the provider
binary from the build instructions above), follow the instructions to
[install it as a plugin](https://www.terraform.io/docs/plugins/basics.html#installing-plugins).
After placing the custom-built provider into your plugins directory,
run `terraform init` to initialize it.

### Provider Documentation

To browse the documentation, you can simply `make website` to serve it locally.

```sh
$ make website
...
```

Additionally, there are a number of templated examples (used in acceptance
testing) inside the `examples/` directory.

## Testing the Provider

In order to test the provider, you can run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run. By
default, acceptance tests for this provider are run in the staging Megaport
environment which does not incur any costs.

```sh
$ make testacc
```

## Contributing

Terraform is the work of thousands of contributors. We appreciate your help!

To contribute, please read the contribution guidelines:
[Contributing to Terraform - Megaport Provider](.github/CONTRIBUTING.md)

GitHub issues are intended for bugs or feature requests related to the Megaport
provider codebase. See https://www.terraform.io/docs/extend/community/index.html
for a list of community resources to ask questions about Terraform.

