# Terraform Provider CloudConnexa

<a href="https://www.terraform.io/" target="_blank">
  <img align="right" src="https://upload.wikimedia.org/wikipedia/commons/thumb/0/04/Terraform_Logo.svg/2560px-Terraform_Logo.svg.png" alt="Terraform" width="120px">
</a>

<a href="https://anna.money/?utm_source=terraform&utm_medium=referral&utm_campaign=docs" target="_blank">
  <img align="right" src="https://upload.wikimedia.org/wikipedia/commons/a/aa/ANNA_Money_Logo_PNG.png" alt="ANNA Money" width="80px">
</a>

<a href="https://openvpn.net/cloud-vpn/?utm_source=terraform&utm_medium=docs" target="_blank">
  <img align="right" src="https://upload.wikimedia.org/wikipedia/commons/thumb/f/f5/OpenVPN_logo.svg/2560px-OpenVPN_logo.svg.png" alt="OpenVPN" width="150px">
</a>

- [Website CloudConnexa](https://openvpn.net/cloud-vpn/?utm_source=terraform&utm_medium=docs)
- [Terraform Registry](https://registry.terraform.io/providers/OpenVPN/cloudconnexa/latest)

## Description

The Terraform provider for [CloudConnexa](https://openvpn.net/cloud-vpn/?utm_source=terraform&utm_medium=docs) allows teams to configure and update CloudConnexa project parameters via their command line.

## Guides

- [Migration from v0.X.X to v1.0.0](https://registry.terraform.io/providers/OpenVPN/cloudconnexa/latest/docs/guides/migration-to-v1)

## Maintainers

This provider plugin is maintained by:

- OpenVPN team at [CloudConnexa](https://openvpn.net/cloud-vpn/?utm_source=terraform&utm_medium=docs)
- SRE Team at [ANNA Money](https://anna.money/?utm_source=terraform&utm_medium=referral&utm_campaign=docs) / [GitHub ANNA Money](http://github.com/anna-money/)
- [@patoarvizu](https://github.com/patoarvizu)

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.12.x
- [Go](https://golang.org/doc/install) 1.18 (to build the provider plugin)

## Building The Provider

Clone repository to: `$GOPATH/src/github.com/OpenVPN/terraform-provider-cloudconnexa`

```sh
mkdir -p $GOPATH/src/github.com/OpenVPN; cd $GOPATH/src/github.com/OpenVPN
git clone git@github.com:OpenVPN/terraform-provider-cloudconnexa.git
```

Enter the provider directory and build the provider

```sh
cd $GOPATH/src/github.com/OpenVPN/terraform-provider-cloudconnexa
make build
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.18+ is _required_). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
make bin
...
$GOPATH/bin/terraform-provider-cloudconnexa
...
```

In order to test the provider, you can simply run `make test`.

```sh
make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

_Note:_ Acceptance tests create real resources, and often cost money to run.

```sh
make testacc
```

_**Please note:** This provider, like CloudConnexa API, is in beta status. Report any problems via issue in this repo._
