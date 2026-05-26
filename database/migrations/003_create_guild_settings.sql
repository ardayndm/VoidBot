CREATE TABLE
    IF NOT EXISTS guild_settings (
        id INT AUTO_INCREMENT PRIMARY KEY,
        guild_id VARCHAR(20) NOT NULL UNIQUE,
        prefix VARCHAR(5) DEFAULT '!',
        log_channel_id VARCHAR(20) DEFAULT NULL,
        welcome_channel_id VARCHAR(20) DEFAULT NULL,
        welcome_message TEXT DEFAULT NULL,
        auto_role_id VARCHAR(20) DEFAULT NULL,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    );