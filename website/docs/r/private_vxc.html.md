---
layout: "megaport"
subcategory: "resources"
page_title: "Megaport: megaport_private_vxc"
description: |-
  Provides a Megaport private Virtual Cross Connect (VXC) resource.
---

# Resource: megaport_private_vxc

Provides a Megaport Virtual Cross Connect (VXC) resource to a private Port.
Allows VXCs to private Ports to be created, updated and deleted.

## Example Usage

```hcl
data "megaport_location" "foo" {
  name_regex = "foo"
}

resource "megaport_port" "foo" {
  name        = "port_a"
  location_id = data.megaport_location.foo.id
  speed       = 1000
  term        = 1
}

data "megaport_location" "bar" {
  name_regex = "bar"
}

resource "megaport_port" "bar" {
  name        = "port_b"
  location_id = data.megaport_location.bar.id
  speed       = 1000
  term        = 1
}

resource "megaport_private_vxc" "foobar" {
  name              = "foobar"
  rate_limit        = 200
  invoice_reference = "foobar"

  a_end {
    product_uid = megaport_port.foo.id
    vlan        = 123
  }

  b_end {
    product_uid = megaport_port.bar.id
    vlan        = 123
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the VXC.
* `rate_limit` - (Required) The rate limit of the VXC (Must not exceed the speed
of the ports at `a_end` and `b_end`)
* `invoice_reference` - (Optional) Used for billing purposes, a reference to
this specific line item.
* `a_end` - (Required) - Points to a port owned by the current account that will
act as one end of the VXC (see [VXC ends](private_vxc.html#vxc-ends)).
* `b_end` - (Required) - Points to a port owned by the current account that will
act as the other end of the VXC (see [VXC ends](private_vxc.html#vxc-ends)).

### VXC ends

The VXC's two ends refer to two ports: A end and B end. Both ports need to be
owned by the current account.

* `product_uid` - (Required, Forces new resource) The product UID of the port.
* `vlan` - (Optional) The VLAN id to use for this connection. If not specified,
Megaport will automatically select an available one.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique product id of the port.
