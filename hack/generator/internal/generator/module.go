package generator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"devinggo/hack/generator/internal/scanner"
	"devinggo/hack/generator/internal/utils"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
)

// ModuleExporter 模块导出器
type ModuleExporter struct {
	ctx context.Context
}

// NewModuleExporter 创建模块导出器
func NewModuleExporter(ctx context.Context) *ModuleExporter {
	return &ModuleExporter{ctx: ctx}
}

// Export 导出模块为zip包
func (e *ModuleExporter) Export(moduleName string) (string, error) {
	// 扫描模块信息
	moduleInfo, err := scanner.ScanModule(moduleName)
	if err != nil {
		return "", err
	}

	// 创建临时目录
	tmpDir := filepath.Join(utils.GetTmpDir(), "module_export_"+moduleName)
	defer os.RemoveAll(tmpDir)

	// 复制文件到临时目录
	for fileType, paths := range moduleInfo.Files {
		for _, path := range paths {
			// 确保源文件存在
			if _, err := os.Stat(path); os.IsNotExist(err) {
				g.Log().Warningf(e.ctx, "[%s]文件不存在，跳过: %s", fileType, path)
				continue
			}

			// 计算目标路径
			dstPath := filepath.Join(tmpDir, path)

			// 创建目标目录
			if err := utils.EnsureDir(filepath.Dir(dstPath)); err != nil {
				return "", fmt.Errorf("创建目录失败: %w", err)
			}

			// 复制文件或目录
			if gfile.IsDir(path) {
				if err := gfile.CopyDir(path, dstPath); err != nil {
					return "", fmt.Errorf("复制目录失败: %w", err)
				}
			} else {
				if err := utils.CopyFile(path, dstPath); err != nil {
					return "", fmt.Errorf("复制文件失败: %w", err)
				}
			}
		}
	}

	// 创建zip文件
	zipFile := fmt.Sprintf("%s.v%s.zip", moduleName, moduleInfo.Version)
	if err := utils.ZipDirectory(e.ctx, tmpDir, zipFile); err != nil {
		return "", fmt.Errorf("创建zip文件失败: %w", err)
	}

	return zipFile, nil
}

// ModuleImporter 模块导入器
type ModuleImporter struct {
	ctx context.Context
}

// NewModuleImporter 创建模块导入器
func NewModuleImporter(ctx context.Context) *ModuleImporter {
	return &ModuleImporter{ctx: ctx}
}

// Import 从zip包导入模块
func (i *ModuleImporter) Import(zipPath string) (string, error) {
	// 检查zip文件是否存在
	if !utils.PathExists(zipPath) {
		return "", fmt.Errorf("模块文件 '%s' 不存在", zipPath)
	}

	// 创建临时解压目录
	tmpDir := filepath.Join(utils.GetTmpDir(), "module_import_"+filepath.Base(zipPath))
	defer os.RemoveAll(tmpDir)

	// 解压zip文件
	if err := utils.UnzipFile(zipPath, tmpDir); err != nil {
		return "", fmt.Errorf("解压模块文件失败: %w", err)
	}

	// 读取模块配置文件
	configPattern := filepath.Join(tmpDir, "modules", "*", "module.json")
	configFiles, err := filepath.Glob(configPattern)
	if err != nil || len(configFiles) == 0 {
		return "", fmt.Errorf("未找到模块配置文件")
	}

	// 读取配置
	config, err := gjson.LoadPath(configFiles[0], gjson.Options{Safe: true})
	if err != nil {
		return "", fmt.Errorf("读取模块配置失败: %w", err)
	}

	moduleName := config.Get("name").String()
	if moduleName == "" {
		return "", fmt.Errorf("模块配置文件中未指定模块名称")
	}

	// 检查模块是否已存在
	modulePath := filepath.Join("modules", moduleName)
	if utils.PathExists(modulePath) {
		return "", fmt.Errorf("模块 '%s' 已存在", moduleName)
	}

	// 复制文件到目标位置
	filesMap := config.Get("files").Map()
	if filesMap != nil {
		for _, v := range filesMap {
			if arr, ok := v.([]interface{}); ok {
				for _, item := range arr {
					path := fmt.Sprintf("%v", item)
					srcPath := filepath.Join(tmpDir, path)
					dstPath := path

					// 确保源文件存在
					if !utils.PathExists(srcPath) {
						g.Log().Warningf(i.ctx, "文件不存在，跳过: %s", srcPath)
						continue
					}

					// 创建目标目录
					if err := utils.EnsureDir(filepath.Dir(dstPath)); err != nil {
						return "", fmt.Errorf("创建目录失败: %w", err)
					}

					// 复制文件或目录
					if gfile.IsDir(srcPath) {
						if err := gfile.CopyDir(srcPath, dstPath); err != nil {
							return "", fmt.Errorf("复制目录失败: %w", err)
						}
					} else {
						if err := utils.CopyFile(srcPath, dstPath); err != nil {
							return "", fmt.Errorf("复制文件失败: %w", err)
						}
					}
				}
			}
		}
	}

	return moduleName, nil
}
