---
layout: "megaport"
subcategory: "resources"
page_title: "Megaport: megaport_mcr"
description: |-
  Provides a Megaport MCR resource.
---

# Resource: megaport_mcr

Provides a Megaport Cloud Router (MCR) resource. Allows MCRs to be created,
updated and deleted.

## Example Usage

```hcl
data "megaport_location" "foo" {
  name_regex = "foobar"
}

resource "megaport_mcr" "foo" {
  name        = "foo"
  location_id = data.megaport_location.foo.id
  speed       = 1000
}
```

## Argument Reference

The following arguments are supported:

* `location_id` - (Required, Forces new resource) The numeric id of the location
where this MCR should be created in.
* `name` - (Required) The name of the MCR.
* `rate_limit` - (Required, Forces new resource) The speed of the MCR. Please
check with the Megaport documentation for available speeds.
* `asn` - (Optional, Forces new resource) The Autonomous System Number (ASN) to
use for BGP peering sessions on VXCs connected to this MCR. If not configured,
the Megaport supplied public ASN will be used.
* `invoice_reference` - (Optional) Used for billing purposes, a reference to
this specific line item.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique product id of the MCR.
