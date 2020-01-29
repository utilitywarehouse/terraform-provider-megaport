data "megaport_location" "foo" {
  name_regex = "{{ .locationA }}"
}

resource "megaport_port" "foo" {
  name        = "terraform_acctest_a_{{ .uid }}"
  location_id = data.megaport_location.foo.id
  speed       = 1000
  term        = 1
}

data "megaport_location" "bar" {
  name_regex = "{{ .locationB }}"
}

resource "megaport_port" "bar" {
  name        = "terraform_acctest_b_{{ .uid }}"
  location_id = data.megaport_location.bar.id
  speed       = 1000
  term        = 1
}

resource "megaport_private_vxc" "foobar" {
  name              = "terraform_acctest_{{ .uid }}"
  rate_limit        = 200
  invoice_reference = "{{ .uid }}"

  a_end {
    product_uid = megaport_port.foo.id
    vlan        = {{ .vlanA }}
  }

  b_end {
    product_uid = megaport_port.bar.id
    vlan        = {{ .vlanB }}
  }
}
