---
layout: "megaport"
subcategory: "datasources"
page_title: "Megaport: megaport_location"
description: |-
  Get information on a Megaport location.
---

# Data Source: megaport_location

Use this datasource to retrieve the id of a Megaport location for use in other
resources.

## Example Usage

```hcl
data "megaport_location" "foo" {
  name_regex = "foobar"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Required, Forces new resource) A regex string filter to apply
to the location list returned by Megaport.

~> **Note:** If more or less than a single match is returned by the search,
Terraform will fail. Ensure that your search is specific enough to return a
single location.

## Attribute Reference

The `id` of the datasource is set to the id of the found location.
