data "megaport_location" "foo" {
  name_regex = "{{ .location }}"
}

resource "megaport_port" "foo" {
  name                   = "terraform_acctest_{{ .uid }}"
  location_id            = data.megaport_location.foo.id
  speed                  = 10000
  term                   = 12
  invoice_reference      = "{{ .uid }}"
  marketplace_visibility = "public"
}

