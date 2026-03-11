# DevingGo代码生成工具统一架构改造方案

基于DevingGo项目现状，将分散的代码生成命令（module:create/export/import、worker:create）与新的CRUD生成器统一整合到 `hack/generator` 工具集，形成独立、可复用、功能完整的开发工具链。核心目标：**统一架构、提升效率、独立可维护**。

## 一、现状分析

**当前问题**：
- 代码生成功能分散在 `modules/system/cmd/` 目录，属于业务模块
- 依赖 `devinggo/modules/system/pkg/utils`，耦合度高
- 缺少统一的CRUD代码生成工具，存在大量重复模板代码
- 模板文件分散在 `resource/generate/`，与生成器代码分离

**迁移范围**：
- `modules/system/cmd/module.go` - 3个模块管理命令（456行）
- `modules/system/cmd/worker_create.go` - Worker创建命令（418行）
- 新增：CRUD代码生成器
- 模板迁移：`resource/generate/modules/` 的15个模板文件

**模块管理命令优化方向**：
- 统一命令风格：`module <action>` 替代 `module:<action>`
- 增强交互体验：支持交互式创建、配置验证、依赖检查
- 扩展功能集：新增验证、升级、克隆、列表等命令
- 模板系统：支持自定义模块模板和脚手架

## 二、目标架构

**新的目录结构**：
```
hack/generator/
├── main.go                      # CLI入口（gcmd路由）
├── README.md                    # 工具使用文档
├── generator.yaml               # 批量生成配置
├── cmd/                         # 命令层
│   ├── module.go                # module:create/export/import
│   ├── worker.go                # worker:create
│   └── crud.go                  # crud:generate（新增）
├── internal/                    # 核心逻辑层
│   ├── generator/               # 生成器引擎
│   │   ├── module.go            # 模块骨架生成器
│   │   ├── worker.go            # Worker任务生成器
│   │   ├── crud.go              # CRUD代码生成器（新增）
│   │   └── template.go          # 模板渲染引擎
│   ├── scanner/                 # 扫描器
│   │   └── module.go            # 模块发现与校验
│   ├── utils/                   # 独立工具函数
│   │   ├── zip.go               # 压缩/解压（标准库实现）
│   │   ├── file.go              # 文件操作辅助
│   │   └── naming.go            # 命名转换工具
│   └── config/                  # 配置管理
│       └── config.go            # generator.yaml解析
└── templates/                   # 模板文件
    ├── module/                  # 模块骨架模板（15个文件）
    │   ├── module.go.tpl
    │   ├── logic.go.tpl
    │   ├── migration_up.sql.tpl
    │   └── ...
    ├── worker/                  # Worker任务模板
    │   ├── cron.go.tpl
    │   ├── task.go.tpl
    │   └── const.go.tpl
    └── crud/                    # CRUD模板（新增）
        ├── req.go.tpl
        ├── res.go.tpl
        ├── api.go.tpl
        ├── controller.go.tpl
        └── logic.go.tpl
```

## 三、实施步骤

### 阶段一：基础架构搭建（第1天）

**1. 创建项目结构**
- 创建 `hack/generator/` 及所有子目录
- 创建 `main.go` CLI入口：
  ```go
  package main
  
  import (
      "devinggo/hack/generator/cmd"
      "github.com/gogf/gf/v2/os/gcmd"
  )
  
  func main() {
      command := &gcmd.Command{
          Name: "generator",
          Brief: "DevingGo代码生成工具集",
      }
      
      command.AddCommand(
          cmd.CreateModule,   // module:create
          cmd.ExportModule,   // module:export
          cmd.ImportModule,   // module:import
          cmd.CreateWorker,   // worker:create
          cmd.GenerateCrud,   // crud:generate（新增）
      )
      
      command.Run(gcmd.GetCtx())
  }
  ```

**2. 实现独立工具函数**
- `internal/utils/zip.go` - 使用 `archive/zip` 标准库实现：
  - `ZipDirectory(ctx, srcDir, dstZip string) error`
  - `UnzipFile(srcZip, dstDir string) error`
- `internal/utils/naming.go` - 命名转换：
  - `ToSnakeCase(s string) string` - category_name
  - `ToCamelCase(s string) string` - CategoryName
  - `ToConstCase(s string) string` - CATEGORY_NAME
- `internal/utils/file.go` - 文件辅助：
  - `GetTmpDir() string`
  - `ValidateModuleName(name string) error`
  - `FormatGoCode(filePath string) error` - 调用gofmt/goimports

### 阶段二：优化并迁移模块管理命令（第2天）

#### 命令格式说明

**保持原有格式**（使用冒号分隔）：
```bash
go run hack/generator/main.go module:create -name=blog [-template=standard]
go run hack/generator/main.go module:export -name=blog [-output=./exports]
go run hack/generator/main.go module:import -file=blog.zip [-force]
go run hack/generator/main.go module:list [-status=enabled]
go run hack/generator/main.go module:validate -name=blog
go run hack/generator/main.go module:clone -source=system -target=newmodule
go run hack/generator/main.go module:upgrade -name=blog -version=2.0.0
```

**设计理念**：
- 保持与现有命令风格一致（`module:action`）
- 支持参数化调用和交互式调用
- 无参数时自动进入交互式引导模式
- 所有可选参数使用 `-flag=value` 格式

#### 3. 重构 module create（交互式创建）

**源文件**：`modules/system/cmd/module.go#L10-L99`  
**目标**：`hack/generator/cmd/module.go` + `internal/generator/module.go`

**核心优化**：

a) **支持交互式模式**
```go
// 无参数时进入交互式
$ go run hack/generator/main.go module:create

✨ 创建新模块
? 模块名称: blog
? 模块描述: 博客管理模块
? 作者: DevingGo Team
? 选择模块类型:
  ▸ 标准模块（完整功能）
    轻量模块（仅API）
    Worker模块（仅后台任务）
? 是否包含数据库表? Yes
? 表名前缀: blog_
? 是否需要Worker? No
? 是否需要WebSocket? No

✅ 模块 'blog' 创建成功！
```

b) **支持模板选择**
```bash
# 使用预定义模板
make gen-module name=blog template=standard
make gen-module name=notify template=worker-only
make gen-module name=chat template=websocket

# 从现有模块克隆
make gen-module name=shop template=system
```

c) **配置文件驱动**
```yaml
# hack/generator/module-templates/blog.yaml
module:
  name: blog
  description: 博客管理模块
  author: DevingGo Team
  license: MIT
  features:
    api: true
    worker: true
    websocket: false
  tables:
    - name: blog_article
      fields:
        - {name: title, type: string}
        - {name: content, type: text}
```

d) **智能依赖检测**
```go
// 自动检测并提示依赖
if detectDatabaseUsage(config) {
    prompt("检测到数据库使用，是否配置数据库连接？")
}
if detectRedisUsage(config) {
    prompt("检测到Redis使用，是否配置Redis连接？")
}
```

#### 4. 增强 module:export（标准化模块包）

**源文件**：`modules/system/cmd/module.go#L100-L167`

