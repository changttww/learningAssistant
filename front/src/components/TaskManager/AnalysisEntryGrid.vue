<template>
  <div class="grid grid-cols-2 gap-4 mb-6">
    <div
      class="card p-5 cursor-pointer hover:shadow-lg transition-shadow"
      @click="$emit('show-efficiency')"
    >
      <div class="flex items-center mb-4">
        <div
          class="w-12 h-12 rounded-xl bg-gradient-to-br from-purple-500 to-pink-500 flex items-center justify-center mr-3"
        >
          <iconify-icon icon="mdi:chart-line" class="text-2xl text-white" />
        </div>
        <div>
          <h4 class="font-bold text-gray-800">学习效率分析</h4>
          <p class="text-sm text-gray-500">智能分析学习数据</p>
        </div>
      </div>
      <div class="space-y-2">
        <div class="flex justify-between text-sm">
          <span class="text-gray-600">本周学习时长</span>
          <span class="font-medium text-purple-600">
            {{ loading ? "加载中..." : displaySummary.weeklyStudyHours + "小时" }}
          </span>
        </div>
        <div class="flex justify-between text-sm">
          <span class="text-gray-600">专注度评分</span>
          <span class="font-medium text-green-600">
            {{ loading ? "加载中..." : displaySummary.focusScore + "分" }}
          </span>
        </div>
        <div class="flex justify-between text-sm">
          <span class="text-gray-600">任务完成率</span>
          <span class="font-medium text-blue-600">
            {{ loading ? "加载中..." : displaySummary.taskCompletionRate + "%" }}
          </span>
        </div>
      </div>
    </div>

    <div
      class="card p-5 cursor-pointer hover:shadow-lg transition-shadow"
      @click="$emit('show-summary')"
    >
      <div class="flex items-center mb-4">
        <div
          class="w-12 h-12 rounded-xl bg-gradient-to-br from-blue-500 to-cyan-500 flex items-center justify-center mr-3"
        >
          <iconify-icon icon="mdi:brain" class="text-2xl text-white" />
        </div>
        <div>
          <h4 class="font-bold text-gray-800">智能总结复习</h4>
          <p class="text-sm text-gray-500">AI生成复习提纲</p>
        </div>
      </div>
      <div class="space-y-2">
        <div class="flex justify-between text-sm">
          <span class="text-gray-600">待复习内容</span>
          <span class="font-medium text-orange-600">5项</span>
        </div>
        <div class="flex justify-between text-sm">
          <span class="text-gray-600">复习提醒</span>
          <span class="font-medium text-red-600">3条</span>
        </div>
        <div class="flex justify-between text-sm">
          <span class="text-gray-600">知识点掌握</span>
          <span class="font-medium text-green-600">78%</span>
        </div>
      </div>
    </div>

  </div>
</template>

