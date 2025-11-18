<template>
  <aside
    class="w-80 flex-shrink-0 sidebar-shadow rounded-2xl bg-white p-6 flex flex-col"
    :class="{ 'h-full': isChatExpanded }"
  >
    <div class="flex justify-between items-center mb-5">
      <h3 class="text-xl font-bold text-gray-700">学习互动</h3>
      <div class="flex items-center gap-2">
        <div class="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
        <span class="text-sm text-green-600">3人在线</span>
      </div>
    </div>

    <div class="mb-6">
      <div class="flex justify-between items-center mb-3">
        <h4 class="font-medium text-gray-700">学习伙伴</h4>
        <button class="text-blue-600 hover:text-blue-800 flex items-center text-sm">
          <iconify-icon icon="mdi:plus-circle" class="mr-1" />
          添加
        </button>
      </div>

      <div class="space-y-3">
        <div class="friend-card card p-3 flex items-center cursor-pointer hover:bg-blue-50 rounded-xl">
          <div class="relative">
            <div class="w-10 h-10 rounded-full bg-gray-300"></div>
            <div class="absolute bottom-0 right-0 w-3 h-3 rounded-full bg-green-500 border-2 border-white"></div>
          </div>
          <div class="ml-3 flex-1">
            <div class="font-medium text-sm">张伟</div>
            <div class="text-xs text-gray-500">前端开发课程学习中</div>
          </div>
          <div class="flex gap-1">
            <button class="w-7 h-7 flex items-center justify-center bg-gray-100 text-gray-600 rounded-full hover:bg-blue-100 text-sm">
              <iconify-icon icon="mdi:thumb-up" />
            </button>
            <button class="w-7 h-7 flex items-center justify-center bg-blue-100 text-blue-600 rounded-full hover:bg-blue-200 text-sm">
              <iconify-icon icon="mdi:message" />
            </button>
          </div>
        </div>

        <div class="friend-card card p-3 flex items-center cursor-pointer hover:bg-blue-50 rounded-xl">
          <div class="relative">
            <div class="w-10 h-10 rounded-full bg-gray-300"></div>
            <div class="absolute bottom-0 right-0 w-3 h-3 rounded-full bg-gray-300 border-2 border-white"></div>
          </div>
          <div class="ml-3 flex-1">
            <div class="font-medium text-sm">刘燕</div>
            <div class="text-xs text-gray-500">英语学习进行中</div>
          </div>
          <div class="flex gap-1">
            <button class="w-7 h-7 flex items-center justify-center bg-gray-100 text-gray-600 rounded-full hover:bg-blue-100 text-sm">
              <iconify-icon icon="mdi:thumb-up" />
            </button>
            <button class="w-7 h-7 flex items-center justify-center bg-blue-100 text-blue-600 rounded-full hover:bg-blue-200 text-sm">
              <iconify-icon icon="mdi:message" />
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="mb-4">
      <div class="flex justify-between items-center mb-3">
        <h4 class="font-medium text-gray-700">评论提醒</h4>
        <div class="w-5 h-5 bg-red-500 text-white rounded-full flex items-center justify-center text-xs">
          3
        </div>
      </div>
      <div class="space-y-2">
        <div class="p-2 bg-red-50 border-l-4 border-red-400 rounded">
          <div class="text-sm font-medium text-red-800">新评论</div>
          <div class="text-xs text-red-600">王浩评论了你的Vue项目</div>
        </div>
      </div>
    </div>

    <div class="mb-4">
      <div class="flex justify-between items-center mb-3">
        <h4 class="font-medium text-gray-700 flex items-center">
          <iconify-icon icon="mdi:chat" class="mr-2 text-blue-600" />
          互动消息
          <span class="ml-2 bg-red-500 text-white text-xs rounded-full px-2 py-0.5">3</span>
        </h4>
        <button
          @click="$emit('toggle-chat')"
          class="text-blue-600 hover:text-blue-800 text-xs cursor-pointer flex items-center"
        >
          <span>{{ isChatExpanded ? "收起" : "展开" }}</span>
          <iconify-icon
            :icon="isChatExpanded ? 'mdi:chevron-up' : 'mdi:chevron-down'"
            class="ml-1"
          />
        </button>
      </div>

      <div v-if="isChatExpanded" class="mt-6 flex-1 flex flex-col pb-6">
        <div class="flex-1 space-y-3 overflow-y-auto">
          <div
            class="flex items-center p-3 bg-gray-50 rounded-lg hover:bg-gray-100 cursor-pointer transition-colors"
            @click="$emit('open-chat', '陈敏')"
          >
            <div class="relative mr-3">
              <div class="w-10 h-10 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 flex items-center justify-center">
                <span class="text-white font-medium text-sm">陈</span>
              </div>
              <div class="absolute -top-1 -right-1 w-3 h-3 bg-green-500 rounded-full border-2 border-white"></div>
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex justify-between items-center mb-1">
                <h5 class="font-medium text-gray-800 text-sm truncate">陈敏</h5>
                <span class="text-xs text-gray-500">1小时前</span>
              </div>
              <p class="text-xs text-gray-600 truncate">
                你的前端项目太棒了！😄 有时间一起讨论一下吗？
              </p>
            </div>
            <div class="ml-2">
              <div class="w-2 h-2 bg-red-500 rounded-full"></div>
            </div>
          </div>

          <div
            class="flex items-center p-3 bg-gray-50 rounded-lg hover:bg-gray-100 cursor-pointer transition-colors"
            @click="$emit('open-chat', '王小明')"
          >
            <div class="relative mr-3">
              <div class="w-10 h-10 rounded-full bg-gradient-to-br from-green-400 to-green-600 flex items-center justify-center">
                <span class="text-white font-medium text-sm">王</span>
              </div>
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex justify-between items-center mb-1">
                <h5 class="font-medium text-gray-800 text-sm truncate">王小明</h5>
                <span class="text-xs text-gray-500">3小时前</span>
              </div>
              <p class="text-xs text-gray-600 truncate">
                学习小组今晚8点开始，记得参加哦～
              </p>
            </div>
            <div class="ml-2">
              <div class="w-2 h-2 bg-red-500 rounded-full"></div>
            </div>
          </div>

          <div
            class="flex items-center p-3 bg-gray-50 rounded-lg hover:bg-gray-100 cursor-pointer transition-colors"
            @click="$emit('open-chat', '李老师')"
          >
            <div class="relative mr-3">
              <div class="w-10 h-10 rounded-full bg-gradient-to-br from-purple-400 to-purple-600 flex items-center justify-center">
                <span class="text-white font-medium text-sm">李</span>
              </div>
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex justify-between items-center mb-1">
                <h5 class="font-medium text-gray-800 text-sm truncate">李老师</h5>
                <span class="text-xs text-gray-500">昨天</span>
              </div>
              <p class="text-xs text-gray-600 truncate">
                作业完成得很好，继续保持！💪
              </p>
            </div>
            <div class="ml-2">
              <div class="w-2 h-2 bg-red-500 rounded-full"></div>
            </div>
          </div>

          <div
            class="flex items-center p-3 bg-gray-50 rounded-lg hover:bg-gray-100 cursor-pointer transition-colors"
            @click="$emit('open-chat', '刘晓雨')"
          >
            <div class="relative mr-3">
              <div class="w-10 h-10 rounded-full bg-gradient-to-br from-orange-400 to-orange-600 flex items-center justify-center">
                <span class="text-white font-medium text-sm">刘</span>
              </div>
              <div class="absolute -top-1 -right-1 w-3 h-3 bg-green-500 rounded-full border-2 border-white"></div>
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex justify-between items-center mb-1">
                <h5 class="font-medium text-gray-800 text-sm truncate">刘晓雨</h5>
                <span class="text-xs text-gray-500">2天前</span>
              </div>
              <p class="text-xs text-gray-600 truncate">
                Vue3的新特性真的很棒，一起学习吧！🚀
              </p>
            </div>
            <div class="ml-2">
              <div class="w-2 h-2 bg-red-500 rounded-full"></div>
            </div>
          </div>

          <div
            class="flex items-center p-3 bg-gray-50 rounded-lg hover:bg-gray-100 cursor-pointer transition-colors"
            @click="$emit('open-chat', '张小花')"
          >
            <div class="relative mr-3">
              <div class="w-10 h-10 rounded-full bg-gradient-to-br from-pink-400 to-pink-600 flex items-center justify-center">
                <span class="text-white font-medium text-sm">张</span>
              </div>
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex justify-between items-center mb-1">
                <h5 class="font-medium text-gray-800 text-sm truncate">张小花</h5>
                <span class="text-xs text-gray-500">3天前</span>
              </div>
              <p class="text-xs text-gray-600 truncate">
                明天的算法课记得带笔记本哦～📚
              </p>
            </div>
          </div>
        </div>

        <div class="mt-4 p-3 bg-blue-50 rounded-lg">
          <div class="flex items-center">
            <input
              type="text"
              placeholder="快速回复..."
              class="flex-1 bg-white border border-blue-200 focus:ring-2 focus:ring-blue-500 focus:border-transparent rounded-full py-2 px-4 text-sm"
            />
            <button
              class="w-8 h-8 flex items-center justify-center bg-blue-600 text-white rounded-full ml-3 hover:bg-blue-700 transition-colors"
            >
              <iconify-icon icon="mdi:send" class="text-sm" />
            </button>
          </div>
        </div>
      </div>

      <div v-else class="space-y-2">
        <div
          class="flex items-center p-2 bg-gray-50 rounded-lg hover:bg-gray-100 cursor-pointer transition-colors"
          @click="$emit('open-chat', '陈敏')"
        >
          <div class="relative mr-2">
            <div class="w-8 h-8 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 flex items-center justify-center">
              <span class="text-white font-medium text-xs">陈</span>
            </div>
            <div class="absolute -top-0.5 -right-0.5 w-2.5 h-2.5 bg-green-500 rounded-full border border-white"></div>
          </div>
          <div class="flex-1 min-w-0">
            <div class="flex justify-between items-center mb-0.5">
              <h5 class="font-medium text-gray-800 text-xs truncate">陈敏</h5>
              <span class="text-xs text-gray-500">1小时前</span>
            </div>
            <p class="text-xs text-gray-600 truncate">
              你的前端项目太棒了！😄 有时...
            </p>
          </div>
          <div class="ml-1">
            <div class="w-1.5 h-1.5 bg-red-500 rounded-full"></div>
          </div>
        </div>

        <div
          class="flex items-center p-2 bg-gray-50 rounded-lg hover:bg-gray-100 cursor-pointer transition-colors"
          @click="$emit('open-chat', '王小明')"
        >
          <div class="relative mr-2">
            <div class="w-8 h-8 rounded-full bg-gradient-to-br from-green-400 to-green-600 flex items-center justify-center">
              <span class="text-white font-medium text-xs">王</span>
            </div>
          </div>
          <div class="flex-1 min-w-0">
            <div class="flex justify-between items-center mb-0.5">
              <h5 class="font-medium text-gray-800 text-xs truncate">王小明</h5>
              <span class="text-xs text-gray-500">3小时前</span>
            </div>
            <p class="text-xs text-gray-600 truncate">
              学习小组今晚8点开始，记得参...
            </p>
          </div>
          <div class="ml-1">
            <div class="w-1.5 h-1.5 bg-red-500 rounded-full"></div>
          </div>
        </div>

        <div
          class="flex items-center p-2 bg-gray-50 rounded-lg hover:bg-gray-100 cursor-pointer transition-colors"
          @click="$emit('open-chat', '李老师')"
        >
          <div class="relative mr-2">
            <div class="w-8 h-8 rounded-full bg-gradient-to-br from-purple-400 to-purple-600 flex items-center justify-center">
              <span class="text-white font-medium text-xs">李</span>
            </div>
          </div>
          <div class="flex-1 min-w-0">
            <div class="flex justify-between items-center mb-0.5">
              <h5 class="font-medium text-gray-800 text-xs truncate">李老师</h5>
              <span class="text-xs text-gray-500">昨天</span>
            </div>
            <p class="text-xs text-gray-600 truncate">
              作业完成得很好，继续保持！💪
            </p>
          </div>
          <div class="ml-1">
            <div class="w-1.5 h-1.5 bg-red-500 rounded-full"></div>
          </div>
        </div>

        <div
          class="flex items-center p-2 bg-gray-50 rounded-lg hover:bg-gray-100 cursor-pointer transition-colors"
          @click="$emit('open-chat', '刘晓雨')"
        >
          <div class="relative mr-2">
            <div class="w-8 h-8 rounded-full bg-gradient-to-br from-orange-400 to-orange-600 flex items-center justify-center">
              <span class="text-white font-medium text-xs">刘</span>
            </div>
            <div class="absolute -top-0.5 -right-0.5 w-2.5 h-2.5 bg-green-500 rounded-full border border-white"></div>
          </div>
          <div class="flex-1 min-w-0">
            <div class="flex justify-between items-center mb-0.5">
              <h5 class="font-medium text-gray-800 text-xs truncate">刘晓雨</h5>
              <span class="text-xs text-gray-500">2天前</span>
            </div>
            <p class="text-xs text-gray-600 truncate">
              Vue3的新特性真的很棒，一起...
            </p>
          </div>
          <div class="ml-1">
            <div class="w-1.5 h-1.5 bg-red-500 rounded-full"></div>
          </div>
        </div>

        <div
          class="flex items-center p-2 bg-gray-50 rounded-lg hover:bg-gray-100 cursor-pointer transition-colors"
          @click="$emit('open-chat', '张小花')"
        >
          <div class="relative mr-2">
            <div class="w-8 h-8 rounded-full bg-gradient-to-br from-pink-400 to-pink-600 flex items-center justify-center">
              <span class="text-white font-medium text-xs">张</span>
            </div>
          </div>
          <div class="flex-1 min-w-0">
            <div class="flex justify-between items-center mb-0.5">
              <h5 class="font-medium text-gray-800 text-xs truncate">张小花</h5>
              <span class="text-xs text-gray-500">3天前</span>
            </div>
            <p class="text-xs text-gray-600 truncate">
              明天的算法课记得带笔记本哦～📚
            </p>
          </div>
        </div>
      </div>
    </div>

    <div
      v-if="!isChatExpanded"
      class="mt-6 flex-1 flex items-end justify-center pb-4"
    >
      <div class="text-center">
        <p
          class="text-2xl font-bold bg-gradient-to-r from-gray-600 to-gray-800 bg-clip-text text-transparent leading-relaxed font-serif italic"
        >
          {{ currentMotivationalQuote }}
        </p>
      </div>
    </div>
  </aside>
</template>

<script>
  export default {
    name: "InteractionPanel",
    props: {
      isChatExpanded: {
        type: Boolean,
        required: true,
      },
      currentMotivationalQuote: {
        type: String,
        required: true,
      },
    },
    emits: ["toggle-chat", "open-chat"],
  };
</script>

<style scoped>
  .sidebar-shadow {
    box-shadow: 0px 0px 20px rgba(0, 0, 0, 0.05);
  }

  .friend-card:hover {
    background-color: rgba(45, 91, 255, 0.05);
  }
</style>