**模块包标准结构**（参考npm/composer）：
```
blog.v2.1.0.zip
├── module.json              # 模块元信息（必需）
├── INSTALL.md               # 安装说明
├── CHANGELOG.md             # 版本历史
├── LICENSE                  # 许可证
├── .module.yaml             # 导入配置（新增）
├── scripts/                 # 生命周期脚本（新增）
│   ├── pre-install.sh       # 安装前钩子
│   ├── post-install.sh      # 安装后钩子
│   ├── pre-uninstall.sh     # 卸载前钩子
│   └── post-uninstall.sh    # 卸载后钩子
├── migrations/              # 数据库迁移文件
│   ├── 20240101_create_blog_tables.up.sql
│   └── 20240101_create_blog_tables.down.sql
├── static/                  # 静态资源（新增）
│   ├── web/                 # 前端文件
│   │   ├── admin/           # 后台前端
│   │   └── public/          # 公开前端
│   └── uploads/             # 示例上传文件
├── config/                  # 配置文件模板（新增）
│   ├── blog.yaml.template   # 配置模板
│   └── routes.yaml          # 路由配置
├── src/                     # 源代码
│   ├── api/
│   ├── controller/
│   ├── logic/
│   ├── model/
│   └── worker/
└── docs/                    # 文档
    ├── API.md
    └── README.md
```

**核心优化**：

a) **交互式导出（无参数时）**
```bash
$ go run hack/generator/main.go module:export

📦 导出模块包
? 选择要导出的模块:
  ▸ system (系统核心模块)
    api (API接口模块)
    blog (博客管理模块)
    
? 导出类型:
  ▸ 标准包（完整功能，可分发）
    开发包（含测试、文档源码）
    生产包（仅运行时代码）
    模板包（可复制的脚手架）
    
? 包含内容（多选）:
  ☑ 源代码
  ☑ 数据库迁移
  ☑ 静态资源
  ☑ 配置模板
  ☑ 安装脚本
  ☑ 文档
  ☐ 测试代码
  ☐ 示例数据
  
? 输出目录: ./exports
? 包格式: [zip/tar.gz]
? 是否签名: Yes

🔍 验证模块完整性...
✅ 模块结构：正确
✅ 依赖检查：通过
⚠️  敏感信息扫描：发现2处需要模板化
📦 正在打包 'blog.v2.1.0.zip'...
🔐 生成数字签名...
✅ 导出完成！

文件: ./exports/blog.v2.1.0.zip (3.2MB)
签名: ./exports/blog.v2.1.0.zip.sig
SHA256: a1b2c3d4...
```

b) **智能导出前处理**
```go
// 自动处理
1. 敏感信息替换为模板变量
   - 数据库密码 → {{.DB_PASSWORD}}
   - API密钥 → {{.API_KEY}}
   - 域名 → {{.DOMAIN}}

2. 路径规范化
   - 绝对路径 → 相对路径
   - 硬编码路径 → 配置变量

3. 依赖版本锁定
   - 记录当前依赖的精确版本
   - 生成 module.lock 文件

4. 静态资源优化
   - 自动压缩图片
   - 合并CSS/JS（可选）
   - 生成资源清单
```

c) **导出配置文件（.module.yaml）**
```yaml
# .module.yaml - 模块导入配置
module:
  name: blog
  version: 2.1.0
  type: standard  # standard/plugin/theme
  
# 导入时可配置的变量
variables:
  - name: DB_PREFIX
    description: 数据库表前缀
    default: blog_
    required: true
    
  - name: UPLOAD_PATH
    description: 上传文件路径
    default: resource/public/uploads/blog
    required: false
    
  - name: ADMIN_PATH
    description: 管理后台路径前缀
    default: /admin/blog
    required: false

# 安装配置
install:
  # 静态文件部署策略
  static:
    - source: admin-ui
      target: admin-ui/modules/blog
      merge: true  # 合并而非覆盖
      
    - source: static/uploads
      target: resource/public/uploads/blog
      create_if_missing: true
      
  # 配置文件处理
  config:
    - source: config/blog.yaml.template
      target: hack/config/blog.yaml
      action: merge  # merge/replace/skip/prompt
      variables: true  # 替换模板变量
      
  # 数据库迁移
  migrations:
    auto_run: prompt  # always/never/prompt
    backup_before: true
    
  # 路由注册
  routes:
    auto_register: true
    source: config/routes.yaml
    
# 生命周期钩子
hooks:
  pre_install:
    - script: scripts/pre-install.sh
      description: 检查系统依赖
      
  post_install:
    - script: scripts/post-install.sh
      description: 初始化默认数据
      args:
        - "--skip-demo-data"
        
  pre_uninstall:
    - script: scripts/pre-uninstall.sh
      confirm: true
      description: 备份数据
      
  post_uninstall:
    - script: scripts/post-uninstall.sh
      description: 清理临时文件

# 依赖管理
dependencies:
  modules:
    - name: system
      version: ">=1.0.0"
      
  go_packages:
    - github.com/example/pkg: "v1.2.3"
    
  system:
    min_go_version: "1.23"
    min_gf_version: "2.9.0"
    redis: ">=5.0"
    postgres: ">=13"

# 卸载配置
uninstall:
  remove_data: prompt  # always/never/prompt
  backup_before: true
  keep_migrations: false
```

d) **生命周期脚本示例**
```bash
# scripts/post-install.sh
#!/bin/bash
set -e

echo "🔧 初始化博客模块..."

# 1. 创建必要的目录
mkdir -p resource/public/uploads/blog/{articles,covers}

# 2. 复制默认配置
if [ ! -f hack/config/blog.yaml ]; then
    cp config/blog.yaml.template hack/config/blog.yaml
    echo "✅ 配置文件已创建: hack/config/blog.yaml"
fi

# 3. 初始化默认分类（如果指定）
if [ "$INIT_DEFAULT_DATA" = "true" ]; then
    echo "📊 初始化默认数据..."
    # 这里可以调用API或直接插入数据库
fi

# 4. 编译前端资源（如果需要）
if [ -d "admin-ui/src" ]; then
    echo "🎨 编译前端资源..."
    cd admin-ui && npm install && npm run build
fi

echo "✅ 博客模块安装完成！"
echo "📖 配置文档: docs/README.md"
echo "🚀 重启服务: make dev"
```

e) **版本管理与签名**
```bash
# 自动生成版本信息
输出文件: blog.v2.1.0.zip

# 包含的元信息
module.json:
  version: 2.1.0
  build_time: 2026-03-04T10:30:00Z
  builder: generator/v1.0.0
  checksum: sha256:a1b2c3d4...

# 数字签名（可选）
blog.v2.1.0.zip.sig
- 使用私钥签名
- 导入时验证完整性
```

f) **排除规则（增强）**
```yaml
# hack/generator/.moduleexport
# 全局排除
exclude:
  - "*.log"
  - "*.tmp"
  - ".DS_Store"
  - "node_modules/"
  - "vendor/"
  - ".git/"
  - "*.test"
  - "coverage.out"

# 敏感文件警告
warn:
  - "*.key"
  - "*.pem"
  - ".env"
  - "*secret*"
  
# 类型特定规则
types:
  production:
    exclude:
      - "*_test.go"
      - "docs/src/"
      - "*.md.draft"
      
  template:
    replace:  # 替换为模板变量
      - pattern: 'blog_'
        variable: '{{.DB_PREFIX}}'
      - pattern: '/admin/blog'
        variable: '{{.ADMIN_PATH}}'
```

#### 5. 智能 module:import（灵活初始化）

**源文件**：`modules/system/cmd/module.go#L168-L265`

**核心优化**：

