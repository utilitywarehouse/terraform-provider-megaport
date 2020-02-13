---
layout: "megaport"
subcategory: "resources"
page_title: "Megaport: megaport_aws_vxc"
description: |-
  Provides a Megaport AWS Virtual Cross Connect (VXC) resource.
---

# Resource: megaport_aws_vxc

Provides a Megaport Virtual Cross Connect (VXC) resource to AWS. Allows VXCs
to AWS to be created, updated and deleted.

## Example Usage

```hcl
data "megaport_location" "aws" {
  name_regex = "foo"
}

data "megaport_partner_port" "aws" {
  name_regex   = "eu-west-1"
  connect_type = "AWS"
  location_id  = data.megaport_location.aws.id
}

data "megaport_port" "own" {
  name_regex = "bar"
}

resource "megaport_aws_vxc" "foobar" {
  name       = "foobar"
  rate_limit = 100

  a_end {
    product_uid = data.megaport_port.own.id
    vlan        = 567
  }

  b_end {
    product_uid    = data.megaport_partner_port.aws.id
    aws_account_id = "012345678912"
    customer_asn   = "64512"
    type           = "private"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the VXC.
* `rate_limit` - (Required) The rate limit of the VXC (Must not exceed the speed
of the port at `a_end`)
* `invoice_reference` - (Optional) Used for billing purposes, a reference to
this specific line item.
* `a_end` - (Required) - Points to a port owned by the current account that will
act as one end of the VXC (see [VXC ends](aws_vxc.html#vxc-ends)).
* `b_end` - (Required) - Points to an AWS port that will act as the other end of
the VXC (see [VXC ends](aws_vxc.html#vxc-ends)).

### VXC ends

The VXC's two ends refer to two ports: A end is the port owned by the current
account and B end is the port on the AWS side. These have different arguments,
detailed below.

#### A End

* `product_uid` - (Required, Forces new resource) The product UID of the port.
* `vlan` - (Optional) The VLAN id to use for this connection. If not specified,
Megaport will automatically select an available one.

#### B End

* `product_uid` - (Required, Forces new resource) The product UID of the port.
* `aws_connection_name` - (Optional) The name of the virtual interface in AWS.
* `aws_account_id` - (Required) The ID of the AWS account to connect to.
* `aws_ip_address` - (Optional) The IP Address space assigned in the AWS VPC
network to peer with. If unspecified, a private `/30` will be automatically
assigned by Megaport.
* `bgp_auth_key` - (Optional) The BGP auth key for the session. If not
specified, Megaport will automatically generate one.
* `customer_asn` - (Required) Your network's Autonomous System Number. For
Private Direct Connects, this must be a private ASN. For public Direct Connects,
this can be either a Private or Public ASN.
* `customer_ip_address` - (Optional) The IP Address space you will use on your
network to peer with. If unspecified, a private `/30` will be automatically
assigned by Megaport.
* `type` - (Optional, Default: `"private"`) Type of the virtual interface to
AWS.  Accepted values: `"private"`, `"public"`.

Additionally to all arguments above, `b_end` also exports the following
attribute:

* `connected_product_uid` - This is set to the uid of the Port that the VXC is
using for its B End.

~> **Note:** `connected_product_uid` can be different from the supplied
`product_uid` argument. Megaport load balances the VXCs among a pool of AWS
Partner Ports. Effectively, when creating a new VXC to an AWS Partner Port,
Megaport might instead use a different Port from the pool, chosen such that it
will result in a VXC with properties identical to what was requested.
Additionally, Megaport might migrate an existing VXC to a different Port, in
which case `connected_product_uid` will be updated to reflect the change. It is
exported as a separate attribute in order to prevent terraform from producing a
plan when the aforementioned situations occur.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique product id of the port.

## Import

The AWS VXC can be imported using its product uid, like any other resource,
e.g.:

```
$ terraform import megaport_aws_vxc.foobar 1f33ea1d-ecc2-4fc3-a3a4-1e4774b04d76
```

!> **Warning:** When a AWS VXC is imported, any changes to the B End
`product_uid` attribute are ignored. To force an update, you will need to
`taint` the resource. After re-creating the resource, it will start to behave as
expected and will compute the full diff. This is to work around the load
balancing behaviour mentioned in the Note above.
