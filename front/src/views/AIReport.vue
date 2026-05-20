<template>
  <div class="ai-report-page">
    <section class="report-hero">
      <div>
        <div class="eyebrow">AI 学习分析</div>
        <h1>AI 学习报告</h1>
        <p>根据任务完成情况、知识点分布和学习趋势，生成本周学习表现分析与优化建议。</p>
      </div>
      <div class="hero-actions">
        <select v-model.number="days" :disabled="loading" @change="loadReport">
          <option :value="7">最近 7 天</option>
          <option :value="14">最近 14 天</option>
          <option :value="30">最近 30 天</option>
        </select>
        <button :disabled="loading" @click="loadReport">
          <iconify-icon icon="mdi:refresh" width="16" height="16"></iconify-icon>
          {{ loading ? "分析中..." : "重新生成" }}
        </button>
      </div>
    </section>

    <div v-if="error" class="state-card error-card">
      <iconify-icon icon="mdi:alert-circle-outline" width="24" height="24"></iconify-icon>
      <div>
        <strong>报告生成失败</strong>
        <p>{{ error }}</p>
      </div>
    </div>

    <div v-else-if="loading" class="state-card">
      <iconify-icon icon="mdi:loading" width="24" height="24"></iconify-icon>
      <div>
        <strong>AI 正在分析学习数据</strong>
        <p>正在整理学习时长、任务完成率、趋势与复习建议。</p>
      </div>
    </div>

    <template v-else-if="report">
      <section class="metrics-grid">
        <div class="metric-card">
          <span>学习时长</span>
          <strong>{{ summary.weekly_study_hours || 0 }} 小时</strong>
        </div>
        <div class="metric-card">
          <span>专注评分</span>
          <strong>{{ summary.focus_score || 0 }} 分</strong>
        </div>
        <div class="metric-card">
          <span>任务完成率</span>
          <strong>{{ summary.task_completion_rate || 0 }}%</strong>
        </div>
        <div class="metric-card">
          <span>连续学习</span>
          <strong>{{ summary.streak_days || 0 }} 天</strong>
        </div>
      </section>

      <section class="content-grid">
        <div class="panel">
          <div class="panel-head">
            <h2>AI 建议</h2>
            <span>{{ report.model || "AI" }}</span>
          </div>
          <div v-if="recommendations.length" class="item-list">
            <article v-for="(item, index) in recommendations" :key="index" class="advice-item">
              <div class="item-title">{{ item.title || "学习建议" }}</div>
              <p>{{ item.detail }}</p>
              <span class="impact" :class="item.impact">{{ impactLabel(item.impact) }}</span>
            </article>
          </div>
          <p v-else class="empty-text">暂无 AI 建议。</p>
        </div>

        <div class="panel">
          <div class="panel-head">
            <h2>任务洞察</h2>
          </div>
          <div v-if="taskInsights.length" class="item-list">
            <article v-for="(item, index) in taskInsights" :key="index" class="task-item">
              <div class="item-title">{{ item.task_title || "任务" }}</div>
              <p>{{ item.advice || item.risk }}</p>
              <span>{{ item.status }}</span>
            </article>
          </div>
          <p v-else class="empty-text">暂无任务洞察。</p>
        </div>
      </section>

      <section class="content-grid">
        <div class="panel">
          <div class="panel-head">
            <h2>学习趋势</h2>
          </div>
          <div v-if="studyTrend.length" class="trend-list">
            <div v-for="item in studyTrend" :key="item.date" class="trend-row">
              <span>{{ item.date }}</span>
              <div class="trend-bar">
                <i :style="{ width: trendWidth(item.study_hours) }"></i>
              </div>
              <strong>{{ item.study_hours }}h</strong>
            </div>
          </div>
          <p v-else class="empty-text">暂无趋势数据。</p>
        </div>

        <div class="panel">
          <div class="panel-head">
            <h2>复习计划</h2>
          </div>
          <p v-if="reviewPlan.summary" class="review-summary">{{ reviewPlan.summary }}</p>
          <div class="knowledge-map">
            <div>
              <span>已掌握</span>
              <strong>{{ knowledgeMap.mastered || 0 }}%</strong>
            </div>
            <div>
              <span>学习中</span>
              <strong>{{ knowledgeMap.learning || 0 }}%</strong>
            </div>
            <div>
              <span>待学习</span>
              <strong>{{ knowledgeMap.to_learn || 0 }}%</strong>
            </div>
          </div>
          <div v-if="reviewItems.length" class="review-items">
            <div v-for="(item, index) in reviewItems" :key="index">
              <span>{{ item.subject }}</span>
              <em>{{ priorityLabel(item.priority) }}</em>
            </div>
          </div>
        </div>
      </section>
    </template>
  </div>
</template>

<script>
import { fetchEfficiencyAnalysis } from "@/api/modules/analysis";
import { DEFAULT_USER_ID, useCurrentUser } from "@/composables/useCurrentUser";