a) **交互式导入向导（无参数时）**
```bash
$ go run hack/generator/main.go module:import

📥 导入模块包
? 模块来源:
  ▸ 本地文件
    远程URL
    模块仓库
    Git仓库

? 选择文件: blog.v2.1.0.zip

🔍 解析模块包...
📦 模块信息:
   名称: blog
   版本: 2.1.0
   类型: 标准模块
   作者: DevingGo Team
   大小: 3.2MB

📋 依赖检查:
   ✅ Go 1.23+ (当前: 1.23.5)
   ✅ GoFrame 2.9.0+ (当前: 2.9.0)
   ✅ 模块 system@1.0.0 (已安装)
   ✅ PostgreSQL 13+ (已配置)
   ✅ Redis 5+ (已配置)

🔧 配置模块参数:
? 数据库表前缀 [blog_]: blog_
? 上传文件路径 [resource/public/uploads/blog]: 
? 管理后台路径 [/admin/blog]: /admin/blog
? 是否初始化示例数据? No

⚠️  冲突检测:
   ⚠️  表 'blog_article' 已存在
   ⚠️  路由 '/admin/blog' 已被使用

? 冲突解决策略:
  ▸ 交互式处理
    全部跳过
    全部覆盖
    使用前缀 'blog_v2'

=== 交互式处理 ===
⚠️  表 'blog_article' 冲突
? 如何处理:
  ▸ 跳过（保留现有）
    覆盖（备份后替换）
    重命名为 'blog_article_v2'
    查看差异

📁 静态文件部署:
   → admin-ui → admin-ui/modules/blog
   → static/uploads → resource/public/uploads/blog

📝 配置文件:
   → config/blog.yaml.template → hack/config/blog.yaml
   ? 文件已存在，如何处理:
     ▸ 智能合并（推荐）
       覆盖
       保留现有
       查看差异

🗄️  数据库迁移:
   发现 2 个迁移文件
   ? 是否立即执行迁移: Yes
   
   ✅ 20240101_create_blog_tables.up.sql
   
🎣 执行安装钩子:
   ⏳ 运行 pre-install.sh
   ✅ 系统依赖检查通过
   
   ⏳ 运行 post-install.sh
   ✅ 默认目录已创建
   ✅ 配置文件已生成

✅ 模块导入完成！

📊 安装摘要:
   模块: blog v2.1.0
   安装路径: modules/blog/
   静态资源: admin-ui/modules/blog/
   配置文件: hack/config/blog.yaml
   
📋 下一步操作:
   1. 查看配置: cat hack/config/blog.yaml
   2. 运行测试: make test-module name=blog
   3. 重启服务: make dev
   4. 访问后台: http://localhost:8000/admin/blog

💡 文档: modules/blog/docs/README.md
```

b) **变量替换引擎**
```go
// 导入时动态替换模板变量
type VariableReplacer struct {
    Variables map[string]string
}

// 支持的变量类型
Variables: {
    // 用户输入
    "DB_PREFIX": "blog_",
    "ADMIN_PATH": "/admin/blog",
    "UPLOAD_PATH": "resource/public/uploads/blog",
    
    // 自动生成
    "MODULE_NAME": "blog",
    "MODULE_NAME_UPPER": "BLOG",
    "MODULE_NAME_CAMEL": "Blog",
    "INSTALL_TIME": "2026-03-04 10:30:00",
    "INSTALL_USER": "admin",
    
    // 系统信息
    "PROJECT_ROOT": "/path/to/devinggo",
    "GO_VERSION": "1.23.5",
    "GF_VERSION": "2.9.0",
}

// 替换规则
- 文件内容: {{.VARIABLE_NAME}}
- 文件路径: {module} → blog
- SQL语句: ${DB_PREFIX} → blog_
```

c) **静态资源智能部署**
```go
// 静态文件部署策略
type StaticDeployment struct {
    Source string      // 源路径
    Target string      // 目标路径
    Strategy string    // copy/symlink/merge
    Merge MergeRule    // 合并规则
}

// 部署模式
1. 复制模式（copy）
   - 完整复制文件
   - 适用于独立静态资源

2. 符号链接（symlink）
   - 创建软链接
   - 节省空间，便于更新
   - 适用于开发环境

3. 合并模式（merge）
   - 智能合并目录
   - 保留现有文件
   - 适用于共享资源

// 前端资源部署示例
admin-ui/ → admin-ui/modules/blog/
├── components/     # 合并到共享组件
├── views/          # 独立的视图目录
├── assets/         # 合并到全局assets
└── router.js       # 注册到主路由

// 实现
- 检测前端框架（Vue/React/原生）
- 自动注册路由到主应用
- 处理CSS/JS依赖
- 更新构建配置
```

d) **配置文件智能合并**
```go
// YAML配置合并策略
type ConfigMergeStrategy string

const (
    MergeStrategyDeep    ConfigMergeStrategy = "deep"    // 深度合并
    MergeStrategyReplace ConfigMergeStrategy = "replace" // 完全替换
    MergeStrategySkip    ConfigMergeStrategy = "skip"    // 跳过不改
    MergeStrategyPrompt  ConfigMergeStrategy = "prompt"  // 交互式
)

// 示例：合并路由配置
// 现有 hack/config.yaml
server:
  routes:
    - path: /admin/system
      module: system
      
# 模块 config/routes.yaml
server:
  routes:
    - path: /admin/blog
      module: blog
      middlewares:
        - auth
        - permission

# 合并结果（deep策略）
server:
  routes:
    - path: /admin/system
      module: system
    - path: /admin/blog      # 新增
      module: blog
      middlewares:
        - auth
        - permission

// 冲突处理
- 相同key：提示用户选择
- 数组：默认追加
- 对象：递归合并
```

e) **数据库迁移集成**
```go
// 自动执行迁移
type MigrationRunner struct {
    AutoRun bool
    Backup  bool
    Rollback bool
}

// 执行流程
1. 检查迁移历史表
   - 避免重复执行
   - 记录迁移状态

2. 备份数据（可选）
   - 导出当前数据
   - 保存到 backups/

3. 执行迁移文件
   - 按顺序 up 迁移
   - 自动替换变量（表前缀等）
   
4. 验证完整性
   - 检查表是否创建
   - 验证数据完整性

5. 失败回滚
   - 执行 down 迁移
   - 恢复备份数据

// 迁移文件变量替换
-- migrations/20240101_create_blog_tables.up.sql
CREATE TABLE ${DB_PREFIX}article (  -- 替换为 blog_article
    id SERIAL PRIMARY KEY,
    title VARCHAR(200),
    ...
);

-- 自动添加模块标记
COMMENT ON TABLE ${DB_PREFIX}article IS 'module:blog;version:2.1.0';
```

f) **依赖自动安装**
```go
// 依赖解析与安装
type DependencyResolver struct {
    Modules  []ModuleDependency
    Packages []PackageDependency
}

// 流程
1. 解析 .module.yaml 依赖声明

2. 检查模块依赖
   - 已安装：检查版本兼容性
   - 未安装：提示安装或自动拉取

3. 安装Go依赖包
   go get github.com/example/pkg@v1.2.3

4. 检查系统依赖
   - Redis：检查连接
   - PostgreSQL：检查版本
   - 外部API：测试可达性

// 示例
依赖: system@1.0.0
✅ 已安装 system@1.0.5 (兼容)

依赖: payment@2.0.0
❌ 未安装
? 是否自动安装 payment 模块: Yes
⬇️  下载 payment@2.0.0...
📦 安装 payment@2.0.0...
✅ 依赖安装完成
```

