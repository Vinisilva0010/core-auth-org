-- name: CreateRole :one
INSERT INTO roles (organization_id, name, description)
VALUES ($1, $2, $3)
RETURNING id, organization_id, name, description, created_at;

-- name: AssignRoleToMember :exec
UPDATE organization_members
SET role_id = $1
WHERE organization_id = $2 AND user_id = $3;
