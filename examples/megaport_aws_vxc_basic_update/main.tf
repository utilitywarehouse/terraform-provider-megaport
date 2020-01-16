data "megaport_location" "aws" {
  name_regex = "{{ .location }}"
}

data "megaport_partner_port" "aws" {
  name_regex   = "eu-west-1"
  connect_type = "AWS"
  location_id  = data.megaport_location.aws.id
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

resource "megaport_aws_vxc" "foo" {
  name              = "terraform_acctest_{{ .uid }}"
  rate_limit        = 1000
  invoice_reference = "{{ .uid }}"

  a_end {
    product_uid = megaport_port.foo.id
    vlan        = 568
  }

  b_end {
    product_uid    = data.megaport_partner_port.aws.id
    aws_account_id = "{{ .aws_account_id }}"
    customer_asn   = {{ .customer_asn }}
    type           = "private"
  }
}