g) **生命周期钩子执行**
```go
// Hook执行器
type HookExecutor struct {
    Shell     string  // bash/sh/powershell
    Timeout   int     // 超时时间
    CaptureOutput bool
}

// 执行上下文
HookContext: {
    ModuleName: "blog",
    ModuleVersion: "2.1.0",
    InstallPath: "modules/blog",
    Variables: map[string]string,
    DryRun: false,
}

// 执行流程
1. pre-install hook
   - 检查前置条件
   - 创建必要目录
   - 验证权限

2. 文件部署
   - 复制源代码
   - 部署静态资源
   - 生成配置文件

3. post-install hook
   - 初始化数据
   - 编译资源
   - 注册服务

// 失败处理
- 任何hook失败 → 中止安装
- 执行已运行的 rollback hooks
- 清理已复制的文件
- 回滚数据库迁移
```

h) **多种导入来源**
```bash
# 1. 本地文件
go run hack/generator/main.go module:import -file=./blog.v2.1.0.zip

# 2. 远程URL
go run hack/generator/main.go module:import \
  -url=https://modules.devinggo.com/blog/2.1.0

# 3. 模块仓库
go run hack/generator/main.go module:import \
  -repo=official -name=blog -version=2.1.0

# 4. Git仓库
go run hack/generator/main.go module:import \
  -git=https://github.com/devinggo/module-blog.git \
  -branch=v2.1.0

# 5. 从模板创建
go run hack/generator/main.go module:import \
  -template=blog -name=myblog -customize
```

i) **安装验证与测试**
```go
// 安装后自动验证
type InstallValidator struct {
    CheckFiles      bool
    CheckDatabase   bool
    CheckRoutes     bool
    CheckStatic     bool
    RunTests        bool
}

// 验证项目
1. 文件完整性
   ✅ 所有必需文件已复制
   ✅ 文件权限正确
   
2. 数据库完整性
   ✅ 表结构正确
   ✅ 索引已创建
   
3. 路由注册
   ✅ 路由已加载
   ✅ 无冲突
   
4. 静态资源
   ✅ 文件可访问
   ✅ 路径映射正确
   
5. 单元测试（可选）
   ⏳ 运行模块测试...
   ✅ 所有测试通过

// 验证失败处理
- 警告：记录但继续
- 错误：提示修复或回滚
```

j) **卸载与清理**
```bash
$ go run hack/generator/main.go module:uninstall -name=blog

⚠️  卸载模块 'blog'

? 是否备份数据: Yes
📦 导出数据到: backups/blog_backup_20260304.zip

🎣 执行卸载钩子:
   ⏳ 运行 pre-uninstall.sh
   ✅ 数据已备份
   
📁 移除文件:
   ✓ modules/blog/
   ✓ admin-ui/modules/blog/
   
🗄️  数据库清理:
   ? 是否删除数据表: No (保留)
   ? 是否执行 down 迁移: Yes
   
   ⏳ 运行 20240101_create_blog_tables.down.sql
   ✅ 迁移已回滚
   
📝 配置清理:
   ? 删除 hack/config/blog.yaml: No (保留)
   
🎣 post-uninstall hook:
   ⏳ 清理临时文件...
   ✅ 清理完成

✅ 模块 'blog' 已卸载

📊 卸载摘要:
   已删除: 源代码、静态资源
   已保留: 数据表、配置文件
   备份位置: backups/blog_backup_20260304.zip
```

#### 6. 新增 module:list（模块管理）

```bash
$ make list-modules

📦 已安装模块 (3)
┌─────────┬─────────┬────────┬──────────┬────────────┬─────────┐
│ 名称    │ 版本    │ 状态   │ 作者     │ 描述       │ 大小    │
├─────────┼─────────┼────────┼──────────┼────────────┼─────────┤
│ system  │ 1.0.0   │ ✓ 启用 │ devinggo │ 系统核心   │ 5.2MB   │
│ api     │ 1.0.0   │ ✓ 启用 │ devinggo │ API模块    │ 2.1MB   │
│ blog    │ 2.1.0   │ ✗ 禁用 │ custom   │ 博客模块   │ 3.2MB   │
└─────────┴─────────┴────────┴──────────┴────────────┴─────────┘

💡 提示：
  - 启用模块: make enable-module name=blog
  - 升级检查: make check-module-updates
  - 详细信息: make module-info name=blog
```

#### 7. 新增 module:repo（模块仓库管理）

```bash
# 配置模块仓库
$ go run hack/generator/main.go module:repo add official \
  https://modules.devinggo.com/api/v1

# 搜索模块
$ go run hack/generator/main.go module:repo search blog

🔍 搜索结果 (3)
┌──────────────┬─────────┬──────────┬─────────────────┐
│ 名称         │ 版本    │ 下载量   │ 描述            │
├──────────────┼─────────┼──────────┼─────────────────┤
│ blog         │ 2.1.0   │ 1.2k     │ 完整博客系统    │
│ blog-mini    │ 1.0.0   │ 500      │ 轻量级博客      │
│ cms-blog     │ 3.0.0   │ 2.5k     │ CMS博客扩展     │
└──────────────┴─────────┴──────────┴─────────────────┘

# 查看模块详情
$ go run hack/generator/main.go module:repo info blog

📦 blog - 完整博客系统
版本: 2.1.0 (最新)
作者: DevingGo Team
许可: MIT
大小: 3.2MB
下载: 1,234 次
评分: ⭐⭐⭐⭐⭐ (4.8/5.0, 45评价)

📝 描述:
  功能完整的博客模块，支持文章管理、分类、标签、
  评论、SEO优化等功能。

✨ 特性:
  • Markdown编辑器
  • 多语言支持
  • SEO优化
  • 评论系统
  • 标签云
  • 文章统计

📋 依赖:
  - system >= 1.0.0
  - Go >= 1.23
  - PostgreSQL >= 13

🎯 兼容性: ✅ 与当前环境兼容

💾 安装:
  go run hack/generator/main.go module:import -repo=official -name=blog

📖 文档: https://docs.devinggo.com/modules/blog
🐛 问题: https://github.com/devinggo/module-blog/issues
```

#### 8. 新增 module:publish（发布到仓库）

```bash
$ go run hack/generator/main.go module:publish

📤 发布模块到仓库

? 选择要发布的模块: blog
? 目标仓库: official
? 发布类型:
  ▸ 新版本发布
    更新现有版本
    Beta测试版

🔍 验证模块包...
✅ 模块结构完整
✅ 文档齐全
✅ 测试通过
✅ 无安全隐患

📝 版本信息:
   当前版本: 2.1.0
   新版本: 2.2.0
   
? 版本类型:
  ▸ 补丁更新 (2.1.1)
    次要更新 (2.2.0)
    主要更新 (3.0.0)
    
? 更新说明:
  [打开编辑器编写 CHANGELOG]

🔐 认证:
   API Token: ••••••••••••
   ✅ 认证通过
   
📦 打包上传...
⬆️  上传中... [████████████] 100% (3.2MB/3.2MB)
✅ 上传完成

🔍 模块审核:
   ⏳ 等待审核... (通常需要1-3个工作日)
   📧 审核结果将发送到: your@email.com
   
✅ 发布请求已提交！

📊 发布摘要:
   模块: blog v2.2.0
   仓库: official
   状态: 待审核
   追踪: https://modules.devinggo.com/publish/12345
```

#### 9. 新增 module:validate（模块验证）

```bash
$ make validate-module name=blog

🔍 验证模块 'blog'...

✅ 模块结构：正确
✅ module.json：有效
✅ 必需文件：完整
⚠️  缺少文档：README.md
❌ 导入冲突：pkg/utils 循环导入
✅ 代码规范：通过 golangci-lint

📊 验证结果：2 个警告，1 个错误
```

