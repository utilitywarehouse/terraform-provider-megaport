## 0.2.0-rc.1 (October 16, 2020)

NOTES:

* upgrade to terraform-plugin-sdk v2.0.4
* upgrade to go 1.15
* migrate to context-aware CRUD methods for resources and datasources

## 0.1.0 (October 15, 2020)

BREAKING CHANGES:

* Removed support for Megaport Cloud Router v1 ([#6](https://github.com/utilitywarehouse/terraform-provider-megaport/pull/6))

## 0.1.0-rc.3 (October 14, 2020)

NOTES:

Adds resource import tests.

BUG FIXES:

* resource/megaport_mcr: fix perpetual diff after import

## 0.1.0-rc.2 (February 27, 2020)

BUG FIXES:

* data-source/megaport_partner_port: use a randomised key to produce consistent
results ([#1](https://github.com/utilitywarehouse/terraform-provider-megaport/issues/1))

## 0.1.0-rc.1 (February 21, 2020)

NOTES:

This is the first release candidate for the Megaport Provider.

FEATURES:

* **New Resource:** `megaport_port`
* **New Resource:** `megaport_mcr`
* **New Resource:** `megaport_aws_vxc`
* **New Resource:** `megaport_gcp_vxc`
* **New Resource:** `megaport_private_vxc`

* **New Data Source:** `megaport_location`
* **New Data Source:** `megaport_partner_port`
* **New Data Source:** `megaport_port`
