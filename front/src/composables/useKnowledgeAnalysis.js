import { ref, computed, onMounted, onUnmounted } from "vue";
import {
  analyzeUserKnowledge,
  getKnowledgeDistribution,
  getSkillRadarData,
  getLearningTrends,
} from "@/api/modules/knowledge";

// 缓存分析结果
let cachedAnalysis = null;
let cacheTimestamp = 0;
const CACHE_DURATION = 5 * 60 * 1000; // 5分钟缓存

export function useKnowledgeAnalysis() {
  const analysisReport = ref(null);
  const knowledgeDistribution = ref([]);
  const skillRadarData = ref([]);
  const learningTrends = ref([]);
  const loading = ref(false);
  let refreshInterval = null;

  // 获取分析数据（带缓存）
  const fetchAnalysis = async (forceRefresh = false) => {
    // 检查缓存
    if (!forceRefresh && cachedAnalysis && Date.now() - cacheTimestamp < CACHE_DURATION) {
      analysisReport.value = cachedAnalysis;
      return cachedAnalysis;
    }

    loading.value = true;
    try {
      const res = await analyzeUserKnowledge();
      if (res?.data) {
        cachedAnalysis = res.data;
        cacheTimestamp = Date.now();
        analysisReport.value = res.data;
        console.log("[知识库分析] 数据获取成功", res.data);
        return res.data;
      }
    } catch (error) {
      console.error("[知识库分析] 获取失败:", error);
      // 降级处理：返回空数据结构
      analysisReport.value = getDefaultAnalysisReport();
    } finally {
      loading.value = false;
    }
  };

  // 分别获取各部分数据
  const fetchDistribution = async () => {
    try {
      const res = await getKnowledgeDistribution();
      if (res?.data) {
        knowledgeDistribution.value = res.data;
        return res.data;
      }
    } catch (error) {
      console.error("[知识分布] 获取失败:", error);
    }
  };

  const fetchSkillRadar = async () => {
    try {
      const res = await getSkillRadarData();
      if (res?.data) {
        skillRadarData.value = res.data;
        return res.data;
      }
    } catch (error) {
      console.error("[技能雷达] 获取失败:", error);
    }
  };

  const fetchTrends = async () => {
    try {
      const res = await getLearningTrends();
      if (res?.data) {
        learningTrends.value = res.data;
        return res.data;
      }
    } catch (error) {
      console.error("[学习趋势] 获取失败:", error);
    }
  };

  // 获取完整分析报告
  const getAnalysisReport = async (forceRefresh = false) => {
    const report = await fetchAnalysis(forceRefresh);
    if (report) {
      knowledgeDistribution.value = report.knowledge_distribution || [];
      skillRadarData.value = report.skill_radar || [];
      learningTrends.value = report.learning_trends || [];
    }
    return report;
  };

  // 初始化（自动刷新）
  const init = (autoRefreshInterval = 30000) => {
    // 首次加载
    getAnalysisReport();

    // 设置自动刷新
    if (autoRefreshInterval > 0) {
      refreshInterval = setInterval(() => {
        console.log("[知识库] 自动刷新分析数据");
        getAnalysisReport(true);
      }, autoRefreshInterval);
    }
  };

  // 清理资源
  const cleanup = () => {
    if (refreshInterval) {
      clearInterval(refreshInterval);
    }
  };

  onMounted(() => {
    init();
    // 窗口获焦时刷新
    window.addEventListener("focus", () => {
      console.log("[知识库] 窗口获焦，刷新数据");
      getAnalysisReport(true);
    });
  });

  onUnmounted(() => {
    cleanup();
  });

  return {
    analysisReport,
    knowledgeDistribution,
    skillRadarData,
    learningTrends,
    loading,
    fetchAnalysis,
    fetchDistribution,
    fetchSkillRadar,
    fetchTrends,
    getAnalysisReport,
    init,
    cleanup,
  };
}

// 默认分析报告结构
function getDefaultAnalysisReport() {
  return {
    user_id: 0,
    generated_at: new Date().toISOString(),
    knowledge_distribution: [
      {
        category: "编程语言",
        count: 0,
        percentage: 0,
        mastered_count: 0,
        learning_count: 0,
        color: "#3b82f6",
        icon: "code",
      },
      {
        category: "数据科学",
        count: 0,
        percentage: 0,
        mastered_count: 0,
        learning_count: 0,
        color: "#8b5cf6",
        icon: "brain",
      },
    ],
    skill_radar: [
      { skill: "Python", value: 0, max_value: 100, level: "beginner", progress: 0, category: "编程语言" },
      { skill: "JavaScript", value: 0, max_value: 100, level: "beginner", progress: 0, category: "编程语言" },
      { skill: "数据分析", value: 0, max_value: 100, level: "beginner", progress: 0, category: "数据科学" },
    ],
    learning_trends: [],
    learning_insights: ["开始学习以获得分析数据"],
    recommended_topics: ["掌握核心基础概念", "实践项目应用"],
    mastered_skills_count: 0,
    learning_skills_count: 0,
    total_knowledge_points: 0,
    estimated_completion_days: 0,
  };
}
