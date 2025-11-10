package mysql

import (
	"context"
	"testing"
	"time"

	"github.com/ABfry/album-battler/backend/internal/domain/entity"
	"github.com/google/uuid"
)

// TestUserRepository は、UserRepositoryのテストをテーブル駆動で実行します。
func TestUserRepository(t *testing.T) {
	// --- テストのセットアップ ---
	db := newTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Logf("Failed to close db: %v", err)
		}
	}()
	repo := NewUserRepository(db)
	ctx := context.Background()

	// --- テストケースの定義 ---
	// FindByID のテストケース
	findUserID := uuid.New()
	findUserTime := time.Now().Truncate(time.Second)
	findUser := &entity.User{
		ID:             findUserID,
		Name:           "test-user-find",
		IconUrl:        "http://example.com/find.png",
		HashedPassword: "hashed_password_find",
		CreatedAt:      findUserTime,
	}

	// Save(Update) のテストケース
	saveUserID := uuid.New()
	saveUserTime := time.Now().Truncate(time.Second)
	initialUser := &entity.User{
		ID:             saveUserID,
		Name:           "initial-user",
		IconUrl:        "http://example.com/initial.png",
		HashedPassword: "hashed_initial",
		CreatedAt:      saveUserTime,
	}
	updatedUser := &entity.User{
		ID:             saveUserID,
		Name:           "updated-user",
		IconUrl:        "http://example.com/updated.png",
		HashedPassword: "hashed_updated",
		CreatedAt:      saveUserTime, // Save処理ではCreatedAtは更新されない
	}

	// --- テストの実行 ---
	t.Run("FindByID", func(t *testing.T) {
		// FindByIDのテスト用のデータを準備
		if err := repo.Save(ctx, findUser); err != nil {
			t.Fatalf("Failed to save user for find test: %v", err)
		}
		t.Cleanup(func() { cleanupTestUser(db, "users", findUserID, t) })

		// テーブル駆動テストの定義
		testCases := []struct {
			name        string
			inputID     uuid.UUID
			wantUser    *entity.User
			expectError bool
		}{
			{
				name:        "存在するユーザーをIDで取得できる",
				inputID:     findUserID,
				wantUser:    findUser,
				expectError: false,
			},
			{
				name:        "存在しないユーザーをIDで取得するとnilが返る",
				inputID:     uuid.New(),
				wantUser:    nil,
				expectError: false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				gotUser, err := repo.FindByID(ctx, tc.inputID)

				if tc.expectError {
					if err == nil {
						t.Errorf("Expected an error, but got none")
					}
				} else {
					if err != nil {
						t.Errorf("Unexpected error: %v", err)
					}
					assertEqual(t, tc.wantUser, gotUser)
				}
			})
		}
	})

	t.Run("Save", func(t *testing.T) {
		// テーブル駆動テストの定義
		testCases := []struct {
			name     string
			user     *entity.User
			wantUser *entity.User
		}{
			{
				name:     "新しいユーザーを保存できる",
				user:     initialUser,
				wantUser: initialUser,
			},
			{
				name:     "既存のユーザー情報を更新できる",
				user:     updatedUser,
				wantUser: updatedUser,
			},
		}

		t.Cleanup(func() { cleanupTestUser(db, "users", saveUserID, t) })

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := repo.Save(ctx, tc.user)
				if err != nil {
					t.Fatalf("Failed to save user: %v", err)
				}

				t.Logf("Table 'users' after saving user '%s':", tc.user.Name)
				showTable(t, db, "users")

				// 保存/更新したデータを取得して検証
				foundUser, err := repo.FindByID(ctx, tc.user.ID)
				if err != nil {
					t.Fatalf("Failed to find user after save: %v", err)
				}
				assertEqual(t, tc.wantUser, foundUser)
			})
		}
	})

	t.Run("SaveMultipleUsers", func(t *testing.T) {
		// テーブル駆動テストの定義
		testCases := []struct {
			name  string
			users []*entity.User
		}{
			{
				name: "複数のユーザーを一度に保存できる",
				users: []*entity.User{
					{
						ID:             uuid.New(),
						Name:           "multi-user-1",
						IconUrl:        "http://example.com/multi1.png",
						HashedPassword: "hashed_multi_1",
						CreatedAt:      time.Now().Truncate(time.Second),
					},
					{
						ID:             uuid.New(),
						Name:           "multi-user-2",
						IconUrl:        "http://example.com/multi2.png",
						HashedPassword: "hashed_multi_2",
						CreatedAt:      time.Now().Truncate(time.Second),
					},
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				for _, user := range tc.users {
					// ループ変数をキャプチャするため、ローカル変数にコピーする
					u := user
					t.Cleanup(func() { cleanupTestUser(db, "users", u.ID, t) })
					err := repo.Save(ctx, u)
					if err != nil {
						t.Fatalf("Failed to save user %s: %v", u.Name, err)
					}
				}

				t.Logf("Table 'users' after saving multiple users:")
				showTable(t, db, "users")

				// 全員が正しく保存されたか検証
				for _, u := range tc.users {
					foundUser, err := repo.FindByID(ctx, u.ID)
					if err != nil {
						t.Fatalf("Failed to find user %s: %v", u.Name, err)
					}
					assertEqual(t, u, foundUser)
				}
			})
		}
	})
}
