---
layout: "megaport"
subcategory: "resources"
page_title: "Megaport: megaport_port"
description: |-
  Provides a Megaport port resource.
---

# Resource: megaport_port

Provides a Megaport port resource. Allows ports to be created, updated and
deleted.

## Example Usage

```hcl
data "megaport_location" "foo" {
  name_regex = "foobar"
}

resource "megaport_port" "foo" {
  name        = "foo"
  location_id = data.megaport_location.foo.id
  speed       = 1000
  term        = 1
}
```

## Argument Reference

The following arguments are supported:

* `location_id` - (Required, Forces new resource) The numeric id of the location
where this port should be created.
* `name` - (Required) The name of the port.
* `speed` - (Required, Forces new resource) The speed of the port (`1000`,
`10000`, or `100000` Mbps, subject to availability).
* `term` - (Required, Forces new resource) Length of the contract (`1`, `12`,
`24` or `36` months).
* `invoice_reference` - (Optional) Used for billing purposes, a reference to
this specific line item.
* `marketplace_visibility` - (Optional, Default: `"private"`) Whether this port
will be listed on the Megaport Marketplace.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique product id of the port.
* `associated_vxcs` - A list of all the VXCs associated with this port.
