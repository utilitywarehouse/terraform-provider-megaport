---
layout: "megaport"
page_title: "Provider: Megaport"
description: |-
  The Megaport provider is used to interact with the various Megaport resources.
---

# Megaport Provider

The Megaport provider is used to interact with the various Megaport resources.
The provider needs to be configured with a valid Megaport API token.

## Example Usage

It is recommended that the API token is provided through an environment variable
(`MEGAPORT_TOKEN`) and not hardcoded in the provider configuration.

```hcl
provider "megaport" {}
```

## Argument Reference

* `token` - (Optional) This is the Megaport API token. It must be provided but
it can also be sourced from the `MEGAPORT_TOKEN` environment variable.

* `api_endpoint` - (Optional) This is the Megaport API endpoint. It can also be
sourced from the `MEGAPORT_API_ENDPOINT` environment variable and can be used to
point the provider to an alternative Megaport environment. It defaults to the
production environment and is primarily used for testing.

