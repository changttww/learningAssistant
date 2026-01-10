<template>
  <div class="knowledge-base-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <button @click="$router.push({ name: 'TeamTasks', query: { teamId: teamId } })" class="btn-back">
          <iconify-icon icon="mdi:arrow-left"></iconify-icon>
          返回
        </button>
        <div class="title-group">
          <h1 class="page-title">
            📚 团队知识库
          </h1>
          <p class="page-subtitle">共收录 {{ totalCount }} 条知识点</p>
        </div>
      </div>
      <div class="header-right">
        <router-link to="/knowledge-graph" class="btn-feature">
          🔗 知识图谱
        </router-link>
        <router-link to="/knowledge-chat" class="btn-feature chat">
          💬 智能问答
        </router-link>
        <button class="btn-sync" @click="handleSyncKnowledgeBase" :disabled="syncing">
          {{ syncing ? '⏳ 同步中...' : '🔄 同步知识库' }}
        </button>
        <button class="btn-refresh" @click="fetchKnowledgeList">
          🔃 刷新
        </button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-row">
      <div class="stat-card mastered">
        <div class="stat-icon">🎯</div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.mastered }}</span>
          <span class="stat-label">已掌握</span>
        </div>
      </div>
      <div class="stat-card learning">
        <div class="stat-icon">📚</div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.learning }}</span>
          <span class="stat-label">学习中</span>
        </div>
      </div>
      <div class="stat-card unfamiliar">
        <div class="stat-icon">🔍</div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.unfamiliar }}</span>
          <span class="stat-label">待巩固</span>
        </div>
      </div>
      <div class="stat-card review">
        <div class="stat-icon">⏰</div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.needReview }}</span>
          <span class="stat-label">待复习</span>
        </div>
      </div>
    </div>

    <!-- 筛选栏 -->
    <div class="filter-bar">
      <div class="search-box">
        <span class="search-icon">🔍</span>
        <input 
          v-model="searchQuery" 
          type="text" 
          placeholder="搜索团队知识点..."
          @keyup.enter="handleSearch"
        />
        <button class="search-btn" @click="handleSearch">
          搜索
        </button>
        <button v-if="searchQuery" class="clear-btn" @click="searchQuery = ''; fetchKnowledgeList()">
          ✕
        </button>
      </div>
      <div class="filter-group">
        <select v-model="filterCategory" @change="currentPage = 1; fetchKnowledgeList()">
          <option value="">📁 全部分类</option>
          <option v-for="cat in categories" :key="cat" :value="cat">{{ cat }}</option>
        </select>
        <select v-model="filterLevel" @change="currentPage = 1; fetchKnowledgeList()">
          <option value="">📊 全部等级</option>
          <option value="3">✅ 已掌握</option>
          <option value="2">📖 熟悉</option>
          <option value="1">👀 了解</option>
          <option value="0">📝 待学习</option>
        </select>
      </div>
    </div>

    <!-- 搜索结果提示 -->
    <div v-if="searchQuery && !loading" class="search-result-hint">
      <span>🔍 搜索 "{{ searchQuery }}" 的结果：共 {{ totalCount }} 条</span>
      <button class="clear-search-btn" @click="searchQuery = ''; fetchKnowledgeList()">
        清除搜索
      </button>
    </div>

    <!-- 知识列表 -->
    <div class="knowledge-list" v-if="!loading && knowledgeList.length > 0">
      <div 
        class="knowledge-card" 
        v-for="item in knowledgeList" 
        :key="item.id"
        @click="showKnowledgeDetail(item)"
        :style="{ borderLeftColor: getCategoryColor(item.category) }"
      >
        <div class="card-header">
          <span class="knowledge-source" :class="getSourceClass(item.source_type)">
            {{ getSourceLabel(item.source_type) }}
          </span>
          <span 
            class="knowledge-category-tag" 
            :style="{ 
              backgroundColor: getCategoryBgColor(item.category), 
              color: getCategoryColor(item.category),
              borderColor: getCategoryColor(item.category)
            }"
          >
            {{ getCategoryIcon(item.category) }} {{ item.category || '未分类' }}
          </span>
          <span class="knowledge-level" :class="getLevelClass(item.level)">
            {{ getLevelLabel(item.level) }}
          </span>
        </div>
        <h3 class="knowledge-title">{{ item.title }}</h3>
        <p class="knowledge-summary">{{ truncateText(item.summary || item.content, 120) }}</p>
        <div class="card-footer">
          <span class="knowledge-date">
            {{ formatDate(item.created_at) }}
          </span>
        </div>
        <div class="card-actions">
          <button 
            class="action-btn upgrade" 
            @click.stop="upgradeLevel(item)" 
            :title="item.level >= 3 ? '已达最高等级' : '提升熟练度'"
            :disabled="item.level >= 3"
          >
            {{ item.level >= 3 ? '✓' : '↑' }}
          </button>
          <button class="action-btn danger" @click.stop="confirmDelete(item)" title="删除">
            ×
          </button>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div class="empty-state" v-if="!loading && knowledgeList.length === 0">
      <div class="empty-icon">📚</div>
      <h3>暂无团队知识点</h3>
      <p>团队成员协作过程中产生的知识点将汇聚于此</p>
    </div>

    <!-- 加载状态 -->
    <div class="loading-state" v-if="loading">
      <div class="loading-spinner"></div>
      <p>加载中...</p>
    </div>

    <!-- 分页 -->
    <div class="pagination" v-if="totalCount > pageSize">
      <button 
        class="page-btn" 
        :disabled="currentPage === 1"
        @click="changePage(currentPage - 1)"
      >
        上一页
      </button>
      <span class="page-info">{{ currentPage }} / {{ totalPages }}</span>
      <button 
        class="page-btn" 
        :disabled="currentPage >= totalPages"
        @click="changePage(currentPage + 1)"
      >
        下一页
      </button>
    </div>

    <!-- 详情弹窗 -->
    <div class="modal-overlay" v-if="showDetail" @click="closeDetail">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h2>{{ selectedItem?.title }}</h2>
          <button class="close-btn" @click="closeDetail">×</button>
        </div>
        <div class="modal-body" v-if="selectedItem">
          <div class="detail-meta">
            <span class="meta-item">
              <i class="icon-source"></i>
              来源: {{ getSourceLabel(selectedItem.source_type) }}
            </span>
            <span class="meta-item">
              <i class="icon-category"></i>
              分类: {{ selectedItem.category || '未分类' }}
            </span>
            <span class="meta-item">
              <i class="icon-level"></i>
              等级: {{ getLevelLabel(selectedItem.level) }}
            </span>
            <span class="meta-item">
              <i class="icon-time"></i>
              创建: {{ formatDate(selectedItem.created_at) }}
            </span>
          </div>
          <div class="detail-content">
            <h4>内容</h4>
            <div class="content-text">{{ selectedItem.content }}</div>
          </div>
          <div class="detail-summary" v-if="selectedItem.summary">
            <h4>摘要</h4>
            <p>{{ selectedItem.summary }}</p>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn-secondary" @click="closeDetail">关闭</button>
          <button 
            class="btn-primary" 
            @click="upgradeLevel(selectedItem)"
            :disabled="selectedItem?.level >= 3"
          >
            提升掌握等级
          </button>
        </div>
      </div>
    </div>

    <!-- 删除确认弹窗 -->
    <div class="modal-overlay" v-if="showDeleteConfirm" @click="closeDeleteConfirm">
      <div class="modal-content delete-modal" @click.stop>
        <div class="modal-header">
          <h2>确认删除</h2>
          <button class="close-btn" @click="closeDeleteConfirm">×</button>
        </div>
        <div class="modal-body">
          <p>确定要删除知识点 "<strong>{{ itemToDelete?.title }}</strong>" 吗？</p>
          <p class="warning-text">此操作不可恢复</p>
        </div>
        <div class="modal-footer">
          <button class="btn-secondary" @click="closeDeleteConfirm">取消</button>
          <button class="btn-danger" @click="deleteKnowledge">删除</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { 
  listTeamKnowledge, 
  searchKnowledge, 
  updateKnowledgeLevel, 
  deleteKnowledgeEntry,
  getTeamKnowledgeStats,
  syncTeamKnowledgeBase
} from '@/api/modules/knowledge';

