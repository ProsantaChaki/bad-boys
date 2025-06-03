package services

import (
	"bad_boyes/internal/models"
	"bad_boyes/internal/repository"
)

type RoleService struct {
	roleRepo *repository.RoleRepository
}

func NewRoleService(roleRepo *repository.RoleRepository) *RoleService {
	return &RoleService{
		roleRepo: roleRepo,
	}
}

// CreateRole creates a new role
func (s *RoleService) CreateRole(name, description string) (*models.Role, error) {
	role := &models.Role{
		Name:        name,
		Description: description,
	}
	return s.roleRepo.CreateRole(role)
}

// AssignRoleToUser assigns a role to a user
func (s *RoleService) AssignRoleToUser(userID, roleID uint) error {
	return s.roleRepo.AssignRoleToUser(userID, roleID)
}

// RemoveRoleFromUser removes a role from a user
func (s *RoleService) RemoveRoleFromUser(userID, roleID uint) error {
	return s.roleRepo.RemoveRoleFromUser(userID, roleID)
}

// GetUserRoles returns all roles assigned to a user
func (s *RoleService) GetUserRoles(userID uint) ([]models.Role, error) {
	return s.roleRepo.GetUserRoles(userID)
}

// CreatePermission creates a new permission
func (s *RoleService) CreatePermission(name, description, resource, action string) (*models.Permission, error) {
	permission := &models.Permission{
		Name:        name,
		Description: description,
		Resource:    resource,
		Action:      action,
	}
	return s.roleRepo.CreatePermission(permission)
}

// AssignPermissionToRole assigns a permission to a role
func (s *RoleService) AssignPermissionToRole(roleID, permissionID uint) error {
	return s.roleRepo.AssignPermissionToRole(roleID, permissionID)
}

// RemovePermissionFromRole removes a permission from a role
func (s *RoleService) RemovePermissionFromRole(roleID, permissionID uint) error {
	return s.roleRepo.RemovePermissionFromRole(roleID, permissionID)
}

// GetRolePermissions returns all permissions assigned to a role
func (s *RoleService) GetRolePermissions(roleID uint) ([]models.Permission, error) {
	return s.roleRepo.GetRolePermissions(roleID)
}

// CheckPermission checks if a user has a specific permission
func (s *RoleService) CheckPermission(userID uint, resource, action string) (bool, error) {
	roles, err := s.GetUserRoles(userID)
	if err != nil {
		return false, err
	}

	for _, role := range roles {
		permissions, err := s.GetRolePermissions(role.ID)
		if err != nil {
			return false, err
		}

		for _, permission := range permissions {
			if permission.Resource == resource && permission.Action == action {
				return true, nil
			}
		}
	}

	return false, nil
}
