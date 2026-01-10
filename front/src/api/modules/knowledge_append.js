
/**
 * 获取团队知识库统计
 * @param {string|number} teamId
 */
export function getTeamKnowledgeStats(teamId) {
  return request.get("/knowledge-base/team/stats", { team_id: teamId });
}

/**
 * 列表团队知识库
 * @param {number} page
 * @param {number} pageSize
 * @param {string} category
 * @param {number|string} level
 * @param {string|number} teamId
 */
export function listTeamKnowledge(page, pageSize, category, level, teamId) {
  const safePage = page ?? 1;
  const safePageSize = pageSize ?? 20;

  const params = { page: safePage, page_size: safePageSize, team_id: teamId };
  if (category && category !== '') {
    params.category = category;
  }
  if (level !== undefined && level !== null && level !== '') {
    params.level = level;
  }
  return request.get("/knowledge-base/team/list", params);
}
