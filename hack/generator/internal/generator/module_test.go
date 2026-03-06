package generator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestModuleCreator_Create 测试创建新模块
func TestModuleCreator_Create(t *testing.T) {
	// 注意：由于模板文件位于项目目录，测试需要在项目根目录运行
	// 检查是否在正确的目录
	if _, err := os.Stat("hack/generator/templates"); os.IsNotExist(err) {
		t.Skip("跳过测试：模板文件不在当前目录，需要在项目根目录运行")
	}

	// 创建临时模块名（在modules目录下）
	moduleName := "testmodule_" + fmt.Sprintf("%d", os.Getpid())
	modulePath := filepath.Join("modules", moduleName)
	
	// 清理：测试结束后删除
	defer os.RemoveAll(modulePath)

	ctx := context.Background()
	creator := NewModuleCreator(ctx)

	// 创建测试模块
	err := creator.Create(moduleName)
	require.NoError(t, err, "创建模块应该成功")

	// 验证目录结构
	expectedDirs := []string{
		modulePath,
		filepath.Join(modulePath, "api"),
		filepath.Join(modulePath, "controller"),
		filepath.Join(modulePath, "logic"),
		filepath.Join(modulePath, "logic", "hook"),
		filepath.Join(modulePath, "logic", "middleware"),
		filepath.Join(modulePath, "logic", moduleName),
		filepath.Join(modulePath, "service"),
		filepath.Join(modulePath, "worker"),
	}

	for _, dir := range expectedDirs {
		assert.DirExists(t, dir, "目录应该存在: %s", dir)
	}

	// 验证关键文件存在
	expectedFiles := []string{
		filepath.Join(modulePath, "module.go"),
		filepath.Join(modulePath, "module.json"),
		filepath.Join(modulePath, "logic", "logic.go"),
	}

	for _, file := range expectedFiles {
		assert.FileExists(t, file, "文件应该存在: %s", file)
	}

	// 验证module.json内容
	configPath := filepath.Join(modulePath, "module.json")
	jsonData, err := gjson.LoadPath(configPath, gjson.Options{Safe: true})
	require.NoError(t, err, "读取module.json应该成功")
	assert.Equal(t, moduleName, jsonData.Get("name").String(), "模块名应该正确")
	assert.Equal(t, "1.0.0", jsonData.Get("version").String(), "版本应该是1.0.0")

	t.Logf("✓ 模块创建测试通过: %d个目录, %d个文件", len(expectedDirs), len(expectedFiles))
}

// TestModuleCreator_CreateExisting 测试创建已存在的模块
func TestModuleCreator_CreateExisting(t *testing.T) {
	if _, err := os.Stat("hack/generator/templates"); os.IsNotExist(err) {
		t.Skip("跳过测试：模板文件不在当前目录")
	}

	ctx := context.Background()
	creator := NewModuleCreator(ctx)

	moduleName := "existingmodule_" + fmt.Sprintf("%d", os.Getpid())
	modulePath := filepath.Join("modules", moduleName)
	defer os.RemoveAll(modulePath)

	// 第一次创建应该成功
	err := creator.Create(moduleName)
	require.NoError(t, err)

	// 第二次创建同名模块应该失败
	err = creator.Create(moduleName)
	assert.Error(t, err, "创建已存在的模块应该返回错误")
	assert.Contains(t, err.Error(), "已存在", "错误消息应该提及模块已存在")

	t.Logf("✓ 模块重复创建测试通过，错误: %v", err)
}