export default {
  name: 'TeamKnowledge',
  data() {
    return {
      knowledgeList: [],
      loading: false,
      syncing: false,
      currentPage: 1,
      pageSize: 12,
      totalCount: 0,
      searchQuery: '',
      filterCategory: '',
      filterLevel: '',
      teamId: null,
      categories: ['数学', '物理', '化学', '生物', '语文', '英语', '历史', '地理', '政治', '编程', '计算机', '经济', '法律', '心理学', '艺术', '音乐', '体育', '通用', '其他'],
      stats: {
        mastered: 0,
        learning: 0,
        unfamiliar: 0,
        needReview: 0
      },
      showDetail: false,
      selectedItem: null,
      showDeleteConfirm: false,
      itemToDelete: null
    };
  },
  computed: {
    totalPages() {
      // 避免除以0的情况
      if (this.pageSize <= 0) return 1;
      return Math.ceil((this.totalCount || 0) / this.pageSize);
    }
  },
  async mounted() {
    this.teamId = this.$route.params.teamId || this.$route.query.teamId || sessionStorage.getItem("currentTeamId");
    // 安全加载，防止初次渲染出错
    try {
        await this.fetchKnowledgeList();
        await this.fetchStats();
    } catch (e) {
        console.error("Mount error:", e);
    }
  },
  methods: {
    async fetchKnowledgeList() {
      this.loading = true;
      try {
        if (!this.teamId) {
            console.warn("TeamID missing");
            this.loading = false;
            return;
        }
        console.log('[知识库] 开始获取团队知识列表...', {
          page: this.currentPage,
          pageSize: this.pageSize,
          category: this.filterCategory,
          level: this.filterLevel,
          teamId: this.teamId
        });
        const res = await listTeamKnowledge(
          this.currentPage, 
          this.pageSize, 
          this.filterCategory, 
          this.filterLevel,
          this.teamId
        );
        console.log('[知识库] 获取知识列表响应:', res);
        
        // 兼容多种响应格式
        if (res && (res.code === 0 || res.code === undefined)) {
          const data = res.data || res;
          this.knowledgeList = data.items || data || [];
          this.totalCount = data.total || this.knowledgeList.length || 0;
          console.log('[知识库] 解析完成，列表数量:', this.knowledgeList.length, '总数:', this.totalCount);
        } else {
          console.warn('[知识库] 响应code非0:', res?.code, res?.msg);
          this.knowledgeList = [];
          this.totalCount = 0;
        }
      } catch (error) {
        console.error('[知识库] 获取知识列表失败:', error);
        this.knowledgeList = [];
        this.totalCount = 0;
      } finally {
        this.loading = false;
      }
    },

    async fetchStats() {
      try {
        if (!this.teamId) return;
        console.log('[知识库] 开始获取团队统计数据...');
        const res = await getTeamKnowledgeStats(this.teamId);
        // 兼容多种响应格式
        if (res && (res.code === 0 || res.code === undefined)) {
          const data = res.data || res;
          this.stats = {
            mastered: data.level_3_count || 0,
            learning: data.level_2_count || 0,
            unfamiliar: data.level_1_count || 0,
            needReview: data.review_needed || 0
          };
        }
      } catch (error) {
        console.error('[知识库] 获取统计失败:', error);
      }
    },

    async handleSearch() {
      const query = this.searchQuery.trim();
      if (query) {
        this.loading = true;
        try {
          const res = await searchKnowledge(query, 50);
          if (res && (res.code === 0 || res.code === undefined)) {
            const data = res.data || res;
            this.knowledgeList = data.results || data.items || data || [];
            this.totalCount = data.total || this.knowledgeList.length;
          } else {
            this.knowledgeList = [];
            this.totalCount = 0;
          }
        } catch (error) {
          console.error('[知识库] 搜索失败:', error);
          this.knowledgeList = [];
          this.totalCount = 0;
        } finally {
          this.loading = false;
        }
      } else {
        // 清空搜索时重新获取列表
        this.currentPage = 1;
        this.fetchKnowledgeList();
      }
    },

    changePage(page) {
      if (page >= 1 && page <= this.totalPages) {
        this.currentPage = page;
        this.fetchKnowledgeList();
      }
    },

    showKnowledgeDetail(item) {
      this.selectedItem = item;
      this.showDetail = true;
    },

    closeDetail() {
      this.showDetail = false;
      this.selectedItem = null;
    },

    confirmDelete(item) {
      this.itemToDelete = item;
      this.showDeleteConfirm = true;
    },

    closeDeleteConfirm() {
      this.showDeleteConfirm = false;
      this.itemToDelete = null;
    },

    async deleteKnowledge() {
      if (!this.itemToDelete) return;
      try {
        const res = await deleteKnowledgeEntry(this.itemToDelete.id);
        if (res && (res.code === 0 || res.code === undefined)) {
          this.knowledgeList = this.knowledgeList.filter(k => k.id !== this.itemToDelete.id);
          this.totalCount--;
          this.closeDeleteConfirm();
          this.fetchStats();
        } else {
          alert('删除失败：' + (res?.msg || '未知错误'));
        }
      } catch (error) {
        console.error('删除失败:', error);
        alert('删除失败：' + (error.message || '请检查网络连接'));
      }
    },

    async upgradeLevel(item) {
      if (item.level >= 3) return;
      try {
        const newLevel = (item.level || 0) + 1;
        const res = await updateKnowledgeLevel(item.id, newLevel);
        if (res && (res.code === 0 || res.code === undefined)) {
          item.level = newLevel;
          this.fetchStats();
        }
      } catch (error) {
        console.error('更新等级失败:', error);
      }
    },

    getSourceLabel(type) {
      const labels = { 1: '任务', 2: '笔记', 3: '测验', 4: '手动' };
      return labels[type] || '未知';
    },

    getSourceClass(type) {
      const classes = { 1: 'source-task', 2: 'source-note', 3: 'source-quiz', 4: 'source-manual' };
      return classes[type] || '';
    },

    getLevelLabel(level) {
      const labels = { 0: '待学习', 1: '了解', 2: '熟悉', 3: '已掌握' };
      return labels[level] || '待学习';
    },

    getLevelClass(level) {
      const classes = { 0: 'level-0', 1: 'level-1', 2: 'level-2', 3: 'level-3' };
      return classes[level] || 'level-0';
    },

    getCategoryConfig(category) {
      // 简化版配置，与 Personal KB 保持一致
      const configs = {
        '数学': { color: '#3b82f6', icon: '🔢', bgColor: '#eff6ff' },
        '语文': { color: '#f59e0b', icon: '📖', bgColor: '#fffbeb' },
        '英语': { color: '#ec4899', icon: '🗣️', bgColor: '#fdf2f8' },
        '编程': { color: '#0ea5e9', icon: '💻', bgColor: '#f0f9ff' },
        '其他': { color: '#64748b', icon: '📁', bgColor: '#f1f5f9' },
      };
      
      if (configs[category]) return configs[category];
      
      // 简单的一级 fallback
      if ((category || '').includes('学')) return configs['数学'];
      if ((category || '').includes('语') || (category || '').includes('文')) return configs['语文'];
      
      return configs['其他'];
    },

    getCategoryIcon(category) { return this.getCategoryConfig(category).icon; },
    getCategoryColor(category) { return this.getCategoryConfig(category).color; },
    getCategoryBgColor(category) { return this.getCategoryConfig(category).bgColor; },

    truncateText(text, maxLength) {
      if (!text) return '';
      return text.length > maxLength ? text.substring(0, maxLength) + '...' : text;
    },

    formatDate(dateStr) {
      if (!dateStr) return '';
      const date = new Date(dateStr);
      return date.toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' });
    },

    async handleSyncKnowledgeBase() {
      if (this.syncing) return;
      if (!this.teamId) {
        alert('缺少团队ID，无法同步');
        return;
      }
      const safeTeamId = Number(this.teamId);
      if (!safeTeamId) {
        alert('团队ID无效，无法同步');
        return;
      }
      this.syncing = true;
      try {
        console.log('[团队知识库] 同步请求 teamId:', safeTeamId);
        const res = await syncTeamKnowledgeBase(safeTeamId);
        if (res && (res.code === 0 || res.code === undefined)) {
          const data = res.data || res;
          alert(`同步完成！已从团队任务构建 ${data.tasks_synced || 0} 条知识。`);
          this.fetchKnowledgeList();
          this.fetchStats();
        } else {
          alert('同步失败');
        }
      } catch (error) {
        console.error('同步失败:', error);
        alert('同步失败');
      } finally {
        this.syncing = false;
      }
    }
  }
};
</script>

