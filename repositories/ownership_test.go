package repositories

import (
	"errors"
	models "moodly/Models"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newRepositoryTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.MoodLog{}, &models.CustomCause{}); err != nil {
		t.Fatalf("migrate test db: %v", err)
	}

	return db
}

func TestMoodLogUpdateScopesByUser(t *testing.T) {
	db := newRepositoryTestDB(t)
	repo := NewMoodLogsRepository(db)

	moodLog := models.MoodLog{
		UserID: 1,
		Mood:   2,
		Note:   "old note",
		Causes: "work",
	}
	if err := db.Create(&moodLog).Error; err != nil {
		t.Fatalf("create mood log: %v", err)
	}

	err := repo.UpdateMoodLog(&models.MoodLog{
		ID:     moodLog.ID,
		UserID: 2,
		Mood:   5,
		Note:   "wrong user",
		Causes: "family",
	})
	if err == nil {
		t.Fatal("expected wrong-user update to fail")
	}

	var unchanged models.MoodLog
	if err := db.First(&unchanged, moodLog.ID).Error; err != nil {
		t.Fatalf("reload mood log: %v", err)
	}
	if unchanged.Mood != 2 || unchanged.Note != "old note" || unchanged.Causes != "work" {
		t.Fatalf("wrong-user update changed mood log: %+v", unchanged)
	}

	err = repo.UpdateMoodLog(&models.MoodLog{
		ID:     moodLog.ID,
		UserID: 1,
		Mood:   4,
		Note:   "updated note",
		Causes: "sleep",
	})
	if err != nil {
		t.Fatalf("same-user update failed: %v", err)
	}

	var updated models.MoodLog
	if err := db.First(&updated, moodLog.ID).Error; err != nil {
		t.Fatalf("reload updated mood log: %v", err)
	}
	if updated.Mood != 4 || updated.Note != "updated note" || updated.Causes != "sleep" {
		t.Fatalf("same-user update did not persist: %+v", updated)
	}
}

func TestMoodLogDeleteScopesByUser(t *testing.T) {
	db := newRepositoryTestDB(t)
	repo := NewMoodLogsRepository(db)

	moodLog := models.MoodLog{
		UserID: 1,
		Mood:   3,
		Causes: "work",
	}
	if err := db.Create(&moodLog).Error; err != nil {
		t.Fatalf("create mood log: %v", err)
	}

	if err := repo.DeleteMoodLog(moodLog.ID, 2); err == nil {
		t.Fatal("expected wrong-user delete to fail")
	}

	var existing models.MoodLog
	if err := db.First(&existing, moodLog.ID).Error; err != nil {
		t.Fatalf("wrong-user delete removed mood log: %v", err)
	}

	if err := repo.DeleteMoodLog(moodLog.ID, 1); err != nil {
		t.Fatalf("same-user delete failed: %v", err)
	}

	err := db.First(&models.MoodLog{}, moodLog.ID).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected mood log to be deleted, got %v", err)
	}
}

func TestCustomCauseUpdateScopesByUser(t *testing.T) {
	db := newRepositoryTestDB(t)
	repo := NewCustomCauseRepository(db)

	cause := models.CustomCause{
		UserID: 1,
		Name:   "work",
	}
	if err := db.Create(&cause).Error; err != nil {
		t.Fatalf("create custom cause: %v", err)
	}

	err := repo.Update(&models.CustomCause{
		ID:     cause.ID,
		UserID: 2,
		Name:   "wrong user",
	})
	if err == nil {
		t.Fatal("expected wrong-user update to fail")
	}

	var unchanged models.CustomCause
	if err := db.First(&unchanged, cause.ID).Error; err != nil {
		t.Fatalf("reload custom cause: %v", err)
	}
	if unchanged.Name != "work" || unchanged.UserID != 1 {
		t.Fatalf("wrong-user update changed custom cause: %+v", unchanged)
	}

	err = repo.Update(&models.CustomCause{
		ID:     cause.ID,
		UserID: 1,
		Name:   "sleep",
	})
	if err != nil {
		t.Fatalf("same-user update failed: %v", err)
	}

	var updated models.CustomCause
	if err := db.First(&updated, cause.ID).Error; err != nil {
		t.Fatalf("reload updated custom cause: %v", err)
	}
	if updated.Name != "sleep" {
		t.Fatalf("same-user update did not persist: %+v", updated)
	}
}

func TestCustomCauseDeleteScopesByUser(t *testing.T) {
	db := newRepositoryTestDB(t)
	repo := NewCustomCauseRepository(db)

	cause := models.CustomCause{
		UserID: 1,
		Name:   "work",
	}
	if err := db.Create(&cause).Error; err != nil {
		t.Fatalf("create custom cause: %v", err)
	}

	if err := repo.Delete(cause.ID, 2); err == nil {
		t.Fatal("expected wrong-user delete to fail")
	}

	var existing models.CustomCause
	if err := db.First(&existing, cause.ID).Error; err != nil {
		t.Fatalf("wrong-user delete removed custom cause: %v", err)
	}

	if err := repo.Delete(cause.ID, 1); err != nil {
		t.Fatalf("same-user delete failed: %v", err)
	}

	err := db.First(&models.CustomCause{}, cause.ID).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected custom cause to be deleted, got %v", err)
	}
}
