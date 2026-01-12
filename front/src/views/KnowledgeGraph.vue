<template>
  <div class="knowledge-graph-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">🔗 知识图谱</h1>
        <p class="page-subtitle">AI驱动的知识点关联可视化 · 共 {{ realNodeCount }} 个知识点</p>
      </div>
      <div class="header-right">
        <button class="btn-mine" @click="mineAllRelations" :disabled="mining">
          {{ mining ? '⏳ 挖掘中...' : '🧠 AI挖掘关系' }}
        </button>
        <button class="btn-refresh" @click="fetchGraphData">
          🔄 刷新
        </button>
      </div>
    </div>

    <!-- 图例说明 -->
    <div class="legend-bar">
      <div class="legend-section">
        <div class="legend-title">🔗 关系类型：</div>
        <div class="legend-items">
          <span class="legend-item">
            <span class="legend-line prerequisite"></span> 前置知识
          </span>
          <span class="legend-item">
            <span class="legend-line related"></span> 相关
          </span>
          <span class="legend-item">
            <span class="legend-line extends"></span> 扩展应用
          </span>
        </div>
      </div>
      <div class="legend-section">
        <div class="legend-title">🎨 学科分类：</div>
        <div class="legend-items categories">
          <span class="category-badge" style="background: #dbeafe; color: #1d4ed8;">💻 计算机</span>
          <span class="category-badge" style="background: #fef3c7; color: #b45309;">📚 人文社科</span>
          <span class="category-badge" style="background: #ede9fe; color: #6d28d9;">🔢 数理逻辑</span>
          <span class="category-badge" style="background: #d1fae5; color: #047857;">🔬 自然科学</span>
          <span class="category-badge" style="background: #fee2e2; color: #b91c1c;">💰 经济管理</span>
          <span class="category-badge" style="background: #fce7f3; color: #be185d;">🎨 艺术体育</span>
        </div>
      </div>
      <div class="legend-tip">💡 点击节点查看详情，拖拽调整布局，滚轮缩放</div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <div class="loading-spinner"></div>
      <p>正在加载知识图谱...</p>
    </div>

    <!-- 空状态 -->
    <div v-else-if="!graphData || realNodeCount === 0" class="empty-state">
      <div class="empty-icon">📊</div>
      <h3>暂无知识点数据</h3>
      <p>请先在知识库中添加一些知识点，开始构建你的知识图谱</p>
      <router-link to="/knowledge-base" class="btn-primary">
        前往知识库
      </router-link>
    </div>

    <!-- 图谱容器 -->
    <div v-else class="graph-wrapper">
      <div ref="graphContainer" class="graph-container"></div>
      <!-- 统计信息浮层 -->
      <div class="stats-overlay">
        <div class="stat-item">
          <span class="stat-value">{{ realNodeCount }}</span>
          <span class="stat-label">知识点</span>
        </div>
        <div class="stat-item">
          <span class="stat-value">{{ logicLinkCount }}</span>
          <span class="stat-label">逻辑关系</span>
        </div>
        <div class="stat-item">
          <span class="stat-value">{{ categoryCount }}</span>
          <span class="stat-label">学科领域</span>
        </div>
      </div>
    </div>

    <!-- 节点详情弹窗 -->
    <div v-if="selectedNode && !selectedNode.is_virtual" class="node-detail-modal" @click.self="selectedNode = null">
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
            <span class="value level-badge" :class="'level-' + selectedNode.level">
              {{ getLevelText(selectedNode.level) }}
            </span>
          </div>
          <div class="detail-row">
            <span class="label">关联知识：</span>
            <span class="value">{{ getNodeLinkCount(selectedNode.id) }} 个</span>
          </div>
          <div v-if="getNodeRelations(selectedNode.id).length > 0" class="relations-section">
            <div class="relations-title">📎 关联详情：</div>
            <div class="relation-list">
              <div v-for="rel in getNodeRelations(selectedNode.id)" :key="rel.id" class="relation-item">
                <span class="relation-type" :class="'type-' + rel.type">{{ rel.typeLabel }}</span>
                <span class="relation-target">{{ rel.targetName }}</span>
              </div>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn-secondary" @click="mineRelationsForNode(selectedNode.id)">
            🧠 挖掘关系
          </button>
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
import request from '@/utils/request';