<style scoped>
.knowledge-base-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 24px;
  background: #f5f7fa;
  min-height: 100vh;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.btn-back {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  color: #4b5563;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

.btn-back:hover {
  background: #f9fafb;
  color: #1a1a2e;
  border-color: #d1d5db;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: #1a1a2e;
  margin: 0;
  display: flex;
  align-items: center;
  gap: 12px;
}

.page-subtitle {
  color: #666;
  margin: 4px 0 0 0;
  font-size: 14px;
}

.header-right {
  display: flex;
  gap: 12px;
}

.btn-feature {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  text-decoration: none;
  transition: all 0.2s;
}

.btn-feature:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.btn-feature.chat {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.btn-sync {
  background: #059669;
  color: white;
  border: none;
  border-radius: 8px;
  padding: 10px 20px;
  cursor: pointer;
}
.btn-sync:hover:not(:disabled) {
  background: #047857;
}

.btn-refresh {
  background: #4f46e5;
  color: white;
  border: none;
  border-radius: 8px;
  padding: 10px 20px;
  cursor: pointer;
}
.btn-refresh:hover {
  background: #4338ca;
}

/* 统计卡片布局 */
.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-bottom: 24px;
}

.stat-card {
  background: white;
  border-radius: 16px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
}

.stat-card.mastered { border-left: 4px solid #10b981; }
.stat-card.learning { border-left: 4px solid #3b82f6; }
.stat-card.unfamiliar { border-left: 4px solid #f59e0b; }
.stat-card.review { border-left: 4px solid #ef4444; }

.stat-icon { font-size: 32px; }
.stat-info { display: flex; flex-direction: column; }
.stat-value { font-size: 28px; font-weight: 700; color: #1a1a2e; }
.stat-label { font-size: 14px; color: #666; }

/* Filter Bar */
.filter-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: white;
  padding: 16px 20px;
  border-radius: 12px;
  margin-bottom: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}
.search-box {
  display: flex;
  align-items: center;
  background: #f5f7fa;
  border-radius: 10px;
  padding: 4px 8px;
  flex: 1;
  max-width: 450px;
}
.search-box input {
  border: none;
  background: transparent;
  padding: 10px;
  outline: none;
  flex: 1;
}
.search-btn {
  background: #4f46e5;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
}
.filter-group { display: flex; gap: 12px; }
.filter-group select {
  padding: 10px 16px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
}

/* List */
.knowledge-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
}

.knowledge-card {
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.06);
  cursor: pointer;
  position: relative;
  border-left: 4px solid #64748b;
  transition: transform 0.2s;
}
.knowledge-card:hover { transform: translateY(-4px); }

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}
.knowledge-source, .knowledge-category-tag, .knowledge-level {
  font-size: 12px;
  padding: 4px 10px;
  border-radius: 12px;
}
.source-task { background: #dbeafe; color: #2563eb; }
.level-3 { background: #dcfce7; color: #16a34a; }

.knowledge-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 8px;
}
.knowledge-summary {
  font-size: 14px;
  color: #666;
  margin-bottom: 16px;
}
.card-footer {
  font-size: 12px;
  color: #9ca3af;
}

/* Modal */
.modal-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}
.modal-content {
  background: white;
  border-radius: 16px;
  width: 90%;
  max-width: 600px;
  padding: 24px;
  max-height: 80vh;
  overflow-y: auto;
}
.modal-header {
  display: flex; justify-content: space-between; align-items: center;
  margin-bottom: 20px;
}
.close-btn { background: none; border: none; font-size: 24px; cursor: pointer; }
.detail-meta { display: flex; gap: 16px; margin-bottom: 20px; border-bottom: 1px solid #eee; padding-bottom: 10px; }
.content-text { background: #f9fafb; padding: 15px; border-radius: 8px; white-space: pre-wrap; }
.modal-footer { margin-top: 20px; display: flex; justify-content: flex-end; gap: 10px; }

/* Empty State */
.empty-state {
  text-align: center;
  padding: 60px;
  background: white;
  border-radius: 16px;
}
.empty-icon { font-size: 64px; margin-bottom: 16px; }

.action-btn {
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 8px;
  background: #f3f4f6;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  font-weight: 600;
  transition: all 0.2s;
}

.card-actions {
  position: absolute;
  top: 12px;
  right: 12px;
  display: flex;
  gap: 8px;
  opacity: 0;
  transition: opacity 0.2s;
}
.knowledge-card:hover .card-actions {
  opacity: 1;
}

/* Responsive */
@media (max-width: 768px) {
  .stats-row { grid-template-columns: repeat(2, 1fr); }
  .filter-bar { flex-direction: column; }
  .knowledge-list { grid-template-columns: 1fr; }
}
</style>
