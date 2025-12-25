<template>
  <div class="w-full h-full overflow-auto p-4 bg-gray-50">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 flex items-center gap-2">
          <iconify-icon icon="mdi:chart-box" width="28" height="28" class="text-[#2D5BFF]"></iconify-icon>
          AI 智能学习报告
        </h1>
        <p class="text-gray-500 mt-1">基于 AI 的个性化学习数据分析与建议</p>
      </div>
      <div class="flex items-center gap-3">
        <select v-model="reportDays" class="border rounded-lg px-3 py-2 text-sm bg-white">
          <option :value="7">最近 7 天</option>
          <option :value="14">最近 14 天</option>
          <option :value="30">最近 30 天</option>
        </select>
        <button 
          @click="generateReport"
          :disabled="loading"
          class="bg-[#2D5BFF] text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-opacity-90 transition-colors flex items-center gap-2 disabled:opacity-50"
        >
          <iconify-icon v-if="loading" icon="mdi:loading" class="animate-spin" width="18"></iconify-icon>
          <iconify-icon v-else icon="mdi:sparkles" width="18"></iconify-icon>
          {{ loading ? '生成中...' : '生成报告' }}
        </button>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="flex flex-col items-center justify-center py-20">
      <div class="w-16 h-16 border-4 border-blue-200 border-t-[#2D5BFF] rounded-full animate-spin mb-4"></div>
      <p class="text-gray-600">AI 正在分析您的学习数据...</p>
      <p class="text-gray-400 text-sm mt-2">这可能需要几秒钟时间</p>
    </div>

    <!-- 报告内容 -->
    <div v-else-if="report" class="space-y-6">
      <!-- 报告信息栏 -->
      <div class="bg-gradient-to-r from-blue-500 to-purple-600 rounded-2xl p-6 text-white">
        <div class="flex items-center justify-between">
          <div>
            <h2 class="text-xl font-bold">{{ report.period }} 学习报告</h2>
            <p class="text-white/80 text-sm mt-1">生成时间：{{ report.generated_at }}</p>
          </div>
          <div class="text-right">
            <div class="text-4xl font-bold">{{ report.overview?.efficiency_score || 0 }}</div>
            <div class="text-white/80 text-sm">学习效率评分</div>
            <div class="text-sm mt-1">{{ report.overview?.efficiency_level }}</div>
          </div>
        </div>
      </div>

      <!-- 学习概览卡片 -->
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <div class="flex items-center gap-2 text-gray-500 text-sm mb-2">
            <iconify-icon icon="mdi:clock-outline" width="18"></iconify-icon>
            总学习时长
          </div>
          <div class="text-2xl font-bold text-blue-600">{{ report.overview?.total_study_hours?.toFixed(1) || 0 }}h</div>
          <div class="text-xs text-gray-400 mt-1">共 {{ report.overview?.total_study_days || 0 }} 天</div>
        </div>
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <div class="flex items-center gap-2 text-gray-500 text-sm mb-2">
            <iconify-icon icon="mdi:checkbox-marked-circle" width="18"></iconify-icon>
            任务完成
          </div>
          <div class="text-2xl font-bold text-green-600">{{ report.overview?.tasks_completed || 0 }}</div>
          <div class="text-xs text-gray-400 mt-1">完成率 {{ report.overview?.task_completion_rate?.toFixed(0) || 0 }}%</div>
        </div>
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <div class="flex items-center gap-2 text-gray-500 text-sm mb-2">
            <iconify-icon icon="mdi:fire" width="18"></iconify-icon>
            连续学习
          </div>
          <div class="text-2xl font-bold text-orange-600">{{ report.overview?.streak_days || 0 }}天</div>
          <div class="text-xs text-gray-400 mt-1">保持势头！</div>
        </div>
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <div class="flex items-center gap-2 text-gray-500 text-sm mb-2">
            <iconify-icon icon="mdi:lightbulb" width="18"></iconify-icon>
            知识点
          </div>
          <div class="text-2xl font-bold text-purple-600">{{ report.overview?.knowledge_points || 0 }}</div>
          <div class="text-xs text-gray-400 mt-1">笔记 {{ report.overview?.notes_created || 0 }} 篇</div>
        </div>
      </div>

      <!-- 双列布局 -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- 能力雷达图 -->
        <div class="bg-white rounded-xl p-6 shadow-sm">
          <h3 class="font-bold text-lg mb-4 flex items-center gap-2">
            <iconify-icon icon="mdi:radar" width="20" class="text-[#2D5BFF]"></iconify-icon>
            能力雷达图
          </h3>
          <div ref="radarChartRef" class="h-64"></div>
        </div>

        <!-- 学习行为分析 -->
        <div class="bg-white rounded-xl p-6 shadow-sm">
          <h3 class="font-bold text-lg mb-4 flex items-center gap-2">
            <iconify-icon icon="mdi:chart-timeline-variant" width="20" class="text-[#2D5BFF]"></iconify-icon>
            学习行为分析
          </h3>
          <div class="space-y-3">
            <div class="flex items-center justify-between py-2 border-b">
              <span class="text-gray-600">学习高峰时段</span>
              <span class="font-medium">{{ report.behavior_analysis?.peak_study_time || '待分析' }}</span>
            </div>
            <div class="flex items-center justify-between py-2 border-b">
              <span class="text-gray-600">平均单次学习</span>
              <span class="font-medium">{{ report.behavior_analysis?.average_session_time || 0 }} 分钟</span>
            </div>
            <div class="flex items-center justify-between py-2 border-b">
              <span class="text-gray-600">最高效日期</span>
              <span class="font-medium">{{ report.behavior_analysis?.most_productive_day || '待分析' }}</span>
            </div>
            <div class="mt-3">
              <div class="text-gray-600 text-sm mb-2">学习习惯</div>
              <div class="flex flex-wrap gap-2">
                <span 
                  v-for="habit in report.behavior_analysis?.study_habits || []" 
                  :key="habit"
                  class="bg-blue-50 text-blue-700 px-2 py-1 rounded text-sm"
                >
                  {{ habit }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 学习趋势图 -->
      <div class="bg-white rounded-xl p-6 shadow-sm">
        <h3 class="font-bold text-lg mb-4 flex items-center gap-2">
          <iconify-icon icon="mdi:chart-line" width="20" class="text-[#2D5BFF]"></iconify-icon>
          学习趋势
        </h3>
        <div ref="trendChartRef" class="h-64"></div>
      </div>

      <!-- 知识掌握分析 -->
      <div class="bg-white rounded-xl p-6 shadow-sm">
        <h3 class="font-bold text-lg mb-4 flex items-center gap-2">
          <iconify-icon icon="mdi:brain" width="20" class="text-[#2D5BFF]"></iconify-icon>
          知识掌握分析
        </h3>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
          <div class="text-center p-4 bg-green-50 rounded-lg">
            <div class="text-3xl font-bold text-green-600">{{ report.knowledge_analysis?.mastered_count || 0 }}</div>
            <div class="text-sm text-gray-600">已掌握</div>
          </div>
          <div class="text-center p-4 bg-yellow-50 rounded-lg">
            <div class="text-3xl font-bold text-yellow-600">{{ report.knowledge_analysis?.learning_count || 0 }}</div>
            <div class="text-sm text-gray-600">学习中</div>
          </div>
          <div class="text-center p-4 bg-gray-50 rounded-lg">
            <div class="text-3xl font-bold text-gray-600">{{ report.knowledge_analysis?.to_learn_count || 0 }}</div>
            <div class="text-sm text-gray-600">待学习</div>
          </div>
        </div>
        <div v-if="report.knowledge_analysis?.weak_points?.length" class="mt-4">
          <div class="text-sm font-medium text-gray-700 mb-2">🎯 薄弱点分析</div>
          <div class="space-y-2">
            <div 
              v-for="(point, idx) in report.knowledge_analysis.weak_points" 
              :key="idx"
              class="flex items-center gap-2 text-sm text-gray-600 bg-red-50 p-2 rounded"
            >
              <iconify-icon icon="mdi:alert-circle" class="text-red-500"></iconify-icon>
              {{ point }}
            </div>
          </div>
        </div>
      </div>

      <!-- AI 个性化建议 -->
      <div class="bg-white rounded-xl p-6 shadow-sm">
        <h3 class="font-bold text-lg mb-4 flex items-center gap-2">
          <iconify-icon icon="mdi:robot" width="20" class="text-[#2D5BFF]"></iconify-icon>
          AI 个性化建议
        </h3>
        
        <!-- 优势分析 -->
        <div v-if="report.ai_advice?.strength_analysis?.length" class="mb-6">
          <div class="text-sm font-medium text-green-700 mb-2">✅ 你的优势</div>
          <div class="space-y-2">
            <div 
              v-for="(strength, idx) in report.ai_advice.strength_analysis" 
              :key="idx"
              class="flex items-center gap-2 text-sm text-gray-700 bg-green-50 p-3 rounded-lg"
            >
              <iconify-icon icon="mdi:check-circle" class="text-green-500"></iconify-icon>
              {{ strength }}
            </div>
          </div>
        </div>

        <!-- 待提升领域 -->
        <div v-if="report.ai_advice?.improvement_areas?.length" class="mb-6">
          <div class="text-sm font-medium text-orange-700 mb-2">📈 待提升领域</div>
          <div class="space-y-2">
            <div 
              v-for="(area, idx) in report.ai_advice.improvement_areas" 
              :key="idx"
              class="flex items-center gap-2 text-sm text-gray-700 bg-orange-50 p-3 rounded-lg"
            >
              <iconify-icon icon="mdi:trending-up" class="text-orange-500"></iconify-icon>
              {{ area }}
            </div>
          </div>
        </div>

        <!-- 个性化建议卡片 -->
        <div v-if="report.ai_advice?.personalized_tips?.length" class="mb-6">
          <div class="text-sm font-medium text-blue-700 mb-2">💡 个性化建议</div>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div 
              v-for="(tip, idx) in report.ai_advice.personalized_tips" 
              :key="idx"
              class="border rounded-lg p-4"
              :class="{
                'border-red-200 bg-red-50': tip.priority === 'high',
                'border-yellow-200 bg-yellow-50': tip.priority === 'medium',
                'border-gray-200 bg-gray-50': tip.priority === 'low'
              }"
            >
              <div class="flex items-center gap-2 font-medium mb-2">
                <span>{{ tip.icon || '💡' }}</span>
                {{ tip.title }}
                <span 
                  class="text-xs px-2 py-0.5 rounded"
                  :class="{
                    'bg-red-200 text-red-800': tip.priority === 'high',
                    'bg-yellow-200 text-yellow-800': tip.priority === 'medium',
                    'bg-gray-200 text-gray-800': tip.priority === 'low'
                  }"
                >
                  {{ tip.priority === 'high' ? '高优先' : tip.priority === 'medium' ? '中优先' : '低优先' }}
                </span>
              </div>
              <p class="text-sm text-gray-600">{{ tip.description }}</p>
            </div>
          </div>
        </div>

        <!-- 推荐行动 -->
        <div v-if="report.ai_advice?.recommended_actions?.length">
          <div class="text-sm font-medium text-purple-700 mb-2">🚀 推荐行动</div>
          <div class="space-y-3">
            <div 
              v-for="(action, idx) in report.ai_advice.recommended_actions" 
              :key="idx"
              class="bg-purple-50 rounded-lg p-4"
            >
              <div class="font-medium text-purple-800 mb-1">{{ action.action }}</div>
              <div class="text-sm text-gray-600 mb-2">{{ action.reason }}</div>
              <div class="flex items-center gap-4 text-xs">
                <span class="bg-purple-200 text-purple-800 px-2 py-1 rounded">预期效果：{{ action.impact }}</span>
                <span class="bg-gray-200 text-gray-800 px-2 py-1 rounded">难度：{{ action.difficulty === 'easy' ? '简单' : action.difficulty === 'medium' ? '中等' : '困难' }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 激励语 -->
      <div v-if="report.motivation" class="bg-gradient-to-r from-yellow-400 to-orange-500 rounded-xl p-6 text-white text-center">
        <iconify-icon icon="mdi:star-shooting" width="32" class="mb-2"></iconify-icon>
        <p class="text-lg font-medium">{{ report.motivation }}</p>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else class="flex flex-col items-center justify-center py-20 text-gray-500">
      <iconify-icon icon="mdi:chart-box-outline" width="64" class="mb-4 text-gray-300"></iconify-icon>
      <p class="text-lg mb-2">还没有生成学习报告</p>
      <p class="text-sm mb-4">点击上方按钮，让 AI 分析你的学习数据</p>
      <button 
        @click="generateReport"
        class="bg-[#2D5BFF] text-white px-6 py-2 rounded-lg font-medium hover:bg-opacity-90"
      >
        生成第一份报告
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, nextTick } from 'vue'
import { generateAIReport } from '@/api/modules/ai'
import * as echarts from 'echarts'