// TestModuleCloner_Clone 测试克隆模块
func TestModuleCloner_Clone(t *testing.T) {
	if _, err := os.Stat("hack/generator/templates"); os.IsNotExist(err) {
		t.Skip("跳过测试：模板文件不在当前目录")
	}

	ctx := context.Background()

	// 先创建源模块
	sourceModule := "sourcemodule_" + fmt.Sprintf("%d", os.Getpid())
	sourcePath := filepath.Join("modules", sourceModule)
	defer os.RemoveAll(sourcePath)

	creator := NewModuleCreator(ctx)
	err := creator.Create(sourceModule)
	require.NoError(t, err, "创建源模块应该成功")

	// 克隆模块
	targetModule := "targetmodule_" + fmt.Sprintf("%d", os.Getpid())
	targetPath := filepath.Join("modules", targetModule)
	defer os.RemoveAll(targetPath)

	cloner := NewModuleCloner(ctx)
	err = cloner.Clone(sourceModule, targetModule)
	require.NoError(t, err, "克隆模块应该成功")

	// 验证目标模块存在
	assert.DirExists(t, targetPath, "目标模块目录应该存在")

	// 验证关键文件存在
	targetModuleFile := filepath.Join(targetPath, "module.go")
	assert.FileExists(t, targetModuleFile, "module.go应该存在")

	// 验证文件内容中的模块名已替换
	content := gfile.GetContents(targetModuleFile)
	assert.Contains(t, content, targetModule, "module.go应该包含新的模块名")
	assert.NotContains(t, content, sourceModule, "module.go不应该包含旧的模块名")

	// 验证module.json中的名称已更新
	configPath := filepath.Join(targetPath, "module.json")
	jsonData, err := gjson.LoadPath(configPath, gjson.Options{Safe: true})
	require.NoError(t, err)
	assert.Equal(t, targetModule, jsonData.Get("name").String(), "module.json中的名称应该已更新")

	t.Logf("✓ 模块克隆测试通过: %s -> %s", sourceModule, targetModule)
}

// TestModuleCloner_CloneNonExistent 测试克隆不存在的模块
func TestModuleCloner_CloneNonExistent(t *testing.T) {
	ctx := context.Background()
	cloner := NewModuleCloner(ctx)

	// 克隆不存在的模块应该失败
	err := cloner.Clone("nonexistent999999", "target999999")
	assert.Error(t, err, "克隆不存在的模块应该返回错误")
	assert.Contains(t, err.Error(), "不存在", "错误消息应该提及源模块不存在")

	t.Logf("✓ 克隆不存在模块测试通过，错误: %v", err)
}

// TestModuleCloner_CloneToExisting 测试克隆到已存在的目标
func TestModuleCloner_CloneToExisting(t *testing.T) {
	if _, err := os.Stat("hack/generator/templates"); os.IsNotExist(err) {
		t.Skip("跳过测试：模板文件不在当前目录")
	}

	ctx := context.Background()
	creator := NewModuleCreator(ctx)

	// 创建源模块和目标模块
	sourceModule := "source_" + fmt.Sprintf("%d", os.Getpid())
	targetModule := "target_" + fmt.Sprintf("%d", os.Getpid())
	sourcePath := filepath.Join("modules", sourceModule)
	targetPath := filepath.Join("modules", targetModule)
	defer os.RemoveAll(sourcePath)
	defer os.RemoveAll(targetPath)

	err := creator.Create(sourceModule)
	require.NoError(t, err)

	err = creator.Create(targetModule)
	require.NoError(t, err)

	// 克隆到已存在的目标应该失败
	cloner := NewModuleCloner(ctx)
	err = cloner.Clone(sourceModule, targetModule)
	assert.Error(t, err, "克隆到已存在的目标应该返回错误")
	assert.Contains(t, err.Error(), "已存在", "错误消息应该提及目标模块已存在")

	t.Logf("✓ 克隆到已存在目标测试通过，错误: %v", err)
}

// TestModuleExporterImporter_RoundTrip 测试模块导出导入完整流程
func TestModuleExporterImporter_RoundTrip(t *testing.T) {
	if _, err := os.Stat("hack/generator/templates"); os.IsNotExist(err) {
		t.Skip("跳过测试：模板文件不在当前目录")
	}

	ctx := context.Background()

	// 1. 创建测试模块
	moduleName := "exportmodule_" + fmt.Sprintf("%d", os.Getpid())
	modulePath := filepath.Join("modules", moduleName)
	
	creator := NewModuleCreator(ctx)
	err := creator.Create(moduleName)
	require.NoError(t, err, "创建模块应该成功")

	// 2. 导出模块
	exporter := NewModuleExporter(ctx)
	zipFile, err := exporter.Export(moduleName)
	require.NoError(t, err, "导出模块应该成功")
	defer os.Remove(zipFile)

	// 验证zip文件存在且有内容
	assert.FileExists(t, zipFile, "zip文件应该存在")
	stat, err := os.Stat(zipFile)
	require.NoError(t, err)
	assert.Greater(t, stat.Size(), int64(100), "zip文件应该有内容")

	t.Logf("✓ 模块导出成功: %s (大小: %d bytes)", zipFile, stat.Size())

	// 3. 删除原模块
	err = os.RemoveAll(modulePath)
	require.NoError(t, err)

	// 4. 导入模块
	importer := NewModuleImporter(ctx)
	importedName, err := importer.Import(zipFile)
	require.NoError(t, err, "导入模块应该成功")
	assert.Equal(t, moduleName, importedName, "导入的模块名应该正确")
	defer os.RemoveAll(modulePath) // 清理

	// 5. 验证导入的模块
	assert.DirExists(t, modulePath, "导入后模块目录应该存在")

	// 验证关键文件
	expectedFiles := []string{
		filepath.Join(modulePath, "module.go"),
		filepath.Join(modulePath, "module.json"),
		filepath.Join(modulePath, "logic", "logic.go"),
	}

	for _, file := range expectedFiles {
		assert.FileExists(t, file, "文件应该存在: %s", file)
	}

	t.Logf("✓ 模块导入成功: %s", importedName)
	t.Logf("✓ 完整流程测试通过: 创建 -> 导出 -> 删除 -> 导入")
}

