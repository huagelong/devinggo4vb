// Package worker
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package worker_test

import (
	"context"
	"devinggo/modules/system/pkg/worker"
	"testing"
	"time"

	"github.com/hibiken/asynq"
)

// TestData 测试数据结构
type TestData struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

// TestWorker 测试Worker
type TestWorker struct {
	executed bool
}

func (w *TestWorker) GetType() string {
	return "test_task"
}

func (w *TestWorker) Execute(ctx context.Context, t *asynq.Task) error {
	w.executed = true
	data, err := worker.GetParameters[TestData](ctx, t)
	if err != nil {
		return err
	}
	worker.GetLogger().Infof(ctx, "执行任务: %+v", data)
	return nil
}

// TestNewManager 测试创建Manager
func TestNewManager(t *testing.T) {
	ctx := context.Background()
	mgr := worker.New(ctx)
	if mgr == nil {
		t.Fatal("Manager should not be nil")
	}
}

// TestRegisterWorker 测试注册Worker
func TestRegisterWorker(t *testing.T) {
	ctx := context.Background()
	mgr := worker.New(ctx)

	testWorker := &TestWorker{}
	result := mgr.RegisterWorker(testWorker)

	// 测试链式调用
	if result != mgr {
		t.Error("RegisterWorker should return the manager for chaining")
	}
}

// TestTaskBuilderBasic 测试基础TaskBuilder
func TestTaskBuilderBasic(t *testing.T) {
	ctx := context.Background()

	// 测试创建TaskBuilder
	builder := worker.NewTaskBuilder(ctx, "test_task")
	if builder == nil {
		t.Fatal("TaskBuilder should not be nil")
	}

	// 测试设置数据
	builder.WithData(TestData{
		Name:  "test",
		Value: 123,
	})

	// 测试构建Task
	task := builder.Build()
	if task == nil {
		t.Fatal("Task should not be nil")
	}

	if task.Type() != "test_task" {
		t.Errorf("Task type should be 'test_task', got '%s'", task.Type())
	}
}

// TestTaskBuilderChaining 测试链式调用
func TestTaskBuilderChaining(t *testing.T) {
	ctx := context.Background()

	// 测试链式调用
	task := worker.NewTaskBuilder(ctx, "test_task").
		WithData(TestData{Name: "test", Value: 123}).
		WithQueue("critical").
		WithTaskID("unique_id").
		WithDelay(5 * time.Second).
		Build()

	if task == nil {
		t.Fatal("Task should not be nil")
	}
}

// TestTaskBuilderWithQueue 测试设置队列
func TestTaskBuilderWithQueue(t *testing.T) {
	ctx := context.Background()

	builder := worker.NewTaskBuilder(ctx, "test_task").
		WithData(TestData{Name: "test"}).
		WithQueue("critical")

	task := builder.Build()
	if task == nil {
		t.Fatal("Task should not be nil")
	}
}

// TestTaskBuilderWithDelay 测试延迟执行
func TestTaskBuilderWithDelay(t *testing.T) {
	ctx := context.Background()

	delay := 10 * time.Second
	builder := worker.NewTaskBuilder(ctx, "test_task").
		WithData(TestData{Name: "test"}).
		WithDelay(delay)

	task := builder.Build()
	if task == nil {
		t.Fatal("Task should not be nil")
	}
}

// TestTaskBuilderWithProcessAt 测试指定时间执行
func TestTaskBuilderWithProcessAt(t *testing.T) {
	ctx := context.Background()

	processAt := time.Now().Add(1 * time.Hour)
	builder := worker.NewTaskBuilder(ctx, "test_task").
		WithData(TestData{Name: "test"}).
		WithProcessAt(processAt)

	task := builder.Build()
	if task == nil {
		t.Fatal("Task should not be nil")
	}
}

// TestTaskBuilderWithRetention 测试任务保留时间
func TestTaskBuilderWithRetention(t *testing.T) {
	ctx := context.Background()

	builder := worker.NewTaskBuilder(ctx, "test_task").
		WithData(TestData{Name: "test"}).
		WithTaskID("unique_task").
		WithRetention(24 * time.Hour)

	task := builder.Build()
	if task == nil {
		t.Fatal("Task should not be nil")
	}
}

// TestGetParameters 测试参数解析
func TestGetParameters(t *testing.T) {
	ctx := context.Background()

	// 构建任务
	originalData := TestData{
		Name:  "test_name",
		Value: 42,
	}

	task := worker.NewTaskBuilder(ctx, "test_task").
		WithData(originalData).
		Build()

	// 解析参数
	parsedData, err := worker.GetParameters[TestData](ctx, task)
	if err != nil {
		t.Fatalf("GetParameters failed: %v", err)
	}

	// 验证数据
	if parsedData.Name != originalData.Name {
		t.Errorf("Name mismatch: expected '%s', got '%s'", originalData.Name, parsedData.Name)
	}

	if parsedData.Value != originalData.Value {
		t.Errorf("Value mismatch: expected %d, got %d", originalData.Value, parsedData.Value)
	}
}

// TestMultipleWorkerRegistration 测试注册多个Worker
func TestMultipleWorkerRegistration(t *testing.T) {
	ctx := context.Background()
	mgr := worker.New(ctx)

	worker1 := &TestWorker{}
	worker2 := &TestWorker{}

	// 链式注册
	mgr.RegisterWorker(worker1).RegisterWorker(worker2)
}

// BenchmarkTaskBuilderCreate 基准测试：创建TaskBuilder
func BenchmarkTaskBuilderCreate(b *testing.B) {
	ctx := context.Background()
	data := TestData{Name: "benchmark", Value: 100}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		worker.NewTaskBuilder(ctx, "test_task").
			WithData(data).
			Build()
	}
}

// BenchmarkTaskBuilderChaining 基准测试：链式调用
func BenchmarkTaskBuilderChaining(b *testing.B) {
	ctx := context.Background()
	data := TestData{Name: "benchmark", Value: 100}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		worker.NewTaskBuilder(ctx, "test_task").
			WithData(data).
			WithQueue("default").
			WithTaskID("test_id").
			WithDelay(0).
			Build()
	}
}

// ExampleTaskBuilder 示例：使用TaskBuilder
func ExampleTaskBuilder() {
	ctx := context.Background()

	// 创建并发送任务
	worker.NewTaskBuilder(ctx, "send_email").
		WithData(map[string]interface{}{
			"to":      "user@example.com",
			"subject": "欢迎",
			"body":    "欢迎注册！",
		}).
		WithQueue("critical").
		WithDelay(5 * time.Second).
		Send()
}

// ExampleManager 示例：使用Manager
func ExampleManager() {
	ctx := context.Background()

	// 创建Manager
	mgr := worker.New(ctx)

	// 注册Worker
	mgr.RegisterWorker(&TestWorker{})

	// 启动服务器
	go mgr.RunServer()
}
