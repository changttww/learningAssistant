<template>
  <div
    class="constellation-wrapper"
    ref="containerRef"
    @pointermove="handlePointerMove"
    @pointerleave="handlePointerLeave"
  >
    <canvas ref="particleCanvas" class="particle-canvas"></canvas>

    <div class="constellation-layer" :style="layerStyle">
      <svg class="constellation-path" :viewBox="`0 0 ${viewportWidth} ${viewportHeight}`">
        <polyline :points="mainPolyline" />
      </svg>

      <div
        v-for="star in mainStars"
        :key="star.id"
        class="star-node"
        v-show="!detailViewActive || isFocusedStar(star)"
        :class="{ 'is-focused': isFocusedStar(star) }"
        :style="starStyle(star)"
        @click="handleStarClick(star)"
      >
        <span class="star-core"></span>
        <div class="star-dialog">
          <p class="star-name">{{ (star.task && star.task.title) || star.name }}</p>
          <p class="star-meta" v-if="star.task && star.task.due_date">截止：{{ formatDueDate(star.task.due_date) }}</p>
          <small class="star-hint" v-if="!detailViewActive">点击进入详情</small>
        </div>

        <transition name="detail-cluster">
          <div v-if="detailViewActive && isFocusedStar(star)" class="detail-cluster">
            <div v-if="detailLoading" class="detail-loading">星图数据加载中...</div>
            
            <!-- 连线层 -->
            <svg class="detail-connections" viewBox="-500 -500 1000 1000">
              <line 
                v-for="sat in detailSatellites" 
                :key="`${sat.id}-line`"
                :x1="sat.lineStart.x" 
                :y1="sat.lineStart.y" 
                :x2="sat.lineEnd.x" 
                :y2="sat.lineEnd.y" 
              />
            </svg>

            <div
              v-for="satellite in detailSatellites"
              :key="satellite.id"
              class="satellite-node"
              :class="`satellite-node--${satellite.type}`"
              :style="{ transform: `translate3d(${satellite.x}px, ${satellite.y}px, 0)` }"
            >
              <span class="satellite-core"></span>
              <div class="satellite-dialog">
                <p class="satellite-title">{{ satellite.title }}</p>
              </div>
            </div>
          </div>
        </transition>
      </div>
    </div>

    <div class="hud">
      <div class="quote-card">
        <p class="quote-text">{{ currentQuote.text }}</p>
      </div>
      <button
        class="hud-btn"
        @click="detailViewActive ? closeDetail() : router.go(-1)"
      >
        {{ detailViewActive ? '返回星图' : '返回团队任务' }}
      </button>
    </div>

  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onBeforeUnmount, watch } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { getTeamTasks, getTaskDetail } from '@/api/modules/task';

const runtimeWindow = typeof globalThis === 'undefined' ? undefined : globalThis.window;
const router = useRouter();
const route = useRoute();
const containerRef = ref(null);
const particleCanvas = ref(null);
const viewportWidth = ref(runtimeWindow?.innerWidth || 1920);
const viewportHeight = ref(runtimeWindow?.innerHeight || 1080);

const tasks = ref([]);
const detailViewActive = ref(false);
const selectedStar = ref(null);
const detailData = ref({ subtasks: [], attachments: [] });
const detailLoading = ref(false);
const zoomTransform = ref({ x: 0, y: 0, scale: 1 });

const pointer = reactive({ x: 0, y: 0, active: false });

let pointerIdleTimer = null;
let resizeHandler = null;
let particleAnimationId = null;

const fallbackTasks = [];

const currentLayout = ref([]);

const updateLayout = () => {
  const source = tasks.value;
  const count = source.length;

  const newLayout = [];
  
  // 调整分布范围逻辑：
  // 基础范围给大一点，避免挤在一起 (0.6 起步)
  // 随数量增加适当扩大，最大铺满 90% 宽，75% 高
  const baseSpread = 0.5; 
  const extraSpread = Math.min(0.4, count * 0.03); // 每增加一个任务增加一点范围
  const spreadFactor = baseSpread + extraSpread;
  
  const rangeX = 0.9 * spreadFactor;
  const rangeY = 0.75 * spreadFactor;
  
  const minX = 0.5 - rangeX / 2;
  const minY = 0.5 - rangeY / 2;

  // 动态调整最小间距，任务越多间距容忍度越小
  const minDistance = Math.max(0.08, 0.25 - count * 0.015);

  for (let i = 0; i < count; i++) {
    let x, y;
    let attempts = 0;
    let valid = false;

    // 尝试生成不重叠的坐标
    while (attempts < 50 && !valid) {
      x = minX + Math.random() * rangeX;
      y = minY + Math.random() * rangeY;
      
      valid = true;
      for (const existing of newLayout) {
        // 考虑宽高比，让x轴的距离权重稍微大一点（因为屏幕通常是宽屏）
        const dx = (existing.x - x) * 1.2; 
        const dy = existing.y - y;
        if (Math.sqrt(dx * dx + dy * dy) < minDistance) {
          valid = false;
          break;
        }
      }
      attempts++;
    }
    
    newLayout.push({
      id: `star-gen-${i}`,
      name: `节点 ${i + 1}`,
      x: x,
      y: y,
      depth: 0.7 + Math.random() * 0.5,
    });
  }
  newLayout.sort((a, b) => a.x - b.x);
  currentLayout.value = newLayout;
};