// TestModuleExporter_ExportNonExistent 测试导出不存在的模块
func TestModuleExporter_ExportNonExistent(t *testing.T) {
	ctx := context.Background()
	exporter := NewModuleExporter(ctx)

	// 导出不存在的模块应该失败
	_, err := exporter.Export("nonexistent999999")
	assert.Error(t, err, "导出不存在的模块应该返回错误")

	t.Logf("✓ 导出不存在模块测试通过，错误: %v", err)
}

// TestModuleImporter_ImportNonExistent 测试导入不存在的zip文件
func TestModuleImporter_ImportNonExistent(t *testing.T) {
	ctx := context.Background()
	importer := NewModuleImporter(ctx)

	// 导入不存在的zip文件应该失败
	_, err := importer.Import("nonexistent999999.zip")
	assert.Error(t, err, "导入不存在的文件应该返回错误")
	assert.Contains(t, err.Error(), "不存在", "错误消息应该提及文件不存在")

	t.Logf("✓ 导入不存在文件测试通过，错误: %v", err)
}

// TestModuleImporter_ImportToExisting 测试导入已存在的模块
func TestModuleImporter_ImportToExisting(t *testing.T) {
	if _, err := os.Stat("hack/generator/templates"); os.IsNotExist(err) {
		t.Skip("跳过测试：模板文件不在当前目录")
	}

	ctx := context.Background()

	// 创建并导出模块
	moduleName := "dupmodule_" + fmt.Sprintf("%d", os.Getpid())
	modulePath := filepath.Join("modules", moduleName)
	defer os.RemoveAll(modulePath)

	creator := NewModuleCreator(ctx)
	err := creator.Create(moduleName)
	require.NoError(t, err)

	exporter := NewModuleExporter(ctx)
	zipFile, err := exporter.Export(moduleName)
	require.NoError(t, err)
	defer os.Remove(zipFile)

	// 尝试导入到已存在的模块应该失败
	importer := NewModuleImporter(ctx)
	_, err = importer.Import(zipFile)
	assert.Error(t, err, "导入到已存在的模块应该返回错误")
	assert.Contains(t, err.Error(), "已存在", "错误消息应该提及模块已存在")

	t.Logf("✓ 导入到已存在模块测试通过，错误: %v", err)
}

// BenchmarkModuleCreator_Create 性能测试：创建模块
func BenchmarkModuleCreator_Create(b *testing.B) {
	tmpDir, _ := os.MkdirTemp("", "module_benchmark_*")
	defer os.RemoveAll(tmpDir)

	originalWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalWd)

	ctx := context.Background()
	creator := NewModuleCreator(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		moduleName := filepath.Join("bench", "module"+filepath.Base(tmpDir)+fmt.Sprintf("%d", i))
		_ = creator.Create(moduleName)
	}
}

// BenchmarkModuleExporter_Export 性能测试：导出模块
func BenchmarkModuleExporter_Export(b *testing.B) {
	tmpDir, _ := os.MkdirTemp("", "module_benchmark_*")
	defer os.RemoveAll(tmpDir)

	originalWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalWd)

	ctx := context.Background()

	// 创建测试模块
	creator := NewModuleCreator(ctx)
	_ = creator.Create("benchmodule")

	exporter := NewModuleExporter(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		zipFile, _ := exporter.Export("benchmodule")
		os.Remove(zipFile)
	}
}