<script>
  // 直接使用后端返回的结构化 JSON
  export const parseEfficiencyPayload = (response) => {
    const analysis =
      response?.data?.analysis ||
      response?.analysis ||
      response?.data?.data?.analysis;
    if (!analysis) {
      console.warn("未获取到学习效率分析数据", response);
    }
    return analysis || null;
  };

  const clampPercentage = (value) => {
    const numeric = Number(value);
    if (!Number.isFinite(numeric)) return 0;
    if (numeric < 0) return 0;
    if (numeric > 100) return 100;
    return Math.round(numeric);
  };

  const safeNumber = (val, digits = 0) => {
    const num = Number(val);
    if (!Number.isFinite(num)) return 0;
    return digits ? Number(num.toFixed(digits)) : Math.round(num);
  };

  // 统一格式化效率数据，供 TaskManager 消费
  export const transformEfficiencyData = (analysis = {}) => {
    const summary =
      analysis.summary || analysis.metrics_snapshot || analysis.metrics || {};
    const detail = analysis.analysis || analysis.detail || analysis.stats || {};
    const trend = Array.isArray(analysis.study_trend) ? analysis.study_trend : [];
    const recs = Array.isArray(analysis.recommendations)
      ? analysis.recommendations
      : [];
    const insights = [
      ...(Array.isArray(summary.key_strengths) ? summary.key_strengths : []),
      ...(Array.isArray(summary.risks) ? summary.risks : []),
      ...(Array.isArray(analysis.insights) ? analysis.insights.filter(Boolean) : []),
      ...(Array.isArray(analysis.task_insights)
        ? analysis.task_insights
            .map((item) => {
              const title = item.task_title || item.title || "任务";
              const advice =
                item.advice || item.risk || item.message || item.detail;
              if (!title && !advice) return null;
              const status = item.status ? `（${item.status}）` : "";
              return `${title}${status}: ${advice || "请优先跟进"}`;
            })
            .filter(Boolean)
        : []),
    ];

    const weeklyHoursFromSummary =
      summary.weekly_study_hours ??
      summary.weeklyStudyHours ??
      summary.weeklyStudyTime;
    const weeklyHoursFromAvgMinutes = summary.daily_average_minutes
      ? summary.daily_average_minutes * 7
      : undefined;
    const weeklyHoursFromTotalMinutes =
      detail.total_study_minutes ?? detail.total_minutes;
    const weeklyStudyTime = safeNumber(
      weeklyHoursFromSummary ??
        (weeklyHoursFromAvgMinutes ? weeklyHoursFromAvgMinutes / 60 : undefined) ??
        (weeklyHoursFromTotalMinutes ? weeklyHoursFromTotalMinutes / 60 : undefined),
      1
    );

    const focusRaw =
      summary.focus_score ??
      summary.study_consistency_score ??
      summary.focus ??
      summary.focusScore;
    const focusScore = safeNumber(
      typeof focusRaw === "number" && focusRaw > 0 && focusRaw <= 1
        ? focusRaw * 100
        : focusRaw
    );

    const completionRaw =
      summary.task_completion_rate ??
      detail.completion_rate ??
      summary.completion_rate ??
      summary.in_progress_to_done_ratio;
    const normalizedCompletion =
      typeof completionRaw === "number" && completionRaw <= 1
        ? completionRaw * 100
        : completionRaw;
    const taskCompletionRate = clampPercentage(normalizedCompletion);

    let studyTrend = trend.map((item) => safeNumber(item.study_hours, 1));
    let focusTrend = trend.map((item) => safeNumber(item.focus_score));
    let trendLabels = trend.map((item) => item.date || "");
    if (!trend.length && detail.recent_trend) {
      const avgHours = detail.recent_trend.average_daily_minutes
        ? safeNumber(detail.recent_trend.average_daily_minutes / 60, 1)
        : 0;
      const recentHours = detail.recent_trend.last_3_days_total_minutes
        ? safeNumber(detail.recent_trend.last_3_days_total_minutes / 3 / 60, 1)
        : avgHours;
      if (avgHours || recentHours) {
        studyTrend = [recentHours || avgHours];
        focusTrend = [focusScore || 80];
        trendLabels = ["最近3天"];
      }
    }

    const reviewPlanRaw =
      analysis.review_plan || analysis.reviewPlan || analysis.review || {};
    const reviewItems = Array.isArray(reviewPlanRaw.review_items)
      ? reviewPlanRaw.review_items
      : [];
    const remindersRaw = Array.isArray(reviewPlanRaw.reminders)
      ? reviewPlanRaw.reminders
      : [];

    const knowledgeMap = {
      mastered: clampPercentage(
        reviewPlanRaw.knowledge_map?.mastered ?? summary.mastered ?? 0
      ),
      learning: clampPercentage(
        reviewPlanRaw.knowledge_map?.learning ?? summary.learning ?? 0
      ),
      toLearn: clampPercentage(
        reviewPlanRaw.knowledge_map?.to_learn ?? summary.to_learn ?? 0
      ),
    };

    const mappedReviewItems = reviewItems
      .map((item) => ({
        subject: item.subject || item.title || "待复习内容",
        priority: item.priority || "medium",
        progress: clampPercentage(item.progress ?? 0),
        dueDate: item.due_date || item.due || "",
      }))
      .filter(Boolean);

    const mappedReminders = remindersRaw
      .map((item) => ({
        content: item.content || item.title || "",
        time: item.time || item.remind_at || "",
      }))
      .filter((item) => item.content);

    return {
      weeklyStudyTime,
      focusScore,
      taskCompletionRate,
      studyTrend,
      focusTrend,
      trendLabels,
      insights,
      prompt: analysis.prompt,
      reviewPlan: {
        summary: reviewPlanRaw.summary || "",
        knowledgeMap,
        reviewItems: mappedReviewItems,
        reminders: mappedReminders,
      },
      detailStats: {
        streakDays: safeNumber(detail.streak_days ?? detail.streak ?? summary.streak_days),
        totalStudyMinutes: safeNumber(
          detail.total_study_minutes ?? detail.total_minutes ?? summary.total_study_minutes
        ),
        completedTasks: safeNumber(
          detail.completed_tasks ??
            detail.tasks_completed ??
            summary.tasks_completed ??
            (summary.tasks_completed_per_day ? summary.tasks_completed_per_day * 7 : 0)
        ),
        inProgressTasks: safeNumber(
          detail.tasks_in_progress ??
            detail.task_backlog?.in_progress_count ??
            summary.in_progress_count
        ),
        completionRate: taskCompletionRate,
      },
      suggestions: recs
        .map((item) => {
          const title = item.title || "";
          const detailText =
            item.detail ||
            item.expected_outcome ||
            item.action ||
            (Array.isArray(item.action_steps) ? item.action_steps.join("；") : "");
          const description = item.description || item.message || detailText || "";
          const impact = item.expected_outcome || item.impact || item.priority || "";
          const mappedType =
            item.type === "high_impact"
              ? "warning"
              : item.priority === "urgent"
              ? "warning"
              : item.type === "positive"
              ? "positive"
              : item.type === "warning"
              ? "warning"
              : "tip";
          if (!title && !description) return null;
          return {
            type: mappedType,
            title: title || "智能建议",
            description,
            impact,
            message: title && description ? `${title}：${description}` : "",
          };
        })
        .filter(Boolean),
    };
  };

  export default {
    name: "AnalysisEntryGrid",
    props: {
      summary: {
        type: Object,
        default: () => ({
          weeklyStudyHours: 0,
          focusScore: 0,
          taskCompletionRate: 0,
        }),
      },
      loading: {
        type: Boolean,
        default: false,
      },
    },
    emits: ["show-efficiency", "show-summary"],
    computed: {
      displaySummary() {
        const summary = this.summary || {};
        const safeNumber = (value, digits = 0) => {
          if (Number.isFinite(value)) {
            return digits ? value.toFixed(digits) : Math.round(value);
          }
          return 0;
        };
        const weekly =
          summary.weeklyStudyHours ??
          summary.weekly_study_hours ??
          summary.weeklyStudyTime ??
          summary.weekly_study_time;
        const focus =
          summary.focusScore ?? summary.focus_score ?? summary.focus ?? 0;
        const completion =
          summary.taskCompletionRate ??
          summary.task_completion_rate ??
          summary.completionRate ??
          0;
        return {
          weeklyStudyHours: safeNumber(weekly, 1),
          focusScore: safeNumber(focus),
          taskCompletionRate: safeNumber(completion),
        };
      },
    },
  };
</script>
