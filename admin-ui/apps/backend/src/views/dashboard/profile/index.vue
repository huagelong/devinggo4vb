<script setup lang="ts">
import type { LogApi } from '#/api/system/log';

import { onMounted, reactive, ref } from 'vue';

// 瀵煎叆鐢ㄦ埛淇℃伅鐩稿叧鐨?Store
import { useUserStore } from '@vben/stores';

import {
  Button,
  Form,
  FormItem,
  Input,
  MessagePlugin,
  TabPanel,
  Tabs,
  Tag,
  Textarea,
  Timeline,
  TimelineItem,
} from 'tdesign-vue-next';

import {
  getLoginLogListApi,
  getOperationLogListApi,
  modifyPasswordApi,
  updateUserInfoApi,
} from '#/api/core/profile';
import { getSystemInfoApi } from '#/api/core/user';
import { uploadImageFileApi } from '#/api/system/upload';

const userStore = useUserStore();

// 宸︿晶 Tabs
const leftTab = ref('info');
// 鍙充晶 Tabs
const rightTab = ref('loginLog');

// 涓汉璧勬枡琛ㄥ崟
const userInfoForm = reactive({
  username: '',
  nickname: '',
  phone: '',
  email: '',
  signed: '',
  avatar: '',
});

// 瀹夊叏璁剧疆琛ㄥ崟
const securityForm = reactive({
  oldPassword: '',
  newPassword: '',
  newPasswordConfirmation: '',
});

// 鏃ュ織鏁版嵁
const loginLogs = ref<LogApi.LoginLogItem[]>([]);
const operationLogs = ref<LogApi.OperLogItem[]>([]);

// 鑾峰彇涓汉淇℃伅
async function fetchUserInfo() {
  try {
    const res = await getSystemInfoApi();
    if (res && res.user) {
      userInfoForm.username = res.user.username || '';
      userInfoForm.nickname = res.user.nickname || '';
      userInfoForm.phone = res.user.phone || '';
      userInfoForm.email = res.user.email || '';
      userInfoForm.signed = res.user.signed || '';
      userInfoForm.avatar = res.user.avatar || '';
    }
  } catch (error) {
    console.error('鑾峰彇涓汉淇℃伅澶辫触', error);
  }
}

// 鎻愪氦涓汉璧勬枡鏇存柊
async function handleUpdateInfo() {
  try {
    await updateUserInfoApi({
      nickname: userInfoForm.nickname,
      phone: userInfoForm.phone,
      email: userInfoForm.email,
      signed: userInfoForm.signed,
    });
    MessagePlugin.success('涓汉璧勬枡鏇存柊鎴愬姛');
    // 鏇存柊瀹屾垚鍚庨噸鏂拌幏鍙栨暟鎹?    fetchUserInfo();
  } catch {
    MessagePlugin.error('涓汉璧勬枡鏇存柊澶辫触');
  }
}

// 鎻愪氦淇敼瀵嗙爜
async function handleUpdatePassword() {
  if (securityForm.newPassword !== securityForm.newPasswordConfirmation) {
    MessagePlugin.error('两次输入的新密码不一致');
    return;
  }
  try {
    await modifyPasswordApi(securityForm);
    MessagePlugin.success('瀵嗙爜淇敼鎴愬姛');
    // 娓呯┖瀵嗙爜琛ㄥ崟
    securityForm.oldPassword = '';
    securityForm.newPassword = '';
    securityForm.newPasswordConfirmation = '';
  } catch (error) {
    // 閿欒鍦ㄨ姹傛嫤鎴櫒閫氬父鏈夋彁绀?    console.error('瀵嗙爜淇敼澶辫触', error);
  }
}

// 鍥剧墖涓婁紶澶勭悊
function triggerUpload() {
  const fileInput = document.createElement('input');
  fileInput.type = 'file';
  fileInput.accept = 'image/*';
  fileInput.addEventListener('change', async (e: Event) => {
    const file = (e.target as HTMLInputElement).files?.[0];
    if (!file) return;
    try {
      const res: any = await uploadImageFileApi(file);
      // 鏍规嵁鍚庣杩斿洖鏍煎紡鍙栧浘鐗嘦RL
      if (res && res.url) {
        userInfoForm.avatar = res.url;
        await updateUserInfoApi({
          avatar: res.url,
        });
        userStore.setUserInfo({
          ...userStore.userInfo,
          avatar: res.url,
        } as any);
        MessagePlugin.success('澶村儚涓婁紶鎴愬姛');
      }
    } catch (error) {
      console.error('涓婁紶澶辫触', error);
      MessagePlugin.error('澶村儚涓婁紶澶辫触');
    }
  });
  fileInput.click();
}