export default {
  name: 'KnowledgeGraph',
  data() {
    return {
      loading: false,
      mining: false,
      graphData: null,
      chart: null,
      selectedNode: null
    };
  },
  computed: {
    // 真实节点数量（排除虚拟中心节点）
    realNodeCount() {
      if (!this.graphData?.nodes) return 0;
      return this.graphData.nodes.filter(n => !n.is_virtual).length;
    },
    // 逻辑关系数量（排除归属关系）
    logicLinkCount() {
      if (!this.graphData?.links) return 0;
      return this.graphData.links.filter(l => l.relation_type >= 1 && l.relation_type <= 3).length;
    },
    // 涉及的学科数量
    categoryCount() {
      if (!this.graphData?.nodes) return 0;
      const categories = new Set(
        this.graphData.nodes
          .filter(n => !n.is_virtual && n.category)
          .map(n => n.category)
      );
      return categories.size;
    }
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
      const nodes = this.graphData.nodes.map(node => {
        const isVirtual = node.is_virtual;
        return {
          id: String(node.id),
          name: node.name,
          symbolSize: isVirtual ? 100 : Math.min(Math.max(node.symbol_size || 32, 24), 55),
          category: node.category,
          level: node.level,
          color: node.color,
          is_virtual: isVirtual,
          itemStyle: {
            color: node.color,
            borderColor: '#ffffff',
            borderWidth: isVirtual ? 4 : 2,
            shadowBlur: isVirtual ? 35 : 12,
            shadowColor: node.color + '50',
            shadowOffsetY: 2
          },
          label: {
            show: true,
            fontSize: isVirtual ? 14 : 12,
            fontWeight: isVirtual ? 'bold' : '500',
            fontFamily: '"Segoe UI Emoji", "Apple Color Emoji", "Noto Color Emoji", "Segoe UI", system-ui, sans-serif',
            color: isVirtual ? '#ffffff' : '#1e293b',
            position: isVirtual ? 'inside' : 'bottom',
            distance: isVirtual ? 0 : 6,
            backgroundColor: isVirtual ? 'transparent' : 'rgba(255,255,255,0.9)',
            padding: isVirtual ? 0 : [3, 6],
            borderRadius: 4
          }
        };
      });

      // 准备边数据
      const links = this.graphData.links.map(link => {
        const style = this.getLinkStyle(link.relation_type);
        return {
          source: String(link.source),
          target: String(link.target),
          value: link.strength,
          label: link.label,
          relationType: link.relation_type,
          lineStyle: style,
          // 逻辑关系显示箭头
          symbol: link.relation_type <= 3 ? ['none', 'arrow'] : ['none', 'none'],
          symbolSize: link.relation_type <= 3 ? [0, 8] : [0, 0]
        };
      });

      const option = {
        tooltip: {
          trigger: 'item',
          backgroundColor: 'rgba(255,255,255,0.95)',
          borderColor: '#e2e8f0',
          borderWidth: 1,
          padding: [12, 16],
          textStyle: {
            color: '#1e293b'
          },
          formatter: (params) => {
            if (params.dataType === 'node') {
              if (params.data.is_virtual) {
                return `<div style="font-weight:600;font-size:14px;">${params.data.name}</div><div style="color:#64748b;font-size:12px;margin-top:4px;">学科分类中心</div>`;
              }
              return `<div style="font-weight:600;font-size:14px;">${params.data.name}</div>
                <div style="color:#64748b;font-size:12px;margin-top:8px;">
                  <div>分类: ${params.data.category || '未分类'}</div>
                  <div>等级: ${this.getLevelText(params.data.level)}</div>
                </div>`;
            } else if (params.dataType === 'edge') {
              const typeLabels = { 1: '🔵 前置知识', 2: '🟢 相关', 3: '🟣 扩展应用', 6: '归属' };
              return typeLabels[params.data.relationType] || '关联';
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
            repulsion: [250, 500],
            gravity: 0.08,
            edgeLength: [80, 200],
            layoutAnimation: true,
            friction: 0.6
          },
          emphasis: {
            focus: 'adjacency',
            lineStyle: {
              width: 4
            },
            itemStyle: {
              shadowBlur: 25,
              shadowColor: 'rgba(0,0,0,0.25)'
            }
          },
          blur: {
            itemStyle: {
              opacity: 0.25
            },
            lineStyle: {
              opacity: 0.08
            }
          },
          lineStyle: {
            curveness: 0.25,
            opacity: 0.85
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
          if (nodeData && !nodeData.is_virtual) {
            this.selectedNode = nodeData;
          }
        }
      });
    },

    getLinkStyle(relationType) {
      const styles = {
        1: { color: '#2563eb', width: 2.5, type: 'solid', opacity: 0.85 },    // prerequisite - 蓝色实线
        2: { color: '#059669', width: 2, type: 'dashed', opacity: 0.75 },     // related - 绿色虚线
        3: { color: '#7c3aed', width: 2.5, type: 'solid', opacity: 0.85 },    // extends - 紫色实线
        4: { color: '#dc2626', width: 2, type: 'dotted', opacity: 0.8 },      // conflict - 红色点线
        5: { color: '#94a3b8', width: 1, type: 'dashed', opacity: 0.35 },     // same_category
        6: { color: '#cbd5e1', width: 0.8, type: 'solid', opacity: 0.2 }      // 归属中心 - 极细灰线
      };
      return styles[relationType] || styles[2];
    },

    getLevelText(level) {
      const levels = ['待学习', '了解', '熟悉', '已掌握', '精通'];
      return levels[level] || '未知';
    },

    getNodeLinkCount(nodeId) {
      if (!this.graphData) return 0;
      // 只统计逻辑关系（1-3），不统计归属关系（6）
      return this.graphData.links.filter(
        link => (link.source === nodeId || link.target === nodeId) && link.relation_type <= 3
      ).length;
    },

    getNodeRelations(nodeId) {
      if (!this.graphData) return [];
      const relations = [];
      const typeLabels = { 1: '前置', 2: '相关', 3: '扩展' };
      
      this.graphData.links.forEach(link => {
        if (link.relation_type > 3) return; // 跳过非逻辑关系
        
        let targetId = null;
        let direction = '';
        
        if (link.source === nodeId) {
          targetId = link.target;
          direction = link.relation_type === 1 ? '是其前置' : '关联';
        } else if (link.target === nodeId) {
          targetId = link.source;
          direction = link.relation_type === 1 ? '为其提供基础' : '关联';
        }
        
        if (targetId) {
          const targetNode = this.graphData.nodes.find(n => n.id === targetId);
          if (targetNode && !targetNode.is_virtual) {
            relations.push({
              id: `${link.source}-${link.target}`,
              type: link.relation_type,
              typeLabel: typeLabels[link.relation_type] || '关联',
              targetName: targetNode.name,
              direction
            });
          }
        }
      });
      
      return relations.slice(0, 5); // 最多显示5个
    },

    async mineRelationsForNode(entryId) {
      try {
        this.mining = true;
        await request.post('/knowledge-base/mine-relations', { entry_id: entryId });
        this.$nextTick(() => {
          alert('关系挖掘完成！正在刷新图谱...');
          this.fetchGraphData();
        });
      } catch (error) {
        console.error('挖掘关系失败:', error);
        alert('挖掘关系失败，请稍后重试');
      } finally {
        this.mining = false;
        this.selectedNode = null;
      }
    },

    async mineAllRelations() {
      if (!confirm('将为所有知识点挖掘AI逻辑关系，这可能需要几分钟时间。确定继续？')) {
        return;
      }
      
      try {
        this.mining = true;
        await request.post('/knowledge-base/mine-all-relations');
        alert('批量挖掘已在后台开始！请稍后刷新查看结果。');
      } catch (error) {
        console.error('批量挖掘失败:', error);
        alert('启动批量挖掘失败，请稍后重试');
      } finally {
        this.mining = false;
      }
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
  background: linear-gradient(135deg, #f0f4ff 0%, #faf5ff 50%, #f0fdf4 100%);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header-left .page-title {
  font-size: 26px;
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

.btn-refresh, .btn-mine {
  padding: 10px 20px;
  border-radius: 10px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s;
  border: none;
  font-weight: 500;
}

.btn-refresh {
  background: #ffffff;
  color: #4f46e5;
  border: 1px solid #e2e8f0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.btn-refresh:hover {
  background: #f8fafc;
  border-color: #c7d2fe;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.btn-mine {
  background: linear-gradient(135deg, #4f46e5, #7c3aed);
  color: #ffffff;
  box-shadow: 0 4px 15px rgba(79, 70, 229, 0.35);
}

.btn-mine:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(99, 102, 241, 0.4);
}

.btn-mine:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.legend-bar {
  background: white;
  border-radius: 16px;
  padding: 16px 24px;
  display: flex;
  align-items: center;
  gap: 32px;
  margin-bottom: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
  border: 1px solid #e2e8f0;
  flex-wrap: wrap;
}

.legend-section {
  display: flex;
  align-items: center;
  gap: 12px;
}

.legend-title {
  font-weight: 600;
  color: #334155;
  font-size: 13px;
  white-space: nowrap;
}

.legend-items {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.legend-items.categories {
  gap: 8px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #475569;
}

.legend-line {
  width: 28px;
  height: 3px;
  border-radius: 2px;
}

.legend-line.prerequisite { background: #2563eb; }
.legend-line.related { border: 2px dashed #059669; background: transparent; }
.legend-line.extends { background: #7c3aed; }

.category-badge {
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 500;
}

.legend-tip {
  margin-left: auto;
  font-size: 12px;
  color: #64748b;
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
  width: 50px;
  height: 50px;
  border: 3px solid #e2e8f0;
  border-top-color: #6366f1;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.empty-state .empty-icon {
  font-size: 80px;
  margin-bottom: 20px;
}

.empty-state h3 {
  font-size: 20px;
  color: #1e293b;
  margin: 0 0 8px 0;
}

.empty-state p {
  margin: 0 0 24px 0;
  color: #64748b;
}

.btn-primary {
  padding: 12px 28px;
  background: linear-gradient(135deg, #4f46e5, #7c3aed);
  color: #fff;
  border-radius: 10px;
  text-decoration: none;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s;
  box-shadow: 0 4px 15px rgba(79, 70, 229, 0.35);
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(79, 70, 229, 0.45);
}

.btn-secondary {
  padding: 10px 20px;
  background: #f1f5f9;
  color: #4f46e5;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-secondary:hover {
  background: #e2e8f0;
  border-color: #c7d2fe;
}

.graph-wrapper {
  background: #ffffff;
  border-radius: 20px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.08);
  overflow: hidden;
  position: relative;
  border: 1px solid #e2e8f0;
}

.graph-container {
  width: 100%;
  height: calc(100vh - 240px);
  min-height: 500px;
}

.stats-overlay {
  position: absolute;
  bottom: 20px;
  left: 20px;
  display: flex;
  gap: 20px;
  background: #ffffff;
  padding: 12px 20px;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #4f46e5;
  line-height: 1;
}

.stat-label {
  font-size: 11px;
  color: #64748b;
  margin-top: 4px;
}

/* 节点详情弹窗 */
.node-detail-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.node-detail-modal .modal-content {
  background: #ffffff;
  border-radius: 20px;
  width: 400px;
  max-width: 90vw;
  box-shadow: 0 25px 50px rgba(0, 0, 0, 0.15);
  border: 1px solid #e2e8f0;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid #e2e8f0;
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
  color: #1e293b;
  font-weight: 600;
}

.close-btn {
  background: none;
  border: none;
  font-size: 28px;
  color: #94a3b8;
  cursor: pointer;
  padding: 0;
  line-height: 1;
  transition: color 0.2s;
}

.close-btn:hover {
  color: #475569;
}

.modal-body {
  padding: 24px;
}

.detail-row {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
}

.detail-row:last-child {
  margin-bottom: 0;
}

.detail-row .label {
  width: 90px;
  color: #64748b;
  font-size: 14px;
}

.detail-row .value {
  font-size: 14px;
  color: #1e293b;
}

.category-tag {
  padding: 5px 12px;
  border-radius: 16px;
  font-size: 12px;
  font-weight: 500;
}

.level-badge {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
}

.level-badge.level-0 { background: #f1f5f9; color: #475569; }
.level-badge.level-1 { background: #dbeafe; color: #1d4ed8; }
.level-badge.level-2 { background: #fef3c7; color: #b45309; }
.level-badge.level-3 { background: #d1fae5; color: #047857; }
.level-badge.level-4 { background: #ede9fe; color: #6d28d9; }

.relations-section {
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px solid #e2e8f0;
}

.relations-title {
  font-size: 13px;
  color: #64748b;
  margin-bottom: 12px;
}

.relation-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.relation-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
}

.relation-type {
  padding: 2px 8px;
  border-radius: 8px;
  font-size: 11px;
  font-weight: 500;
}

.relation-type.type-1 { background: #dbeafe; color: #1d4ed8; }
.relation-type.type-2 { background: #d1fae5; color: #047857; }
.relation-type.type-3 { background: #ede9fe; color: #6d28d9; }

.relation-target {
  color: #334155;
}

.modal-footer {
  padding: 16px 24px;
  border-top: 1px solid #e2e8f0;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
