package repository

import (
	"bad_boyes/internal/models"

	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) CreatePost(post *models.Post) error {
	return r.db.Create(post).Error
}

func (r *PostRepository) GetPostByID(id uint) (*models.Post, error) {
	var post models.Post
	err := r.db.Preload("User").First(&post, id).Error
	return &post, err
}

func (r *PostRepository) UpdatePost(post *models.Post) error {
	// Create backup in post_history
	history := &models.PostHistory{
		PostID:        post.ID,
		UserID:        post.UserID,
		Title:         post.Title,
		Description:   post.Description,
		Address:       post.Address,
		ContactName:   post.ContactName,
		MobileNumber:  post.MobileNumber,
		IncidentDate:  post.IncidentDate,
		Status:        post.Status,
		IsAnonymous:   post.IsAnonymous,
		Visibility:    post.Visibility,
		AllowComments: post.AllowComments,
	}
	if err := r.db.Create(history).Error; err != nil {
		return err
	}

	return r.db.Save(post).Error
}

func (r *PostRepository) DeletePost(id uint) error {
	return r.db.Delete(&models.Post{}, id).Error
}

func (r *PostRepository) ListPosts(page, pageSize int, userID *uint) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	query := r.db.Model(&models.Post{}).Where("visibility = ?", "public")
	if userID != nil {
		query = query.Or("user_id = ?", userID)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("User").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&posts).Error

	return posts, total, err
}

func (r *PostRepository) CreateReport(report *models.Report) error {
	return r.db.Create(report).Error
}

func (r *PostRepository) GetReportByID(id uint) (*models.Report, error) {
	var report models.Report
	err := r.db.Preload("Post").Preload("Reporter").First(&report, id).Error
	return &report, err
}

func (r *PostRepository) UpdateReportStatus(id uint, status string) error {
	return r.db.Model(&models.Report{}).Where("id = ?", id).Update("status", status).Error
}

func (r *PostRepository) ListReports(page, pageSize int) ([]models.Report, int64, error) {
	var reports []models.Report
	var total int64

	err := r.db.Model(&models.Report{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Post").Preload("Reporter").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&reports).Error

	return reports, total, err
}

func (r *PostRepository) GetPostHistory(postID uint) ([]models.PostHistory, error) {
	var history []models.PostHistory
	err := r.db.Where("post_id = ?", postID).Order("created_at DESC").Find(&history).Error
	return history, err
}
