<template>
  <div class="note-enhance-panel">
    <!-- 增强功能选择 -->
    <div class="flex items-center gap-2 mb-4 flex-wrap">
      <button 
        v-for="btn in enhanceButtons" 
        :key="btn.type"
        @click="handleEnhance(btn.type)"
        :disabled="loading"
        class="px-3 py-2 rounded-lg text-sm font-medium flex items-center gap-1 transition-colors"
        :class="activeType === btn.type ? 'bg-[#2D5BFF] text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'"
      >
        <iconify-icon :icon="btn.icon" width="16"></iconify-icon>
        {{ btn.label }}
      </button>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="flex items-center justify-center py-8">
      <iconify-icon icon="mdi:loading" class="animate-spin text-[#2D5BFF]" width="24"></iconify-icon>
      <span class="ml-2 text-gray-600">AI 正在处理...</span>
    </div>

    <!-- 结果展示 -->
    <div v-else-if="result" class="space-y-4">
      <!-- 摘要结果 -->
      <div v-if="result.summary" class="bg-white rounded-lg p-4 border">
        <h4 class="font-medium text-gray-900 mb-3 flex items-center gap-2">
          <iconify-icon icon="mdi:text-box-outline" width="18" class="text-blue-500"></iconify-icon>
          智能摘要
        </h4>
        <div class="space-y-3">
          <div class="bg-blue-50 p-3 rounded-lg">
            <div class="text-sm text-gray-500 mb-1">一句话总结</div>
            <p class="text-gray-800">{{ result.summary.brief }}</p>
          </div>
          <div>
            <div class="text-sm text-gray-500 mb-2">核心要点</div>
            <ul class="space-y-2">
              <li 
                v-for="(point, idx) in result.summary.key_points" 
                :key="idx"
                class="flex items-start gap-2 text-gray-700"
              >
                <iconify-icon icon="mdi:checkbox-marked-circle" class="text-green-500 mt-0.5"></iconify-icon>
                {{ point }}
              </li>
            </ul>
          </div>
          <div v-if="result.summary.conclusion" class="bg-gray-50 p-3 rounded-lg">
            <div class="text-sm text-gray-500 mb-1">总结</div>
            <p class="text-gray-700">{{ result.summary.conclusion }}</p>
          </div>
          <div class="text-xs text-gray-400 flex items-center gap-1">
            <iconify-icon icon="mdi:clock-outline"></iconify-icon>
            预计阅读时间：{{ result.summary.reading_time }} 分钟
          </div>
        </div>
      </div>

      <!-- 关键词结果 -->
      <div v-if="result.keywords" class="bg-white rounded-lg p-4 border">
        <h4 class="font-medium text-gray-900 mb-3 flex items-center gap-2">
          <iconify-icon icon="mdi:tag-multiple" width="18" class="text-purple-500"></iconify-icon>
          关键词提取
        </h4>
        <div class="space-y-4">
          <div>
            <div class="text-sm text-gray-500 mb-2">主要关键词</div>
            <div class="flex flex-wrap gap-2">
              <span 
                v-for="kw in result.keywords.main_keywords" 
                :key="kw.word"
                class="inline-flex items-center gap-1 px-3 py-1 rounded-full text-sm"
                :style="{ 
                  backgroundColor: `rgba(99, 102, 241, ${0.1 + kw.weight * 0.3})`,
                  color: '#4f46e5'
                }"
              >
                {{ kw.word }}
                <span class="text-xs opacity-70">({{ (kw.weight * 100).toFixed(0) }}%)</span>
              </span>
            </div>
          </div>
          <div v-if="result.keywords.related_concepts?.length">
            <div class="text-sm text-gray-500 mb-2">相关概念</div>
            <div class="flex flex-wrap gap-2">
              <span 
                v-for="concept in result.keywords.related_concepts" 
                :key="concept"
                class="bg-gray-100 text-gray-700 px-2 py-1 rounded text-sm"
              >
                {{ concept }}
              </span>
            </div>
          </div>
          <div class="flex items-center gap-4">
            <div v-if="result.keywords.category" class="flex items-center gap-1 text-sm">
              <iconify-icon icon="mdi:folder" class="text-orange-500"></iconify-icon>
              分类：{{ result.keywords.category }}
            </div>
            <div v-if="result.keywords.tags?.length" class="flex items-center gap-2">
              <iconify-icon icon="mdi:tag" class="text-gray-400"></iconify-icon>
              <span v-for="tag in result.keywords.tags" :key="tag" class="text-sm text-blue-600">#{{ tag }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 思维导图结果 -->
      <div v-if="result.mindmap" class="bg-white rounded-lg p-4 border">
        <h4 class="font-medium text-gray-900 mb-3 flex items-center gap-2">
          <iconify-icon icon="mdi:sitemap" width="18" class="text-green-500"></iconify-icon>
          思维导图
        </h4>
        <div class="bg-gray-50 rounded-lg p-4 overflow-x-auto">
          <div class="mindmap-tree">
            <MindmapNode :node="result.mindmap.root" :level="0" />
          </div>
        </div>
        <div v-if="result.mindmap.markdown" class="mt-4">
          <button 
            @click="showMarkdown = !showMarkdown"
            class="text-sm text-blue-600 hover:underline flex items-center gap-1"
          >
            <iconify-icon :icon="showMarkdown ? 'mdi:chevron-up' : 'mdi:chevron-down'"></iconify-icon>
            {{ showMarkdown ? '收起' : '查看' }} Markdown 大纲
          </button>
          <pre v-if="showMarkdown" class="mt-2 bg-gray-900 text-gray-100 p-3 rounded text-sm overflow-x-auto">{{ result.mindmap.markdown }}</pre>
        </div>
      </div>

      <!-- 复习问题结果 -->
      <div v-if="result.questions" class="bg-white rounded-lg p-4 border">
        <h4 class="font-medium text-gray-900 mb-3 flex items-center gap-2">
          <iconify-icon icon="mdi:help-circle" width="18" class="text-orange-500"></iconify-icon>
          复习问题
        </h4>
        <div class="space-y-4">
          <div v-for="(q, idx) in result.questions.review_questions" :key="idx" class="border rounded-lg p-3">
            <div class="flex items-start justify-between gap-2 mb-2">
              <div class="font-medium text-gray-800">
                {{ idx + 1 }}. {{ q.question }}
              </div>
              <div class="flex gap-1">
                <span 
                  class="text-xs px-2 py-0.5 rounded"
                  :class="{
                    'bg-green-100 text-green-700': q.difficulty === 'easy',
                    'bg-yellow-100 text-yellow-700': q.difficulty === 'medium',
                    'bg-red-100 text-red-700': q.difficulty === 'hard'
                  }"
                >
                  {{ q.difficulty === 'easy' ? '简单' : q.difficulty === 'medium' ? '中等' : '困难' }}
                </span>
                <span class="text-xs px-2 py-0.5 rounded bg-gray-100 text-gray-600">
                  {{ q.type === 'concept' ? '概念' : q.type === 'application' ? '应用' : '分析' }}
                </span>
              </div>
            </div>
            <div class="bg-green-50 p-2 rounded text-sm text-gray-700">
              <span class="text-green-600 font-medium">参考答案：</span>{{ q.answer }}
            </div>
          </div>
          <div v-if="result.questions.thinking_questions?.length">
            <div class="text-sm font-medium text-gray-700 mb-2">💭 思考题</div>
            <ul class="space-y-2">
              <li 
                v-for="(tq, idx) in result.questions.thinking_questions" 
                :key="idx"
                class="text-gray-600 text-sm bg-yellow-50 p-2 rounded"
              >
                {{ tq }}
              </li>
            </ul>
          </div>
        </div>
      </div>

      <!-- 润色结果 -->
      <div v-if="result.polish" class="bg-white rounded-lg p-4 border">
        <h4 class="font-medium text-gray-900 mb-3 flex items-center gap-2">
          <iconify-icon icon="mdi:auto-fix" width="18" class="text-pink-500"></iconify-icon>
          笔记润色
        </h4>
        <div class="space-y-4">
          <div v-if="result.polish.improvements?.length">
            <div class="text-sm text-gray-500 mb-2">改进说明</div>
            <ul class="space-y-1">
              <li 
                v-for="(imp, idx) in result.polish.improvements" 
                :key="idx"
                class="flex items-center gap-2 text-sm text-green-700"
              >
                <iconify-icon icon="mdi:check"></iconify-icon>
                {{ imp }}
              </li>
            </ul>
          </div>
          <div>
            <div class="text-sm text-gray-500 mb-2">润色后的内容</div>
            <div 
              class="prose prose-sm max-w-none bg-gray-50 p-4 rounded-lg" 
              v-html="renderMarkdown(result.polish.polished_content)"
            ></div>
          </div>
          <div v-if="result.polish.suggestions?.length">
            <div class="text-sm text-gray-500 mb-2">进一步建议</div>
            <ul class="space-y-1">
              <li 
                v-for="(sug, idx) in result.polish.suggestions" 
                :key="idx"
                class="flex items-center gap-2 text-sm text-blue-700"
              >
                <iconify-icon icon="mdi:lightbulb"></iconify-icon>
                {{ sug }}
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else class="text-center py-8 text-gray-500">
      <iconify-icon icon="mdi:magic-staff" width="48" class="mb-2 text-gray-300"></iconify-icon>
      <p>选择上方功能，让 AI 增强你的笔记</p>
    </div>
  </div>
