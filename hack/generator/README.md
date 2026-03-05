# DevingGo 代码生成工具

DevingGo代码生成工具集，提供模块管理、Worker任务和CRUD代码生成功能。

## 功能特性

### 1. 模块管理 (module) ✅
- `module:create` - 创建新模块 ✅
- `module:clone` - 克隆现有模块 ✅
- `module:export` - 导出模块包 ✅
- `module:import` - 导入模块包 ✅
- `module:list` - 列出已安装模块 ✅
- `module:validate` - 验证模块完整性 ✅

### 2. Worker任务生成 (worker) ⏳
- `worker:create` - 创建Worker任务（task/cron/both）⏳

### 3. CRUD代码生成 (crud) ⏳
- `crud:generate` - 生成CRUD代码 ⏳
- 支持单表生成和批量生成 ⏳

## 目录结构

```
hack/generator/
├── main.go                      # CLI入口
├── README.md                    # 本文档
├── generator.yaml               # 批量生成配置
├── cmd/                         # 命令层
│   ├── module.go                # 模块管理命令
│   ├── worker.go                # Worker任务命令
│   └── crud.go                  # CRUD生成命令
├── internal/                    # 核心逻辑层
│   ├── generator/               # 生成器引擎
│   │   ├── module.go            # 模块生成器
│   │   ├── worker.go            # Worker生成器
│   │   ├── crud.go              # CRUD生成器
│   │   └── template.go          # 模板渲染引擎
│   ├── scanner/                 # 扫描器
│   │   └── module.go            # 模块扫描器
│   ├── utils/                   # 工具函数
│   │   ├── zip.go               # 压缩/解压
│   │   ├── file.go              # 文件操作
│   │   └── naming.go            # 命名转换
│   └── config/                  # 配置管理
│       └── config.go            # 配置解析
└── templates/                   # 模板文件
    ├── module/                  # 模块模板
    ├── worker/                  # Worker模板
    └── crud/                    # CRUD模板
```

## 使用方法

### 运行工具

```bash
# 查看帮助
go run hack/generator/main.go -h

# 或通过Makefile（需配置）
make gen-module name=blog
make gen-worker name=test_task type=task
make gen-crud module=system table=users business=User
```

### 模块管理示例

```bash
# 创建新模块
go run hack/generator/main.go module:create -name blog

# 克隆现有模块
go run hack/generator/main.go module:clone -source blog -target news

# 导出模块
go run hack/generator/main.go module:export -name blog

# 导入模块
go run hack/generator/main.go module:import -file blog.v1.0.0.zip

# 列出模块
go run hack/generator/main.go module:list

# 验证模块
go run hack/generator/main.go module:validate -name blog
```

### Worker任务示例

```bash
# 创建异步任务
go run hack/generator/main.go worker:create -name=send_email -type=task

# 创建定时任务
go run hack/generator/main.go worker:create -name=clean_cache -type=cron

# 创建混合任务
go run hack/generator/main.go worker:create -name=data_sync -type=both
```

### CRUD生成示例

```bash
# 单表生成
go run hack/generator/main.go crud:generate \
  -module=system \
  -table=users \
  -business=User

# 批量生成（使用配置文件）
go run hack/generator/main.go crud:generate -config=generator.yaml
```

## 配置文件

### generator.yaml

批量生成CRUD代码的配置文件：

```yaml
# 目标模块
module: system

# 生成配置列表
tables:
  - table: users
    business: User
    description: 用户管理
    
  - table: roles
    business: Role
    description: 角色管理
    
  - table: permissions
    business: Permission
    description: 权限管理
```

## 开发指南

### 添加新命令

1. 在 `cmd/` 目录创建命令文件
2. 实现 `gcmd.Command` 结构
3. 在 `main.go` 中注册命令
4. 在 `internal/generator/` 实现核心逻辑

### 添加新模板

1. 在 `templates/` 对应目录创建模板文件
2. 使用 `.tpl` 后缀
3. 使用 `{{.Variable}}` 语法定义变量
4. 在生成器中调用模板渲染

## 工具函数

### naming.go - 命名转换

```go
// snake_case
ToSnakeCase("CategoryName") // "category_name"

// camelCase
ToCamelCase("category_name") // "categoryName"

// PascalCase
ToPascalCase("category_name") // "CategoryName"

// CONST_CASE
ToConstCase("CategoryName") // "CATEGORY_NAME"

// kebab-case
ToKebabCase("CategoryName") // "category-name"
```

### file.go - 文件操作

```go
// 检查路径
PathExists(path) bool
IsDir(path) bool
IsFile(path) bool

// 目录操作
EnsureDir(dir) error
GetProjectRoot() (string, error)
GetModuleName() (string, error)

// 文件操作
CopyFile(src, dst) error
WriteFile(path, content, overwrite) error

// 代码格式化
FormatGoCode(filePath) error
FormatGoCodeInDir(dir) error
```

### zip.go - 压缩/解压

```go
// 压缩目录
ZipDirectory(ctx, srcDir, dstZip) error

// 解压文件
UnzipFile(srcZip, dstDir) error
```

## 注意事项

1. **独立性**: 工具不依赖项目业务代码，可独立编译和分发
2. **安全性**: 解压文件时会检查路径穿越攻击
3. **代码格式化**: 生成的Go代码会自动格式化（gofmt/goimports）
4. **模板变量**: 使用 `{{.Variable}}` 格式，遵循Go template语法

## 版本历史

- **v1.0.0** (2026-03-04)
  - ✅ 完成基础架构搭建
  - ✅ 完成模块管理功能（创建/克隆/导入/导出/列表/验证）
  - ✅ 实现.module.yaml配置系统
  - ✅ 实现变量替换引擎
  - ✅ 创建17个模块模板文件
  - ⏳ Worker任务生成（待开发）
  -.module.yaml 配置规范](docs/MODULE_YAML_SPEC.md) ✅
- [ ⏳ CRUD代码生成（待开发）

## 相关文档

- [技术方案](../../docs/PLAN-CodeGeneratorUnification.md)
- [任务清单](../../docs/TODO-CodeGeneratorUnification.md)
- [模块开发指南](../../docs/MODULE.md) - 待创建
- [Worker开发指南](../../docs/WORKER.md) - 待创建

## 许可证

MIT License
