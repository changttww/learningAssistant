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
          <p class="star-meta" v-if="star.task && star.task.owner_name">负责人：{{ star.task.owner_name }}</p>
          <p class="star-meta" v-if="star.task && star.task.due_date">截止：{{ star.task.due_date }}</p>
          <small class="star-hint">点击进入详情</small>
        </div>

        <transition name="detail-cluster">
          <div v-if="detailViewActive && isFocusedStar(star)" class="detail-cluster">
            <div v-if="detailLoading" class="detail-loading">星图数据加载中...</div>
            <div
              v-for="satellite in detailSatellites"
              :key="`${satellite.id}-link`"
              class="satellite-link"
              :style="satelliteLinkStyle(satellite)"
            ></div>
            <div
              v-for="satellite in detailSatellites"
              :key="satellite.id"
              class="satellite-node"
              :class="`satellite-node--${satellite.type}`"
              :style="satelliteStyle(satellite)"
            >
              <span class="satellite-core"></span>
              <div class="satellite-dialog">
                <p class="satellite-title">{{ satellite.title }}</p>
                <small>{{ satellite.meta }}</small>
              </div>
            </div>
          </div>
        </transition>
      </div>
    </div>

    <div class="hud">
      <div class="status-card">
        <p class="status-title">交互状态</p>
        <p class="status-value">{{ interactionStatus }}</p>
        <p class="status-meta">{{ interactionMeta }}</p>
        <p class="status-tip">移动鼠标划开星尘 · 点击节点查看详情</p>
      </div>
      <button
        class="hud-btn"
        @click="detailViewActive ? closeDetail() : router.push('/team-tasks')"
      >
        {{ detailViewActive ? '返回星图' : '返回团队任务' }}
      </button>
    </div>

  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onBeforeUnmount } from 'vue';
import { useRouter } from 'vue-router';
import { getTeamTasks, getTaskDetail } from '@/api/modules/task';

const runtimeWindow = typeof globalThis === 'undefined' ? undefined : globalThis.window;
const router = useRouter();
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

const dipperLayout = [
  { id: 'star-1', name: '天枢', x: 0.12, y: 0.62, depth: 0.85 },
  { id: 'star-2', name: '天璇', x: 0.24, y: 0.46, depth: 0.7 },
  { id: 'star-3', name: '天玑', x: 0.36, y: 0.34, depth: 0.95 },
  { id: 'star-4', name: '天权', x: 0.5, y: 0.4, depth: 1.1 },
  { id: 'star-5', name: '玉衡', x: 0.63, y: 0.29, depth: 0.9 },
  { id: 'star-6', name: '开阳', x: 0.76, y: 0.43, depth: 1.05 },
  { id: 'star-7', name: '摇光', x: 0.88, y: 0.58, depth: 0.8 },
];

const fallbackTasks = [
  { id: 1, title: '登录功能开发', description: '完善登录注册与安全策略', owner_name: '王同学', due_date: '2025-12-30' },
  { id: 2, title: '支付模块设计', description: '设计支付流程与风控', owner_name: '陈同学', due_date: '2025-12-12' },
  { id: 3, title: '学习看板重构', description: '实现全新任务看板交互', owner_name: '李同学', due_date: '2026-01-05' },
  { id: 4, title: '通知中心优化', description: '重构实时通知分发链路', owner_name: '赵同学', due_date: '2025-12-03' },
  { id: 5, title: '学习房间稳定性', description: '优化 WebRTC 链路和重连策略', owner_name: '周同学', due_date: '2026-02-01' },
];

const mainStars = computed(() => {
  const source = tasks.value.length ? tasks.value : fallbackTasks;
  return dipperLayout.map((layout, index) => ({
    ...layout,
    task: source[index % source.length] || null,
  }));
});

const satelliteLayout = {
  subtasks: { spread: 110, baseAngle: -10, radius: 190 },
  attachments: { spread: 90, baseAngle: 155, radius: 210 },
};

const detailSatellites = computed(() => {
  if (!detailViewActive.value || !selectedStar.value) return [];
  const payload = detailData.value || {};
  const buildNodes = (type) => {
    const items = Array.isArray(payload[type]) ? payload[type].slice(0, 6) : [];
    if (!items.length) return [];
    const config = satelliteLayout[type];
    const spread = config.spread;
    const base = config.baseAngle - spread / 2;
    const step = items.length > 1 ? spread / (items.length - 1) : spread / 2;
    return items.map((item, index) => ({
      id: `${selectedStar.value.id}-${type}-${item.id || index}`,
      type,
      title: type === 'subtasks' ? item.title : item.name,
      meta: type === 'subtasks' ? item.status : item.size,
      angle: base + step * index,
      radius: config.radius + index * 6,
    }));
  };

  return [...buildNodes('subtasks'), ...buildNodes('attachments')];
});

