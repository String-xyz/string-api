resource "aws_cloudfront_distribution" "this" {
  enabled         = true
  is_ipv6_enabled = true
  aliases         = ["string-api.${local.root_domain}", "www.string-api.${local.root_domain}"]

  origin {
    domain_name = aws_alb.alb.dns_name
    origin_id   = local.origin_id
    custom_origin_config {
      http_port              = 80
      https_port             = 443
      origin_protocol_policy = "https-only"
      origin_ssl_protocols   = ["TLSv1", "TLSv1.1", "TLSv1.2"]
    }
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
      locations        = []
    }
  }


  default_cache_behavior {
    target_origin_id = local.origin_id
    compress         = true
    allowed_methods  = ["DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT"]
    cached_methods   = ["GET", "HEAD", "OPTIONS"]

    forwarded_values {
      query_string = true
      headers      = ["X-Forwarded-For", "Host", "X-Api-Key"]
      cookies {
        forward = "all"
      }
    }

    viewer_protocol_policy = "redirect-to-https"
    min_ttl                = 0
    default_ttl            = 7200
    max_ttl                = 86400
  }

  viewer_certificate {
    ssl_support_method             = "sni-only"
    acm_certificate_arn            = module.acm.arn
    minimum_protocol_version       = "TLSv1.1_2016"
    cloudfront_default_certificate = false
  }
}
