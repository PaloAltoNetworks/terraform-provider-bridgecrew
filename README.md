# terraform-provider-bridgecrew

[![Maintained by Bridgecrew.io](https://img.shields.io/badge/maintained%20by-bridgecrew.io-blueviolet)](https://bridgecrew.io/?utm_source=github&utm_medium=organic_oss&utm_campaign=terraform-provider-bridgecrew)
[![release](https://github.com/PaloAltoNetworks/terraform-provider-bridgecrew/actions/workflows/release.yml/badge.svg)](https://github.com/PaloAltoNetworks/terraform-provider-bridgecrew/actions/workflows/release.yml)
[![slack-community](https://img.shields.io/badge/slack-bridgecrew-blueviolet.svg?logo=slack)](https://codifiedsecurity.slack.com/)
[![Go Report Card](https://goreportcard.com/badge/github.com/PaloAltoNetworks/terraform-provider-bridgecrew)](https://goreportcard.com/report/github.com/PaloAltoNetworks/terraform-provider-bridgecrew)
[![Go Reference](https://pkg.go.dev/badge/github.com/PaloAltoNetworks/terraform-provider-bridgecrew.svg)](https://pkg.go.dev/github.com/jameswoolfenden/terraform-provider-bridgecrew)
[![GitHub All Releases](https://img.shields.io/github/downloads/PaloAltoNetworks/terraform-provider-bridgecrew/total)](https://github.com/PaloAltoNetworks/terraform-provider-bridgecrew/releases)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/PaloAltoNetworks/terraform-provider-bridgecrew)
![GitHub issues](https://img.shields.io/github/issues/PaloAltoNetworks/terraform-provider-bridgecrew)

This guide is to help you develop/debug the Terraform Bridgecrew provider, to get started you need to obtain and add your secret - your Bridgecrew API key, as an env var,
BRIDGECREW_API, or it won't work.

First obtain your API key here: <https://www.bridgecrew.cloud/integrations/api-token>

If this is your first time using this provider you will need to build
and then run it with:

```bash
make check
```

This will build and install the provider locally, and run a test template.

If you're not using a Mac you will have to change OS_ARCH=darwin_amd64 to the value for your platform.

Terraform examples live in a sub-folder Terraform:
The example tf gets all the repositories you have in Bridgecrew and lists them, also included it's a sample policy that can be created, updated and destroyed by the platform.

Once installed you can use the Provider via the normal Terraform workflow:

```bash
terraform init
terraform plan
terraform apply
...

Changes to Outputs:
  + repos       = {
      + id           = "1627304954"
      + repositories = [
          + {
              + creationdate  = "2021-05-19T06:23:36.966Z"
              + defaultbranch = "master"
              + id            = "d56e6193-82b7-44ce-ba5f-2751bedc3842"
              + ispublic      = false
              + owner         = "JamesWoolfenden"
              + repository    = "shift-left"
              + source        = "Github"
            },
          + {
```

The Terraform config is in main.tf.
Currently, there is only support for 8 data sources:

- bridgecrew_apitokens
- bridgecrew_errors
- bridgecrew_integrations
- bridgecrew_policies
- bridgecrew_repositories
- bridgecrew_repository_branches
- bridgecrew_suppressions
- bridgecrew_users

and two resources:

- bridgecrew_policy
- bridgecrew_simple_policy

More will follow.

## Examples

For more detailed examples see:  <https://github.com/JamesWoolfenden/terraform-bridgecrew-examples>, each example has a video for you to follow.
There is also a published module that uses the Provider here: <https://registry.terraform.io/modules/JamesWoolfenden/simplepolicy/bridgecrew/latest>.

## Debugging

To see the debug output for a provider set:

```bash
export TF_LOG_CORE=""
```

and

```bash
export TF_LOG_PROVIDER="DEBUG"
```

I have added the ability to debug the provider,
![debug options](debug.png "debug options")

Run in debug (with some breakpoints) and the output shows:

```shell
 TF_REATTACH_PROVIDERS='{"jameswoolfenden/dev/bridgecrew":{"Protocol":"grpc","ProtocolVersion":5,"Pid":68437,"Test":true,"Addr":{"Network":"unix","String":"/var/folders/pd/vb2p7hln0gj15hjkn6g31_j00000gn/T/plugin752084824"}}}'
```

 export that to your shell, and use Terraform in debug!

## Building The Documentation

The documentation is built from components (go templates) stored in the `templates` folder.
Building the documentation, copies the full markdown into the `docs` folder, ready for deployment to Hashicorp.

> NOTE: you'll need the [`tfplugindocs`](https://github.com/hashicorp/terraform-plugin-docs) tool for generating the Markdown to be deployed to Hashicorp. For more information on generating documentation, refer to https://www.terraform.io/docs/registry/providers/docs.html

- To validate the `/template` directory structure:

```bash
make validate-docs
```

- To build the `/docs` documentation Markdown files:

```bash
make docs
```

If you add new resources you will need to add a template for it in the template folder and update *scripts/generate-docs.go*, once built you will need to add the generated markdown file.

- To view the documentation:
The provider has online documentation here:<https://registry.terraform.io/providers/PaloAltoNetworks/bridgecrew/latest/docs>
If you want to preview your modified docs you can paste your `/docs` folder Markdown file content into <https://registry.terraform.io/tools/doc-preview>

## Contributing

The repository uses the pre-commit framework to format and test code prior to checkin, pre-commit is installed via pip and then the config is installed (from the root)after you initially clone the repo:

```bash
git clone git@github.com:bridgecrewio/terraform-provider-bridgecrew.git
pip3 install pre-commit
pre-commit install
```

For details on the hooks used see the config: .pre-commit-config.yaml.

## Building a release

This repository uses GitHub actions in conjunction with goreleaser, pushing a tag will invoke a matrix build of goreleaser.

## Checkov/Bridgecrew

The Terraform you create for this provider is already supported by Checkov and the Bridgecrew platform.
For example If you run the Checkov cli over this repository, you'll see there's a security check on the bridgecrew provider:

```bash
checkov -d .

       _               _
   ___| |__   ___  ___| | _______   __
  / __| '_ \ / _ \/ __| |/ / _ \ \ / /
 | (__| | | |  __/ (__|   < (_) \ V /
  \___|_| |_|\___|\___|_|\_\___/ \_/

By bridgecrew.io | version: 2.0.413

terraform scan results:

Passed checks: 1, Failed checks: 0, Skipped checks: 0

Check: CKV_BCW_1: "Ensure no hard coded API token exist in the provider"
    PASSED for resource: bridgecrew.default
    File: /terraform/provider.bridgecrew.tf:1-4

```
