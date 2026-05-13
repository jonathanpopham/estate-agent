# AWS Deployment Target

The first AWS deployment should be intentionally boring:

- API Gateway or ALB for HTTPS webhook intake
- ECS Fargate service running `estate-agent serve`
- SQS queue between intake and worker once builder mode exists
- Secrets Manager for GitHub credentials and webhook secret
- CloudWatch logs and alarms
- ECR for container image storage

Terraform will live here once the local MVP is stable.

Initial variables:

- `github_webhook_secret`
- `github_token_secret_arn`
- `default_repo_owner`
- `default_repo_name`
- `dry_run`