#### 10. 新增 module:clone（快速克隆）

```bash
# 从现有模块快速创建新模块
$ make clone-module source=system target=crm

🔄 正在克隆模块 'system' → 'crm'...

自动处理：
  ✓ 替换模块名称（System → Crm）
  ✓ 更新包导入路径
  ✓ 重命名文件和目录
  ✓ 更新 module.json
  ✓ 清理业务数据

✅ 模块克隆完成！新模块位置：modules/crm/
```

#### 11. 新增 module:upgrade（版本升级）

```bash
$ make upgrade-module name=blog version=3.0.0

📦 升级模块 'blog' 从 2.1.0 到 3.0.0

检查兼容性：
  ✓ Go版本：1.23+ → 1.23+ 兼容
  ✓ 依赖模块：system@1.0.0 兼容
  ⚠️  破坏性变更：API路径调整

? 是否继续升级？Yes

⬇️  下载升级包...
🔄 备份当前版本...
📝 执行迁移脚本...
✅ 升级完成！

⚠️  请注意：
  - API路径已变更，需更新前端代码
  - 新增配置项：blog.cache.ttl
```

#### 核心逻辑提取

- `createModuleFiles()` → `Generator.GenerateModuleFiles()`
- `createModuleMigrationFiles()` → `Generator.GenerateMigration()`
- `createModuleConfigFile()` → `Generator.GenerateModuleConfig()`
- 新增：`Generator.ValidateModule()` - 模块验证
- 新增：`Generator.CloneModule()` - 模块克隆
- 新增：`Generator.UpgradeModule()` - 版本升级
- 新增：`Exporter.PackageModule()` - 标准化打包
- 新增：`Importer.InstallModule()` - 智能安装
- 新增：`Importer.UninstallModule()` - 安全卸载
- 新增：`VariableReplacer.Replace()` - 变量替换
- 新增：`ConfigMerger.Merge()` - 配置合并
- 新增：`HookExecutor.Run()` - 钩子执行
- 新增：`StaticDeployer.Deploy()` - 静态资源部署
- 新增：`RepositoryClient.Search()` - 仓库搜索
- 新增：`RepositoryClient.Publish()` - 发布模块

#### 模块包管理系统优势

**🎯 对比传统zip打包方式**：

| 特性 | 传统方式 | 增强方案 |
|------|---------|---------|
| 包结构 | 随意 | 标准化（.module.yaml） |
| 初始化 | 手动 | 自动（生命周期钩子） |
| 配置 | 硬编码 | 模板变量+交互式 |
| 静态资源 | 手动复制 | 智能部署+合并 |
| 数据库 | 手动执行SQL | 自动迁移+回滚 |
| 依赖 | 未管理 | 自动检测+安装 |
| 冲突 | 覆盖或失败 | 智能检测+多种策略 |
| 卸载 | 手动删除 | 安全卸载+备份 |
| 分发 | 手动传播 | 模块仓库+版本管理 |
| 更新 | 重新安装 | 就地升级+数据保留 |

**🚀 关键改进**：

1. **标准化打包格式**
   - 类似npm/composer的模块包标准
   - 包含完整的元信息和依赖声明
   - 支持数字签名验证

2. **灵活的初始化**
   - 生命周期钩子（pre/post install/uninstall）
   - 模板变量系统（数据库前缀、路径等可配置）
   - 交互式配置向导

3. **智能静态资源管理**
   - 自动部署前端资源到正确位置
   - 支持符号链接（开发环境）
   - 智能合并共享资源

4. **配置文件友好**
   - YAML配置深度合并
   - 冲突智能检测
   - 多种合并策略

5. **完整的生命周期管理**
   - 安装验证
   - 依赖自动解析
   - 安全卸载+数据备份
   - 版本升级+迁移

6. **模块生态系统**
   - 官方模块仓库
   - 版本管理和更新检查
   - 社区分享和发现

**💡 使用场景**：

```bash
# 场景1：开发者发布模块
make export-module name=blog   # 标准化打包
make publish-module name=blog  # 发布到仓库

# 场景2：用户安装模块（交互式）
make import-module             # 选择来源→配置参数→自动安装

# 场景3：批量部署（配置文件）
# deployment.yaml
modules:
  - name: blog
    source: official
    version: 2.1.0
    config:
      db_prefix: blog_
      admin_path: /admin/blog

make deploy-modules -f deployment.yaml

# 场景4：模块升级
make upgrade-module name=blog version=2.2.0  # 保留数据+配置

# 场景5：开发环境同步
make export-modules-list > modules.lock
# 在另一台机器
make import-modules-list < modules.lock
```

#### 关键变更

**代码组织**：
- 移除 `devinggo/modules/system/pkg/utils` 依赖
- 模板路径改为 `hack/generator/templates/module`
- 生成的文件路径中 `modules/_/` → `modules/bootstrap/`
- 使用 `internal/utils/zip.ZipDirectory()` 替换 `utils.ZipDirectory()`
- 使用 `internal/utils/zip.UnzipFile()` 替换 `utils.UnzipFile()`

**模块包标准化**：
- 采用 `.module.yaml` 配置文件定义模块元信息
- 标准化包结构（src/、static/、config/、migrations/、scripts/）
- 支持生命周期钩子（pre/post install/uninstall）
- 支持模板变量系统（{{.VARIABLE}}）

**安装增强**：
- 实现变量替换引擎（用户可配置参数）
- 智能配置合并（YAML深度合并）
- 静态资源自动部署（多种策略）
- 数据库迁移自动执行（带回滚）
- 依赖自动检测与安装
- 安装后验证机制

**模块生态**：
- 新增模块仓库客户端（搜索、安装、发布）
- 支持多种导入来源（本地、URL、Git、仓库）
- 版本管理与更新检查
- 模块市场集成

**12. 实现模块包标准化工具**
- 创建 `.module.yaml` 解析器：`internal/config/module_config.go`
- 实现变量替换引擎：`internal/generator/variable_replacer.go`
- 实现配置合并器：`internal/utils/config_merger.go`
- 实现生命周期钩子执行器：`internal/generator/hook_executor.go`
- 实现静态资源部署器：`internal/generator/static_deployer.go`

**13. 实现模块仓库客户端（可选，后期扩展）**
- 创建仓库API客户端：`internal/repository/client.go`
- 实现模块搜索：`SearchModules(keyword string)`
- 实现模块下载：`DownloadModule(name, version string)`
- 实现模块发布：`PublishModule(packagePath string)`
- 配置管理：`hack/generator/repositories.yaml`

**14. 迁移模板文件**
- 复制 `resource/generate/modules/` 所有模板到 `hack/generator/templates/module/`
- 模板重命名：`.html` → `.tpl`（如 `module.go.html` → `module.go.tpl`）
- 统一模板变量：移除 `{%.moduleName%}` 特殊标记，改为标准Go template `{{.ModuleName}}`
- 更新路径引用：`modules/_/` → `modules/bootstrap/`

**15. 迁移SQL迁移文件生成**
- 提取 `modules/system/cmd/module.go#L367-L418` 的 `createModuleMigrationFiles` 逻辑
- 移到 `hack/generator/internal/generator/module.go` 的 `GenerateMigration()` 方法
- 模板文件：`sql/module_up_postgres.html` → `templates/module/migration_up.sql.tpl`

