<template>
  <div class="mindmap-node" :style="{ marginLeft: level * 24 + 'px' }">
    <div 
      class="node-content flex items-center gap-2 py-1 px-2 rounded cursor-pointer hover:bg-white transition-colors"
      :class="nodeClass"
      @click="toggleExpand"
    >
      <!-- 展开/收起图标 -->
      <iconify-icon 
        v-if="node.children?.length"
        :icon="expanded ? 'mdi:chevron-down' : 'mdi:chevron-right'"
        class="text-gray-400"
        width="16"
      ></iconify-icon>
      <span v-else class="w-4"></span>
      
      <!-- 节点图标 -->
      <iconify-icon 
        :icon="nodeIcon"
        :class="iconClass"
        width="16"
      ></iconify-icon>
      
      <!-- 节点文本 -->
      <span :class="textClass">{{ node.text }}</span>
    </div>
    
    <!-- 子节点 -->
    <div v-if="expanded && node.children?.length" class="children">
      <MindmapNode 
        v-for="(child, idx) in node.children" 
        :key="idx" 
        :node="child" 
        :level="level + 1"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, defineProps } from 'vue'

const props = defineProps({
  node: { type: Object, required: true },
  level: { type: Number, default: 0 }
})

const expanded = ref(true)

const toggleExpand = () => {
  if (props.node.children?.length) {
    expanded.value = !expanded.value
  }
}

const nodeClass = computed(() => {
  const level = props.level
  if (level === 0) return 'bg-blue-100'
  if (level === 1) return 'bg-green-50'
  return ''
})

const nodeIcon = computed(() => {
  const level = props.level
  if (level === 0) return 'mdi:lightbulb'
  if (level === 1) return 'mdi:folder'
  return 'mdi:circle-small'
})

const iconClass = computed(() => {
  const level = props.level
  if (level === 0) return 'text-blue-500'
  if (level === 1) return 'text-green-500'
  return 'text-gray-400'
})

const textClass = computed(() => {
  const level = props.level
  if (level === 0) return 'font-bold text-blue-800'
  if (level === 1) return 'font-medium text-gray-800'
  return 'text-gray-600 text-sm'
})
</script>

<style scoped>
.mindmap-node {
  position: relative;
}

.mindmap-node::before {
  content: '';
  position: absolute;
  left: 8px;
  top: 24px;
  bottom: 0;
  width: 1px;
  background: #e5e7eb;
}

.mindmap-node:last-child::before {
  display: none;
}

.children {
  border-left: 1px dashed #d1d5db;
  margin-left: 8px;
  padding-left: 8px;
}
</style>
