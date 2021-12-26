package form

import "gitlab.com/fibocloud/medtech/gin/databases"

// PermissionParams create body params
type PermissionParams struct {
	Path        string `json:"path" binding:"required"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

// RoleCreateParams create body params
type RoleCreateParams struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

// RoleUpdateParams create body params
type RoleUpdateParams struct {
	Name        string                           `json:"name" binding:"required"`
	Description string                           `json:"description"`
	IsActive    bool                             `json:"is_active"`
	Permissions []*databases.MedSystemPermission `json:"permissions"`
}
