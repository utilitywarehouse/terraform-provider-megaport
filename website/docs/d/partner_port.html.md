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
data "megaport_location" "foo" {
  name_regex = "Telehouse North"
}

data "megaport_partner_port" "foo" {
  name_regex   = "eu-west-1"

  marketplace {
    location_id = data.megaport_location.foo.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Required, Forces new resource) A regex string filter to apply
to the Port list returned by Megaport.
* `aws` - (Optional, Conflicts with `marketplace` and `gcp`) Search Ports that
are suitable to use for connections to AWS. See below for supported attributes.
* `gcp` - (Optional, Conflicts with `marketplace` and `aws`) Search Ports that
are suitable to use for connections to GCP. See below for supported attributes.
* `marketplace` - (Optional, Conflicts with `aws` and `gcp`) Search Ports from
the Megaport marketplace.

The `aws` and `marketplace` blocks support:

* `location_id` - (Required, Forces new resource) Filter Ports based on a
location id, as returned by the [megaport_location](/docs/providers/megaport/d/location.html)
datasource.
* `vxc_permitted` - (Optional, Forces new resource, Default: `true`) Limit
search to Ports that have the `vxcPermitted` flag set. This is true by default
since Megaport will only accept VXCs to these Ports.

The `gcp` block supports:

* `pairing_key` - (Required, Forces new resource) The GCP Partner Interconnect
pairing key that will be used for the VXC.

~> **Note:** If more or less than a single match is returned by the search,
Terraform will fail. Ensure that your search is specific enough to return a
single Port.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Product UID of the selected Port.
* `bandwidths` - A list of bandwidths supported for VXCs to this Port. This is
only populated when using the `gcp` block.

