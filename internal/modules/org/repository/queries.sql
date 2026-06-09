-- name: CreateOrganization :one
INSERT INTO organizations (name, slug)
VALUES ($1, $2)
RETURNING id, name, slug, created_at, updated_at;

-- name: AddUserToOrganization :exec
INSERT INTO organization_members (organization_id, user_id)
VALUES ($1, $2);

-- name: GetOrganizationByID :one
SELECT id, name, slug, created_at, updated_at
FROM organizations
WHERE id = $1 LIMIT 1;

-- name: CreateOrganizationUnit :one
INSERT INTO organization_units (organization_id, name)
VALUES ($1, $2)
RETURNING id, organization_id, name, created_at;

-- name: ListOrganizationUnits :many
SELECT id, organization_id, name, created_at
FROM organization_units
WHERE organization_id = $1
ORDER BY name;
