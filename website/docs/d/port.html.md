---
layout: "megaport"
subcategory: "datasources"
page_title: "Megaport: megaport_port"
description: |-
  Get information on a Megaport Port resource.
---

# Data Source: megaport_port

Use this datasource to retrieve the uid of a Megaport Port for use in other
resources.

## Example Usage

```hcl
data "megaport_port" "foo" {
  name_regex = "foobar"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Required, Forces new resource) A regex string filter to apply
to the Port list returned by Megaport.

~> **Note:** If more or less than a single match is returned by the search, Terraform will
fail. Ensure that your search is specific enough to return a single Port.

## Attribute Reference

The `id` of the datasource is set to the uid of the found Port.
