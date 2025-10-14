/**
 * 统计相关API
 */

import { request } from '@/utils/request'

/**
 * 获取用户学习统计概览
 */
export function getUserStatisticsOverview(userId) {
  return request.get(`/statistics/users/${userId}/overview`)
}

/**
 * 获取用户学习时长统计
 */
export function getUserStudyTimeStats(userId, params = {}) {
  return request.get(`/statistics/users/${userId}/study-time`, params)
}

/**
 * 获取用户任务完成统计
 */
export function getUserTaskCompletionStats(userId, params = {}) {
  return request.get(`/statistics/users/${userId}/task-completion`, params)
}

/**
 * 获取用户技能雷达图数据
 */
export function getUserSkillRadarData(userId) {
  return request.get(`/statistics/users/${userId}/skill-radar`)
}

/**
 * 获取团队统计概览
 */
export function getTeamStatisticsOverview(teamId) {
  return request.get(`/statistics/teams/${teamId}/overview`)
}

/**
 * 获取团队成员活跃度统计
 */
export function getTeamMemberActivityStats(teamId, params = {}) {
  return request.get(`/statistics/teams/${teamId}/member-activity`, params)
}

/**
 * 获取团队任务进度统计
 */
export function getTeamTaskProgressStats(teamId) {
  return request.get(`/statistics/teams/${teamId}/task-progress`)
}

/**
 * 获取团队积分排行榜
 */
export function getTeamPointsRanking(teamId, params = {}) {
  return request.get(`/statistics/teams/${teamId}/points-ranking`, params)
}

/**
 * 获取全局统计数据
 */
export function getGlobalStatistics() {
  return request.get('/statistics/global')
}

/**
 * 获取学习趋势数据
 */
export function getStudyTrendData(params = {}) {
  return request.get('/statistics/study-trend', params)
}

/**
 * 获取知识点掌握情况
 */
export function getKnowledgePointMastery(params = {}) {
  return request.get('/statistics/knowledge-mastery', params)
}

/**
 * 获取学习效率分析
 */
export function getStudyEfficiencyAnalysis(params = {}) {
  return request.get('/statistics/study-efficiency', params)
}

/**
 * 获取学习目标达成情况
 */
export function getStudyGoalAchievement(params = {}) {
  return request.get('/statistics/goal-achievement', params)
}

/**
 * 获取学习习惯分析
 */
export function getStudyHabitAnalysis(params = {}) {
  return request.get('/statistics/study-habits', params)
}

/**
 * 获取学习报告
 */
export function getStudyReport(params = {}) {
  return request.get('/statistics/study-report', params)
}

/**
 * 导出统计数据
 */
export function exportStatisticsData(params = {}) {
  return request.download('/statistics/export', params)
}

/**
 * 获取实时统计数据
 */
export function getRealTimeStatistics() {
  return request.get('/statistics/real-time')
}

/**
 * 获取历史统计数据对比
 */
export function getHistoricalComparison(params = {}) {
  return request.get('/statistics/historical-comparison', params)
}