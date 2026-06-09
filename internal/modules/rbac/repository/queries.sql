-- name: CreateRole :one
INSERT INTO roles (organization_id, name, description)
VALUES ($1, $2, $3)
RETURNING id, organization_id, name, description, created_at;

-- name: AssignRoleToMember :exec
UPDATE organization_members
SET role_id = $1
WHERE organization_id = $2 AND user_id = $3;

-- name: CreatePermission :one
INSERT INTO permissions (name)
VALUES ($1)
RETURNING id, name;

-- name: AssignPermissionToRole :exec
INSERT INTO role_permissions (role_id, permission_id)
VALUES ($1, $2);
