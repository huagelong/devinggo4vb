-- 删除代码生成相关的菜单权限
DELETE FROM system_menu WHERE id IN (4200, 4201, 4202, 4203, 4204, 4205, 4206);

-- 删除代码生成业务字段信息表
DROP TABLE IF EXISTS setting_generate_columns CASCADE;

-- 删除代码生成业务信息表
DROP TABLE IF EXISTS setting_generate_tables CASCADE;