const layerStyle = computed(() => ({
  transform:
    detailViewActive.value && selectedStar.value
      ? `translate3d(${zoomTransform.value.x}px, ${zoomTransform.value.y}px, 0) scale(${zoomTransform.value.scale})`
      : 'translate3d(0, 0, 0) scale(1)',
}));

const interactionStatus = computed(() => (pointer.active ? '划开星尘' : '静待星海'));

const interactionMeta = computed(() => `粒子${pointer.active ? '受扰动' : '平衡'}`);

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
      return `${x},${y}`;
    })
    .join(' ');

const mainPolyline = computed(() => polylineFor(dipperLayout));

const loadTasks = async () => {
  try {
    const res = await getTeamTasks();
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
    console.warn('获取团队任务失败，使用示例数据：', error);
    tasks.value = fallbackTasks;
  }
};

const ensureTaskDetail = async (task) => {
  if (!task) return null;
  try {
    const res = await getTaskDetail(task.id);
    const raw = res?.data?.data || res?.data || res;
    return normalizeDetail(raw, task);
  } catch (error) {
    console.warn('获取任务详情失败，使用示例：', error);
    return buildMockDetail(task);
  }
};

const normalizeDetail = (raw, task) => {
  if (!raw) return buildMockDetail(task);
  const mock = buildMockDetail(task);
  return {
    subtasks: Array.isArray(raw.subtasks) && raw.subtasks.length ? raw.subtasks : mock.subtasks,
    attachments: Array.isArray(raw.attachments) && raw.attachments.length ? raw.attachments : mock.attachments,
  };
};

const buildMockDetail = (task) => {
  const baseTitle = task?.title || '任务';
  return {
    subtasks: [
      { id: `${task?.id}-s1`, title: `${baseTitle} · 方案评审`, status: '进行中' },
      { id: `${task?.id}-s2`, title: `${baseTitle} · 开发排期`, status: '未开始' },
      { id: `${task?.id}-s3`, title: `${baseTitle} · 验收联调`, status: '待验收' },
    ],
    attachments: [
      { id: `${task?.id}-a1`, name: '需求说明书.pdf', size: '2.4M' },
      { id: `${task?.id}-a2`, name: '交互稿.fig', size: '4.1M' },
    ],
  };
};

const isFocusedStar = (star) => Boolean(selectedStar.value && selectedStar.value.id === star.id);

const satelliteStyle = (node) => {
  const angleRad = (node.angle * Math.PI) / 180;
  const x = Math.cos(angleRad) * node.radius;
  const y = Math.sin(angleRad) * node.radius;
  return { transform: `translate3d(${x}px, ${y}px, 0)` };
};

const satelliteLinkStyle = (node) => ({
  width: `${node.radius}px`,
  transform: `rotate(${node.angle}deg)`,
});

const focusOnStar = (star) => {
  if (!star) return;
  const { x, y } = projectPoint(star);
  const centerX = viewportWidth.value / 2;
  const centerY = viewportHeight.value / 2;
  const scale = 2.6;
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

      const buffer = 40;
      if (particle.x < -buffer) particle.x = canvas.width + buffer;
      if (particle.x > canvas.width + buffer) particle.x = -buffer;
      if (particle.y < -buffer) particle.y = canvas.height + buffer;
      if (particle.y > canvas.height + buffer) particle.y = -buffer;

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
  loadTasks();
  initParticles();
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
  pointer-events: auto;
  transform: translate3d(0, 0, 0);
}

.star-core {
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
  position: relative;
  margin-left: 28px;
  padding: 12px 16px;
  min-width: 220px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  background: rgba(6, 10, 28, 0.85);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.4);
}

.star-dialog::before {
  content: '';
  position: absolute;
  left: -60px;
  top: 50%;
  width: 52px;
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
}

.status-card {
  padding: 14px 18px;
  border-radius: 16px;
  background: rgba(8, 14, 46, 0.75);
  border: 1px solid rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(6px);
}

.status-title {
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: rgba(255, 255, 255, 0.55);
}

.status-value {
  font-size: 18px;
  font-weight: 600;
}

.status-meta,
.status-tip {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.6);
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

.satellite-link {
  position: absolute;
  top: 0;
  left: 0;
  height: 2px;
  background: linear-gradient(90deg, rgba(255, 255, 255, 0.65), rgba(86, 171, 255, 0.35));
  box-shadow: 0 0 12px rgba(126, 166, 255, 0.5);
  opacity: 0.85;
  transform-origin: 0 50%;
  pointer-events: none;
}

.satellite-node {
  position: absolute;
  display: flex;
  align-items: center;
  pointer-events: auto;
}

.satellite-core {
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
  margin-left: 18px;
  padding: 8px 12px;
  border-radius: 10px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  background: rgba(5, 9, 24, 0.92);
  min-width: 160px;
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