### 阶段三：迁移Worker创建命令（第3天）

**16. 重构 worker:create（交互式创建）**
- 源文件：`modules/system/cmd/worker_create.go` 全部418行
- 目标：`hack/generator/cmd/worker.go` + `internal/generator/worker.go`

**核心优化**：

a) **交互式模式（无参数时）**
```bash
$ go run hack/generator/main.go worker:create

⚙️  创建Worker任务
? 任务名称: send_email
? 任务描述: 发送邮件通知
? 所属模块:
  ▸ system
    api
    blog
? 任务类型:
  ▸ 异步任务（Task）- 实时触发
    定时任务（Cron）- 定时执行
    两者都要（Both）

# 如果选择Cron或Both
? Cron表达式: */5 * * * *  (每5分钟)

? 是否需要参数结构体? Yes
? 参数字段（回车结束）:
  字段1名称: email
  字段1类型: string
  字段2名称: subject
  字段2类型: string
  字段3名称: 

✅ Worker任务创建成功！
生成文件:
  - modules/system/worker/server/send_email_worker.go
  - modules/system/worker/consts/const.go (已更新)
```

b) **保持所有现有功能**：
  - 支持 `-type` 参数（task/cron/both）
  - 支持 `-module` 参数（默认system）
  - 自动创建目录结构
  - 智能更新 `consts/const.go` 文件
  - 数据结构复用（both类型时）

c) **命令行快捷方式（有参数时）**
```bash
# 快速创建（跳过交互）
make gen-worker name=send_email type=task module=system desc="发送邮件"
```

**17. 创建Worker模板**
- `templates/worker/cron.go.tpl` - 定时任务模板：
  ```go
  package cron
  
  import (
      "context"
      "devinggo/modules/{{.ModuleName}}/worker/consts"
      "devinggo/modules/system/pkg/worker"
      glob2 "devinggo/modules/system/pkg/worker/glob"
  )
  
  type {{.StructName}} struct {
      // TODO: 定义参数字段
  }
  
  func init() {
      worker.RegisterCronFunc(consts.{{.ConstName}}_CRON, "{{.Description}}", {{.HandlerName}})
  }
  
  func {{.HandlerName}}(ctx context.Context, payload *glob2.Payload, params *gjson.Json) {
      // 任务逻辑
  }
  ```
- `templates/worker/task.go.tpl` - 异步任务模板
- `templates/worker/const.go.tpl` - 常量文件模板

**18. 优化常量文件更新逻辑**
- 现有实现使用字符串查找最后的 `)` 插入常量
- 改进：使用Go AST解析，安全地插入新常量
- 重复检测更准确，避免语法错误

**19. 增强Worker生成功能**
- 支持从配置文件批量生成多个Worker（`-config=hack/generator.yaml`）
- 为Worker添加单元测试模板（可选）
- 生成示例数据结构和注释

### 阶段四：新增CRUD生成器（第4-5天）

**20. 交互式CRUD生成（无参数时）**

```bash
$ go run hack/generator/main.go crud:generate

🔧 生成CRUD代码

? 生成模式:
  ▸ 单表生成
    批量生成（配置文件）

# === 单表模式 ===
? 所属模块:
  ▸ system
    api
    blog

? 选择数据表:
  ▸ system_category
    system_dept
    system_menu
    system_post
    ...

🔍 分析表结构 'system_category'...
检测到字段: id, name, code, status, sort, remark, created_at, updated_at

? 业务名称: Category
? 中文名称: 分类管理
? 权限前缀: system:category

? 配置搜索字段（多选）:
  ☑ name (模糊搜索)
  ☑ code (精确匹配)
  ☑ status (精确匹配)
  ☑ created_at (范围搜索)
  ☐ sort
  ☐ remark

? 生成的CRUD方法（多选）:
  ☑ Index (分页列表)
  ☑ List (完整列表)
  ☑ Save (新增)
  ☑ Update (更新)
  ☑ Delete (软删除)
  ☑ Read (详情)
  ☑ ChangeStatus (状态切换)
  ☐ BatchDelete (批量删除)
  ☐ Export (导出Excel)
  ☐ Import (导入Excel)

✅ 配置完成！开始生成代码...

生成文件:
  ✓ modules/system/model/req/system_category.go
  ✓ modules/system/model/res/system_category.go
  ✓ modules/system/api/system_category.go
  ✓ modules/system/logic/system_category.go
  ✓ modules/system/controller/system_category.go

🔄 运行 gf gen service...
✅ Service接口已生成

🔄 运行 gf gen ctrl...
✅ Controller注册已更新

✅ CRUD代码生成完成！

📝 下一步:
  1. 检查生成的代码: modules/system/
  2. 补充业务逻辑: logic/system_category.go
  3. 测试接口: make dev
```

**21. 设计配置文件**
- 创建 `hack/generator.yaml`：
```yaml
generator:
  global:
    author: "DevingGo Generator"
    license: "Apache-2.0"

  modules:
    - name: "system"
      path: "modules/system"
      tables:
        - table: "system_category"
          business: "Category"
          cnName: "分类"
          searchFields:
            - field: "name"
              type: "like"            # like/eq/between/in
            - field: "status"
              type: "eq"
            - field: "created_at"
              type: "between"
          methods:                    # 生成的CRUD方法
            - "Index"                 # 列表（分页）
            - "List"                  # 列表（不分页）
            - "Save"                  # 新增
            - "Update"                # 更新
            - "Delete"                # 软删除
            - "Read"                  # 详情
            - "ChangeStatus"          # 修改状态
          permission: "system:category"  # 权限前缀
```

**22. 实现CRUD生成器核心**
- `internal/generator/crud.go` 实现：
  - `GenerateModel()` - 生成 req/res 模型
  - `GenerateAPI()` - 生成 API 定义
  - `GenerateLogic()` - 生成 Logic 实现
  - `GenerateController()` - 生成 Controller 实现
- 数据源：读取 `internal/model/entity/` 获取表字段
- 生成顺序：Model → API → Logic → Controller

**23. 创建CRUD模板**
- `templates/crud/req.go.tpl` - 请求模型：
  ```go
  type {{.ModuleName}}{{.BusinessName}}Search struct {
      {{range .SearchFields}}
      {{.Name}} {{.Type}} `json:"{{.JsonName}}"`
      {{end}}
  }
  
  type {{.ModuleName}}{{.BusinessName}}Save struct {
      Id int64 `json:"id"`
      {{range .Fields}}
      {{.Name}} {{.Type}} `json:"{{.JsonName}}" v:"{{.Validation}}"`
      {{end}}
  }
  ```
- `templates/crud/res.go.tpl` - 响应模型
- `templates/crud/api.go.tpl` - API定义（含g.Meta标签）
- `templates/crud/logic.go.tpl` - Logic实现（继承GenericService）
- `templates/crud/controller.go.tpl` - Controller实现

**24. 实现 crud:generate 命令**
- `cmd/crud.go` 命令入口：
  ```go
  var GenerateCrud = &gcmd.Command{
      Name: "crud:generate",
      Brief: "生成CRUD业务代码",
      Func: func(ctx context.Context, parser *gcmd.Parser) error {
          // 模式1：单表生成
          if parser.GetOpt("table").String() != "" {
              return generateSingleTable(ctx, parser)
          }
          // 模式2：批量生成
          if parser.GetOpt("config").String() != "" {
              return generateBatch(ctx, parser)
          }
      },
  }
  ```
- 生成后自动运行 `gofmt` 和 `goimports` 格式化代码

