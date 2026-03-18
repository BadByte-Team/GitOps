CREATE DATABASE IF NOT EXISTS curso_db;
USE curso_db;

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role ENUM('admin', 'student') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE modules (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    is_hidden BOOLEAN DEFAULT FALSE
);

CREATE TABLE episodes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    module_id INT,
    title VARCHAR(100) NOT NULL,
    video_url VARCHAR(255) NOT NULL,
    is_hidden BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (module_id) REFERENCES modules(id) ON DELETE CASCADE
);

-- Contraseñas hasheadas con bcrypt (password original: admin123)
INSERT INTO users (username, password, role) VALUES 
('Leo', '$2a$10$0RAUOqn4rlh4JLn6CqPmAOCj/EDoN6E.oeJ00oUKtYB5eYGE0h55O', 'admin'),
('Gael', '$2a$10$0RAUOqn4rlh4JLn6CqPmAOCj/EDoN6E.oeJ00oUKtYB5eYGE0h55O', 'admin'),
('Jocelyn', '$2a$10$0RAUOqn4rlh4JLn6CqPmAOCj/EDoN6E.oeJ00oUKtYB5eYGE0h55O', 'admin'),
('Sergio', '$2a$10$0RAUOqn4rlh4JLn6CqPmAOCj/EDoN6E.oeJ00oUKtYB5eYGE0h55O', 'admin');

INSERT INTO modules (title) VALUES ('Modulo 1: VS Code'), ('Modulo 2: Git y GitHub');
INSERT INTO episodes (module_id, title, video_url) VALUES 
(1, 'EP 01: Instalacion', 'https://www.youtube.com/embed/dQw4w9WgXcQ'),
(2, 'EP 04: Llaves SSH', 'https://www.youtube.com/embed/dQw4w9WgXcQ');
