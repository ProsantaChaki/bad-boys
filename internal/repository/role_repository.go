package repository

import (
	"bad_boyes/internal/models"

	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) CreateRole(role *models.Role) (*models.Role, error) {
	err := r.db.Create(role).Error
	return role, err
}

func (r *RoleRepository) AssignRoleToUser(userID, roleID uint) error {
	return r.db.Create(&models.UserRole{
		UserID: userID,
		RoleID: roleID,
	}).Error
}

func (r *RoleRepository) RemoveRoleFromUser(userID, roleID uint) error {
	return r.db.Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&models.UserRole{}).Error
}

func (r *RoleRepository) GetUserRoles(userID uint) ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Model(&models.User{}).Where("id = ?", userID).Association("Roles").Find(&roles)
	return roles, err
}

func (r *RoleRepository) CreatePermission(permission *models.Permission) (*models.Permission, error) {
	err := r.db.Create(permission).Error
	return permission, err
}

func (r *RoleRepository) AssignPermissionToRole(roleID, permissionID uint) error {
	return r.db.Create(&models.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
	}).Error
}

func (r *RoleRepository) RemovePermissionFromRole(roleID, permissionID uint) error {
	return r.db.Where("role_id = ? AND permission_id = ?", roleID, permissionID).Delete(&models.RolePermission{}).Error
}

func (r *RoleRepository) GetRolePermissions(roleID uint) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Model(&models.Role{}).Where("id = ?", roleID).Association("Permissions").Find(&permissions)
	return permissions, err
}
