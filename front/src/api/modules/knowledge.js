/**
 * 知识库 API
 */

import { request } from "@/utils/request";

/**
 * 添加知识库条目
 */
export function addKnowledgeEntry(data) {
  return request.post("/knowledge-base/add", data);
}

/**
 * 从任务创建知识库条目
 */
export function addKnowledgeFromTask(data) {
  return request.post("/knowledge-base/add-from-task", data);
}

/**
 * 从笔记创建知识库条目
 */
export function addKnowledgeFromNote(data) {
  return request.post("/knowledge-base/add-from-note", data);
}

/**
 * 搜索知识库
 */
export function searchKnowledge(query, limit = 10) {
  console.log('[API] searchKnowledge 请求:', { q: query, limit });
  // request.get 的第二个参数直接是 params 对象，不需要再包裹一层
  return request.get("/knowledge-base/search", { q: query, limit });
}

/**
 * 获取单个知识库条目
 */
export function getKnowledgeEntry(id) {
  return request.get(`/knowledge-base/entry/${id}`);
}

/**
 * 更新知识点掌握等级
 */
export function updateKnowledgeLevel(id, level) {
  return request.put(`/knowledge-base/entry/${id}/level`, { level });
}

/**
 * 删除知识库条目
 */
export function deleteKnowledgeEntry(id) {
  return request.delete(`/knowledge-base/entry/${id}`);
}

/**
 * 获取用户知识库统计
 */
export function getUserKnowledgeStats() {
  return request.get("/knowledge-base/stats");
}

/**
 * 列表用户知识库
 * @param {number} page
 * @param {number} pageSize
 * @param {string} category
 * @param {number|string} level
 */
export function listUserKnowledge(page, pageSize, category, level) {
  // 兼容旧调用：若未传 page/pageSize，使用默认值
  const safePage = page ?? 1;
  const safePageSize = pageSize ?? 20;

  const params = { page: safePage, page_size: safePageSize };
  // 只有当 category 和 level 有值时才添加到参数中
  if (category && category !== '') {
    params.category = category;
  }
  if (level !== undefined && level !== null && level !== '') {
    params.level = level;
  }
  console.log('[API] listUserKnowledge 请求参数:', params);
  // request.get 的第二个参数直接是 params 对象，不需要再包裹一层
  return request.get("/knowledge-base/list", params);
}

/**
 * 获取完整的AI分析报告
 */
export function analyzeUserKnowledge() {
  return request.get("/knowledge-base/analysis");
}

/**
 * 获取知识点分布
 */
export function getKnowledgeDistribution() {
  return request.get("/knowledge-base/distribution");
}

/**
 * 获取技能雷达数据
 */
export function getSkillRadarData() {
  return request.get("/knowledge-base/skill-radar");
}

/**
 * 获取学习趋势
 * @param {"30"|"90"|"year"} range
 */
export function getLearningTrends(range = "30") {
  return request.get("/knowledge-base/trends", { range });
}

/**
 * 获取知识关系
 */
export function getKnowledgeRelations(id) {
  return request.get(`/knowledge-base/relations/${id}`);
}

/**
 * 创建知识关系
 */
export function createKnowledgeRelation(data) {
  return request.post("/knowledge-base/relations", data);
}

/**
 * 同步用户知识库（从任务和笔记构建）
 */
export function syncUserKnowledgeBase() {
  return request.post("/knowledge-base/sync-all");
}

/**
 * 仅同步任务到知识库
 */
export function syncTasksToKnowledge() {
  return request.post("/knowledge-base/sync-tasks");
}

/**
 * 仅同步笔记到知识库
 */
export function syncNotesToKnowledge() {
  return request.post("/knowledge-base/sync-notes");
}

/**
 * 获取知识图谱数据
 */
export function getKnowledgeGraph() {
  return request.get("/knowledge-base/graph");
}

/**
 * RAG问答（带引用溯源）
 * @param {string} query - 用户问题
 * @param {number} limit - 最大引用数量
 */
export function ragChat(query, limit = 5) {
  return request.post("/knowledge-base/chat", { query, limit });
}