const loading = ref(false)
const report = ref(null)
const reportDays = ref(7)
const radarChartRef = ref(null)
const trendChartRef = ref(null)

let radarChart = null
let trendChart = null

// 生成报告
const generateReport = async () => {
  loading.value = true
  try {
    const res = await generateAIReport({
      days: reportDays.value,
      report_type: reportDays.value <= 7 ? 'weekly' : 'monthly'
    })
    if (res.code === 200) {
      report.value = res.data
      await nextTick()
      renderCharts()
    } else {
      console.error('生成报告失败:', res.message)
    }
  } catch (err) {
    console.error('生成报告错误:', err)
  } finally {
    loading.value = false
  }
}

// 渲染图表
const renderCharts = () => {
  renderRadarChart()
  renderTrendChart()
}

// 渲染雷达图
const renderRadarChart = () => {
  if (!radarChartRef.value || !report.value?.ability_radar?.dimensions) return
  
  if (radarChart) {
    radarChart.dispose()
  }
  radarChart = echarts.init(radarChartRef.value)

  const dimensions = report.value.ability_radar.dimensions
  const indicator = dimensions.map(d => ({ name: d.name, max: d.max }))
  const values = dimensions.map(d => d.value)

  const option = {
    radar: {
      indicator,
      shape: 'polygon',
      splitNumber: 4,
      axisName: {
        color: '#666',
        fontSize: 12
      },
      splitLine: {
        lineStyle: { color: '#ddd' }
      },
      splitArea: {
        areaStyle: { color: ['#f5f5f5', '#fff'] }
      }
    },
    series: [{
      type: 'radar',
      data: [{
        value: values,
        name: '能力分布',
        areaStyle: {
          color: 'rgba(45, 91, 255, 0.3)'
        },
        lineStyle: {
          color: '#2D5BFF',
          width: 2
        },
        itemStyle: {
          color: '#2D5BFF'
        }
      }]
    }]
  }

  radarChart.setOption(option)
}