watch(tasks, updateLayout, { deep: true });

const mainStars = computed(() => {
  const source = tasks.value;
  
  // 按照截止时间排序（从早到晚）
  const sortedSource = [...source].sort((a, b) => {
    const dateA = new Date(a.due_date || a.due_at || 0).getTime();
    const dateB = new Date(b.due_date || b.due_at || 0).getTime();
    return dateA - dateB;
  });

  // 确保布局更新后再映射，避免越界（虽然 updateLayout 是同步的）
  if (currentLayout.value.length !== sortedSource.length) {
    updateLayout();
  }
  return currentLayout.value.map((layout, index) => ({
    ...layout,
    task: sortedSource[index] || null,
  }));
});

const formatDueDate = (dateStr) => {
  if (!dateStr) return '';
  const date = new Date(dateStr);
  if (Number.isNaN(date.getTime())) return dateStr;
  
  const year = date.getFullYear();
  const month = date.getMonth() + 1;
  const day = date.getDate();
  const hours = date.getHours();
  const minutes = date.getMinutes();
  
  const pad = (n) => n < 10 ? `0${n}` : n;
  
  return `${year}年${month}月${day}日 ${pad(hours)}:${pad(minutes)}`;
};

const detailSatellites = computed(() => {
  if (!detailViewActive.value || !selectedStar.value) return [];
  const payload = detailData.value || {};
  const items = Array.isArray(payload.subtasks) ? payload.subtasks : [];
  if (!items.length) return [];

  const satellites = [];
  const minDistance = 50; // 卫星之间的最小间距

  for (const [index, item] of items.entries()) {
    let dx, dy;
    let attempts = 0;
    let valid = false;

    // 兼容字符串类型的子任务（后端返回的是字符串数组）
    const title = typeof item === 'string' ? item : item.title;
    const id = (typeof item === 'object' && item.id) ? item.id : index;

    // 尝试生成不重叠的随机位置
    while (attempts < 50 && !valid) {
      // 随机角度
      const angle = Math.random() * Math.PI * 2;
      // 随机半径 (80px - 180px)
      const radius = 80 + Math.random() * 100;
      
      dx = Math.cos(angle) * radius;
      dy = Math.sin(angle) * radius;

      valid = true;
      for (const existing of satellites) {
        const dist = Math.hypot(existing.x - dx, existing.y - dy);
        if (dist < minDistance) {
          valid = false;
          break;
        }
      }
      attempts++;
    }

    satellites.push({
      id: `${selectedStar.value.id}-sub-${id}`,
      type: 'subtasks',
      title: title,
      x: dx,
      y: dy,
      // 连线坐标：从主星中心(0,0)到卫星中心(dx, dy)
      lineStart: { x: 0, y: 0 },
      lineEnd: { x: dx, y: dy }
    });
  }

  return satellites;
});

const layerStyle = computed(() => ({
  transform:
    detailViewActive.value && selectedStar.value
      ? `translate3d(${zoomTransform.value.x}px, ${zoomTransform.value.y}px, 0) scale(${zoomTransform.value.scale})`
      : 'translate3d(0, 0, 0) scale(1)',
}));

const quotes = [
  { text: '天行健，君子以自强不息' },
  { text: '路漫漫其修远兮，吾将上下而求索' },
  { text: '学而不思则罔，思而不学则殆' },
  { text: '博观而约取，厚积而薄发' },
  { text: '业精于勤，荒于嬉；行成于思，毁于随' },
  { text: '锲而不舍，金石可镂' },
  { text: '不积跬步，无以至千里' },
  { text: '会当凌绝顶，一览众山小' },
  { text: '长风破浪会有时，直挂云帆济沧海' },
  { text: '非淡泊无以明志，非宁静无以致远' }
];
const currentQuote = ref(quotes[0]);

