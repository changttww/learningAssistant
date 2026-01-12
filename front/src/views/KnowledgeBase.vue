<template>
  <div class="knowledge-base-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">
          📚 我的知识库
        </h1>
        <p class="page-subtitle">共收录 {{ totalCount }} 条知识点</p>
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
          placeholder="搜索知识点..."
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
      <h3>暂无知识记录</h3>
      <p>完成任务或创建笔记后，知识将自动收录到这里</p>
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
  listUserKnowledge, 
  searchKnowledge, 
  updateKnowledgeLevel, 
  deleteKnowledgeEntry,
  getUserKnowledgeStats,
  syncUserKnowledgeBase
} from '@/api/modules/knowledge';

export default {
  name: 'KnowledgeBase',
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
      categories: [
        '计算机', '人文社科', '数理逻辑', '自然科学', '经济管理', '艺术体育'
      ],
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
      return Math.ceil(this.totalCount / this.pageSize);
    }
  },
  mounted() {
    this.fetchKnowledgeList();
    this.fetchStats();
  },
  methods: {
    async fetchKnowledgeList() {
      this.loading = true;
      try {
        console.log('[知识库] 开始获取知识列表...', {
          page: this.currentPage,
          pageSize: this.pageSize,
          category: this.filterCategory,
          level: this.filterLevel
        });
        const res = await listUserKnowledge(
          this.currentPage, 
          this.pageSize, 
          this.filterCategory, 
          this.filterLevel
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
        console.log('[知识库] 开始获取统计数据...');
        const res = await getUserKnowledgeStats();
        console.log('[知识库] 获取统计响应:', res);
        
        // 兼容多种响应格式
        if (res && (res.code === 0 || res.code === undefined)) {
          const data = res.data || res;
          this.stats = {
            mastered: data.level_3_count || 0,
            learning: data.level_2_count || 0,
            unfamiliar: data.level_1_count || 0,
            needReview: data.review_needed || 0
          };
          console.log('[知识库] 统计数据:', this.stats);
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
          console.log('[知识库] 搜索关键词:', query);
          const res = await searchKnowledge(query, 200);
          console.log('[知识库] 搜索结果:', res);
          if (res && (res.code === 0 || res.code === undefined)) {
            const data = res.data || res;
            this.knowledgeList = data.results || data.items || data || [];
            this.totalCount = data.total || this.knowledgeList.length;
            console.log('[知识库] 搜索到', this.knowledgeList.length, '条结果');
          } else {
            console.warn('[知识库] 搜索响应异常:', res);
            this.knowledgeList = [];
            this.totalCount = 0;
          }
        } catch (error) {
          console.error('[知识库] 搜索失败:', error);
          alert('搜索失败：' + (error.message || '请检查网络'));
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
        console.log('[知识库] 删除:', this.itemToDelete.id);
        const res = await deleteKnowledgeEntry(this.itemToDelete.id);
        console.log('[知识库] 删除结果:', res);
        if (res && (res.code === 0 || res.code === undefined)) {
          this.knowledgeList = this.knowledgeList.filter(k => k.id !== this.itemToDelete.id);
          this.totalCount--;
          this.closeDeleteConfirm();
          this.fetchStats();
        } else {
          alert('删除失败：' + (res?.msg || res?.error || '未知错误'));
        }
      } catch (error) {
        console.error('删除失败:', error);
        alert('删除失败：' + (error.message || '请检查网络连接'));
      }
    },

    async upgradeLevel(item) {
      if (item.level >= 3) {
        console.log('[知识库] 已达最高等级');
        return;
      }
      try {
        const newLevel = (item.level || 0) + 1;
        console.log('[知识库] 提升等级:', item.id, newLevel);
        const res = await updateKnowledgeLevel(item.id, newLevel);
        console.log('[知识库] 提升等级结果:', res);
        if (res && (res.code === 0 || res.code === undefined)) {
          item.level = newLevel;
          this.fetchStats();
        } else {
          console.error('提升等级失败:', res?.msg || res?.error);
        }
      } catch (error) {
        console.error('更新等级失败:', error);
      }
    },

    getSourceLabel(type) {
      const labels = {
        1: '任务',
        2: '笔记',
        3: '测验',
        4: '手动'
      };
      return labels[type] || '未知';
    },

    getSourceClass(type) {
      const classes = {
        1: 'source-task',
        2: 'source-note',
        3: 'source-quiz',
        4: 'source-manual'
      };
      return classes[type] || '';
    },

    getLevelLabel(level) {
      const labels = {
        0: '待学习',
        1: '了解',
        2: '熟悉',
        3: '已掌握'
      };
      return labels[level] || '待学习';
    },

    getLevelClass(level) {
      const classes = {
        0: 'level-0',
        1: 'level-1',
        2: 'level-2',
        3: 'level-3'
      };
      return classes[level] || 'level-0';
    },

    // 获取分类的显示配置（6大学科分类体系）
    getCategoryConfig(category) {
      const configs = {
        // 6大学科分类
        '计算机': { color: '#3b82f6', icon: '�', gradient: 'linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%)', bgColor: '#dbeafe' },
        '人文社科': { color: '#f59e0b', icon: '📚', gradient: 'linear-gradient(135deg, #f59e0b 0%, #d97706 100%)', bgColor: '#fef3c7' },
        '数理逻辑': { color: '#8b5cf6', icon: '🔢', gradient: 'linear-gradient(135deg, #8b5cf6 0%, #6d28d9 100%)', bgColor: '#ede9fe' },
        '自然科学': { color: '#10b981', icon: '🔬', gradient: 'linear-gradient(135deg, #10b981 0%, #059669 100%)', bgColor: '#d1fae5' },
        '经济管理': { color: '#ef4444', icon: '💰', gradient: 'linear-gradient(135deg, #ef4444 0%, #dc2626 100%)', bgColor: '#fee2e2' },
        '艺术体育': { color: '#ec4899', icon: '🎨', gradient: 'linear-gradient(135deg, #ec4899 0%, #db2777 100%)', bgColor: '#fce7f3' },
        // 其他/未分类
        '未分类': { color: '#64748b', icon: '�', gradient: 'linear-gradient(135deg, #94a3b8 0%, #64748b 100%)', bgColor: '#f1f5f9' },
        '其他': { color: '#64748b', icon: '📁', gradient: 'linear-gradient(135deg, #94a3b8 0%, #64748b 100%)', bgColor: '#f1f5f9' },
      };
      
      // 直接匹配
      if (configs[category]) {
        return configs[category];
      }
      
      // 模糊匹配（兼容旧数据）
      const lowerCat = (category || '').toLowerCase();
      const keywordMap = {
        // 计算机类
        '编程': '计算机', '代码': '计算机', '开发': '计算机', 'programming': '计算机', 'computer': '计算机',
        // 人文社科类
        '文学': '人文社科', '历史': '人文社科', '语文': '人文社科', '英语': '人文社科', '政治': '人文社科', '哲学': '人文社科',
        // 数理逻辑类
        '数学': '数理逻辑', '物理': '数理逻辑', '逻辑': '数理逻辑', 'math': '数理逻辑', 'physics': '数理逻辑',
        // 自然科学类
        '化学': '自然科学', '生物': '自然科学', '地理': '自然科学', 'chemistry': '自然科学', 'biology': '自然科学',
        // 经济管理类
        '经济': '经济管理', '金融': '经济管理', '管理': '经济管理', '会计': '经济管理',
        // 艺术体育类
        '艺术': '艺术体育', '音乐': '艺术体育', '美术': '艺术体育', '体育': '艺术体育', '运动': '艺术体育',
      };
      
      for (const [keyword, subject] of Object.entries(keywordMap)) {
        if (lowerCat.includes(keyword)) {
          return configs[subject];
        }
      }
      
      return configs['其他'];
    },

    // 获取分类的图标
    getCategoryIcon(category) {
      return this.getCategoryConfig(category).icon;
    },

    // 获取分类的颜色
    getCategoryColor(category) {
      return this.getCategoryConfig(category).color;
    },

    // 获取分类的背景色
    getCategoryBgColor(category) {
      return this.getCategoryConfig(category).bgColor;
    },

    truncateText(text, maxLength) {
      if (!text) return '';
      return text.length > maxLength ? text.substring(0, maxLength) + '...' : text;
    },

    formatDate(dateStr) {
      if (!dateStr) return '';
      const date = new Date(dateStr);
      return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit'
      });
    },

    async handleSyncKnowledgeBase() {
      if (this.syncing) return;
      this.syncing = true;
      try {
        console.log('[知识库] 开始同步知识库...');
        const res = await syncUserKnowledgeBase();
        console.log('[知识库] 同步结果:', res);
        
        if (res && (res.code === 0 || res.code === undefined)) {
          const msg = res.msg || '同步请求已提交';
          alert(`同步成功！${msg}`);
          // 刷新列表和统计
          this.fetchKnowledgeList();
          this.fetchStats();
        } else {
          alert('同步失败：' + (res?.msg || '未知错误'));
        }
      } catch (error) {
        console.error('[知识库] 同步失败:', error);
        alert('同步失败：' + (error.message || '请检查网络连接'));
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

.btn-feature.chat:hover {
  box-shadow: 0 4px 12px rgba(245, 87, 108, 0.4);
}

.btn-sync {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  background: #059669;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

.btn-sync:hover:not(:disabled) {
  background: #047857;
  transform: translateY(-1px);
}

.btn-sync:disabled {
  background: #9ca3af;
  cursor: not-allowed;
}

.btn-refresh {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  background: #4f46e5;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

.btn-refresh:hover {
  background: #4338ca;
  transform: translateY(-1px);
}

/* 统计卡片 */
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
  transition: transform 0.2s;
}

.stat-card:hover {
  transform: translateY(-2px);
}

.stat-card.mastered {
  border-left: 4px solid #10b981;
}

.stat-card.learning {
  border-left: 4px solid #3b82f6;
}

.stat-card.unfamiliar {
  border-left: 4px solid #f59e0b;
}

.stat-card.review {
  border-left: 4px solid #ef4444;
}

.stat-icon {
  font-size: 32px;
}

.stat-info {
  display: flex;
  flex-direction: column;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #1a1a2e;
}

.stat-label {
  font-size: 14px;
  color: #666;
}

/* 筛选栏 */
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
  border: 2px solid transparent;
  transition: all 0.2s;
}

.search-box:focus-within {
  border-color: #4f46e5;
  background: white;
}

.search-icon {
  font-size: 16px;
  margin-left: 8px;
}

.search-box input {
  border: none;
  background: transparent;
  padding: 10px 12px;
  font-size: 14px;
  flex: 1;
  outline: none;
}

.search-btn {
  background: #4f46e5;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s;
}

.search-btn:hover {
  background: #4338ca;
}

.clear-btn {
  background: #e5e7eb;
  border: none;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  cursor: pointer;
  margin-left: 8px;
  font-size: 12px;
  color: #6b7280;
  transition: all 0.2s;
}

.clear-btn:hover {
  background: #ef4444;
  color: white;
}

.filter-group {
  display: flex;
  gap: 12px;
}

.filter-group select {
  padding: 10px 16px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
  background: white;
  cursor: pointer;
  min-width: 120px;
  transition: all 0.2s;
}

.filter-group select:hover {
  border-color: #4f46e5;
}

.filter-group select:focus {
  outline: none;
  border-color: #4f46e5;
  box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
}

/* 搜索结果提示 */
.search-result-hint {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  border-radius: 10px;
  padding: 12px 16px;
  margin-bottom: 20px;
  color: #1e40af;
  font-size: 14px;
}

.clear-search-btn {
  background: #2563eb;
  color: white;
  border: none;
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.clear-search-btn:hover {
  background: #1d4ed8;
}

/* 知识列表 */
.knowledge-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
}

.knowledge-card {
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
  border-left: 4px solid #64748b;
}

.knowledge-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
}

.knowledge-card:hover .card-actions {
  opacity: 1;
}

.card-header {
  display: flex;
  justify-content: flex-start;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}

.knowledge-source {
  font-size: 12px;
  padding: 4px 10px;
  border-radius: 12px;
  font-weight: 500;
}

.knowledge-category-tag {
  font-size: 12px;
  padding: 4px 10px;
  border-radius: 12px;
  font-weight: 500;
  border: 1px solid;
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.source-task {
  background: #dbeafe;
  color: #2563eb;
}

.source-note {
  background: #dcfce7;
  color: #16a34a;
}

.source-quiz {
  background: #fef3c7;
  color: #d97706;
}

.source-manual {
  background: #e5e7eb;
  color: #4b5563;
}

.knowledge-level {
  font-size: 12px;
  padding: 4px 10px;
  border-radius: 12px;
  font-weight: 500;
}

.level-0 {
  background: #f3f4f6;
  color: #6b7280;
}

.level-1 {
  background: #fef3c7;
  color: #d97706;
}

.level-2 {
  background: #dbeafe;
  color: #2563eb;
}

.level-3 {
  background: #dcfce7;
  color: #16a34a;
}

.knowledge-title {
  font-size: 16px;
  font-weight: 600;
  color: #1a1a2e;
  margin: 0 0 8px 0;
  line-height: 1.4;
}

.knowledge-summary {
  font-size: 14px;
  color: #666;
  line-height: 1.6;
  margin: 0 0 16px 0;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: #9ca3af;
}

.knowledge-category {
  display: flex;
  align-items: center;
  gap: 4px;
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

.action-btn:hover:not(:disabled) {
  transform: scale(1.1);
}

.action-btn.upgrade {
  background: #dbeafe;
  color: #2563eb;
}

.action-btn.upgrade:hover:not(:disabled) {
  background: #2563eb;
  color: white;
}

.action-btn.upgrade:disabled {
  background: #dcfce7;
  color: #16a34a;
  cursor: default;
}

.action-btn.danger {
  background: #fee2e2;
  color: #ef4444;
}

.action-btn.danger:hover {
  background: #ef4444;
  color: white;
}

/* 空状态和加载 */
.empty-state,
.loading-state {
  text-align: center;
  padding: 60px 20px;
  background: white;
  border-radius: 16px;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
}

.empty-state h3 {
  font-size: 20px;
  color: #1a1a2e;
  margin: 0 0 8px 0;
}

.empty-state p {
  color: #666;
  margin: 0;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #e5e7eb;
  border-top-color: #4f46e5;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 16px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* 分页 */
.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 16px;
  margin-top: 32px;
}

.page-btn {
  padding: 10px 20px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  background: white;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

.page-btn:hover:not(:disabled) {
  background: #4f46e5;
  color: white;
  border-color: #4f46e5;
}

.page-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.page-info {
  font-size: 14px;
  color: #666;
}

/* 弹窗 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
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
  max-height: 80vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.delete-modal {
  max-width: 400px;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid #e5e7eb;
}

.modal-header h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #1a1a2e;
}

.close-btn {
  background: none;
  border: none;
  font-size: 24px;
  color: #9ca3af;
  cursor: pointer;
  line-height: 1;
}

.close-btn:hover {
  color: #1a1a2e;
}

.modal-body {
  padding: 24px;
  overflow-y: auto;
  flex: 1;
}

.detail-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid #e5e7eb;
}

.meta-item {
  font-size: 14px;
  color: #666;
  display: flex;
  align-items: center;
  gap: 4px;
}

.detail-content h4,
.detail-summary h4 {
  font-size: 14px;
  font-weight: 600;
  color: #1a1a2e;
  margin: 0 0 12px 0;
}

.content-text {
  background: #f9fafb;
  padding: 16px;
  border-radius: 8px;
  font-size: 14px;
  line-height: 1.6;
  color: #374151;
  white-space: pre-wrap;
  word-break: break-word;
  max-height: 300px;
  overflow-y: auto;
}

.detail-summary {
  margin-top: 20px;
}

.detail-summary p {
  font-size: 14px;
  color: #666;
  line-height: 1.6;
  margin: 0;
}

.warning-text {
  color: #ef4444;
  font-size: 13px;
  margin-top: 8px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-top: 1px solid #e5e7eb;
}

.btn-secondary,
.btn-primary,
.btn-danger {
  padding: 10px 20px;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-secondary {
  background: #f3f4f6;
  border: none;
  color: #374151;
}

.btn-secondary:hover {
  background: #e5e7eb;
}

.btn-primary {
  background: #4f46e5;
  border: none;
  color: white;
}

.btn-primary:hover {
  background: #4338ca;
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-danger {
  background: #ef4444;
  border: none;
  color: white;
}

.btn-danger:hover {
  background: #dc2626;
}

/* 响应式 */
@media (max-width: 768px) {
  .stats-row {
    grid-template-columns: repeat(2, 1fr);
  }

  .filter-bar {
    flex-direction: column;
    gap: 16px;
  }

  .search-box {
    max-width: 100%;
  }

  .filter-group {
    width: 100%;
    justify-content: space-between;
  }

  .filter-group select {
    flex: 1;
  }

  .knowledge-list {
    grid-template-columns: 1fr;
  }
}
</style>
