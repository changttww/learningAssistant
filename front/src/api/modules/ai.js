/**
 * AI 功能相关 API
 */

import { request } from "@/utils/request";

/**
 * 使用 AI 解析自然语言任务
 * @param {string} input - 自然语言任务描述
 */
export function parseTaskWithAI(input) {
  return request.post("/ai/parse-task", { input });
}

/**
 * 获取任务学习指导
 * @param {Object} data - 任务信息
 * @param {string} data.title - 任务标题
 * @param {string} data.description - 任务描述
 * @param {string} data.category - 任务分类
 */
export function getTaskGuidance(data) {
  return request.post("/ai/task-guidance", data);
}

/**
 * 生成智能测验
 * @param {Object} data - 测验配置
 * @param {string} data.topic - 测验主题
 * @param {string} data.content - 学习内容/笔记（可选）
 * @param {string} data.difficulty - 难度：easy/medium/hard
 * @param {number} data.quizCount - 题目数量
 * @param {boolean} data.includeEssay - 是否包含简答题
 */
export function generateQuiz(data) {
  return request.post("/ai/generate-quiz", data);
}

/**
 * 通用聊天
 * @param {Object} data
 * @param {string} data.message
 * @param {Array} data.history
 */
export function chatWithAI(data) {
  return request.post("/ai/chat", data);
}

/**
 * 生成学习计划
 * @param {Object} data
 */
export function generateStudyPlan(data) {
  return request.post("/ai/study-plan", data);
}

/**
 * 生成房间创意
 * @param {Object} data
 * @param {string} data.prompt
 */
export function generateRoomIdea(data) {
  return request.post("/ai/room-idea", data);
}

/**
 * 提交测验答案并加入知识库
 * @param {Object} data - 测验结果
 * @param {number} data.task_id - 任务ID
 * @param {string} data.topic - 测验主题
 * @param {Array} data.questions - 题目列表
 * @param {Object} data.answers - 用户答案
 * @param {number} data.score - 得分
 */
export function submitQuizToKnowledge(data) {
  return request.post("/ai/submit-quiz", data);
}

/**
 * 智能笔记增强
 * @param {Object} data - 笔记增强配置
 * @param {number} data.note_id - 笔记ID（可选）
 * @param {string} data.content - 笔记内容
 * @param {string} data.title - 笔记标题
 * @param {string} data.type - 增强类型：all/summary/keywords/mindmap/questions/polish
 */
export function enhanceNote(data) {
  return request.post("/notes/enhance", data);
}

/**
 * 生成笔记摘要
 * @param {Object} data - 笔记内容
 */
export function generateNoteSummary(data) {
  return request.post("/notes/generate-summary", data);
}

/**
 * 提取笔记关键词
 * @param {Object} data - 笔记内容
 */
export function extractNoteKeywords(data) {
  return request.post("/notes/extract-keywords", data);
}

/**
 * 生成思维导图
 * @param {Object} data - 笔记内容
 */
export function generateNoteMindmap(data) {
  return request.post("/notes/generate-mindmap", data);
}

/**
 * 生成复习问题
 * @param {Object} data - 笔记内容
 */
export function generateNoteQuestions(data) {
  return request.post("/notes/generate-questions", data);
}

/**
 * 润色笔记
 * @param {Object} data - 笔记内容
 */
export function polishNote(data) {
  return request.post("/notes/polish", data);
}

export default {
  parseTaskWithAI,
  getTaskGuidance,
  generateQuiz,
  chatWithAI,
  generateStudyPlan,
  generateRoomIdea,
  submitQuizToKnowledge,
  enhanceNote,
  generateNoteSummary,
  extractNoteKeywords,
  generateNoteMindmap,
  generateNoteQuestions,
  polishNote,
};
