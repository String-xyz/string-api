data "aws_route53_zone" "root" {
  name = local.root_domain
}

resource "aws_route53_record" "domain" {
  name = "string-api.${local.root_domain}"
  type    = "A"
  zone_id = data.aws_route53_zone.root.zone_id
  alias {
    evaluate_target_health = false
    name                   = aws_cloudfront_distribution.this.domain_name
    zone_id                = aws_cloudfront_distribution.this.hosted_zone_id
  }
}

resource "aws_route53_record" "www_domain" {
  name    = "www.string-api.${local.root_domain}"
  type    = "A"
  zone_id = data.aws_route53_zone.root.zone_id
  alias {
    evaluate_target_health = false
    name                   = aws_cloudfront_distribution.this.domain_name
    zone_id                = aws_cloudfront_distribution.this.hosted_zone_id
  }
}

module "acm" {
  source      = "../acm"
  domain_name = "string-api.${local.root_domain}"
  alternative_names = ["www.string-api.${local.root_domain}"]
  aws_region  = "us-east-1"
  zone_id     = data.aws_route53_zone.root.zone_id
  tags = {
    Environment = local.env
    Name = "string-api.${local.root_domain}"
  }
}