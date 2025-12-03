<template>
  <div class="relative" ref="containerRef">
    <!-- Trigger Button -->
    <button
      @click="toggle"
      class="relative p-2 text-gray-600 hover:bg-gray-100 rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500/20"
    >
      <iconify-icon
        icon="mdi:bell-outline"
        width="24"
        height="24"
      ></iconify-icon>
      <span
        v-if="unreadCount > 0"
        class="absolute top-1 right-1 w-2.5 h-2.5 bg-red-500 rounded-full border border-white animate-pulse"
      ></span>
    </button>

    <!-- Dropdown Panel -->
    <div
      v-if="isOpen"
      class="absolute right-0 mt-2 w-96 bg-white rounded-xl shadow-2xl border border-gray-100 z-50 overflow-hidden origin-top-right transition-all animate-in fade-in zoom-in-95 duration-200"
    >
      <!-- Header -->
      <div
        class="p-4 border-b bg-gray-50/50 backdrop-blur-sm flex justify-between items-center sticky top-0"
      >
        <h3 class="font-bold text-gray-800 flex items-center gap-2">
          通知中心
          <span
            v-if="unreadCount > 0"
            class="bg-red-100 text-red-600 text-xs px-2 py-0.5 rounded-full font-medium"
            >{{ unreadCount }}</span
          >
        </h3>
        <div class="flex items-center gap-3">
          <button
            v-if="unreadCount > 0"
            @click="markAllRead"
            class="text-xs text-blue-600 hover:text-blue-700 hover:underline font-medium"
          >
            全部已读
          </button>
          <button
            @click="isOpen = false"
            class="text-gray-400 hover:text-gray-600 transition-colors"
          >
            <iconify-icon icon="mdi:close" width="20"></iconify-icon>
          </button>
        </div>
      </div>

      <!-- Tabs -->
      <div class="flex border-b px-4 bg-white">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          @click="currentTab = tab.id"
          :class="[
            'flex-1 py-3 text-sm font-medium border-b-2 transition-colors relative',
            currentTab === tab.id
              ? 'border-blue-600 text-blue-600'
              : 'border-transparent text-gray-500 hover:text-gray-700',
          ]"
        >
          {{ tab.label }}
        </button>
      </div>

      <!-- List -->
      <div class="max-h-[400px] overflow-y-auto custom-scrollbar bg-gray-50/30">
        <div
          v-if="filteredNotifications.length === 0"
          class="py-12 text-center text-gray-500 text-sm flex flex-col items-center"
        >
          <div
            class="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mb-3 text-gray-400"
          >
            <iconify-icon icon="mdi:bell-off-outline" width="32"></iconify-icon>
          </div>
          <p>暂无{{ currentTabLabel }}通知</p>
        </div>

        <div v-else class="divide-y divide-gray-100">
          <div
            v-for="note in filteredNotifications"
            :key="note.id"
            @click="openDetail(note)"
            :class="[
              'p-4 hover:bg-white transition-colors cursor-pointer relative group',
              !note.is_read ? 'bg-blue-50/40 hover:bg-blue-50/60' : 'bg-white',
            ]"
          >
            <div class="flex gap-3">
              <!-- Icon based on type -->
              <div class="mt-1 flex-shrink-0">
                <div
                  :class="[
                    'w-8 h-8 rounded-full flex items-center justify-center',
                    getTypeStyle(note.type),
                  ]"
                >
                  <iconify-icon
                    :icon="getTypeIcon(note.type)"
                    width="16"
                  ></iconify-icon>
                </div>
              </div>

              <div class="flex-1 min-w-0">
                <div class="flex justify-between items-start mb-1">
                  <h4
                    :class="[
                      'text-sm font-semibold truncate pr-2',
                      !note.is_read ? 'text-gray-900' : 'text-gray-600',
                    ]"
                  >
                    {{ note.title }}
                  </h4>
                  <span
                    class="text-xs text-gray-400 whitespace-nowrap flex-shrink-0"
                    >{{ formatTime(note.created_at) }}</span
                  >
                </div>
                <p
                  class="text-xs text-gray-500 line-clamp-2 mb-2 leading-relaxed"
                >
                  {{ note.content }}
                </p>

                <!-- Quick Actions -->
                <div
                  v-if="canQuickAction(note)"
                  class="flex gap-2 mt-2"
                  @click.stop
                >
                  <button
                    @click="handleAction(note, 'REJECT')"
                    class="px-3 py-1 text-xs border border-gray-200 bg-white rounded hover:bg-gray-50 text-gray-600 transition-colors"
                  >
                    拒绝
                  </button>
                  <button
                    @click="handleAction(note, 'ACCEPT')"
                    class="px-3 py-1 text-xs bg-blue-600 text-white rounded hover:bg-blue-700 shadow-sm transition-colors"
                  >
                    {{ getAcceptText(note.type) }}
                  </button>
                </div>
              </div>

              <!-- Unread Dot -->
              <div
                v-if="!note.is_read"
                class="absolute top-4 right-4 w-2 h-2 bg-blue-500 rounded-full"
              ></div>
            </div>
          </div>
        </div>
      </div>

      <!-- Footer -->
      <div class="p-2 border-t bg-gray-50 text-center">
        <button
          @click="viewHistory"
          class="text-xs text-gray-500 hover:text-blue-600 py-1 w-full transition-colors"
        >
          查看历史记录
        </button>
      </div>
    </div>

    <!-- Detail Modal -->
    <div
      v-if="showModal"
      class="fixed inset-0 z-[100] flex items-center justify-center p-4"
    >
      <div
        class="absolute inset-0 bg-black/40 backdrop-blur-sm transition-opacity"
        @click="closeModal"
      ></div>
      <div
        class="bg-white rounded-xl shadow-2xl w-full max-w-lg relative z-10 overflow-hidden animate-in fade-in zoom-in-95 duration-200"
      >
        <div class="p-6">
          <div class="flex items-start justify-between mb-4">
            <div class="flex items-center gap-3">
              <div
                :class="[
                  'w-10 h-10 rounded-full flex items-center justify-center',
                  getTypeStyle(selectedNote?.type),
                ]"
              >
                <iconify-icon
                  :icon="getTypeIcon(selectedNote?.type)"
                  width="20"
                ></iconify-icon>
              </div>
              <div>
                <h3 class="text-lg font-bold text-gray-900">
                  {{ selectedNote?.title }}
                </h3>
                <span class="text-xs text-gray-500">{{
                  formatFullTime(selectedNote?.created_at)
                }}</span>
              </div>
            </div>
            <button
              @click="closeModal"
              class="text-gray-400 hover:text-gray-600 transition-colors"
            >
              <iconify-icon icon="mdi:close" width="24"></iconify-icon>
            </button>
          </div>

          <div
            class="bg-gray-50 p-4 rounded-lg text-gray-700 text-sm leading-relaxed mb-6 whitespace-pre-wrap border border-gray-100"
          >
            {{ selectedNote?.content }}
          </div>

          <div
            v-if="canQuickAction(selectedNote)"
            class="flex justify-end gap-3 border-t pt-4"
          >
            <button
              @click="handleAction(selectedNote, 'REJECT')"
              class="px-4 py-2 text-sm border border-gray-300 rounded-lg hover:bg-gray-50 text-gray-700 transition-colors"
            >
              拒绝
            </button>
            <button
              @click="handleAction(selectedNote, 'ACCEPT')"
              class="px-4 py-2 text-sm bg-blue-600 text-white rounded-lg hover:bg-blue-700 shadow-md transition-colors"
            >
              {{ getAcceptText(selectedNote?.type) }}
            </button>
          </div>
          <div
            v-else-if="getActionStatusText(selectedNote)"
            class="flex justify-end border-t pt-4 text-sm font-medium text-gray-500"
          >
            {{ getActionStatusText(selectedNote) }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
import { getNotificationList, getUnreadNotificationCount, markAllNotificationsAsRead, markNotificationAsRead } from '@/api/modules/notification';
import { handleTeamRequest } from '@/api/modules/team';

const router = useRouter();
const containerRef = ref(null);
const isOpen = ref(false);
const showModal = ref(false);
const notifications = ref([]);
const unreadCount = ref(0);
const selectedNote = ref(null);
const currentTab = ref("ALL");
let pollTimer = null;

const tabs = [
  { id: "ALL", label: "全部" },
  { id: "SYSTEM", label: "系统" },
  { id: "TEAM", label: "团队" },
];

const currentTabLabel = computed(() => {
  const tab = tabs.find((t) => t.id === currentTab.value);
  return tab ? tab.label : "";
});

const filteredNotifications = computed(() => {
  if (currentTab.value === "ALL") return notifications.value;
  if (currentTab.value === "TEAM") {
    return notifications.value.filter((n) =>
      ["TEAM_INVITE", "TEAM_APPLICATION", "TEAM"].includes(n.type)
    );
  }
  return notifications.value.filter(
    (n) =>
      n.type === "SYSTEM" ||
      !["TEAM_INVITE", "TEAM_APPLICATION", "TEAM"].includes(n.type)
  );
});

const toggle = () => {
  isOpen.value = !isOpen.value;
  if (isOpen.value) {
    fetchNotifications();
  }
};

const fetchNotifications = async () => {
  try {
    const res = await getNotificationList({ page: 1, page_size: 20 });
    // Support new pagination structure
    if (res.data && Array.isArray(res.data.items)) {
      notifications.value = res.data.items;
    } else {
      notifications.value = res.data || [];
    }

    // We still want unread count, which is fetched separately or we can infer if needed,
    // but fetchUnreadCount does it separately.
  } catch (e) {
    console.error(e);
  }
};

const fetchUnreadCount = async () => {
  try {
    const res = await getUnreadNotificationCount();
    // Backend returns {code: 0, data: count}
    const count = res.data?.data ?? res.data ?? 0;
    unreadCount.value = typeof count === "number" ? count : 0;
  } catch (e) {
    console.error(e);
  }
};

const markAllRead = async () => {
  try {
    await markAllNotificationsAsRead();
    // Optimistic update
    notifications.value.forEach((n) => (n.is_read = true));
    unreadCount.value = 0;
    fetchNotifications(); // Sync with server
  } catch (e) {
    console.error(e);
  }
};

const markAsRead = async (note) => {
  if (note.is_read) return;
  try {
    await markNotificationAsRead(note.id);
    note.is_read = true;
    unreadCount.value = Math.max(0, unreadCount.value - 1);
  } catch (e) {
    console.error(e);
  }
};

const openDetail = (note) => {
  selectedNote.value = note;
  showModal.value = true;
  isOpen.value = false; // Close dropdown
  markAsRead(note);
};

const closeModal = () => {
  showModal.value = false;
  selectedNote.value = null;
};

const canQuickAction = (note) => {
  if (!note) return false;
  return (
    (note.type === "TEAM_INVITE" || note.type === "TEAM_APPLICATION") &&
    (!note.action_status || note.action_status === "PENDING")
  );
};

const getActionStatusText = (note) => {
  if (!note || !note.action_status || note.action_status === 'PENDING' || note.action_status === 'PROCESSED') return '';
  if (note.action_status === 'REJECTED') return '已拒绝';
  if (note.action_status === 'APPROVED' || note.action_status === 'ACCEPTED') return '已同意';
  return '';
};

const getAcceptText = (type) => {
  return type === "TEAM_APPLICATION" ? "批准申请" : "接受邀请";
};

const getTypeIcon = (type) => {
  switch (type) {
    case "TEAM_INVITE":
    case "TEAM_APPLICATION":
      return "mdi:account-group";
    case "SYSTEM":
    default:
      return "mdi:bell-ring";
  }
};

const getTypeStyle = (type) => {
  switch (type) {
    case "TEAM_INVITE":
    case "TEAM_APPLICATION":
      return "bg-purple-100 text-purple-600";
    case "SYSTEM":
    default:
      return "bg-blue-100 text-blue-600";
  }
};

const handleAction = async (note, action) => {
  try {
    // Use 0 as dummy teamId since backend finds request by ID (RelatedID)
    await handleTeamRequest(0, note.related_id, { action });

    // Show success message (could use a toast library if available)
    const msg =
      action === "ACCEPT"
        ? note.type === "TEAM_APPLICATION"
          ? "已批准申请"
          : "已加入团队"
        : "已拒绝";
    ElMessage.success(msg);

    // Update list
    fetchNotifications();
    if (showModal.value) closeModal();
  } catch (e) {
    const errorMsg = e.response?.data?.error || "操作失败";
    ElMessage.error(errorMsg);
    console.error(e);
  }
};

const formatTime = (str) => {
  if (!str) return "";
  const d = new Date(str);
  const now = new Date();
  const diff = now - d;

  if (diff < 60000) return "刚刚";
  if (diff < 3600000) return `${Math.floor(diff / 60000)}分钟前`;
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}小时前`;
  return `${d.getMonth() + 1}-${d.getDate()}`;
};

const formatFullTime = (str) => {
  if (!str) return "";
  const d = new Date(str);
  return `${d.getFullYear()}-${(d.getMonth() + 1)
    .toString()
    .padStart(2, "0")}-${d.getDate().toString().padStart(2, "0")} ${d
    .getHours()
    .toString()
    .padStart(2, "0")}:${d.getMinutes().toString().padStart(2, "0")}`;
};

const viewHistory = () => {
  isOpen.value = false;
  router.push('/notifications');
};

const handleClickOutside = (event) => {
  if (containerRef.value && !containerRef.value.contains(event.target)) {
    isOpen.value = false;
  }
};

onMounted(() => {
  fetchUnreadCount();
  pollTimer = setInterval(fetchUnreadCount, 10000); // Poll every 10s for better "real-time" feel
  document.addEventListener("click", handleClickOutside);
});

onUnmounted(() => {
  if (pollTimer) clearInterval(pollTimer);
  document.removeEventListener("click", handleClickOutside);
});
</script>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #d1d5db;
  border-radius: 3px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #9ca3af;
}
</style>