const projectPoint = (node) => ({
  x: node.x * viewportWidth.value,
  y: node.y * viewportHeight.value,
});

const starStyle = (node) => {
  const { x, y } = projectPoint(node);
  return { transform: `translate3d(${x}px, ${y}px, 0)` };
};

const polylineFor = (layouts) =>
  layouts
    .map((layout) => {
      const { x, y } = projectPoint(layout);
      // 此时 star-node 已对齐到中心，直接连接中心坐标即可
      return `${x},${y}`;
    })
    .join(' ');

const mainPolyline = computed(() => polylineFor(currentLayout.value));

const loadTasks = async () => {
  try {
    const params = {};
    if (route.query.teamId) {
      params.team_id = route.query.teamId;
    }
    const res = await getTeamTasks(params);
    const data = res?.data?.items || res?.data || res;
    if (Array.isArray(data) && data.length) {
      tasks.value = data.map((item, index) => ({
        id: item.id || runtimeWindow?.crypto?.randomUUID?.() || `${Date.now()}-${index}`,
        title: item.title || item.name || '未命名任务',
        description: item.description || '',
        owner_name: item.owner_name || item.created_by_name || '未知',
        due_date: item.due_at || item.due_date || '',
      }));
    }
  } catch (error) {
    console.warn('获取团队任务失败：', error);
    tasks.value = [];
  }
};

const ensureTaskDetail = async (task) => {
  if (!task) return null;
  try {
    const res = await getTaskDetail(task.id);
    const raw = res?.data?.data || res?.data || res;
    return normalizeDetail(raw, task);
  } catch (error) {
    console.warn('获取任务详情失败：', error);
    return { subtasks: [], attachments: [] };
  }
};

const normalizeDetail = (raw, task) => {
  if (!raw) return { subtasks: [], attachments: [] };
  return {
    subtasks: Array.isArray(raw.subtasks) ? raw.subtasks : [],
    attachments: Array.isArray(raw.attachments) ? raw.attachments : [],
  };
};

const buildMockDetail = null;

const isFocusedStar = (star) => Boolean(selectedStar.value && selectedStar.value.id === star.id);

const satelliteStyle = null;
const satelliteLinkStyle = null;

const focusOnStar = (star) => {
  if (!star) return;
  const { x, y } = projectPoint(star);
  const centerX = viewportWidth.value / 2;
  const centerY = viewportHeight.value / 2;
  // 减小缩放比例，避免子任务超出屏幕 (原 2.6 -> 1.6)
  const scale = 1.6;
  zoomTransform.value = {
    x: centerX - x * scale,
    y: centerY - y * scale,
    scale,
  };
};

const handleStarClick = async (star) => {
  if (!star.task) return;
  selectedStar.value = star;
  detailViewActive.value = true;
  focusOnStar(star);
  detailLoading.value = true;
  detailData.value = await ensureTaskDetail(star.task);
  detailLoading.value = false;
};

const closeDetail = () => {
  detailViewActive.value = false;
  selectedStar.value = null;
  detailData.value = { subtasks: [], attachments: [] };
  zoomTransform.value = { x: 0, y: 0, scale: 1 };
};

const handlePointerMove = (event) => {
  const rect = containerRef.value?.getBoundingClientRect();
  if (!rect) return;
  pointer.x = event.clientX - rect.left;
  pointer.y = event.clientY - rect.top;
  pointer.active = true;
  if (pointerIdleTimer) clearTimeout(pointerIdleTimer);
  pointerIdleTimer = runtimeWindow?.setTimeout(() => {
    pointer.active = false;
  }, 800);
};

const handlePointerLeave = () => {
  pointer.active = false;
};

const particleField = [];

const seedParticles = (canvas) => {
  particleField.length = 0;
  const densityBase = Math.floor((canvas.width * canvas.height) / 3500) + 140;
  const total = clamp(densityBase, 420, 880);
  for (let i = 0; i < total; i += 1) {
    particleField.push({
      x: Math.random() * canvas.width,
      y: Math.random() * canvas.height,
      vx: 0,
      vy: 0,
      size: 0.4 + Math.random() * 1.4,
      depth: 0.35 + Math.random() * 0.65,
      twinkle: Math.random() * Math.PI * 2,
      speed: 0.01 + Math.random() * 0.01,
    });
  }
};

const requestFrame = (cb) =>
  (runtimeWindow?.requestAnimationFrame ?? ((fn) => setTimeout(fn, 16)))(cb);