// 鑾峰彇鏃ュ織
async function fetchLogs() {
  try {
    const loginRes: any = await getLoginLogListApi({ page: 1, pageSize: 10 });
    if (loginRes && loginRes.items) {
      loginLogs.value = loginRes.items;
    }

    const opRes: any = await getOperationLogListApi({ page: 1, pageSize: 10 });
    if (opRes && opRes.items) {
      operationLogs.value = opRes.items;
    }
  } catch (error) {
    console.error('鑾峰彇鏃ュ織澶辫触', error);
  }
}

onMounted(() => {
  fetchUserInfo();
  fetchLogs();
});
</script>

<template>
  <div class="h-full p-4 overflow-auto bg-[var(--vben-color-background)]">
    <!-- 椤堕儴 Banner -->
    <div
      class="relative flex flex-col items-center justify-center w-full h-48 overflow-hidden rounded-t-lg bg-blue-50 dark:bg-blue-900/20"
    >
      <!-- 铏氭嫙鑳屾櫙瑁呴グ -->
      <div class="absolute inset-0 pointer-events-none opacity-50">
        <!-- 绫讳技璁捐鍥句腑鐨勫嚑浣曞厓绱?-->
        <div
          class="absolute top-10 left-20 w-12 h-12 bg-teal-300 rounded-full blur-md"
        ></div>
        <div
          class="absolute bottom-10 left-40 w-6 h-6 bg-orange-500 rounded-full blur-sm"
        ></div>
        <div
          class="absolute top-20 right-20 w-16 h-4 bg-indigo-600 rounded rotate-45 blur-sm"
        ></div>
      </div>

      <!-- 澶村儚鍜屼笂浼?-->
      <div class="relative z-10 z-20 mt-4 group">
        <div
          @click="triggerUpload"
          class="flex items-center justify-center w-24 h-24 overflow-hidden border-4 border-white rounded-full shadow-lg bg-gray-100 hover:bg-gray-200 cursor-pointer"
        >
          <img
            v-if="userInfoForm.avatar"
            :src="userInfoForm.avatar"
            class="object-cover w-full h-full"
          />
          <svg
            v-else
            xmlns="http://www.w3.org/2000/svg"
            class="w-10 h-10 text-gray-400"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
            <circle cx="12" cy="7" r="4" />
          </svg>

          <div
            class="absolute inset-0 flex flex-col items-center justify-center text-white bg-black bg-opacity-50 opacity-0 group-hover:opacity-100 transition-opacity"
          >
            <span class="text-2xl">+</span>
            <span class="text-xs">鏈湴涓婁紶</span>
          </div>
        </div>
      </div>

      <!-- 瑙掕壊鏍囩 -->
      <div class="z-10 mt-3 mb-2">
        <Tag
          v-for="role in userStore.userInfo?.roles"
          :key="role"
          theme="primary"
          shape="round"
          size="large"
        >
          {{ role === 'superAdmin' ? '超级管理员' : role }}
        </Tag>
      </div>
    </div>

    <!-- 涓嬫柟涓ゅ垪鍐呭 -->
    <div class="flex flex-col gap-4 mt-4 md:flex-row">
      <!-- 宸︽爮锛氫釜浜鸿祫鏂?瀹夊叏璁剧疆 -->
      <div
        class="flex-1 p-4 bg-white rounded shadow-sm dark:bg-[var(--vben-color-background-elevated)] min-h-[500px]"
      >
        <Tabs v-model="leftTab" class="h-full">
          <TabPanel value="info" label="涓汉璧勬枡">
            <div class="pt-6 mt-4">
              <Form
                :data="userInfoForm"
                label-align="left"
                label-width="100px"
                @submit="handleUpdateInfo"
              >
                <FormItem label="璐︽埛鍚?" name="username">
                  <Input v-model="userInfoForm.username" disabled />
                </FormItem>
                <FormItem label="鏄电О" name="nickname">
                  <Input
                    v-model="userInfoForm.nickname"
                    placeholder="璇疯緭鍏ユ樀绉?"
                  />
                </FormItem>
                <FormItem label="鎵嬫満" name="phone">
                  <Input
                    v-model="userInfoForm.phone"
                    placeholder="璇疯緭鍏ユ墜鏈哄彿"
                  />
                </FormItem>
                <FormItem label="閭" name="email">
                  <Input
                    v-model="userInfoForm.email"
                    placeholder="璇疯緭鍏ラ偖绠?"
                  />
                </FormItem>
                <FormItem label="涓汉绛惧悕" name="signed">
                  <Textarea
                    v-model="userInfoForm.signed"
                    placeholder="璇疯緭鍏ヤ釜浜虹鍚?"
                    :maxlength="255"
                    :autosize="{ minRows: 3, maxRows: 5 }"
                  />
                </FormItem>
                <FormItem>
                  <Button
                    theme="default"
                    type="submit"
                    class="bg-gray-800 text-white hover:bg-gray-700"
                  >
                    淇濆瓨
                  </Button>
                </FormItem>
              </Form>
            </div>
          </TabPanel>

          <TabPanel value="security" label="瀹夊叏璁剧疆">
            <div class="pt-6 mt-4">
              <Form
                :data="securityForm"
                label-align="left"
                label-width="100px"
                @submit="handleUpdatePassword"
              >
                <FormItem label="鏃у瘑鐮?" name="oldPassword" required-mark>
                  <Input
                    type="password"
                    v-model="securityForm.oldPassword"
                    placeholder="璇疯緭鍏ユ棫瀵嗙爜"
                  />
                </FormItem>
                <FormItem label="鏂板瘑鐮?" name="newPassword" required-mark>
                  <Input
                    type="password"
                    v-model="securityForm.newPassword"
                    placeholder="璇疯緭鍏ユ柊瀵嗙爜"
                  />
                </FormItem>
                <FormItem
                  label="纭瀵嗙爜"
                  name="newPasswordConfirmation"
                  required-mark
                >
                  <Input
                    type="password"
                    v-model="securityForm.newPasswordConfirmation"
                    placeholder="璇峰啀娆¤緭鍏ユ柊瀵嗙爜"
                  />
                </FormItem>
                <FormItem>
                  <Button
                    theme="default"
                    type="submit"
                    class="bg-gray-800 text-white hover:bg-gray-700"
                  >
                    淇濆瓨
                  </Button>
                </FormItem>
              </Form>
            </div>
          </TabPanel>
        </Tabs>
      </div>

      <!-- 鍙虫爮锛氭棩蹇?-->
      <div
        class="flex-1 p-4 bg-white rounded shadow-sm dark:bg-[var(--vben-color-background-elevated)] min-h-[500px]"
      >
        <Tabs v-model="rightTab" class="h-full">
          <TabPanel value="loginLog" label="鐧诲綍鏃ュ織">
            <div class="pt-6 mt-4 overflow-y-auto max-h-[400px]">
              <Timeline>
                <TimelineItem
                  v-for="log in loginLogs"
                  :key="log.id"
                  theme="primary"
                >
                  <div
                    class="text-sm font-medium text-gray-800 dark:text-gray-200"
                  >
                    鎮ㄤ簬 {{ log.login_time || log.created_at }} 鐧诲綍绯荤粺锛寋{
                      log.status === 1 ? '鐧诲綍鎴愬姛' : '鐧诲綍澶辫触'
                    }}
                  </div>
                  <div class="mt-1 text-xs text-gray-500">
                    鍦扮悊浣嶇疆: {{ log.ip_location || '鏈煡' }}锛屾搷浣滅郴缁?
                    {{ log.os || '鏈煡' }}
                  </div>
                </TimelineItem>
                <div
                  v-if="loginLogs.length === 0"
                  class="text-center text-gray-400 py-10"
                >
                  鏆傛棤鏃ュ織
                </div>
              </Timeline>
            </div>
          </TabPanel>

          <TabPanel value="opLog" label="鎿嶄綔鏃ュ織">
            <div class="pt-6 mt-4 overflow-y-auto max-h-[400px]">
              <Timeline>
                <TimelineItem
                  v-for="log in operationLogs"
                  :key="log.id"
                  theme="primary"
                >
                  <div
                    class="text-sm font-medium text-gray-800 dark:text-gray-200"
                  >
                    鎮ㄤ簬 {{ log.created_at }} 鎵ц浜?                    {{ log.service_name || '鎿嶄綔' }}
                  </div>
                  <div class="mt-1 text-xs text-gray-500">
                    鍦扮悊浣嶇疆: {{ log.ip_location || '鏈煡' }}锛屾柟寮?
                    {{ log.method }}锛岃矾鐢? {{ log.router }}
                  </div>
                </TimelineItem>
                <div
                  v-if="operationLogs.length === 0"
                  class="text-center text-gray-400 py-10"
                >
                  鏆傛棤鏃ュ織
                </div>
              </Timeline>
            </div>
          </TabPanel>
        </Tabs>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* 閬垮厤瑕嗙洊鍏ㄥ眬琛ㄥ崟鏍峰紡 */
</style>
