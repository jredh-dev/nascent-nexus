-- Example database initialization
-- This creates sample tables for the database tool to query

CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS notes (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    tags VARCHAR(255)[],
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Sample data
INSERT INTO tasks (title, description, status) VALUES
    ('Setup toolbox', 'Initialize the toolbox project with Go and Docker', 'completed'),
    ('Implement database tool', 'Create tool for querying PostgreSQL', 'in_progress'),
    ('Add web UI', 'Build simple web interface', 'pending');

INSERT INTO notes (content, tags) VALUES
    ('Toolbox uses AGPL-3.0 to keep improvements open source', ARRAY['license', 'legal']),
    ('Tools define capability boundaries - no unrestricted system access', ARRAY['security', 'architecture']),
    ('Start with read-only tools, add write capabilities carefully', ARRAY['security', 'best-practice']);