export default {
  name: "AIReport",
  setup() {
    const { profile, loadCurrentUser } = useCurrentUser();
    return { profile, loadCurrentUser };
  },
  data() {
    return {
      days: 7,
      loading: false,
      error: "",
      report: null,
    };
  },
  computed: {
    summary() {
      return this.report?.summary || {};
    },
    recommendations() {
      return Array.isArray(this.report?.recommendations) ? this.report.recommendations : [];
    },
    taskInsights() {
      return Array.isArray(this.report?.task_insights) ? this.report.task_insights : [];
    },
    studyTrend() {
      return Array.isArray(this.report?.study_trend) ? this.report.study_trend : [];
    },
    reviewPlan() {
      return this.report?.review_plan || {};
    },
    knowledgeMap() {
      return this.reviewPlan.knowledge_map || {};
    },
    reviewItems() {
      return Array.isArray(this.reviewPlan.review_items) ? this.reviewPlan.review_items : [];
    },
  },
  mounted() {
    this.loadReport();
  },
  methods: {
    async loadReport() {
      this.loading = true;
      this.error = "";
      try {
        let userId = this.profile?.id || DEFAULT_USER_ID;
        if (!this.profile?.id) {
          const loaded = await this.loadCurrentUser().catch(() => null);
          userId = loaded?.id || DEFAULT_USER_ID;
        }
        const res = await fetchEfficiencyAnalysis({
          user_id: userId,
          days: this.days,
          force_refresh: true,
        });
        const analysis = res?.data?.analysis || res?.analysis || res?.data?.data?.analysis;
        if (!analysis) {
          throw new Error("接口未返回分析数据");
        }
        this.report = analysis;
      } catch (error) {
        this.report = null;
        this.error = error?.message || "网络连接失败或 AI 分析接口不可用";
      } finally {
        this.loading = false;
      }
    },
    trendWidth(hours) {
      const max = Math.max(...this.studyTrend.map((item) => Number(item.study_hours) || 0), 1);
      const value = Number(hours) || 0;
      return `${Math.max(6, Math.round((value / max) * 100))}%`;
    },
    impactLabel(impact) {
      return { high: "高影响", medium: "中影响", low: "低影响" }[impact] || "建议";
    },
    priorityLabel(priority) {
      return { high: "高优先级", medium: "中优先级", low: "低优先级" }[priority] || "中优先级";
    },
  },
};
</script>

<style scoped>
.ai-report-page {
  width: 100%;
  padding: 24px 0 40px;
}

.report-hero,
.state-card,
.panel,
.metric-card {
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  background: #ffffff;
}

.report-hero {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 20px;
  padding: 26px;
  background: linear-gradient(135deg, rgba(45, 91, 255, 0.08), rgba(16, 185, 129, 0.08)), #ffffff;
}

.eyebrow {
  color: #2d5bff;
  font-size: 13px;
  font-weight: 700;
}

h1,
h2,
p {
  margin: 0;
}

h1 {
  margin: 8px 0;
  color: #111827;
  font-size: 28px;
  font-weight: 800;
}

p {
  color: #4b5563;
  line-height: 1.7;
}

.hero-actions {
  display: flex;
  gap: 10px;
}

select,
button {
  min-height: 38px;
  border-radius: 10px;
  border: 1px solid #d1d5db;
  padding: 0 12px;
  font-size: 13px;
}

button {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: #ffffff;
  border-color: #2d5bff;
  background: #2d5bff;
  font-weight: 600;
}

button:disabled,
select:disabled {
  opacity: 0.65;
}

.state-card {
  margin-top: 18px;
  padding: 18px;
  display: flex;
  gap: 12px;
  align-items: flex-start;
}

.error-card {
  border-color: #fecaca;
  background: #fff7f7;
  color: #dc2626;
}

.metrics-grid {
  margin-top: 18px;
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
}

.metric-card {
  padding: 18px;
}

.metric-card span,
.panel-head span,
.task-item span,
.knowledge-map span {
  color: #6b7280;
  font-size: 13px;
}

.metric-card strong {
  display: block;
  margin-top: 8px;
  color: #111827;
  font-size: 24px;
}

.content-grid {
  margin-top: 18px;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 18px;
}

.panel {
  padding: 20px;
}

.panel-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 14px;
}

h2 {
  color: #111827;
  font-size: 18px;
  font-weight: 800;
}

.item-list {
  display: grid;
  gap: 12px;
}

.advice-item,
.task-item {
  padding: 14px;
  border-radius: 10px;
  background: #f9fafb;
}

.item-title {
  color: #111827;
  font-weight: 700;
  margin-bottom: 6px;
}

.impact {
  display: inline-block;
  margin-top: 10px;
  padding: 3px 8px;
  border-radius: 999px;
  color: #1d4ed8;
  background: #dbeafe;
  font-size: 12px;
}

.impact.high {
  color: #b91c1c;
  background: #fee2e2;
}

.impact.low {
  color: #047857;
  background: #d1fae5;
}

.trend-list {
  display: grid;
  gap: 10px;
}

.trend-row {
  display: grid;
  grid-template-columns: 86px 1fr 52px;
  gap: 10px;
  align-items: center;
  color: #374151;
  font-size: 13px;
}

.trend-bar {
  height: 9px;
  overflow: hidden;
  border-radius: 999px;
  background: #eef2ff;
}

.trend-bar i {
  display: block;
  height: 100%;
  border-radius: inherit;
  background: #2d5bff;
}

.knowledge-map {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
  margin-top: 14px;
}

.knowledge-map div {
  padding: 12px;
  border-radius: 10px;
  background: #f9fafb;
}

.knowledge-map strong {
  display: block;
  margin-top: 6px;
  color: #111827;
}

.review-items {
  display: grid;
  gap: 8px;
  margin-top: 14px;
}

.review-items div {
  display: flex;
  justify-content: space-between;
  gap: 10px;
  padding: 10px 0;
  border-top: 1px solid #f3f4f6;
}

.review-items em {
  color: #2d5bff;
  font-style: normal;
  font-size: 13px;
  white-space: nowrap;
}

.empty-text {
  color: #9ca3af;
}

@media (max-width: 900px) {
  .metrics-grid,
  .content-grid {
    grid-template-columns: 1fr;
  }

  .report-hero {
    flex-direction: column;
  }

  .hero-actions {
    width: 100%;
    flex-wrap: wrap;
  }
}
</style>
