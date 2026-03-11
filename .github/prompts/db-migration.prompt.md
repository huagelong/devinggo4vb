---
name: db-migration
description: "在 DevingGo-Light 项目中创建和管理 PostgreSQL 数据库迁移文件。适用于：新建数据表、修改表结构（加字段/改字段/加索引）、插入初始化数据、回滚迁移等操作。"
---

# 数据库迁移

迁移文件存放于 `resource/migrations/`，命名格式：`{timestamp}_{name}.up.sql` / `.down.sql`。

## 命令

```bash
go run main.go migrate:create -name ${name}  # 创建迁移文件
go run main.go migrate:up                    # 执行迁移
go run main.go migrate:down                  # 回滚
```

## 新建表模板

**up.sql**：
```sql
CREATE TABLE ${module}_${resource} (
    id          bigint       NOT NULL,
    name        varchar(100) NOT NULL DEFAULT '',
    status      int2         NOT NULL DEFAULT 1,
    sort        int4         NOT NULL DEFAULT 0,
    remark      varchar(500) NOT NULL DEFAULT '',
    created_at  timestamp    WITHOUT TIME ZONE,
    updated_at  timestamp    WITHOUT TIME ZONE,
    deleted_at  timestamp    WITHOUT TIME ZONE,
    created_by  bigint,
    updated_by  bigint,
    CONSTRAINT ${module}_${resource}_pkey PRIMARY KEY (id)
);
COMMENT ON TABLE ${module}_${resource} IS '${中文说明}';
```

**down.sql**：`DROP TABLE IF EXISTS ${module}_${resource};`

## 常用 DDL 片段

```sql
-- 添加字段
ALTER TABLE ${table} ADD COLUMN IF NOT EXISTS ${col} ${type} NOT NULL DEFAULT ${default};
COMMENT ON COLUMN ${table}.${col} IS '${说明}';

-- 删除字段
ALTER TABLE ${table} DROP COLUMN IF EXISTS ${col};

-- 添加索引
CREATE INDEX IF NOT EXISTS idx_${table}_${col} ON ${table} (${col});
CREATE UNIQUE INDEX IF NOT EXISTS uk_${table}_${col} ON ${table} (${col}) WHERE deleted_at IS NULL;

-- 插入初始数据
INSERT INTO ${table} (...) VALUES (...) ON CONFLICT (id) DO NOTHING;
```

## 新建表后

```bash
make dao  # 重新生成 DAO/Entity/DO
```

