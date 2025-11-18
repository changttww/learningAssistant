<template>
  <div class="grid grid-cols-3 gap-4 mb-6">
    <div
      class="card p-5 cursor-pointer hover:shadow-lg transition-shadow"
      @click="$emit('show-efficiency')"
    >
      <div class="flex items-center mb-4">
        <div
          class="w-12 h-12 rounded-xl bg-gradient-to-br from-purple-500 to-pink-500 flex items-center justify-center mr-3"
        >
          <iconify-icon icon="mdi:chart-line" class="text-2xl text-white" />
        </div>
        <div>
          <h4 class="font-bold text-gray-800">学习效率分析</h4>
          <p class="text-sm text-gray-500">智能分析学习数据</p>
        </div>
      </div>
      <div class="space-y-2">
        <div class="flex justify-between text-sm">
          <span class="text-gray-600">本周学习时长</span>
          <span class="font-medium text-purple-600">28.5小时</span>
        </div>
        <div class="flex justify-between text-sm">
          <span class="text-gray-600">专注度评分</span>
          <span class="font-medium text-green-600">85分</span>
        </div>
        <div class="flex justify-between text-sm">
          <span class="text-gray-600">任务完成率</span>
          <span class="font-medium text-blue-600">92%</span>
        </div>
      </div>
    </div>

    <div
      class="card p-5 cursor-pointer hover:shadow-lg transition-shadow"
      @click="$emit('show-summary')"
    >
      <div class="flex items-center mb-4">
        <div
          class="w-12 h-12 rounded-xl bg-gradient-to-br from-blue-500 to-cyan-500 flex items-center justify-center mr-3"
        >
          <iconify-icon icon="mdi:brain" class="text-2xl text-white" />
        </div>
        <div>
          <h4 class="font-bold text-gray-800">智能总结复习</h4>
          <p class="text-sm text-gray-500">AI生成复习提纲</p>
        </div>
      </div>
      <div class="space-y-2">
        <div class="flex justify-between text-sm">
          <span class="text-gray-600">待复习内容</span>
          <span class="font-medium text-orange-600">5项</span>
        </div>
        <div class="flex justify-between text-sm">
          <span class="text-gray-600">复习提醒</span>
          <span class="font-medium text-red-600">3条</span>
        </div>
        <div class="flex justify-between text-sm">
          <span class="text-gray-600">知识点掌握</span>
          <span class="font-medium text-green-600">78%</span>
        </div>
      </div>
    </div>

    <div
      class="card p-5 cursor-pointer hover:shadow-lg transition-shadow"
      @click="$emit('show-check-in')"
    >
      <div class="flex items-center mb-4">
        <div
          class="w-12 h-12 rounded-xl bg-gradient-to-br from-green-500 to-teal-500 flex items-center justify-center mr-3"
        >
          <iconify-icon
            icon="mdi:calendar-check"
            class="text-2xl text-white"
          />
        </div>
        <div>
          <h4 class="font-bold text-gray-800">打卡记录分析</h4>
          <p class="text-sm text-gray-500">学习习惯与建议</p>
        </div>
      </div>
      <div class="space-y-2">
        <div class="flex justify-between text-sm">
          <span class="text-gray-600">连续打卡</span>
          <span class="font-medium text-green-600">28天</span>
        </div>
        <div class="flex justify-between text-sm">
          <span class="text-gray-600">平均打字速度</span>
          <span class="font-medium text-blue-600">65WPM</span>
        </div>
        <div class="flex justify-between text-sm">
          <span class="text-gray-600">建议状态</span>
          <span class="font-medium text-purple-600">继续保持</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  export default {
    name: "AnalysisEntryGrid",
    emits: ["show-efficiency", "show-summary", "show-check-in"],
  };
</script>