const cancelFrame = (id) => {
  const cancel = runtimeWindow?.cancelAnimationFrame ?? clearTimeout;
  cancel(id);
};

const initParticles = () => {
  const canvas = particleCanvas.value;
  const container = containerRef.value;
  if (!canvas || !container) return;
  const ctx = canvas.getContext('2d');

  const syncSize = () => {
    const rect = container.getBoundingClientRect();
    canvas.width = rect.width;
    canvas.height = rect.height;
    viewportWidth.value = rect.width;
    viewportHeight.value = rect.height;
    seedParticles(canvas);
    if (detailViewActive.value && selectedStar.value) {
      focusOnStar(selectedStar.value);
    }
  };

  syncSize();

  const render = () => {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    ctx.globalCompositeOperation = 'lighter';
    const influenceRadius = 200;

    for (const particle of particleField) {
      const dx = particle.x - pointer.x;
      const dy = particle.y - pointer.y;
      const distance = Math.max(Math.hypot(dx, dy), 0.001);

      if (pointer.active && distance < influenceRadius) {
        const force = (1 - distance / influenceRadius) * 7;
        particle.vx += (dx / distance) * force;
        particle.vy += (dy / distance) * force;
      }

      particle.vx *= 0.92;
      particle.vy *= 0.92;
      particle.twinkle += particle.speed;
      particle.x += particle.vx + Math.cos(particle.twinkle) * 0.25;
      particle.y += particle.vy + Math.sin(particle.twinkle) * 0.2;

      // 边界检查与重置
      const buffer = 40;
      const width = canvas.width;
      const height = canvas.height;

      // 防止 NaN 或无限大导致粒子消失
      if (!Number.isFinite(particle.x) || !Number.isFinite(particle.y)) {
        particle.x = Math.random() * width;
        particle.y = Math.random() * height;
        particle.vx = 0;
        particle.vy = 0;
      }

      if (particle.x < -buffer) {
        particle.x = width + buffer;
      } else if (particle.x > width + buffer) {
        particle.x = -buffer;
      }

      if (particle.y < -buffer) {
        particle.y = height + buffer;
      } else if (particle.y > height + buffer) {
        particle.y = -buffer;
      }

      const drawX = particle.x;
      const drawY = particle.y;
      const alpha = 0.25 + (Math.sin(particle.twinkle) + 1) * 0.15;

      ctx.beginPath();
      ctx.arc(drawX, drawY, particle.size, 0, Math.PI * 2);
      ctx.fillStyle = `rgba(255,255,255,${alpha})`;
      ctx.shadowColor = 'rgba(130,160,255,0.8)';
      ctx.shadowBlur = 12 * particle.depth;
      ctx.fill();
    }

    particleAnimationId = requestFrame(render);
  };

  particleAnimationId = requestFrame(render);
  resizeHandler = () => syncSize();
  runtimeWindow?.addEventListener('resize', resizeHandler);
};

const clamp = (value, min, max) => Math.min(max, Math.max(min, value));

onMounted(() => {
  updateLayout();
  loadTasks();
  initParticles();
  
  // 随机展示古文
  const idx = Math.floor(Math.random() * quotes.length);
  currentQuote.value = quotes[idx];
});

onBeforeUnmount(() => {
  if (particleAnimationId) cancelFrame(particleAnimationId);
  if (resizeHandler && runtimeWindow) runtimeWindow.removeEventListener('resize', resizeHandler);
  if (pointerIdleTimer) clearTimeout(pointerIdleTimer);
});
</script>

