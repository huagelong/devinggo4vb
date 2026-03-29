# Admin 扩展动作模板（Week 3 产物）

更新时间：2026-03-29
适用目录：`admin-ui/apps/backend/src/views/system/role/`

---

## 1. 模板概述

扩展动作模板适用于**列表 CRUD + 特殊业务动作**的场景，典型场景：
- 角色管理（列表 + 权限分配）
- 用户管理（列表 + 关联设置）

核心特征：
- 基础 CRUD 能力（使用 useCrudPage）
- 业务保护规则（如 superAdmin 不可操作）
- 独立 Modal 执行扩展动作（如菜单权限、数据权限）
- 扩展动作完成后刷新列表

---

## 2. 标准目录结构

```text
src/views/system/<module>/
  index.vue                              # 列表页
  components/
    <module>-modal.vue                  # 主表新增/编辑弹窗
    <module>-<action>-modal.vue         # 扩展动作弹窗（如 menu-permission、data-permission）
  model.ts                               # 页面类型定义
  schemas.ts                             # 列定义、表单默认值
  use-<module>-crud.ts                   # 主表 CRUD Hook
```

---

## 3. index.vue 模板

```vue
<script setup>
// 基础 CRUD
const {
  fetchTableData,
  pagination,
  tableData,
  selectedRowKeys,
  isRecycleBin,
  // ...
} = useXxxCrud();

// 扩展动作 Modal ref
const permissionModalRef = ref();

// 扩展动作触发
function handlePermission(row: XxxItem) {
  permissionModalRef.value?.open(row);
}

// 扩展动作完成回调
function handleSuccess() {
  void fetchTableData();
}
</script>

<template>
  <Page>
    <Table :data="tableData" :pagination="pagination">
      <!-- 业务保护示例：超级管理员行特殊处理 -->
      <template #action="{ row }">
        <div v-if="row.code === 'superAdmin'">
          <Button disabled>编辑</Button>
          <Button disabled>删除</Button>
        </div>
        <div v-else>
          <Button @click="handleEdit(row)">编辑</Button>
          <Button @click="handlePermission(row)">权限分配</Button>
          <Button @click="handleDelete(row)">删除</Button>
        </div>
      </template>

      <!-- 状态列特殊处理 -->
      <template #status="{ row }">
        <Switch
          :disabled="isRecycleBin || row.code === 'superAdmin'"
          :value="Number(row.status) === 1"
          @change="handleStatusChange(row, value)"
        />
      </template>
    </Table>

    <!-- 主表弹窗 -->
    <XxxModal ref="modalRef" @success="handleSuccess" />

    <!-- 扩展动作弹窗 -->
    <XxxPermissionModal ref="permissionModalRef" @success="handleSuccess" />
  </Page>
</template>
```

---

## 4. 业务保护模式

### 4.1 保护规则示例

```ts
// 超级管理员不可编辑/删除/禁用
function handleEdit(row: RoleListItem) {
  if (row.id === 1 || row.code === 'superAdmin') {
    message.error('超级管理员角色不可编辑');
    return;
  }
  roleModalRef.value?.open(row);
}

function handleDelete(row: RoleListItem) {
  if (row.id === 1 || row.code === 'superAdmin') {
    message.error('超级管理员角色不可删除');
    return;
  }
  // 执行删除
}

async function handleStatusChange(row: RoleListItem, checked: boolean) {
  if (row.code === 'superAdmin') {
    message.info('超级管理员角色不能禁用');
    return;
  }
  await changeRoleStatus({ id: row.id, status: checked ? 1 : 2 });
  await fetchTableData();
}
```

### 4.2 批量操作保护

```ts
async function handleBatchDelete() {
  const ids = toIds(selectedRowKeys.value);
  if (ids.includes(1)) {
    message.error('超级管理员角色不可删除');
    return;
  }
  // 执行批量删除
}
```

---

## 5. 扩展动作 Modal 实现

### 5.1 权限分配 Modal 模板

```vue
<script setup>
// 接收主行数据
const currentRow = ref<RoleListItem>();

// 获取当前角色的权限数据
async function fetchPermission() {
  const data = await getMenuByRole(currentRow.value!.id);
  // 处理数据，转换格式
}

// 提交权限更新
async function handleSubmit() {
  await updateRoleMenuPermission(currentRow.value!.id, {
    menu_ids: selectedMenuIds.value,
  });
  MessagePlugin.success('权限更新成功');
  emit('success');
  modalApi.close();
}

async function open(row: RoleListItem) {
  currentRow.value = row;
  modalApi.setState({ title: `菜单权限 - ${row.name}` });
  modalApi.open();
  await fetchPermission();
}
</script>

<template>
  <Modal>
    <Tree
      v-model:checked="selectedMenuIds"
      :data="menuTreeData"
      checkable
    />
    <template #footer>
      <Button @click="modalApi.close()">取消</Button>
      <Button theme="primary" @click="handleSubmit">保存</Button>
    </template>
  </Modal>
</template>
```

### 5.2 数据权限 Modal

类似结构，处理 `dept_ids` 和 `data_scope`。

---

## 6. 与普通 CRUD 的差异

| 特性 | 普通 CRUD | 扩展动作 |
|------|---------|---------|
| 列表 | 基础 CRUD | 基础 CRUD + 业务保护 |
| 操作列 | 编辑/删除 | 编辑/删除 + 扩展动作入口 |
| 状态列 | 正常切换 | 可受业务规则约束 |
| 批量操作 | 基础保护 | 需检查是否包含受保护行 |
| 扩展 Modal | 无 | 多个独立扩展 Modal |

---

## 7. 接入示例：role

**API**：`api/system/role.ts` - 包含 `updateRoleMenuPermission` / `updateRoleDataPermission`

**扩展动作 Modal**：
- `role-menu-permission-modal.vue` - 菜单权限
- `role-data-permission-modal.vue` - 数据权限

**业务保护**：
```ts
// superAdmin 行不可编辑/删除/禁用
// id=1 的角色受保护
```

---

## 8. 验收清单

- [ ] 列表 CRUD 正常
- [ ] 业务保护规则生效（superAdmin/id=1）
- [ ] 扩展动作 Modal 正常打开
- [ ] 扩展动作数据正确加载
- [ ] 扩展动作提交后列表刷新
- [ ] 批量操作包含受保护行时正确拦截
- [ ] typecheck 通过
