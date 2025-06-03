-- Create post_history table for backup
CREATE TABLE post_history (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    post_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    address TEXT NOT NULL,
    contact_name VARCHAR(255) NOT NULL,
    mobile_number VARCHAR(20) NOT NULL,
    incident_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL,
    is_anonymous BOOLEAN,
    visibility VARCHAR(20) NOT NULL,
    allow_comments BOOLEAN,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
); 