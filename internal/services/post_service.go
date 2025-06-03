package services

import (
	"bad_boyes/internal/models"
	"bad_boyes/internal/repository"
	"errors"
	"log"
	"time"
)

type PostService struct {
	postRepo  *repository.PostRepository
	auditRepo *repository.AuditRepository
}

func NewPostService(postRepo *repository.PostRepository, auditRepo *repository.AuditRepository) *PostService {
	return &PostService{
		postRepo:  postRepo,
		auditRepo: auditRepo,
	}
}

func (s *PostService) CreatePost(userID uint, req models.CreatePostRequest) (*models.Post, error) {
	log.Printf("Starting post creation for user ID: %d", userID)
	log.Printf("Request data: %+v", req)

	post := &models.Post{
		UserID:        userID,
		Title:         req.Title,
		Description:   req.Description,
		Address:       req.Address,
		ContactName:   req.ContactName,
		MobileNumber:  req.MobileNumber,
		IncidentDate:  req.IncidentDate,
		IsAnonymous:   req.IsAnonymous,
		Visibility:    req.Visibility,
		AllowComments: req.AllowComments,
		Status:        "active",
	}
	log.Printf("Created post object: %+v", post)

	log.Printf("Attempting to save post to database")
	if err := s.postRepo.CreatePost(post); err != nil {
		log.Printf("Failed to create post: %v", err)
		return nil, err
	}
	log.Printf("Post saved successfully with ID: %d", post.ID)

	// Create audit log
	log.Printf("Creating audit log for post creation")
	auditLog := &models.AuditLog{
		UserID:    &userID,
		Action:    "create",
		TableName: "posts",
		RecordID:  post.ID,
		NewValues: models.JSON{
			"title":          post.Title,
			"description":    post.Description,
			"address":        post.Address,
			"contact_name":   post.ContactName,
			"mobile_number":  post.MobileNumber,
			"incident_date":  time.Time(post.IncidentDate).Format("2006-01-02"),
			"is_anonymous":   post.IsAnonymous,
			"visibility":     post.Visibility,
			"allow_comments": post.AllowComments,
			"status":         post.Status,
		},
	}
	log.Printf("Audit log object created: %+v", auditLog)

	log.Printf("Attempting to save audit log")
	if err := s.auditRepo.CreateLog(auditLog); err != nil {
		log.Printf("Failed to create audit log: %v", err)
		return nil, err
	}
	log.Printf("Audit log saved successfully")

	log.Printf("Post creation completed successfully for user ID: %d", userID)
	return post, nil
}

func (s *PostService) GetPost(id uint) (*models.Post, error) {
	return s.postRepo.GetPostByID(id)
}

func (s *PostService) UpdatePost(userID uint, postID uint, req models.UpdatePostRequest) (*models.Post, error) {
	log.Printf("Starting post update for post ID: %d by user ID: %d", postID, userID)
	log.Printf("Update request data: %+v", req)

	post, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		log.Printf("Failed to fetch post: %v", err)
		return nil, err
	}
	log.Printf("Found existing post: %+v", post)

	if post.UserID != userID {
		log.Printf("Unauthorized update attempt: post belongs to user %d, but user %d is trying to update", post.UserID, userID)
		return nil, errors.New("unauthorized")
	}

	// Store old values for audit
	oldValues := models.JSON{
		"title":          post.Title,
		"description":    post.Description,
		"address":        post.Address,
		"contact_name":   post.ContactName,
		"mobile_number":  post.MobileNumber,
		"incident_date":  time.Time(post.IncidentDate).Format("2006-01-02"),
		"is_anonymous":   post.IsAnonymous,
		"visibility":     post.Visibility,
		"allow_comments": post.AllowComments,
	}
	log.Printf("Stored old values for audit: %+v", oldValues)

	// Update fields if provided
	if req.Title != "" {
		log.Printf("Updating title from '%s' to '%s'", post.Title, req.Title)
		post.Title = req.Title
	}
	if req.Description != "" {
		log.Printf("Updating description")
		post.Description = req.Description
	}
	if req.Address != "" {
		log.Printf("Updating address from '%s' to '%s'", post.Address, req.Address)
		post.Address = req.Address
	}
	if req.ContactName != "" {
		log.Printf("Updating contact name from '%s' to '%s'", post.ContactName, req.ContactName)
		post.ContactName = req.ContactName
	}
	if req.MobileNumber != "" {
		log.Printf("Updating mobile number from '%s' to '%s'", post.MobileNumber, req.MobileNumber)
		post.MobileNumber = req.MobileNumber
	}
	if req.IncidentDate != (models.Date{}) {
		log.Printf("Updating incident date")
		post.IncidentDate = req.IncidentDate
	}
	post.IsAnonymous = req.IsAnonymous
	if req.Visibility != "" {
		log.Printf("Updating visibility from '%s' to '%s'", post.Visibility, req.Visibility)
		post.Visibility = req.Visibility
	}
	post.AllowComments = req.AllowComments

	log.Printf("Attempting to save updated post")
	if err := s.postRepo.UpdatePost(post); err != nil {
		log.Printf("Failed to update post: %v", err)
		return nil, err
	}
	log.Printf("Post updated successfully")

	// Create audit log
	log.Printf("Creating audit log for post update")
	auditLog := &models.AuditLog{
		UserID:    &userID,
		Action:    "update",
		TableName: "posts",
		RecordID:  post.ID,
		OldValues: oldValues,
		NewValues: models.JSON{
			"title":          post.Title,
			"description":    post.Description,
			"address":        post.Address,
			"contact_name":   post.ContactName,
			"mobile_number":  post.MobileNumber,
			"incident_date":  time.Time(post.IncidentDate).Format("2006-01-02"),
			"is_anonymous":   post.IsAnonymous,
			"visibility":     post.Visibility,
			"allow_comments": post.AllowComments,
		},
	}
	log.Printf("Audit log object created: %+v", auditLog)

	log.Printf("Attempting to save audit log")
	if err := s.auditRepo.CreateLog(auditLog); err != nil {
		log.Printf("Failed to create audit log: %v", err)
		return nil, err
	}
	log.Printf("Audit log saved successfully")

	log.Printf("Post update completed successfully for post ID: %d", postID)
	return post, nil
}

