package services

import (
	"errors"
	"fmt"
	"time"

	"carbuyer/internal/db/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ThreadService struct {
	db *gorm.DB
}

func NewThreadService(db *gorm.DB) *ThreadService {
	return &ThreadService{
		db: db,
	}
}

// CreateThread creates a new thread for a user
func (s *ThreadService) CreateThread(userID uuid.UUID, sellerName string, sellerType models.SellerType) (*models.Thread, error) {
	// Validate input
	if sellerName == "" {
		return nil, errors.New("seller name is required")
	}

	if sellerType != models.SellerTypePrivate && sellerType != models.SellerTypeDealership && sellerType != models.SellerTypeOther {
		return nil, errors.New("invalid seller type")
	}

	// Create thread
	thread := &models.Thread{
		UserID:     userID,
		SellerName: sellerName,
		SellerType: sellerType,
	}

	if err := s.db.Create(thread).Error; err != nil {
		return nil, fmt.Errorf("failed to create thread: %w", err)
	}

	return thread, nil
}

// ThreadWithCounts represents a thread with calculated counts and preview
type ThreadWithCounts struct {
	models.Thread
	MessageCount      int64  `json:"messageCount"`
	UnreadCount       int64  `json:"unreadCount"`
	LastMessagePreview string `json:"lastMessagePreview"`
	DisplayName       string `json:"displayName"`
}

// GetUserThreads retrieves all threads for a user with calculated counts and previews
func (s *ThreadService) GetUserThreads(userID uuid.UUID) ([]ThreadWithCounts, error) {
	var threads []models.Thread
	if err := s.db.Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("COALESCE(last_message_at, created_at) DESC").
		Find(&threads).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve threads: %w", err)
	}

	result := make([]ThreadWithCounts, len(threads))
	for i, thread := range threads {
		// Calculate message count
		var messageCount int64
		if err := s.db.Model(&models.Message{}).
			Where("thread_id = ? AND deleted_at IS NULL", thread.ID).
			Count(&messageCount).Error; err != nil {
			return nil, fmt.Errorf("failed to count messages: %w", err)
		}

		// Calculate unread count
		var unreadCount int64
		query := s.db.Model(&models.Message{}).
			Where("thread_id = ? AND deleted_at IS NULL", thread.ID)
		
		if thread.LastReadAt != nil {
			query = query.Where("timestamp > ?", thread.LastReadAt)
		}
		
		if err := query.Count(&unreadCount).Error; err != nil {
			return nil, fmt.Errorf("failed to count unread messages: %w", err)
		}

		// Get last message preview
		var lastMessage models.Message
		lastMessagePreview := ""
		if err := s.db.Where("thread_id = ? AND deleted_at IS NULL", thread.ID).
			Order("timestamp DESC").
			First(&lastMessage).Error; err == nil {
			// Truncate to ~50 characters
			preview := lastMessage.Content
			if len(preview) > 50 {
				preview = preview[:47] + "..."
			}
			lastMessagePreview = preview
		}

		// Calculate display name
		displayName := thread.SellerName
		if thread.SellerName == thread.Phone || thread.SellerName == "" {
			displayName = thread.Phone
		}

		result[i] = ThreadWithCounts{
			Thread:            thread,
			MessageCount:      messageCount,
			UnreadCount:       unreadCount,
			LastMessagePreview: lastMessagePreview,
			DisplayName:       displayName,
		}
	}

	return result, nil
}

// GetThreadByID retrieves a specific thread
func (s *ThreadService) GetThreadByID(threadID, userID uuid.UUID) (*models.Thread, error) {
	var thread models.Thread
	if err := s.db.Where("id = ? AND user_id = ? AND deleted_at IS NULL", threadID, userID).First(&thread).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("thread not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &thread, nil
}

// DeleteThread deletes a thread (soft delete in future)
func (s *ThreadService) DeleteThread(threadID, userID uuid.UUID) error {
	result := s.db.Where("id = ? AND user_id = ?", threadID, userID).Delete(&models.Thread{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete thread: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("thread not found")
	}

	return nil
}

// ArchiveThread soft deletes a thread by setting deleted_at
func (s *ThreadService) ArchiveThread(threadID, userID uuid.UUID) error {
	// Verify the thread exists, belongs to the user, and is not already archived
	var thread models.Thread
	if err := s.db.Where("id = ? AND user_id = ? AND deleted_at IS NULL", threadID, userID).First(&thread).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("thread not found")
		}
		return fmt.Errorf("failed to verify thread: %w", err)
	}

	// Soft delete the thread by setting deleted_at
	now := time.Now()
	thread.DeletedAt = &now
	if err := s.db.Save(&thread).Error; err != nil {
		return fmt.Errorf("failed to archive thread: %w", err)
	}

	return nil
}

// UpdateThreadName updates a thread's seller name
func (s *ThreadService) UpdateThreadName(threadID, userID uuid.UUID, sellerName string) (*models.Thread, error) {
	// Validate input
	if sellerName == "" {
		return nil, errors.New("seller name is required")
	}

	// Verify thread exists, belongs to user, and is not archived
	var thread models.Thread
	if err := s.db.Where("id = ? AND user_id = ? AND deleted_at IS NULL", threadID, userID).First(&thread).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("thread not found")
		}
		return nil, fmt.Errorf("failed to find thread: %w", err)
	}

	// Update seller name
	thread.SellerName = sellerName
	thread.UpdatedAt = time.Now()

	if err := s.db.Save(&thread).Error; err != nil {
		return nil, fmt.Errorf("failed to update thread: %w", err)
	}

	return &thread, nil
}

