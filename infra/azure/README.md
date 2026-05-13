# Azure Deployment Target

The first Azure deployment should map cleanly to the same shape as AWS:

- Azure Container Apps or App Service for `estate-agent serve`
- Azure Queue Storage or Service Bus for async work
- Key Vault for GitHub credentials and webhook secret
- Application Insights for logs and traces
- Container Registry for images

Terraform will live here once the local MVP is stable.

Initial variables:

- `github_webhook_secret`
- `github_token_secret_name`
- `default_repo_owner`
- `default_repo_name`
- `dry_run`

