<template>
  <div class="p-6 max-w-4xl mx-auto">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-gray-800">通知历史</h1>
      <div class="flex gap-3">
        <button @click="markAllRead" class="text-sm text-blue-600 hover:underline">全部已读</button>
        <button @click="clearAll" class="text-sm text-red-600 hover:underline">清空通知</button>
      </div>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
      <div v-if="loading" class="p-8 text-center text-gray-500">加载中...</div>
      <div v-else-if="notifications.length === 0" class="p-12 text-center text-gray-500">
        <iconify-icon icon="mdi:bell-off-outline" width="48" class="mx-auto mb-4 text-gray-300"></iconify-icon>
        <p>暂无通知记录</p>
      </div>
      <div v-else class="divide-y divide-gray-100">
        <div 
          v-for="note in notifications" 
          :key="note.id" 
          :class="['p-4 hover:bg-gray-50 transition-colors', !note.is_read ? 'bg-blue-50/30' : '']"
        >
          <div class="flex gap-4">
            <div class="mt-1">
               <div :class="['w-10 h-10 rounded-full flex items-center justify-center', getTypeStyle(note.type)]">
                  <iconify-icon :icon="getTypeIcon(note.type)" width="20"></iconify-icon>
               </div>
            </div>
            <div class="flex-1">
              <div class="flex justify-between items-start mb-1">
                <h3 :class="['text-sm font-semibold', !note.is_read ? 'text-gray-900' : 'text-gray-600']">{{ note.title }}</h3>
                <span class="text-xs text-gray-400">{{ formatTime(note.created_at) }}</span>
              </div>
              <p class="text-sm text-gray-600 mb-2 whitespace-pre-wrap">{{ note.content }}</p>
              
               <!-- Actions if applicable -->
              <div v-if="canQuickAction(note)" class="flex gap-2 mt-2">
                 <button 
                    @click="handleAction(note, 'REJECT')" 
                    class="px-3 py-1 text-xs border border-gray-200 bg-white rounded hover:bg-gray-50 text-gray-600"
                  >
                    拒绝
                  </button>
                  <button 
                    @click="handleAction(note, 'ACCEPT')" 
                    class="px-3 py-1 text-xs bg-blue-600 text-white rounded hover:bg-blue-700"
                  >
                    {{ getAcceptText(note.type) }}
                  </button>
              </div>
              <div v-else-if="getActionStatusText(note)" class="mt-2 text-sm text-gray-500 font-medium">
                 {{ getActionStatusText(note) }}
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Pagination -->
      <div class="p-4 border-t bg-gray-50 flex justify-between items-center" v-if="total > pageSize">
        <button 
          :disabled="page <= 1"
          @click="changePage(page - 1)"
          class="px-3 py-1 text-sm bg-white border rounded hover:bg-gray-100 disabled:opacity-50"
        >
          上一页
        </button>
        <span class="text-sm text-gray-600">第 {{ page }} 页 / 共 {{ Math.ceil(total / pageSize) }} 页</span>
        <button 
          :disabled="page * pageSize >= total"
          @click="changePage(page + 1)"
          class="px-3 py-1 text-sm bg-white border rounded hover:bg-gray-100 disabled:opacity-50"
        >
          下一页
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { ElMessage } from 'element-plus';
import { getNotificationList, markAllNotificationsAsRead, clearAllNotifications } from '@/api/modules/notification';
import { handleTeamRequest } from '@/api/modules/team';

const notifications = ref([]);
const loading = ref(false);
const page = ref(1);
const pageSize = ref(20);
const total = ref(0);

const fetchList = async () => {
  loading.value = true;
  try {
    const res = await getNotificationList({ page: page.value, page_size: pageSize.value });
    if (res.data && Array.isArray(res.data.items)) {
      notifications.value = res.data.items;
      total.value = res.data.total;
    } else {
      // Fallback
      notifications.value = res.data || [];
      total.value = notifications.value.length;
    }
  } catch (e) {
    console.error(e);
  } finally {
    loading.value = false;
  }
};

const changePage = (newPage) => {
  page.value = newPage;
  fetchList();
  window.scrollTo(0, 0);
};

const markAllRead = async () => {
  if (!confirm("确定要标记所有通知为已读吗？")) return;
  try {
    await markAllNotificationsAsRead();
    fetchList();
  } catch (e) {
    console.error(e);
  }
};

const clearAll = async () => {
  if (!confirm("确定要清空所有通知吗？此操作不可恢复。")) return;
  try {
    await clearAllNotifications();
    fetchList();
  } catch (e) {
    console.error(e);
  }
};

const formatTime = (str) => {
  if (!str) return '';
  return new Date(str).toLocaleString();
};

// Reusing helper logic
const canQuickAction = (note) => {
  if (!note) return false;
  return (note.type === 'TEAM_INVITE' || note.type === 'TEAM_APPLICATION') && (!note.action_status || note.action_status === 'PENDING');
};

const getActionStatusText = (note) => {
  if (!note || !note.action_status || note.action_status === 'PENDING' || note.action_status === 'PROCESSED') return '';
  if (note.action_status === 'REJECTED') return '已拒绝';
  if (note.action_status === 'APPROVED' || note.action_status === 'ACCEPTED') return '已同意';
  return '';
};

const getAcceptText = (type) => {
  return type === 'TEAM_APPLICATION' ? '批准申请' : '接受邀请';
};

const getTypeIcon = (type) => {
  switch(type) {
    case 'TEAM_INVITE':
    case 'TEAM_APPLICATION':
      return 'mdi:account-group';
    case 'SYSTEM':
    default:
      return 'mdi:bell-ring';
  }
};

const getTypeStyle = (type) => {
  switch(type) {
    case 'TEAM_INVITE':
    case 'TEAM_APPLICATION':
      return 'bg-purple-100 text-purple-600';
    case 'SYSTEM':
    default:
      return 'bg-blue-100 text-blue-600';
  }
};

const handleAction = async (note, action) => {
  try {
    await handleTeamRequest(0, note.related_id, { action });
    ElMessage.success("操作成功");
    fetchList();
  } catch (e) {
    ElMessage.error(e.response?.data?.error || "操作失败");
  }
};

onMounted(() => {
  fetchList();
});
</script>
