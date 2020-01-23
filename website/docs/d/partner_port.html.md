---
layout: "megaport"
subcategory: "datasources"
page_title: "Megaport: megaport_partner_port"
description: |-
  Get information on a Megaport Partner Port resource.
---

# Data Source: megaport_partner_port

Use this datasource to retrieve the uid of a Megaport Partner Port from the
Megaport Markeplace for use in other resources.

## Example Usage

```hcl
data "megaport_partner_port" "foo" {
  name_regex   = "eu-west-1"
  connect_type = "AWS"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Required, Forces new resource) A regex string filter to apply
to the Port list returned by Megaport.
* `connect_type` - (Optional, Forces new resource) The type of Partner Port is
a filter applied to the Port list returned by Megaport. Supported values:
`"AWS"`.
* `location_id` - (Optional, Forces new resource) Limit the Port search by a
location id.
* `vxc_permitted` - (Optional, Default: `true`) Limit search to Ports that have
the `vxcPermitted` flag set. This is true by default since Megaport will only
accept VXCs to these Ports.

~> **Note:** If more or less than a single match is returned by the search, Terraform will
fail. Ensure that your search is specific enough to return a single Port.

## Attribute Reference

The `id` of the datasource is set to the uid of the found Port.
