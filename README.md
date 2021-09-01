# terraform-provider-bridgecrew

First you need to add your secret - the Bridgecrew API key, as an env var,
BRIDGECREW_API, or it won't work.

If this is you first time using this provider you will need to build
and then run with:

```bash
make check
```

This will build and install the provider locally, and run a test template.

If you're not using a Mac you will have to change OS_ARCH=darwin_amd64 to your value.

The example tf gets all the repositories you have in Bridgecrew and lists them.

Once installed you can use the provider via the normal Terraform workflow:

```bash
terraform init
terraform plan

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
 Currently there uis only support for 2 data sources:

- bridgecrew_repositories
- bridgecrew_suppressions
- bridgecrew_policies

export TF_LOG_CORE=""
and
export TF_LOG_PROVIDER="DEBUG"

## Building The Documentation

The documentation is built from components (go templates) stored in the `templates` folder.
Building the documentation copies the full markdown into the `docs` folder, ready for deployment to Hashicorp.

> NOTE: you'll need the [`tfplugindocs`](https://github.com/hashicorp/terraform-plugin-docs) tool for generating the Markdown to be deployed to Hashicorp. For more information on generating documentation, refer to https://www.terraform.io/docs/registry/providers/docs.html

- To validate the `/template` directory structure:

```bash
make validate-docs
```

- To build the `/docs` documentation Markdown files:

```bash
make generate-docs
```

- To view the documentation:
Paste `/docs` Markdown file content into https://registry.terraform.io/tools/doc-preview
