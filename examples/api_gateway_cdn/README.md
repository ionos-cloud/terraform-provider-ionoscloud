# API Gateway + CDN Example
Demo setup for API Gateway in combination with CDN Distribution, DNS records and a certificate for usage with a custom domain. The Terraform manifests provision following infrastructure:

- Creates auto certificate provider with Let's Encrypt and requests SSL certs for a configurable subdomain
- Sets up an API Gateway and routes to send requests to specific backends
- Establishes CDN distribution with HTTPS, caching, WAF, rate limiting, and geo-restrictions
- Creates DNS A and AAAA records for the CDN distribution
- Sets up a virtual server with nginx as the backend server via cloud-init
