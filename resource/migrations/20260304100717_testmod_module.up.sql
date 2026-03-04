-- Create testmod module tables
-- Author: devinggo
-- Date: 2026-03-04 18:07:17

-- Example table for testmod module
CREATE TABLE IF NOT EXISTS testmod_example (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL DEFAULT '',
    description TEXT,
    status SMALLINT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create index
CREATE INDEX idx_testmod_example_status ON testmod_example(status);
CREATE INDEX idx_testmod_example_deleted_at ON testmod_example(deleted_at);

-- Add comment
COMMENT ON TABLE testmod_example IS 'Testmod模块示例表';
COMMENT ON COLUMN testmod_example.id IS 'ID';
COMMENT ON COLUMN testmod_example.name IS '名称';
COMMENT ON COLUMN testmod_example.description IS '描述';
COMMENT ON COLUMN testmod_example.status IS '状态:1=启用,0=禁用';
COMMENT ON COLUMN testmod_example.created_at IS '创建时间';
COMMENT ON COLUMN testmod_example.updated_at IS '更新时间';
COMMENT ON COLUMN testmod_example.deleted_at IS '删除时间';
