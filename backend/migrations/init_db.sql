SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

CREATE DATABASE IF NOT EXISTS album_battler CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE album_battler;

-- 開発用
DROP TABLE IF EXISTS images;
DROP TABLE IF EXISTS battles;
DROP TABLE IF EXISTS rooms;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id CHAR(36) PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    icon_url TEXT NOT NULL,
    hashed_password VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE rooms (
    id CHAR(36) PRIMARY KEY,
    room_number INT NOT NULL,
    host_user_id CHAR(36),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expired_at DATETIME,
    status ENUM('waiting', 'full', 'battling', 'result', 'closed') NOT NULL DEFAULT 'waiting',
    max_users INT NOT NULL DEFAULT 5,
    CHECK (room_number BETWEEN 0 AND 9999),
    CONSTRAINT fk_rooms_host_user FOREIGN KEY (host_user_id) REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    UNIQUE KEY uq_rooms_room_number_status (room_number, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE battles (
    id CHAR(36) PRIMARY KEY,
    room_id CHAR(36) NOT NULL,
    started_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_battles_room FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE images (
    id CHAR(36) PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    battle_id CHAR(36) NOT NULL,
    image_url TEXT NOT NULL,
    uploaded_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ai_score FLOAT,
    user_score INT,
    CHECK (ai_score IS NULL OR (ai_score >= 0 AND ai_score <= 100)),
    CHECK (user_score IS NULL OR user_score >= 0),
    CONSTRAINT fk_images_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_images_battle FOREIGN KEY (battle_id) REFERENCES battles(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE battle_users (
    battle_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    PRIMARY KEY (battle_id, user_id),
    CONSTRAINT fk_battle_users_battle FOREIGN KEY (battle_id) REFERENCES battles(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_battle_users_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE room_users (
    room_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    joined_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (room_id, user_id),
    CONSTRAINT fk_room_users_room FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_room_users_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_rooms_host_user ON rooms (host_user_id);
CREATE INDEX idx_battles_room ON battles (room_id);
CREATE INDEX idx_images_user ON images (user_id);
CREATE INDEX idx_images_battle ON images (battle_id);
CREATE INDEX idx_battle_users_user ON battle_users (user_id);
CREATE INDEX idx_room_users_user ON room_users (user_id);

-- 仮ユーザーデータの挿入 (開発用)
-- パスワードはすべて "password123" (bcryptでハッシュ化: $2a$10$YourHashedPasswordHere)
INSERT INTO users (id, name, icon_url, hashed_password, created_at) VALUES
    (
        '550e8400-e29b-41d4-a716-446655440001',
        'テストユーザー1',
        'https://api.dicebear.com/7.x/avataaars/svg?seed=test1',
        '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
        '2025-01-01 10:00:00'
    ),
    (
        '550e8400-e29b-41d4-a716-446655440002',
        'テストユーザー2',
        'https://api.dicebear.com/7.x/avataaars/svg?seed=test2',
        '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
        '2025-01-02 11:30:00'
    ),
    (
        '550e8400-e29b-41d4-a716-446655440003',
        'テストユーザー3',
        'https://api.dicebear.com/7.x/avataaars/svg?seed=test3',
        '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
        '2025-01-03 14:45:00'
    );
