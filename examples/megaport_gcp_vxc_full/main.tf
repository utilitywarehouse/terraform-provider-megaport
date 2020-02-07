data "megaport_partner_port" "gcp" {
  name_regex = "{{ .nameRegex }}"

  gcp {
    pairing_key = "{{ .pairingKey }}"
  }
}

data "megaport_location" "foo" {
  name_regex = "Telehouse North"
}

resource "megaport_port" "foo" {
  name        = "terraform_acctest_{{ .uid }}"
  location_id = data.megaport_location.foo.id
  speed       = 1000
  term        = 1
}

resource "megaport_gcp_vxc" "foo" {
  name              = "terraform_acctest_{{ .uid }}"
  rate_limit        = {{ .rateLimit }}
  invoice_reference = "{{ .uid }}"

  a_end {
    product_uid = megaport_port.foo.id
    vlan        = {{ .vlan }}
  }

  b_end {
    product_uid = data.megaport_partner_port.gcp.id
    pairing_key = "{{ .pairingKey }}"
  }
}
