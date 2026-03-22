import { readFileSync, writeFileSync } from 'fs';

const filePath = 'admin-ui/apps/backend/src/views/system/user/index.vue';
let content = readFileSync(filePath, 'utf8');

// 1. Add DialogPlugin to tdesign-vue-next imports
const old1 = `import {
  Button,
  Dropdown,
  MessagePlugin,
  Popconfirm,
  Switch,
} from 'tdesign-vue-next';`;
const new1 = `import {
  Button,
  DialogPlugin,
  Dropdown,
  MessagePlugin,
  Popconfirm,
  Switch,
} from 'tdesign-vue-next';`;
content = content.replace(old1, new1);
console.log('1. DialogPlugin import:', content.includes('DialogPlugin') ? 'OK' : 'FAIL');

// 2. Fix icon imports (remove RefreshIcon, sort alphabetically)
const old2 = `import {
  AddIcon,
  MoreIcon,
  DeleteIcon,
  DownloadIcon,
  EditIcon,
  RefreshIcon,
  UploadIcon,
} from 'tdesign-icons-vue-next';`;
const new2 = `import {
  AddIcon,
  DeleteIcon,
  DownloadIcon,
  EditIcon,
  MoreIcon,
  UploadIcon,
} from 'tdesign-icons-vue-next';`;
content = content.replace(old2, new2);
console.log('2. RefreshIcon removed:', !content.includes('RefreshIcon') ? 'OK' : 'FAIL');

// 3. Reduce action column width from 260 to 200
const old3 = `        field: 'action',
        title: '操作',
        width: 260,`;
const new3 = `        field: 'action',
        title: '操作',
        width: 200,`;
content = content.replace(old3, new3);
console.log('3. Column width 200:', content.includes('width: 200,') ? 'OK' : 'FAIL');

// 4. Update actionDropdownOptions and handleActionDropdownClick
const old4 = `const actionDropdownOptions = [
  { content: '重置密码', value: 'reset_password' },
];

function handleActionDropdownClick(data: any, row: any) {
  if (data.value === 'reset_password') {
    const confirmDia = DialogPlugin.confirm({
      header: '提示',
      body: '确认重置该用户密码吗？',
      onConfirm: () => {
        handleResetPassword(row);
        confirmDia.hide();
      },
      onClose: () => confirmDia.hide(),
    });
  }
}`;
const new4 = `const actionDropdownOptions = [
  { content: '重置密码', value: 'reset_password' },
  { content: '更新缓存', value: 'clear_cache' },
];

function handleActionDropdownClick(data: any, row: any) {
  if (data.value === 'reset_password') {
    const dialog = DialogPlugin.confirm({
      header: '提示',
      body: '确认重置该用户密码吗？',
      onConfirm: () => {
        handleResetPassword(row);
        dialog.hide();
      },
      onClose: () => dialog.hide(),
    });
  } else if (data.value === 'clear_cache') {
    const dialog = DialogPlugin.confirm({
      header: '提示',
      body: '确认更新该用户缓存吗？',
      onConfirm: () => {
        handleClearCache(row);
        dialog.hide();
      },
      onClose: () => dialog.hide(),
    });
  }
}`;
content = content.replace(old4, new4);
console.log('4. actionDropdown updated:', content.includes('clear_cache') ? 'OK' : 'FAIL');

// 5. Remove standalone 更新缓存 Popconfirm button from template
const old5 = `                <Dropdown
                  :options="actionDropdownOptions"
                  trigger="click"
                  @click="
                    (dropdownItem) =>
                      handleActionDropdownClick(dropdownItem, row)
                  "
                >
                  <Button size="small" theme="default" variant="text">
                    <template #icon><MoreIcon /></template>
                    更多
                  </Button>
                </Dropdown>
                <Popconfirm
                  content="确认更新缓存吗？"
                  @confirm="handleClearCache(row)"
                >
                  <Button size="small" theme="default" variant="text">
                    <template #icon><RefreshIcon /></template>
                    更新缓存
                  </Button>
                </Popconfirm>`;
const new5 = `                <Dropdown
                  :options="actionDropdownOptions"
                  trigger="click"
                  @click="(dropdownItem) => handleActionDropdownClick(dropdownItem, row)"
                >
                  <Button size="small" theme="default" variant="text">
                    <template #icon><MoreIcon /></template>
                    更多
                  </Button>
                </Dropdown>`;
content = content.replace(old5, new5);
console.log('5. Standalone 更新缓存 removed:', !content.includes('确认更新缓存吗') ? 'OK' : 'FAIL');

writeFileSync(filePath, content, 'utf8');
console.log('\nAll done! File written.');
