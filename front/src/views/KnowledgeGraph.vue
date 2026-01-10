<template>
  <div class="knowledge-graph-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">🔗 知识图谱</h1>
        <p class="page-subtitle">可视化展示知识点之间的关联关系</p>
      </div>
      <div class="header-right">
        <button class="btn-refresh" @click="fetchGraphData">
          🔄 刷新
        </button>
      </div>
    </div>

    <!-- 图例说明 -->
    <div class="legend-bar">
      <div class="legend-title">关系类型：</div>
      <div class="legend-items">
        <span class="legend-item">
          <span class="legend-line prerequisite"></span> 前置知识
        </span>
        <span class="legend-item">
          <span class="legend-line related"></span> 相关
        </span>
        <span class="legend-item">
          <span class="legend-line extends"></span> 扩展
        </span>
        <span class="legend-item">
          <span class="legend-line same-category"></span> 同分类
        </span>
      </div>
      <div class="legend-tip">💡 点击节点查看详情，拖拽可调整布局</div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <div class="loading-spinner"></div>
      <p>正在加载知识图谱...</p>
    </div>

    <!-- 空状态 -->
    <div v-else-if="!graphData || graphData.nodes.length === 0" class="empty-state">
      <div class="empty-icon">📊</div>
      <h3>暂无知识点数据</h3>
      <p>请先在知识库中添加一些知识点</p>
      <router-link to="/knowledge-base" class="btn-primary">
        前往知识库
      </router-link>
    </div>

    <!-- 图谱容器 -->
    <div v-else class="graph-wrapper">
      <div ref="graphContainer" class="graph-container"></div>
    </div>

    <!-- 节点详情弹窗 -->
    <div v-if="selectedNode" class="node-detail-modal" @click.self="selectedNode = null">
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ selectedNode.name }}</h3>
          <button class="close-btn" @click="selectedNode = null">×</button>
        </div>
        <div class="modal-body">
          <div class="detail-row">
            <span class="label">分类：</span>
            <span class="value category-tag" :style="{ backgroundColor: selectedNode.color + '20', color: selectedNode.color }">
              {{ selectedNode.category || '未分类' }}
            </span>
          </div>
          <div class="detail-row">
            <span class="label">掌握等级：</span>
            <span class="value">{{ getLevelText(selectedNode.level) }}</span>
          </div>
          <div class="detail-row">
            <span class="label">相关连接：</span>
            <span class="value">{{ getNodeLinkCount(selectedNode.id) }} 个</span>
          </div>
        </div>
        <div class="modal-footer">
          <router-link :to="`/knowledge-base?highlight=${selectedNode.id}`" class="btn-primary">
            查看详情
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import * as echarts from 'echarts';
import { getKnowledgeGraph } from '@/api/modules/knowledge';

export default {
  name: 'KnowledgeGraph',
  data() {
    return {
      loading: false,
      graphData: null,
      chart: null,
      selectedNode: null
    };
  },
  mounted() {
    this.fetchGraphData();
    window.addEventListener('resize', this.handleResize);
  },
  beforeUnmount() {
    window.removeEventListener('resize', this.handleResize);
    if (this.chart) {
      this.chart.dispose();
    }
  },
  methods: {
    async fetchGraphData() {
      this.loading = true;
      try {
        const res = await getKnowledgeGraph();
        this.graphData = res?.data;
        
        // 使用双重 nextTick 确保 v-else 的 DOM 已经渲染完成
        this.$nextTick(() => {
          this.$nextTick(() => {
            this.renderGraph();
          });
        });
      } catch (error) {
        console.error('获取知识图谱失败:', error);
      } finally {
        this.loading = false;
      }
    },

    renderGraph() {
      if (!this.$refs.graphContainer || !this.graphData) return;

      if (this.chart) {
        this.chart.dispose();
      }

      this.chart = echarts.init(this.$refs.graphContainer);

      // 准备节点数据
      const nodes = this.graphData.nodes.map(node => ({
        id: String(node.id),
        name: node.name,
        symbolSize: Math.min(Math.max(node.value, 20), 60),
        category: node.category,
        level: node.level,
        color: node.color,
        itemStyle: {
          color: node.color
        },
        label: {
          show: true,
          fontSize: 11,
          color: '#333'
        }
      }));

      // 准备边数据
      const links = this.graphData.links.map(link => {
        const lineStyle = this.getLinkStyle(link.relation_type);
        return {
          source: String(link.source),
          target: String(link.target),
          value: link.strength,
          label: link.label,
          lineStyle: lineStyle
        };
      });

      const option = {
        tooltip: {
          trigger: 'item',
          formatter: (params) => {
            if (params.dataType === 'node') {
              return `<strong>${params.data.name}</strong><br/>分类: ${params.data.category || '未分类'}<br/>等级: ${this.getLevelText(params.data.level)}`;
            } else if (params.dataType === 'edge') {
              return `${params.data.label || '关联'}`;
            }
            return '';
          }
        },
        animationDuration: 1500,
        animationEasingUpdate: 'quinticInOut',
        series: [{
          type: 'graph',
          layout: 'force',
          data: nodes,
          links: links,
          roam: true,
          draggable: true,
          focusNodeAdjacency: true,
          force: {
            repulsion: 300,
            gravity: 0.1,
            edgeLength: [80, 200],
            layoutAnimation: true
          },
          emphasis: {
            focus: 'adjacency',
            lineStyle: {
              width: 4
            }
          },
          label: {
            show: true,
            position: 'bottom',
            distance: 5,
            fontSize: 11
          },
          lineStyle: {
            curveness: 0.2,
            opacity: 0.6
          },
          edgeLabel: {
            show: false
          }
        }]
      };

      this.chart.setOption(option);

      // 绑定点击事件
      this.chart.on('click', (params) => {
        if (params.dataType === 'node') {
          const nodeData = this.graphData.nodes.find(n => String(n.id) === params.data.id);
          if (nodeData) {
            this.selectedNode = nodeData;
          }
        }
      });
    },

    getLinkStyle(relationType) {
      const styles = {
        1: { color: '#3b82f6', width: 2, type: 'solid' },      // prerequisite
        2: { color: '#10b981', width: 1.5, type: 'dashed' },   // related
        3: { color: '#8b5cf6', width: 1.5, type: 'solid' },    // extends
        4: { color: '#ef4444', width: 2, type: 'dotted' },     // conflict
        5: { color: '#9ca3af', width: 1, type: 'dashed' },     // same_category
        6: { color: '#f59e0b', width: 1, type: 'dashed' }      // same_tag
      };
      return styles[relationType] || styles[2];
    },

    getLevelText(level) {
      const levels = ['待学习', '了解', '熟悉', '已掌握', '精通'];
      return levels[level] || '未知';
    },

    getNodeLinkCount(nodeId) {
      if (!this.graphData) return 0;
      return this.graphData.links.filter(
        link => link.source === nodeId || link.target === nodeId
      ).length;
    },

    handleResize() {
      if (this.chart) {
        this.chart.resize();
      }
    }
  }
};
</script>

