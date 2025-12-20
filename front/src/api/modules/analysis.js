/**
 * 学习效率分析 API
 */

import { request } from "@/utils/request";

/**
 * 获取学习效率分析报告
 * @param {{user_id: number, days?: number, model?: string}} data
 */
export function fetchEfficiencyAnalysis(data) {
  return request.post("/analysis/efficiency", data);
}

