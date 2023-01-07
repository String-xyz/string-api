data "aws_route53_zone" "root" {
  name = local.domain
}

resource "aws_route53_record" "domain" {
  name = local.domain
  type    = "A"
  zone_id = data.aws_route53_zone.root.zone_id
  alias {
    evaluate_target_health = false
    name                   = aws_cloudfront_distribution.this.domain_name
    zone_id                = aws_cloudfront_distribution.this.hosted_zone_id
  }
}

module "cloudfront" {
  source      = "../acm"
  domain_name = local.domain
  aws_region  = "us-east-1"
  zone_id     = data.aws_route53_zone.root.zone_id
  tags = {
    Environment = local.env
    Name = "cert-${local.domain}"
  }
}

module "alb_acm" {
  source            = "../acm"
  domain_name       = local.domain
  aws_region        = "us-west-2"
  zone_id           = data.aws_route53_zone.root.zone_id
  tags = {
    Name = "cert-${local.domain}-alb"
  }
}