// MarkThreadAsRead marks a thread as read by setting last_read_at to now
func (s *ThreadService) MarkThreadAsRead(threadID, userID uuid.UUID) error {
	// Verify the thread exists and belongs to the user
	var thread models.Thread
	if err := s.db.Where("id = ? AND user_id = ? AND deleted_at IS NULL", threadID, userID).First(&thread).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("thread not found")
		}
		return fmt.Errorf("failed to verify thread: %w", err)
	}

	// Set last_read_at to now
	now := time.Now()
	if err := s.db.Model(&models.Thread{}).Where("id = ?", threadID).Update("last_read_at", now).Error; err != nil {
		return fmt.Errorf("failed to mark thread as read: %w", err)
	}

	return nil
}

// selectParentThread determines which thread should be the parent
// Logic: Prefer thread with displayName (SellerName != Phone), fallback to first in selection order
func selectParentThread(threads []models.Thread, orderedIDs []uuid.UUID) models.Thread {
	// Create map for O(1) lookup
	threadMap := make(map[uuid.UUID]models.Thread)
	for _, t := range threads {
		threadMap[t.ID] = t
	}

	// First pass: find threads with meaningful displayName (not just phone)
	var threadsWithName []models.Thread
	var threadsWithoutName []models.Thread

	for _, id := range orderedIDs {
		thread, exists := threadMap[id]
		if !exists {
			continue
		}

		// Check if has meaningful display name (SellerName exists and is not just the phone)
		hasDisplayName := thread.SellerName != "" && thread.SellerName != thread.Phone

		if hasDisplayName {
			threadsWithName = append(threadsWithName, thread)
		} else {
			threadsWithoutName = append(threadsWithoutName, thread)
		}
	}

	// Return first with displayName, or first in order if none have displayName
	if len(threadsWithName) > 0 {
		return threadsWithName[0]
	}
	if len(threadsWithoutName) > 0 {
		return threadsWithoutName[0]
	}

	// Fallback (shouldn't happen with validation)
	return threads[0]
}

// ConsolidateThreads consolidates multiple threads into a single parent thread
// Parent selection logic:
// 1. Prefer thread with non-empty displayName over empty/phone-only displayName
// 2. If tied, use first thread in selection order
func (s *ThreadService) ConsolidateThreads(userID uuid.UUID, threadIDs []uuid.UUID) (*ThreadWithCounts, error) {
	// Validation
	if len(threadIDs) < 2 {
		return nil, errors.New("at least 2 threads required for consolidation")
	}

	// Begin transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Fetch all threads and verify ownership + existence
	var threads []models.Thread
	if err := tx.Where("id IN ? AND user_id = ? AND deleted_at IS NULL", threadIDs, userID).
		Find(&threads).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to fetch threads: %w", err)
	}

	// Validate all threads were found
	if len(threads) != len(threadIDs) {
		tx.Rollback()
		return nil, errors.New("one or more threads not found or already archived")
	}

	// Determine parent thread using selection logic
	parentThread := selectParentThread(threads, threadIDs)

	// Collect IDs of threads to archive (all except parent)
	var sourceThreadIDs []uuid.UUID
	for _, t := range threads {
		if t.ID != parentThread.ID {
			sourceThreadIDs = append(sourceThreadIDs, t.ID)
		}
	}

	// Move all messages from source threads to parent
	if err := tx.Model(&models.Message{}).
		Where("thread_id IN ? AND deleted_at IS NULL", sourceThreadIDs).
		Update("thread_id", parentThread.ID).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to move messages: %w", err)
	}

	// Move all tracked offers from source threads to parent
	if err := tx.Model(&models.TrackedOffer{}).
		Where("thread_id IN ?", sourceThreadIDs).
		Update("thread_id", parentThread.ID).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to move tracked offers: %w", err)
	}

	// Update parent's LastMessageAt to latest across all threads
	var latestMessageTime *time.Time
	for _, t := range threads {
		if t.LastMessageAt != nil {
			if latestMessageTime == nil || t.LastMessageAt.After(*latestMessageTime) {
				latestMessageTime = t.LastMessageAt
			}
		}
	}
	if latestMessageTime != nil && (parentThread.LastMessageAt == nil || latestMessageTime.After(*parentThread.LastMessageAt)) {
		parentThread.LastMessageAt = latestMessageTime
		if err := tx.Save(&parentThread).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update parent thread: %w", err)
		}
	}

	// Archive source threads (soft delete)
	now := time.Now()
	if err := tx.Model(&models.Thread{}).
		Where("id IN ?", sourceThreadIDs).
		Update("deleted_at", now).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to archive source threads: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return consolidated thread with counts (reuse existing logic)
	// We need to fetch it with the same service instance to get counts
	threadsWithCounts, err := s.GetUserThreads(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated thread: %w", err)
	}

	// Find the parent thread in the results
	for _, t := range threadsWithCounts {
		if t.ID == parentThread.ID {
			return &t, nil
		}
	}

	return nil, errors.New("consolidated thread not found")
}
