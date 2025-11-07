package mysql

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/ABfry/album-battler/backend/internal/domain/entity"
	"github.com/google/uuid"
)

const (
	dsn = "album_user:album_pass@tcp(127.0.0.1:3306)/album_battler?parseTime=true"
)

func newTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("データベースへの接続に失敗しました: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("データベースへのPingに失敗しました: %v", err)
	}
	return db
}

// TestStep1_SaveUser はユーザーの保存をテストします。
func TestStep1_SaveUser(t *testing.T) {
	db := newTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Failed to close database: %v", err)
		}
	}()
	repo := NewUserRepository(db)

	userID := uuid.New()
	user := &entity.User{
		ID:             userID,
		Name:           "test-step1-save",
		IconUrl:        "http://example.com/step1.png",
		HashedPassword: "hashed_password_step1",
		CreatedAt:      time.Now(),
	}

	err := repo.Save(context.Background(), user)
	if err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	// データをクリーンアップします。
	_, err = db.Exec("DELETE FROM users WHERE id = ?", userID.String())
	if err != nil {
		t.Fatalf("Failed to clean up test data: %v", err)
	}
}

// TestStep2_FindUser はユーザーの取得をテストします。
func TestStep2_FindUser(t *testing.T) {
	db := newTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Failed to close database: %v", err)
		}
	}()
	repo := NewUserRepository(db)

	// テスト用のユーザーを作成
	userID := uuid.New()
	user := &entity.User{
		ID:             userID,
		Name:           "test-step2-find",
		IconUrl:        "http://example.com/step2.png",
		HashedPassword: "hashed_password_step2",
		CreatedAt:      time.Now(),
	}
	err := repo.Save(context.Background(), user)
	if err != nil {
		t.Fatalf("Failed to save user for find test: %v", err)
	}
	// テスト終了後にデータを削除
	defer func() {
		_, err := db.Exec("DELETE FROM users WHERE id = ?", userID.String())
		if err != nil {
			t.Fatalf("Failed to clean up test data: %v", err)
		}
	}()

	foundUser, err := repo.FindByID(context.Background(), userID)
	if err != nil {
		t.Fatalf("Failed to find user: %v", err)
	}
	if foundUser == nil {
		t.Fatal("User not found")
	}
	if foundUser.Name != "test-step2-find" {
		t.Errorf("Expected user name 'test-step2-find', but got %s", foundUser.Name)
	}
}

// TestStep3_UpdateUser はユーザーの更新をテストします。
func TestStep3_UpdateUser(t *testing.T) {
	db := newTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Failed to close database: %v", err)
		}
	}()
	repo := NewUserRepository(db)

	// テスト用のユーザーを作成
	userID := uuid.New()
	user := &entity.User{
		ID:             userID,
		Name:           "test-step3-update",
		IconUrl:        "http://example.com/step3.png",
		HashedPassword: "hashed_password_step3",
		CreatedAt:      time.Now(),
	}
	err := repo.Save(context.Background(), user)
	if err != nil {
		t.Fatalf("Failed to save user for update test: %v", err)
	}
	// テスト終了後にデータを削除
	defer func() {
		_, err := db.Exec("DELETE FROM users WHERE id = ?", userID.String())
		if err != nil {
			t.Fatalf("Failed to clean up test data: %v", err)
		}
	}()

	updatedUser := &entity.User{
		ID:             userID,
		Name:           "updated-name-step3",
		IconUrl:        "http://example.com/step3.png",
		HashedPassword: "hashed_password_step3",
		CreatedAt:      user.CreatedAt,
	}
	err = repo.Save(context.Background(), updatedUser)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	foundUser, err := repo.FindByID(context.Background(), userID)
	if err != nil {
		t.Fatalf("Failed to find updated user: %v", err)
	}
	if foundUser.Name != "updated-name-step3" {
		t.Errorf("Expected updated user name 'updated-name-step3', but got %s", foundUser.Name)
	}
}
