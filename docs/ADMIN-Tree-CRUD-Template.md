# Admin 树形 CRUD 模板（Week 3 产物）

更新时间：2026-03-29
适用目录：`admin-ui/apps/backend/src/views/system/{menu,dept}/`

---

## 1. 模板概述

树形 CRUD 模板适用于**层级结构数据**的展示与维护，典型场景：
- 菜单管理（menu）
- 部门管理（dept）

核心特征：
- 树形表格展示（可展开/折叠）
- 新增子级节点
- 排序直接在行内编辑
- 状态切换独立于行操作

---

## 2. 标准目录结构

```text
src/views/system/<module>/
  index.vue                  # 树形表格页（页面编排）
  components/
    <module>-modal.vue      # 新增/编辑弹窗
  model.ts                   # 页面类型定义
  schemas.ts                 # 列定义、表单默认值工厂函数
  use-<module>-page.ts       # 列表逻辑 Hook（独立实现）
```

---

## 3. index.vue 模板职责

```vue
<!-- 关键结构 -->
<template>
  <Page auto-content-height>
    <div class="flex h-full flex-col gap-3">
      <!-- 搜索区 -->
      <SearchForm />

      <!-- 工具栏 + 树形表格 -->
      <Table
        :tree="{ childrenKey: 'children', defaultExpandAll: true }"
        :columns="columns"
        :data="tableData"
      >
        <!-- 排序列：行内编辑 -->
        <template #sort="{ row }">
          <InputNumber :value="row.sort" @change="handleSortChange(value, row)" />
        </template>

        <!-- 状态列：独立切换 -->
        <template #status="{ row }">
          <Switch :value="Number(row.status) === 1" @change="handleStatusChange(row, value)" />
        </template>

        <!-- 操作列：子级新增 + 编辑 + 删除 -->
        <template #action="{ row }">
          <Button v-if="row.type === 'M'" @click="handleAdd(row.id)">新增子级</Button>
          <Button @click="handleEdit(row)">编辑</Button>
          <Popconfirm @confirm="handleDelete(row)">删除</Popconfirm>
        </template>
      </Table>
    </div>

    <ModuleModal ref="modalRef" @success="fetchTableData" />
  </Page>
</template>
```

---

## 4. use-<module>-page.ts 实现要点

```ts
// 关键实现
export function useMenuPage() {
  // 1. 数据类型用 TreeItem
  const tableData = ref<MenuApi.TreeItem[]>([]);

  // 2. 树形 API 直接返回树结构数组（不分页）
  async function fetchTableData() {
    const params = buildParams();
    tableData.value = isRecycleBin.value
      ? await getRecycleMenuTreeList(params)
      : await getMenuTreeList(params);
  }

  // 3. 排序/状态切换独立接口
  async function handleSortChange(value: number, row: MenuItem) {
    await updateMenuNumber({ id: row.id, numberName: 'sort', numberValue: value });
    await fetchTableData();
  }

  // 4. 新增时传递 parent_id
  function handleAdd(parentId = 0) {
    modalRef.value?.open({ parent_id: parentId });
  }
}
```

---

## 5. Modal 组件要点

```vue
<script setup>
// 上级选项使用 TreeSelect
const parentMenuOptions = ref<TreeOptionItem[]>([]);

async function fetchParentOptions() {
  // 调用 tree 接口获取选项
  const tree = await getMenuTreeOptions({ onlyMenu: true });
  parentMenuOptions.value = [
    { id: 0, label: '顶级菜单', value: 0, children: tree }
  ];
}

// 打开时注入 parent_id
async function open(data?: Partial<SubmitPayload>) {
  modalApi.setState({ title: data?.id ? '编辑' : '新增' });
  await fetchParentOptions();
  formApi.setValues({ parent_id: data?.parent_id ?? 0, ... });
}
</script>

<template>
  <Modal>
    <Form>
      <!-- TreeSelect 上级选择 -->
      <FormItem label="上级菜单">
        <TreeSelect v-model="form.parent_id" :data="parentMenuOptions" />
      </FormItem>
      <!-- 根据类型动态显隐字段 -->
      <FormItem label="图标" v-if="isFieldVisible('icon', form.type)">
        <Input />
      </FormItem>
    </Form>
  </Modal>
</template>
```

---

## 6. 树形 vs 普通列表差异

| 特性 | 普通 CRUD | 树形 CRUD |
|------|---------|----------|
| 数据结构 | 扁平列表 + 分页 | 树形 JSON + 不分页 |
| 新增 | 单一新增 | 顶级新增 + 子级新增 |
| parent_id | 无或隐藏字段 | TreeSelect 上级选择 |
| 排序 | 表格列排序 | 行内 InputNumber |
| 展开/折叠 | 无 | 工具栏或全展开 |

---

## 7. 接入示例：menu

**API**：`api/system/menu.ts` - `getMenuTreeList` / `getRecycleMenuTreeList`

**Model**：
```ts
export interface MenuTreeItem {
  children?: MenuTreeItem[];
  id: number;
  name: string;
  parent_id?: IdType;
  type?: MenuTypeValue;
  // ...
}
```

**Schema 工厂**：
```ts
export function createMenuTableColumns(): MenuTableColumn[] {
  return [
    { colKey: 'name', title: '菜单名称' },
    { colKey: 'type', title: '类型' },  // Tag 显示
    { colKey: 'sort', title: '排序' },  // InputNumber
    { colKey: 'status', title: '状态' }, // Switch
  ];
}
```

---

## 8. 验收清单

- [ ] 树形表格正确渲染（childrenKey）
- [ ] 新增顶级菜单正常
- [ ] 新增子级菜单正常（传递正确 parent_id）
- [ ] 编辑回填正确
- [ ] 排序行内编辑生效
- [ ] 状态切换生效
- [ ] 删除/回收站/恢复正常
- [ ] 工具栏刷新正常
- [ ] typecheck 通过
