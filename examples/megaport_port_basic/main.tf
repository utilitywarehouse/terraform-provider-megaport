data "megaport_location" "port" {
  name_regex = "{{ .location }}"
}

resource "megaport_port" "test" {
  name        = "terraform_acctest_{{ .uid }}"
  location_id = data.megaport_location.port.id
  speed       = 1000
  term        = 1
}

