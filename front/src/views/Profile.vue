<template>
  <div class="w-full h-full overflow-auto px-4">
    <div class="card">
      <div class="flex items-center justify-between mb-6">
        <h1 class="text-2xl font-bold">个人资料</h1>
      </div>

      <div v-if="loading" class="py-10 text-center text-gray-500">
        正在加载个人资料...
      </div>

      <div v-else-if="error" class="py-10 text-center text-red-500">
        {{ error }}
      </div>

      <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- 左侧个人信息 -->
        <div class="lg:col-span-2">
          <div class="flex items-start gap-6 mb-6">
            <div
              class="w-24 h-24 rounded-full bg-gray-200 overflow-hidden flex items-center justify-center text-3xl font-semibold text-[#2D5BFF]"
            >
              <img
                v-if="profile.avatar_url"
                :src="profile.avatar_url"
                :alt="profile.display_name"
                class="w-full h-full object-cover"
              />
              <span v-else>{{ profile.display_name?.slice(0, 1) || "U" }}</span>
            </div>
            <div class="flex-1">
              <h2 class="text-xl font-bold mb-2">{{ profile.display_name }}</h2>
              <p class="text-gray-600 mb-3">
                {{ profile.bio || "这个用户还没有填写简介。" }}
              </p>
              <div class="flex items-center gap-4 flex-wrap">
                <div
                  v-if="profile.role"
                  class="bg-blue-100 text-blue-800 px-3 py-1 rounded-full text-sm font-medium"
                >
                  {{ profile.role }}
                </div>
                <div
                  class="bg-green-100 text-green-800 px-3 py-1 rounded-full text-sm font-medium flex items-center gap-1"
                >
                  <span
                    class="w-2 h-2 rounded-full"
                    :class="statusDotClass"
                  ></span>
                  {{ statusLabel }}
                </div>
                <div
                  v-for="badge in profileBadges"
                  :key="badge"
                  class="bg-indigo-100 text-indigo-700 px-3 py-1 rounded-full text-sm font-medium"
                >
                  {{ badge }}
                </div>
              </div>
            </div>
          </div>

          <!-- 基本信息 -->
          <div class="space-y-4">
            <h3 class="font-bold text-lg">基本信息</h3>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <ProfileReadonlyInput label="姓名" :value="profile.basic_info?.real_name" />
              <ProfileReadonlyInput label="邮箱" :value="profile.basic_info?.email" />
              <ProfileReadonlyInput label="学校" :value="profile.basic_info?.school" />
              <ProfileReadonlyInput label="专业" :value="profile.basic_info?.major" />
              <ProfileReadonlyInput label="地区" :value="profile.basic_info?.location" />
              <ProfileReadonlyInput label="加入时间" :value="profile.basic_info?.join_date" />
            </div>
          </div>

          <!-- 技能标签 -->
          <div class="mt-6">
            <h3 class="font-bold text-lg mb-3">技能标签</h3>
            <div v-if="allSkills.length" class="flex flex-wrap gap-2">
              <span
                v-for="skill in allSkills"
                :key="skill"
                class="bg-blue-100 text-blue-800 px-3 py-1 rounded-full text-sm"
              >
                {{ skill }}
              </span>
            </div>
            <div v-else class="text-sm text-gray-500">暂无技能标签</div>
          </div>

          <!-- 账户操作 -->
          <div class="mt-8 border border-gray-100 rounded-xl p-5 bg-gray-50">
            <h3 class="font-bold text-lg mb-1">账户与安全</h3>
            <p class="text-sm text-gray-500">
              管理您的账户，退出登录后需要重新验证才能访问受限页面。
            </p>
            <div class="flex flex-wrap gap-3 mt-4">
              <button
                class="px-4 py-2 bg-red-500 text-white rounded-lg text-sm hover:bg-red-600 transition disabled:opacity-60 disabled:cursor-not-allowed"
                @click="handleLogout"
                :disabled="loggingOut"
              >
                {{ loggingOut ? "退出中..." : "退出登录" }}
              </button>
            </div>
            <p v-if="logoutError" class="text-sm text-red-500 mt-2">
              {{ logoutError }}
            </p>
          </div>
        </div>

        <!-- 右侧统计信息 -->
        <div class="space-y-4">
          <div class="bg-blue-50 p-4 rounded-lg">
            <div class="text-center">
              <div class="text-2xl font-bold text-blue-600">
                {{ studyStats.level_label }}
              </div>
              <div class="text-sm text-gray-600 mt-1">当前等级</div>
              <div class="mt-3">
                <div class="progress-bar">
                  <div
                    class="progress-fill"
                    :style="{ width: `${studyStats.progress_percent || 0}%` }"
                  ></div>
                </div>
                <div class="text-xs text-gray-600 mt-1">
                  距离下一级还需 {{ studyStats.distance_to_next || 0 }} 积分
                </div>
              </div>
            </div>
          </div>

          <div class="space-y-3">
            <ProfileStatCard
              icon="mdi:clock"
              icon-class="text-blue-600"
              label="总学习时长"
              :value="`${studyStats.total_study_hours || 0}h`"
              value-class="text-blue-600"
            />
            <ProfileStatCard
              icon="mdi:check-circle"
              icon-class="text-green-600"
              label="完成任务"
              :value="studyStats.tasks_completed || 0"
              value-class="text-green-600"
            />
            <ProfileStatCard
              icon="mdi:certificate"
              icon-class="text-purple-600"
              label="获得证书"
              :value="studyStats.certificates_count || 0"
              value-class="text-purple-600"
            />
            <ProfileStatCard
              icon="mdi:account-group"
              icon-class="text-orange-600"
              label="学习小组"
              :value="studyStats.study_groups || 0"
              value-class="text-orange-600"
            />
          </div>

        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from "vue";
