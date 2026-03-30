-- Initialize pgvector extension and create embeddings table
-- This script runs automatically when the PostgreSQL container starts for the first time
--
-- Tetrix CE uses 1024-dimension embeddings (Voyage AI voyage-code-3)

-- Enable pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- Create embeddings table
CREATE TABLE IF NOT EXISTS embeddings (
    id TEXT PRIMARY KEY,
    source_id TEXT NOT NULL DEFAULT '',
    symbol_id TEXT NOT NULL,
    symbol_name TEXT NOT NULL,
    symbol_kind TEXT NOT NULL,
    file_path TEXT NOT NULL,
    repository_id TEXT NOT NULL,
    line_start INTEGER NOT NULL,
    line_end INTEGER NOT NULL,
    content TEXT NOT NULL,
    embedding vector(1024) NOT NULL,
    metadata JSONB DEFAULT '{}',
    embedding_type TEXT DEFAULT 'technical',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_embeddings_symbol_id ON embeddings(symbol_id);
CREATE INDEX IF NOT EXISTS idx_embeddings_repository_id ON embeddings(repository_id);
CREATE INDEX IF NOT EXISTS idx_embeddings_file_path ON embeddings(file_path);
CREATE INDEX IF NOT EXISTS idx_embeddings_source_id ON embeddings(source_id);
CREATE INDEX IF NOT EXISTS idx_embeddings_embedding_type ON embeddings(embedding_type);

-- Create HNSW vector similarity search index
CREATE INDEX IF NOT EXISTS idx_embeddings_vector ON embeddings
USING hnsw (embedding vector_cosine_ops)
WITH (m = 16, ef_construction = 64);

-- Create function to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_embeddings_updated_at
    BEFORE UPDATE ON embeddings
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Grant permissions to the default app user
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'tetrix') THEN
        GRANT ALL PRIVILEGES ON TABLE embeddings TO tetrix;
        GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO tetrix;
    END IF;
END
$$;