### 阶段五：工作流集成（第6天）

**25. 更新Makefile命令**
- 在 `hack/hack.mk` 添加新命令：
```makefile
# ========================================
# 代码生成工具命令
# ========================================

# 模块管理
.PHONY: gen-module
gen-module: cli.install
	@go run hack/generator/main.go module:create -name=$(name)

.PHONY: export-module
export-module:
	@go run hack/generator/main.go module:export -name=$(name)

.PHONY: import-module
import-module:
	@go run hack/generator/main.go module:import -file=$(file)

# Worker任务
.PHONY: gen-worker
gen-worker:
	@go run hack/generator/main.go worker:create -name=$(name) -module=$(or $(module),system) -type=$(or $(type),both) -desc="$(desc)"

# CRUD生成
.PHONY: gen-crud
gen-crud: cli.install dao
	@go run hack/generator/main.go crud:generate -module=$(module) -table=$(table) -business=$(business)
	@gf gen service -s=modules/$(module)/logic -d=modules/$(module)/service
	@gf gen ctrl

.PHONY: gen-batch
gen-batch: cli.install dao
	@go run hack/generator/main.go crud:generate -config=hack/generator.yaml
	@$(MAKE) service
	@gf gen ctrl

# 开发工作流
.PHONY: dev
dev: gen-batch
	@gf run main.go

# 帮助信息
.PHONY: help-gen
help-gen:
	@echo "代码生成工具命令:"
	@echo "  make gen-module name=模块名                    创建新模块"
	@echo "  make export-module name=模块名                 导出模块zip"
	@echo "  make import-module file=模块.zip               导入模块"
	@echo "  make gen-worker name=任务名 [module=] [type=]  创建Worker任务"
	@echo "  make gen-crud module=模块 table=表 business=   生成CRUD代码"
	@echo "  make gen-batch                                 批量生成（配置文件）"
	@echo ""
	@echo "示例:"
	@echo "  make gen-module name=blog"
	@echo "  make gen-worker name=send_email desc='发送邮件'"
	@echo "  make gen-crud module=system table=system_post business=Post"
```

**26. 清理旧代码**
- 删除 `modules/system/cmd/module.go`
- 删除 `modules/system/cmd/worker_create.go`
- 如果 `modules/system/cmd/` 目录变空，删除整个目录
- 检查是否有其他地方引用这些命令，一并更新

**27. 重命名模块加载器**
- 将 `modules/_/` 重命名为 `modules/bootstrap/`
- 更新 `main.go` 导入路径：`_ "devinggo/modules/bootstrap/logic"`
- 更新所有模板中的路径引用
- 更新代码生成器中的路径常量

### 阶段六：文档与测试（第7天）

**28. 创建工具文档**
- `hack/generator/README.md`：
```markdown
# DevingGo 代码生成工具集

统一的代码生成工具，支持模块管理、Worker任务、CRUD代码生成。

## 🚀 快速开始

### 创建新模块
```bash
make gen-module name=blog
# 或
go run hack/generator/main.go module:create -name=blog
```

### 创建Worker任务
```bash
make gen-worker name=send_email desc="发送邮件"
# 指定模块和类型
make gen-worker name=clean_logs module=system type=cron desc="清理日志"
```

### 生成CRUD代码
```bash
# 单表生成
make gen-crud module=system table=system_post business=Post

# 批量生成（配置文件）
make gen-batch
```

## 📦 命令详解

| 命令 | 功能 | 参数 |
|------|------|------|
| `module:create` | 创建模块骨架 | `-name` 模块名 |
| `module:export` | 导出模块zip | `-name` 模块名 |
| `module:import` | 导入模块zip | `-file` zip文件路径 |
| `worker:create` | 创建Worker任务 | `-name` 任务名<br>`-module` 模块名（默认system）<br>`-type` 类型（task/cron/both）<br>`-desc` 描述 |
| `crud:generate` | 生成CRUD代码 | `-module` 模块名<br>`-table` 表名<br>`-business` 业务名<br>或 `-config` 配置文件 |

## 🔧 配置文件

编辑 `hack/generator.yaml` 进行批量生成配置。

## 📁 目录结构

生成器完全独立，不依赖业务模块：
- `cmd/` - 命令实现
- `internal/generator/` - 生成器引擎
- `internal/utils/` - 独立工具函数
- `templates/` - 模板文件

## 🧪 开发流程

1. 设计数据库表
2. `make dao` - 生成DAO/Entity（GoFrame）
3. `make gen-crud module=system table=表名 business=业务名` - 生成CRUD代码
4. `gf gen service` 自动生成Service接口
5. 补充业务逻辑
6. 测试

## 📖 相关文档

- [模块开发指南](../../docs/MODULE.md)
- [Worker开发指南](../../docs/WORKER.md)
- [CRUD生成指南](../../docs/CRUD.md)
```

**29. 添加目录说明**
- `internal/dao/internal/README.md`：
  ```markdown
  # ⚠️ 自动生成目录
  
  此目录由 `gf gen dao` 自动生成，**禁止手动编辑**。
  
  如需扩展DAO方法，请在上级目录 `internal/dao/` 中添加。
  ```

- `modules/README.md`：
  ```markdown
  # 模块目录
  
  ## 模块架构
  
  DevingGo采用模块化架构，每个模块独立开发、独立部署。
  
  ### 标准模块结构
  
  ```
  modules/{module}/
    ├── module.go           # 模块入口
    ├── module.json         # 模块元信息
    ├── api/                # API定义
    ├── controller/         # HTTP处理器
    ├── logic/              # 业务逻辑
    ├── service/            # 接口定义（自动生成）
    ├── model/
    │   ├── req/            # 请求模型
    │   └── res/            # 响应模型
    ├── router/             # 路由定义（可选）
    └── worker/             # 后台任务（可选）
  ```
  
  ### 创建新模块
  
  ```bash
  make gen-module name=模块名
  ```
  
  详见：[模块开发指南](../docs/MODULE.md)
  ```

**30. 编写单元测试**
- `hack/generator/internal/utils/naming_test.go` - 测试命名转换
- `hack/generator/internal/utils/zip_test.go` - 测试压缩解压
- `hack/generator/internal/generator/module_test.go` - 测试模块生成

**31. 集成测试脚本**
- 创建 `hack/generator/test.sh`：
```bash
#!/bin/bash
set -e

echo "🧪 开始集成测试..."

# 测试模块创建
echo "1️⃣  测试模块创建"
go run main.go module:create -name=testmod
test -d ../../modules/testmod || exit 1

# 测试Worker创建
echo "2️⃣  测试Worker创建"
go run main.go worker:create -name=test_worker -module=testmod -type=both -desc="测试任务"
test -f ../../modules/testmod/worker/cron/test_worker_cron.go || exit 1
test -f ../../modules/testmod/worker/server/test_worker_worker.go || exit 1

# 测试模块导出
echo "3️⃣  测试模块导出"
go run main.go module:export -name=testmod
test -f testmod.v1.0.0.zip || exit 1

# 清理
echo "🧹 清理测试数据"
rm -rf ../../modules/testmod
rm -f testmod.v1.0.0.zip

echo "✅ 所有测试通过！"
```

## 四、验证清单

### 功能验证

**模块管理**：
- [ ] `make gen-module name=testmod` 创建模块，生成15个文件
- [ ] 模块包含标准目录结构（api、controller、logic等）
- [ ] 生成的 `module.json` 配置正确
- [ ] 生成的SQL迁移文件可执行
- [ ] `make export-module name=testmod` 生成zip文件
- [ ] `make import-module file=testmod.v1.0.0.zip` 成功导入
- [ ] 不依赖 `modules/system/pkg/utils`

