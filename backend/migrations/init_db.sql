
CREATE DATABASE IF NOT EXISTS album_battler CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE album_battler;

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
    host_user_id CHAR(36) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expired_at DATETIME,
    status ENUM('waiting', 'full', 'battling', 'result', 'closed') NOT NULL DEFAULT 'waiting',
    CHECK (room_number BETWEEN 0 AND 9999),
    CONSTRAINT fk_rooms_host_user FOREIGN KEY (host_user_id) REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    UNIQUE KEY uq_rooms_room_number (room_number)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE battles (
    id CHAR(36) PRIMARY KEY,
    room_id CHAR(36) NOT NULL,
    user_1_id CHAR(36) NOT NULL,
    user_2_id CHAR(36) NOT NULL,
    winner_id CHAR(36),
    started_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_battles_room FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_battles_user1 FOREIGN KEY (user_1_id) REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT fk_battles_user2 FOREIGN KEY (user_2_id) REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT fk_battles_winner FOREIGN KEY (winner_id) REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE,
    CONSTRAINT uq_battles_participants UNIQUE (room_id, user_1_id, user_2_id),
    CHECK (user_1_id <> user_2_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE images (
    id CHAR(36) PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    battle_id CHAR(36) NOT NULL,
    image_url TEXT NOT NULL,
    uploaded_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ai_score DECIMAL(5,2),
    audience_score INT,
    CHECK (ai_score IS NULL OR (ai_score >= 0 AND ai_score <= 100)),
    CHECK (audience_score IS NULL OR audience_score >= 0),
    CONSTRAINT fk_images_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_images_battle FOREIGN KEY (battle_id) REFERENCES battles(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_rooms_host_user ON rooms (host_user_id);
CREATE INDEX idx_battles_room ON battles (room_id);
CREATE INDEX idx_battles_user1 ON battles (user_1_id);
CREATE INDEX idx_battles_user2 ON battles (user_2_id);
CREATE INDEX idx_battles_winner ON battles (winner_id);
CREATE INDEX idx_images_user ON images (user_id);
CREATE INDEX idx_images_battle ON images (battle_id);
