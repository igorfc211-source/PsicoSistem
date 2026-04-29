# Infrastructure Guide

## Storage drivers

The backend can run with two persistence drivers:

- `json`: local development and low-friction demos
- `postgres`: production-style storage with tenant-safe SQL repositories

### Local JSON

- `STORAGE_DRIVER=json`
- `DATA_FILE=data/app_state.json`

### PostgreSQL

- `STORAGE_DRIVER=postgres`
- `DATABASE_URL=postgres://...`
- `DATABASE_AUTO_MIGRATE=true`

When `postgres` is enabled, the bootstrap opens a `pgxpool`, ensures the base schema, seeds the default plans and swaps the repositories used by `auth`, `tenant`, `subscription` and `user`.

## AWS profile

The backend validates a runtime profile through:

- `CLOUD_PROVIDER=aws`
- `AWS_REGION`
- `AWS_S3_BUCKET`
- `AWS_SECRETS_PREFIX`
- `AWS_USE_IAM_ROLE`
- `AWS_CLOUDWATCH_NAMESPACE`

### Recommended AWS architecture

- API runtime: ECS Fargate or EKS
- Load balancing: Application Load Balancer
- Database: Amazon RDS PostgreSQL
- Documents and AI artifacts: Amazon S3
- Secrets: AWS Secrets Manager or SSM Parameter Store
- Logs and metrics: Amazon CloudWatch

## Permission model

The permission engine is account-scoped and designed so future domain modules can reuse the same semantics:

- `all`: query the entire tenant scope
- `own`: always filter by the authenticated professional/user id
- `none`: deny access

For future modules like patients, services, calendar, finance and AI history, the rule is:

- `all`: `WHERE tenant_id = $1`
- `own`: `WHERE tenant_id = $1 AND responsible_user_id = $2`
- `none`: deny before reaching the repository