// 渲染趋势图
const renderTrendChart = () => {
  if (!trendChartRef.value || !report.value?.behavior_analysis?.daily_trend) return
  
  if (trendChart) {
    trendChart.dispose()
  }
  trendChart = echarts.init(trendChartRef.value)

  const trend = report.value.behavior_analysis.daily_trend
  const dates = trend.map(t => t.date)
  const hours = trend.map(t => t.study_hours)
  const scores = trend.map(t => t.focus_score)

  const option = {
    tooltip: {
      trigger: 'axis'
    },
    legend: {
      data: ['学习时长', '专注度'],
      top: 0
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      top: '15%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: dates,
      axisLabel: { fontSize: 11 }
    },
    yAxis: [
      {
        type: 'value',
        name: '小时',
        axisLabel: { fontSize: 11 }
      },
      {
        type: 'value',
        name: '专注度',
        max: 100,
        axisLabel: { fontSize: 11 }
      }
    ],
    series: [
      {
        name: '学习时长',
        type: 'bar',
        data: hours,
        itemStyle: { color: '#2D5BFF' }
      },
      {
        name: '专注度',
        type: 'line',
        yAxisIndex: 1,
        data: scores,
        smooth: true,
        itemStyle: { color: '#10B981' }
      }
    ]
  }

  trendChart.setOption(option)
}

// 监听窗口变化
onMounted(() => {
  window.addEventListener('resize', () => {
    radarChart?.resize()
    trendChart?.resize()
  })
})

// 监听天数变化
watch(reportDays, () => {
  if (report.value) {
    generateReport()
  }
})
</script>

<style scoped>
.card {
  @apply bg-white rounded-xl p-6 shadow-sm;
}
</style>