**Worker任务**：
- [ ] `make gen-worker name=send_email type=task` 仅生成Task文件
- [ ] `make gen-worker name=clean_logs type=cron` 仅生成Cron文件
- [ ] `make gen-worker name=notify type=both` 同时生成Task和Cron，数据结构复用
- [ ] 常量文件正确更新，无重复
- [ ] 生成的代码编译通过
- [ ] 支持指定模块：`make gen-worker name=test module=api`

**CRUD生成**：
- [ ] `make gen-crud module=system table=system_test business=Test` 生成5类文件
- [ ] 生成的Model包含正确的字段和标签
- [ ] 生成的API定义包含完整的g.Meta标签
- [ ] 生成的Logic继承GenericService，包含搜索方法
- [ ] 生成的Controller实现完整CRUD方法
- [ ] 执行 `gf gen service` 后Service接口自动生成
- [ ] 批量生成：`make gen-batch` 处理配置文件中所有表

### 工作流验证

- [ ] `make help-gen` 显示所有代码生成命令
- [ ] 生成的代码运行 `go build` 编译通过
- [ ] 启动服务，Swagger显示新生成的接口
- [ ] 测试接口功能正常（增删改查）

### 工具独立性验证

- [ ] `hack/generator` 可独立编译：`cd hack/generator && go build`
- [ ] 不导入任何 `modules/` 下的业务代码
- [ ] 只使用标准库和GoFrame核心库
- [ ] 可以复制到其他DevingGo项目使用

### E2E测试

1. 创建新模块：`make gen-module name=blog`
2. 设计表：`blog_article`（id, title, content, status, created_at等）
3. 生成DAO：`make dao`
4. 生成CRUD：`make gen-crud module=blog table=blog_article business=Article`
5. 启动服务：`make dev`
6. 测试接口：创建文章、查询列表、更新、删除
7. 验证：所有功能正常，代码质量良好

## 五、技术决策

**1. 工具独立性**
- ✅ 采用：hack/generator 完全独立，不依赖业务模块
- ❌ 拒绝：继续放在 modules/system/cmd，耦合度高
- 理由：工具代码应独立维护，便于复用和测试

**2. 模板系统**
- ✅ 采用：迁移到 hack/generator/templates，自包含
- ❌ 拒绝：共享 resource/generate 目录
- 理由：生成器应自包含所有依赖，提高可移植性

**3. 工具函数实现**
- ✅ 采用：使用标准库重新实现（archive/zip）
- ❌ 拒绝：继续依赖 modules/system/pkg/utils
- 理由：移除对业务模块的依赖，提高独立性

**4. 与GoFrame工具链协作**
- ✅ 采用：生成业务层代码，调用 `gf gen service` 生成接口
- ❌ 拒绝：修改 `gf gen dao` 或手动管理Service接口
- 理由：尊重GoFrame工具链，职责分离

**5. 配置驱动生成**
- ✅ 采用：支持 generator.yaml 批量生成
- ❌ 拒绝：仅支持命令行参数
- 理由：批量生成提升效率，配置可版本控制

**6. 命令风格**
- ✅ 采用：保持 gcmd.Command 框架，`module:create` 格式
- ❌ 拒绝：改为 Cobra 等其他CLI框架
- 理由：与现有命令风格一致，降低学习成本

**7. 模块加载器命名**
- ✅ 采用：重命名为 modules/bootstrap
- ❌ 拒绝：保持 modules/_
- 理由：下划线包名违反Go规范，bootstrap语义更清晰

## 六、预期收益

**开发效率提升**：
- 模块创建：从手动30分钟 → 自动化1分钟（提升97%）
- 模块导入：从手动配置20分钟 → 交互式安装3分钟（提升85%）
- Worker创建：从手动15分钟 → 自动化30秒（提升97%）
- CRUD开发：从手动2小时 → 自动化5分钟（提升96%）

**代码质量提升**：
- 统一代码风格，消除人工差异
- 减少模板代码重复率80%+
- 自动生成的代码经过充分测试，错误率低
- 标准化模块结构，提高可维护性

**部署与分发**：
- 模块打包：从随意打包 → 标准化+验证（安全性提升100%）
- 模块部署：从手动复制+配置 → 自动安装（速度提升90%）
- 静态资源：从手动管理 → 自动部署+合并（错误率降低95%）
- 配置管理：从覆盖冲突 → 智能合并（冲突减少80%）

**维护性提升**：
- 工具代码独立，易于测试和升级
- 模板集中管理，修改一处即可
- 清晰的职责划分，降低理解成本
- 生命周期钩子支持，扩展性强

**生态系统建设**：
- 模块可复用、可分发、可共享
- 官方模块市场，降低重复开发
- 版本管理，支持升级和回滚
- 社区贡献，加速功能迭代

**可扩展性**：
- 新增生成器类型容易（如：生成API文档、测试代码）
- 支持自定义模板覆盖
- 可发布为独立CLI工具
- 模块仓库系统，支持企业私有部署

## 七、实施时间表

| 阶段 | 任务 | 预计时间 | 产出 |
|------|------|----------|------|
| 第1天 | 基础架构搭建 + 独立工具函数 | 1天 | 项目结构、utils工具 |
| 第2-3天 | 模块管理命令迁移与增强 | 2天 | module:create/export/import/list/validate/clone/upgrade<br>+ 模块包标准化（.module.yaml、钩子系统） |
| 第4天 | 模块仓库客户端（可选） | 1天 | module:repo命令、模块搜索/发布 |
| 第5天 | 迁移Worker创建命令 | 1天 | worker:create（交互式增强） |
| 第6-7天 | 新增CRUD生成器 | 2天 | crud:generate（交互式+批量） |
| 第8天 | 工作流集成 + Makefile更新 | 1天 | make命令、清理旧代码 |
| 第9天 | 文档编写 + 测试验证 | 1天 | README、测试脚本、使用文档 |

**总计：9个工作日**（含模块仓库客户端）  
**核心功能：7个工作日**（不含模块仓库，可后期扩展）

**分阶段交付**：
- **阶段1（3天）**：基础工具 + 标准化模块包导入导出
- **阶段2（2天）**：Worker生成器 + 完整的模块管理命令
- **阶段3（2天）**：CRUD生成器
- **阶段4（2天）**：集成测试 + 文档 + （可选）模块仓库

## 八、风险与应对

**风险1：模板迁移后路径错误**
- 应对：创建测试脚本验证所有生成的文件路径正确
- 验证：生成测试模块，确保编译通过

**风险2：常量文件更新破坏现有代码**
- 应对：使用AST解析而非字符串操作
- 验证：测试重复创建场景，确保错误提示友好

**风险3：CRUD生成器与现有代码风格不一致**
- 应对：基于现有代码提取模板，保持一致性
- 验证：生成的代码与手写代码对比，确保无差异

**风险4：GoFrame工具链版本兼容性**
- 应对：明确依赖的GoFrame版本（v2.9.0）
- 验证：在干净环境测试完整工作流

**风险5：工具独立性不足**
- 应对：代码审查确保不导入 `modules/` 下的代码
- 验证：独立编译 `hack/generator` 成功

---

**方案批准**：待确认后开始实施

**联系人**：开发团队负责人

**最后更新**：2026年3月4日