<style scoped>
.constellation-wrapper {
  position: relative;
  width: 100%;
  height: 100vh;
  overflow: hidden;
  background: radial-gradient(circle at 20% -20%, #1d234d 0%, #070a1b 55%, #010109 100%);
  color: #fff;
  cursor: default;
}

.particle-canvas {
  position: absolute;
  inset: 0;
}

.constellation-layer {
  position: relative;
  width: 100%;
  height: 100%;
  pointer-events: none;
  transition: transform 0.65s cubic-bezier(0.22, 1, 0.36, 1);
  transform-origin: 0 0;
}

.constellation-path,
.detail-path {
  position: absolute;
  width: 100%;
  height: 100%;
  opacity: 0.35;
}

.constellation-path polyline,
.detail-path polyline {
  fill: none;
  stroke: rgba(255, 255, 255, 0.45);
  stroke-width: 2;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.star-node {
  position: absolute;
  top: 0;
  left: 0;
  width: 0;
  height: 0;
  pointer-events: auto;
  transform: translate3d(0, 0, 0);
  overflow: visible;
}

.star-core {
  position: absolute;
  left: -8px;
  top: -8px;
  display: block;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: linear-gradient(145deg, #fefefe, #8ab8ff);
  box-shadow: 0 0 18px rgba(138, 184, 255, 0.9), 0 0 32px rgba(255, 255, 255, 0.5);
  animation: pulse 2.8s ease-in-out infinite;
}

.star-node:hover .star-core,
.star-node.is-focused .star-core {
  transform: scale(1.2);
}

.star-dialog {
  position: absolute;
  left: 24px;
  top: -24px;
  padding: 12px 16px;
  width: max-content;
  max-width: 300px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  background: rgba(6, 10, 28, 0.85);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.4);
}

.star-dialog::before {
  content: '';
  position: absolute;
  left: -24px;
  top: 24px;
  width: 24px;
  height: 2px;
  background: linear-gradient(90deg, rgba(255, 255, 255, 0), rgba(255, 255, 255, 0.7));
}

.star-name {
  font-size: 18px;
  font-weight: 600;
}

.star-meta {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.8);
}

.star-hint {
  display: block;
  margin-top: 6px;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.6);
}

.hud {
  position: absolute;
  top: 32px;
  right: 32px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  z-index: 5;
  align-items: flex-end;
}

.quote-card {
  padding: 14px 18px;
  border-radius: 16px;
  background: rgba(8, 14, 46, 0.75);
  border: 1px solid rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(6px);
  max-width: 320px;
}

.quote-text {
  font-size: 14px;
  line-height: 1.5;
  color: rgba(255, 255, 255, 0.8);
  letter-spacing: 0.05em;
}

.hud-btn {
  padding: 10px 16px;
  border-radius: 999px;
  border: 1px solid rgba(255, 255, 255, 0.25);
  background: transparent;
  color: #fff;
  font-size: 14px;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  cursor: pointer;
  transition: background 0.2s ease, color 0.2s ease, border-color 0.2s ease;
}

.hud-btn:hover {
  background: rgba(255, 255, 255, 0.12);
  border-color: rgba(255, 255, 255, 0.5);
}

.hud-btn:active {
  background: rgba(255, 255, 255, 0.2);
}

@keyframes pulse {
  0% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.2);
  }
  100% {
    transform: scale(1);
  }
}


.detail-cluster {
  position: absolute;
  top: 0;
  left: 0;
  width: 0;
  height: 0;
  pointer-events: none;
  overflow: visible; /* 关键：允许内容溢出 */
}

.detail-loading {
  position: absolute;
  top: -120px;
  left: -40px;
  padding: 8px 14px;
  border-radius: 999px;
  background: rgba(6, 10, 28, 0.8);
  border: 1px solid rgba(255, 255, 255, 0.2);
  font-size: 12px;
  letter-spacing: 0.08em;
  pointer-events: auto;
}

.detail-connections {
  position: absolute;
  top: -500px;
  left: -500px;
  width: 1000px;
  height: 1000px;
  pointer-events: none;
}

.detail-connections line {
  stroke: rgba(255, 255, 255, 0.65);
  stroke-width: 1.5;
  stroke-linecap: round;
  filter: drop-shadow(0 0 4px rgba(126, 166, 255, 0.5));
}

.satellite-node {
  position: absolute;
  top: 0;
  left: 0;
  width: 0;
  height: 0;
  overflow: visible;
  pointer-events: auto;
}

.satellite-core {
  position: absolute;
  left: -5px;
  top: -5px;
  display: block;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: linear-gradient(145deg, #e8f0ff, #88a9ff);
  box-shadow: 0 0 14px rgba(136, 169, 255, 0.85);
}

.satellite-node--attachments .satellite-core {
  background: linear-gradient(145deg, #ffd9c2, #ff9f6b);
  box-shadow: 0 0 14px rgba(255, 174, 127, 0.9);
}

.satellite-dialog {
  position: absolute;
  left: 12px;
  top: -16px;
  padding: 8px 12px;
  border-radius: 10px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  background: rgba(5, 9, 24, 0.92);
  width: max-content;
  max-width: 240px;
  box-shadow: 0 10px 26px rgba(0, 0, 0, 0.35);
}

.satellite-title {
  font-size: 13px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.9);
  margin-bottom: 2px;
}

.satellite-dialog small {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.65);
}

.detail-cluster-enter-active,
.detail-cluster-leave-active {
  transition: opacity 0.25s ease;
}

.detail-cluster-enter-from,
.detail-cluster-leave-to {
  opacity: 0;
}
</style>
