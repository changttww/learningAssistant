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

export default {
  parseTaskWithAI,
  getTaskGuidance,
  generateQuiz,
  submitQuizToKnowledge,
};