import { useRouter } from "vue-router";
import {
  getUserProfile,
  getUserStudyStats,
  getUserSkills,
} from "@/api/modules/user";
import { DEFAULT_USER_ID, useCurrentUser } from "@/composables/useCurrentUser";
import { logout as logoutApi } from "@/api/modules/auth";
import { clearAuth } from "@/utils/auth";

const loading = ref(false);
const error = ref("");
const profile = reactive({});
const studyStats = reactive({});
const skills = reactive({
  primary: [],
  secondary: [],
});
const loggingOut = ref(false);
const logoutError = ref("");

const statusLabel = computed(() => {
  if (profile.status === "offline") return "离线";
  if (profile.status === "busy") return "忙碌";
  return "在线";
});

const statusDotClass = computed(() => {
  switch (profile.status) {
    case "offline":
      return "bg-gray-400";
    case "busy":
      return "bg-yellow-500";
    default:
      return "bg-green-500";
  }
});

const allSkills = computed(() => [
  ...(skills.primary || []),
  ...(skills.secondary || []),
]);

const profileBadges = computed(() => profile.badges || []);

const router = useRouter();
const { profile: currentUser, clearCurrentUser } = useCurrentUser();

const ProfileReadonlyInput = {
  name: "ProfileReadonlyInput",
  props: {
    label: {
      type: String,
      required: true,
    },
    value: {
      type: [String, Number],
      default: "",
    },
  },
  template: `
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1">{{ label }}</label>
      <input
        type="text"
        :value="value || '未填写'"
        class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-blue-500 focus:border-blue-500 bg-gray-50"
        readonly
      />
    </div>
  `,
};

const ProfileStatCard = {
  name: "ProfileStatCard",
  props: {
    icon: {
      type: String,
      required: true,
    },
    iconClass: {
      type: String,
      default: "",
    },
    label: {
      type: String,
      required: true,
    },
    value: {
      type: [String, Number],
      default: "",
    },
    valueClass: {
      type: String,
      default: "",
    },
  },
  template: `
    <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
      <div class="flex items-center gap-3">
        <iconify-icon :icon="icon" width="20" height="20" :class="iconClass"></iconify-icon>
        <span class="font-medium">{{ label }}</span>
      </div>
      <span class="font-bold" :class="valueClass">{{ value }}</span>
    </div>
  `,
};

async function fetchProfileData() {
  loading.value = true;
  error.value = "";
  try {
    const userId = currentUser.value?.id || DEFAULT_USER_ID;
    const [profileRes, statsRes, skillsRes] = await Promise.all([
      getUserProfile(userId),
      getUserStudyStats(userId),
      getUserSkills(userId),
    ]);

    Object.assign(profile, profileRes.data || {});
    Object.assign(studyStats, statsRes.data || {});
    Object.assign(skills, skillsRes.data || { primary: [], secondary: [] });
  } catch (err) {
    console.error("加载个人资料失败:", err);
    error.value = err?.message || "加载个人资料失败，请稍后重试。";
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  fetchProfileData();
});

async function handleLogout() {
  if (loggingOut.value) return;
  loggingOut.value = true;
  logoutError.value = "";
  try {
    await logoutApi();
  } catch (err) {
    console.error("退出登录失败:", err);
    logoutError.value = err?.message || "退出登录失败，请稍后重试。";
  } finally {
    clearAuth();
    clearCurrentUser();
    loggingOut.value = false;
    router.push({
      name: "Login",
      query: { redirect: router.currentRoute.value.fullPath },
    });
  }
}

defineOptions({
  name: "Profile",
});
</script>
