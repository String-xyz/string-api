data "aws_route53_zone" "root" {
  name = local.root_domain
}

resource "aws_route53_record" "domain" {
  name = "admin.${local.root_domain}"
  type    = "A"
  zone_id = data.aws_route53_zone.root.zone_id
  alias {
    evaluate_target_health = false
    name                   = aws_alb.alb.dns_name
    zone_id                = aws_alb.alb.zone_id
  }
}
