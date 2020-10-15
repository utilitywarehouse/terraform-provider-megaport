data "megaport_location" "foo" {
  name_regex = "{{ .location }}"
}

resource "megaport_mcr" "foo" {
  name              = "terraform_acctest_{{ .uid }}"
  location_id       = data.megaport_location.foo.id
  rate_limit        = {{ .rate_limit }}
  invoice_reference = "{{ .uid }}"
{{- if .asn }}
  asn               = {{ .asn }}
{{- end }}
}

