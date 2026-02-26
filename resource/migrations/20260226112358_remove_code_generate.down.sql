-- 代码生成业务信息表
CREATE TABLE IF NOT EXISTS setting_generate_tables (
    id bigint PRIMARY KEY,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp DEFAULT NULL,
    source_name varchar DEFAULT ''::varchar NOT NULL,
    table_name varchar DEFAULT ''::varchar NOT NULL,
    table_comment varchar DEFAULT ''::varchar NOT NULL,
    namespace varchar DEFAULT ''::varchar NOT NULL,
    module_name varchar DEFAULT ''::varchar NOT NULL,
    package_name varchar DEFAULT ''::varchar NOT NULL,
    business_name varchar DEFAULT ''::varchar NOT NULL,
    function_name varchar DEFAULT ''::varchar NOT NULL,
    function_author varchar DEFAULT ''::varchar NOT NULL,
    gen_type varchar DEFAULT 'zip'::varchar NOT NULL,
    gen_path varchar DEFAULT '/'::varchar NOT NULL,
    tpl_type varchar DEFAULT 'default'::varchar NULL,
    PRIMARY KEY (id)
);
COMMENT ON TABLE setting_generate_tables IS '代码生成业务信息表';

-- 代码生成业务字段信息表
CREATE TABLE IF NOT EXISTS setting_generate_columns (
    id bigint PRIMARY KEY,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp DEFAULT NULL,
    table_id bigint DEFAULT 0 NOT NULL,
    column_name varchar DEFAULT ''::varchar NOT NULL,
    column_comment varchar DEFAULT ''::varchar NOT NULL,
    column_type varchar DEFAULT ''::varchar NOT NULL,
    go_column_type varchar DEFAULT ''::varchar NOT NULL,
    go_field_type varchar DEFAULT ''::varchar NOT NULL,
    is_pk int DEFAULT 0 NOT NULL,
    is_increment int DEFAULT 0 NOT NULL,
    is_required int DEFAULT 0 NOT NULL,
    is_insert int DEFAULT 1 NOT NULL,
    is_edit int DEFAULT 1 NOT NULL,
    is_list int DEFAULT 1 NOT NULL,
    is_query int DEFAULT 0 NOT NULL,
    query_type varchar DEFAULT 'EQ'::varchar NOT NULL,
    html_type varchar DEFAULT 'input'::varchar NOT NULL,
    dict_type varchar DEFAULT ''::varchar NOT NULL,
    PRIMARY KEY (id)
);
COMMENT ON TABLE setting_generate_columns IS '代码生成业务字段信息表';

-- 恢复代码生成相关菜单权限
INSERT INTO system_menu (id, parent_id, parent_ids, menu_name, perms, menu_icon, menu_type, menu_url, route_name, route_path, route_component, route_redirect, sort_num, is_visible, status, is_cache, is_affix, created_by, created_at, updated_by, updated_at, remark, deleted_at) VALUES
       (4200, 4000, ',0,1000,4000,', '代码生成器', 'system:code', 'ma-icon-code', 'M', 'code', 'system/code/index', '/system/code', '/system/code/index', '', 2, 1, 1, 0, 1, 1, '2024-08-19 11:29:36', NULL, '2024-08-19 15:56:57', NULL, NULL),
       (4201, 4200, ',0,1000,4000,4200,', '预览代码', 'system:code:preview', '', 'B', '', '', '', '', 1, 1, 1, 0, 0, 1, '2024-08-19 11:29:36', NULL, '2024-08-19 11:29:36', NULL, NULL),
       (4202, 4200, ',0,1000,4000,4200,', '装载数据表', 'system:code:loadTable', '', 'B', '', '', '', '', 1, 1, 1, 0, 0, 1, '2024-08-19 11:29:36', NULL, '2024-08-19 11:29:36', NULL, NULL),
       (4203, 4200, ',0,1000,4000,4200,', '删除业务表', 'system:code:delete', '', 'B', '', '', '', '', 1, 1, 1, 0, 0, 1, '2024-08-19 11:29:36', NULL, '2024-08-19 11:29:36', NULL, NULL),
       (4204, 4200, ',0,1000,4000,4200,', '同步业务表', 'system:code:sync', '', 'B', '', '', '', '', 1, 1, 1, 0, 0, 1, '2024-08-19 11:29:36', NULL, '2024-08-19 11:29:36', NULL, NULL),
       (4205, 4200, ',0,1000,4000,4200,', '生成代码', 'system:code:generate', '', 'B', '', '', '', '', 1, 1, 1, 0, 0, 1, '2024-08-19 11:29:36', NULL, '2024-08-19 11:29:36', NULL, NULL),
       (4206, 4200, ',0,1000,4000,4200,', '代码生成列表', 'system:code:index', '', 'B', '', '', '', '', 1, 1, 1, 0, 0, 1, '2024-08-19 11:29:36', NULL, '2024-08-19 11:29:36', NULL, NULL);
