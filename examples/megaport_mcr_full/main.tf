data "megaport_location" "foo" {
  name_regex    = "{{ .location }}"
  mcr_available = {{ .mcr_version }}
}

resource "megaport_mcr" "foo" {
  mcr_version       = {{ .mcr_version }}
  name              = "terraform_acctest_{{ .uid }}"
  location_id       = data.megaport_location.foo.id
  rate_limit        = {{ .rate_limit }}
  invoice_reference = "{{ .uid }}"
{{- if .term }}
  term              = {{ .term }}
{{- end }}
{{- if .asn }}
  asn               = {{ .asn }}
{{- end }}
}