func (s *PostService) DeletePost(userID uint, postID uint) error {
	post, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		return err
	}

	if post.UserID != userID {
		return errors.New("unauthorized")
	}

	if err := s.postRepo.DeletePost(postID); err != nil {
		return err
	}

	// Create audit log
	auditLog := &models.AuditLog{
		UserID:    &userID,
		Action:    "delete",
		TableName: "posts",
		RecordID:  postID,
		OldValues: models.JSON{
			"title":          post.Title,
			"description":    post.Description,
			"address":        post.Address,
			"contact_name":   post.ContactName,
			"mobile_number":  post.MobileNumber,
			"incident_date":  post.IncidentDate,
			"is_anonymous":   post.IsAnonymous,
			"visibility":     post.Visibility,
			"allow_comments": post.AllowComments,
			"status":         post.Status,
		},
	}
	return s.auditRepo.CreateLog(auditLog)
}

func (s *PostService) ListPosts(page, pageSize int, userID *uint) ([]models.Post, int64, error) {
	return s.postRepo.ListPosts(page, pageSize, userID)
}

func (s *PostService) CreateReport(userID uint, postID uint, req models.CreateReportRequest) error {
	report := &models.Report{
		PostID:     postID,
		ReporterID: userID,
		Reason:     req.Reason,
		Status:     "pending",
	}

	if err := s.postRepo.CreateReport(report); err != nil {
		return err
	}

	// Create audit log
	auditLog := &models.AuditLog{
		UserID:    &userID,
		Action:    "create_report",
		TableName: "reports",
		RecordID:  report.ID,
		NewValues: models.JSON{
			"post_id":    postID,
			"reason":     req.Reason,
			"status":     "pending",
			"created_at": time.Now(),
		},
	}
	return s.auditRepo.CreateLog(auditLog)
}

func (s *PostService) UpdateReportStatus(userID uint, reportID uint, status string) error {
	report, err := s.postRepo.GetReportByID(reportID)
	if err != nil {
		return err
	}

	oldStatus := report.Status
	if err := s.postRepo.UpdateReportStatus(reportID, status); err != nil {
		return err
	}

	// Create audit log
	auditLog := &models.AuditLog{
		UserID:    &userID,
		Action:    "update_report_status",
		TableName: "reports",
		RecordID:  reportID,
		OldValues: models.JSON{
			"status": oldStatus,
		},
		NewValues: models.JSON{
			"status": status,
		},
	}
	return s.auditRepo.CreateLog(auditLog)
}

func (s *PostService) ListReports(page, pageSize int) ([]models.Report, int64, error) {
	return s.postRepo.ListReports(page, pageSize)
}

func (s *PostService) GetPostHistory(postID uint) ([]models.PostHistory, error) {
	return s.postRepo.GetPostHistory(postID)
}