</template>

<script setup>
import { ref, defineProps } from 'vue'
import { enhanceNote } from '@/api/modules/ai'
import { marked } from 'marked'

// 思维导图节点组件
import MindmapNode from './MindmapNode.vue'

const props = defineProps({
  noteId: { type: Number, default: 0 },
  content: { type: String, default: '' },
  title: { type: String, default: '' }
})

const loading = ref(false)
const result = ref(null)
const activeType = ref('')
const showMarkdown = ref(false)

const enhanceButtons = [
  { type: 'all', label: '全面增强', icon: 'mdi:auto-awesome' },
  { type: 'summary', label: '生成摘要', icon: 'mdi:text-box-outline' },
  { type: 'keywords', label: '提取关键词', icon: 'mdi:tag-multiple' },
  { type: 'mindmap', label: '思维导图', icon: 'mdi:sitemap' },
  { type: 'questions', label: '复习问题', icon: 'mdi:help-circle' },
  { type: 'polish', label: '润色优化', icon: 'mdi:auto-fix' }
]

const handleEnhance = async (type) => {
  if (!props.content && !props.noteId) {
    alert('请先输入笔记内容')
    return
  }

  loading.value = true
  activeType.value = type
  result.value = null

  try {
    const res = await enhanceNote({
      note_id: props.noteId || undefined,
      content: props.content,
      title: props.title,
      type
    })
    if (res.code === 200) {
      result.value = res.data
    } else {
      console.error('增强失败:', res.message)
    }
  } catch (err) {
    console.error('增强错误:', err)
  } finally {
    loading.value = false
  }
}

const renderMarkdown = (content) => {
  if (!content) return ''
  return marked(content)
}
</script>

<style scoped>
.note-enhance-panel {
  @apply bg-gray-50 rounded-xl p-4;
}

.mindmap-tree {
  min-width: fit-content;
}
</style>
