-- name: CreateAuditLog :one
INSERT INTO audit_logs (organization_id, user_id, action, resource, details, ip_address)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, organization_id, user_id, action, resource, details, ip_address, created_at;

-- name: ListAuditLogsByOrg :many
SELECT id, organization_id, user_id, action, resource, details, ip_address, created_at
FROM audit_logs
WHERE organization_id = $1
ORDER BY created_at DESC;