<style scoped>
.knowledge-graph-container {
  padding: 24px;
  min-height: 100vh;
  background: linear-gradient(135deg, #f0f4ff 0%, #fef3f2 100%);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header-left .page-title {
  font-size: 24px;
  font-weight: 700;
  color: #1e293b;
  margin: 0;
}

.header-left .page-subtitle {
  font-size: 14px;
  color: #64748b;
  margin: 4px 0 0 0;
}

.header-right {
  display: flex;
  gap: 12px;
}

.btn-refresh {
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
  text-decoration: none;
}

.btn-refresh {
  background: #2D5BFF;
  color: white;
  border: none;
}

.btn-refresh:hover {
  background: #1e40af;
}

.legend-bar {
  background: white;
  border-radius: 12px;
  padding: 12px 20px;
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.legend-title {
  font-weight: 600;
  color: #374151;
}

.legend-items {
  display: flex;
  gap: 16px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #6b7280;
}

.legend-line {
  width: 24px;
  height: 3px;
  border-radius: 2px;
}

.legend-line.prerequisite { background: #3b82f6; }
.legend-line.related { background: #10b981; border-style: dashed; }
.legend-line.extends { background: #8b5cf6; }
.legend-line.same-category { background: #9ca3af; border-style: dashed; }

.legend-tip {
  margin-left: auto;
  font-size: 12px;
  color: #9ca3af;
}

.loading-container, .empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 60vh;
  color: #64748b;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #e2e8f0;
  border-top-color: #2D5BFF;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.empty-state .empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
}

.empty-state h3 {
  font-size: 18px;
  color: #374151;
  margin: 0 0 8px 0;
}

.empty-state p {
  margin: 0 0 20px 0;
}

.btn-primary {
  padding: 10px 24px;
  background: #2D5BFF;
  color: white;
  border-radius: 8px;
  text-decoration: none;
  font-size: 14px;
  transition: background 0.2s;
}

.btn-primary:hover {
  background: #1e40af;
}

.graph-wrapper {
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  overflow: hidden;
}

.graph-container {
  width: 100%;
  height: calc(100vh - 220px);
  min-height: 500px;
}

/* 节点详情弹窗 */
.node-detail-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.node-detail-modal .modal-content {
  background: white;
  border-radius: 16px;
  width: 360px;
  max-width: 90vw;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #f1f5f9;
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
  color: #1e293b;
}

.close-btn {
  background: none;
  border: none;
  font-size: 24px;
  color: #94a3b8;
  cursor: pointer;
  padding: 0;
  line-height: 1;
}

.close-btn:hover {
  color: #475569;
}

.modal-body {
  padding: 20px;
}

.detail-row {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
}

.detail-row:last-child {
  margin-bottom: 0;
}

.detail-row .label {
  width: 80px;
  color: #64748b;
  font-size: 14px;
}

.detail-row .value {
  font-size: 14px;
  color: #1e293b;
}

.category-tag {
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 12px;
}

.modal-footer {
  padding: 16px 20px;
  border-top: 1px solid #f1f5f9;
  display: flex;
  justify-content: flex-end;
}
</style>
