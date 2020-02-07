---
layout: "megaport"
subcategory: "resources"
page_title: "Megaport: megaport_gcp_vxc"
description: |-
  Provides a Megaport GCP Virtual Cross Connect (VXC) resource.
---

# Resource: megaport_gcp_vxc

Provides a Megaport Virtual Cross Connect (VXC) resource to GCP. Allows VXCs
to GCP to be created, updated and deleted.

## Example Usage

```hcl
data "megaport_partner_port" "gcp" {
  name_regex = "London"

  gcp {
    pairing_key = var.pairing_key
  }
}

data "megaport_port" "own" {
  name_regex = "bar"
}

resource "megaport_gcp_vxc" "foobar" {
  name       = "foobar"
  rate_limit = data.megaport_partner_port.gcp.bandwidths[1]

  a_end {
    product_uid = data.megaport_port.own.id
  }

  b_end {
    product_uid = data.megaport_partner_port.gcp.id
    pairing_key = var.pairing_key
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the VXC.
* `rate_limit` - (Required) The rate limit of the VXC (Must be one of the
available bandwidths, as exported by the
[`megaport_partner_port`](/docs/providers/megaport/d/partner_port.html)
datasource.)
* `invoice_reference` - (Optional) Used for billing purposes, a reference to
this specific line item.
* `a_end` - (Required) - Points to a port owned by the current account that will
act as one end of the VXC (see [VXC ends](gcp_vxc.html#vxc-ends)).
* `b_end` - (Required) - Points to an GCP port that will act as the other end of
the VXC (see [VXC ends](gcp_vxc.html#vxc-ends)).

### VXC ends

The VXC's two ends refer to two ports: A end is the port owned by the current
account and B end is the port on the GCP side. These have different arguments,
detailed below.

#### A End

* `product_uid` - (Required, Forces new resource) The product UID of the port.
* `vlan` - (Optional) The VLAN id to use for this connection. If not specified,
Megaport will automatically select an available one.

#### B End

* `product_uid` - (Required, Forces new resource) The product UID of the port.
* `pairing_key` - (Required, Forces new resource) The GCP Partner Interconnect
[Pairing Key](https://cloud.google.com/interconnect/docs/concepts/terminology#pairingkey)
to use for this connection.

Additionally to all arguments above, `b_end` also exports the following
attribute:

* `connected_product_uid` - This is set to the uid of the Port that the VXC is
using for its B End.

~> **Note:** `connected_product_uid` can be different from the supplied
`product_uid` argument. Megaport load balances the VXCs among a pool of GCP
Partner Ports. Effectively, when creating a new VXC to an GCP Partner Port,
Megaport might instead use a different Port from the pool, chosen such that it
will result in a VXC with properties identical to what was requested.
Additionally, Megaport might migrate an existing VXC to a different Port, in
which case `connected_product_uid` will be updated to reflect the change. It is
exported as a separate attribute in order to prevent terraform from producing a
plan when the aforementioned situations occur.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique product id of the port.
