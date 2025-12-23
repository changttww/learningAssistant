/**
 * 学科分类显示配置 - 面向学习助手
 * 为每个学科提供特色化的颜色、图标和样式配置
 */

// 学科显示配置表
export const subjectConfigs = {
  // 理科类
  '数学': {
    color: '#3b82f6',
    gradient: 'linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%)',
    icon: 'mdi:calculator-variant',
    emoji: '🔢',
    lightBg: '#eff6ff',
    textColor: '#1e40af',
  },
  '物理': {
    color: '#8b5cf6',
    gradient: 'linear-gradient(135deg, #8b5cf6 0%, #6d28d9 100%)',
    icon: 'mdi:atom',
    emoji: '⚛️',
    lightBg: '#f5f3ff',
    textColor: '#5b21b6',
  },
  '化学': {
    color: '#06b6d4',
    gradient: 'linear-gradient(135deg, #06b6d4 0%, #0891b2 100%)',
    icon: 'mdi:flask-outline',
    emoji: '🧪',
    lightBg: '#ecfeff',
    textColor: '#0e7490',
  },
  '生物': {
    color: '#10b981',
    gradient: 'linear-gradient(135deg, #10b981 0%, #059669 100%)',
    icon: 'mdi:dna',
    emoji: '🧬',
    lightBg: '#ecfdf5',
    textColor: '#047857',
  },
  // 文科类
  '语文': {
    color: '#f59e0b',
    gradient: 'linear-gradient(135deg, #f59e0b 0%, #d97706 100%)',
    icon: 'mdi:book-open-page-variant',
    emoji: '📖',
    lightBg: '#fffbeb',
    textColor: '#b45309',
  },
  '英语': {
    color: '#ec4899',
    gradient: 'linear-gradient(135deg, #ec4899 0%, #db2777 100%)',
    icon: 'mdi:alphabetical',
    emoji: '🗣️',
    lightBg: '#fdf2f8',
    textColor: '#be185d',
  },
  '历史': {
    color: '#92400e',
    gradient: 'linear-gradient(135deg, #92400e 0%, #78350f 100%)',
    icon: 'mdi:castle',
    emoji: '🏛️',
    lightBg: '#fef3c7',
    textColor: '#78350f',
  },
  '地理': {
    color: '#16a34a',
    gradient: 'linear-gradient(135deg, #22c55e 0%, #16a34a 100%)',
    icon: 'mdi:earth',
    emoji: '🌍',
    lightBg: '#f0fdf4',
    textColor: '#15803d',
  },
  '政治': {
    color: '#dc2626',
    gradient: 'linear-gradient(135deg, #ef4444 0%, #dc2626 100%)',
    icon: 'mdi:bank',
    emoji: '⚖️',
    lightBg: '#fef2f2',
    textColor: '#b91c1c',
  },
  // 技能类
  '编程': {
    color: '#0ea5e9',
    gradient: 'linear-gradient(135deg, #0ea5e9 0%, #0284c7 100%)',
    icon: 'mdi:code-braces',
    emoji: '💻',
    lightBg: '#f0f9ff',
    textColor: '#0369a1',
  },
  '计算机': {
    color: '#6366f1',
    gradient: 'linear-gradient(135deg, #6366f1 0%, #4f46e5 100%)',
    icon: 'mdi:laptop',
    emoji: '🖥️',
    lightBg: '#eef2ff',
    textColor: '#4338ca',
  },
  '艺术': {
    color: '#f472b6',
    gradient: 'linear-gradient(135deg, #f472b6 0%, #ec4899 100%)',
    icon: 'mdi:palette',
    emoji: '🎨',
    lightBg: '#fdf2f8',
    textColor: '#db2777',
  },
  '音乐': {
    color: '#a855f7',
    gradient: 'linear-gradient(135deg, #a855f7 0%, #9333ea 100%)',
    icon: 'mdi:music',
    emoji: '🎵',
    lightBg: '#faf5ff',
    textColor: '#7e22ce',
  },
  '体育': {
    color: '#f97316',
    gradient: 'linear-gradient(135deg, #f97316 0%, #ea580c 100%)',
    icon: 'mdi:basketball',
    emoji: '⚽',
    lightBg: '#fff7ed',
    textColor: '#c2410c',
  },
  // 通识类
  '学习方法': {
    color: '#14b8a6',
    gradient: 'linear-gradient(135deg, #14b8a6 0%, #0d9488 100%)',
    icon: 'mdi:lightbulb-on',
    emoji: '💡',
    lightBg: '#f0fdfa',
    textColor: '#0f766e',
  },
  '考试技巧': {
    color: '#eab308',
    gradient: 'linear-gradient(135deg, #eab308 0%, #ca8a04 100%)',
    icon: 'mdi:school',
    emoji: '📝',
    lightBg: '#fefce8',
    textColor: '#a16207',
  },
  '阅读': {
    color: '#84cc16',
    gradient: 'linear-gradient(135deg, #84cc16 0%, #65a30d 100%)',
    icon: 'mdi:book-open-variant',
    emoji: '📚',
    lightBg: '#f7fee7',
    textColor: '#4d7c0f',
  },
  '写作': {
    color: '#0891b2',
    gradient: 'linear-gradient(135deg, #22d3ee 0%, #0891b2 100%)',
    icon: 'mdi:pencil',
    emoji: '✍️',
    lightBg: '#ecfeff',
    textColor: '#0e7490',
  },
  '思维训练': {
    color: '#7c3aed',
    gradient: 'linear-gradient(135deg, #a78bfa 0%, #7c3aed 100%)',
    icon: 'mdi:brain',
    emoji: '🧠',
    lightBg: '#f5f3ff',
    textColor: '#6d28d9',
  },
  '项目': {
    color: '#2563eb',
    gradient: 'linear-gradient(135deg, #3b82f6 0%, #2563eb 100%)',
    icon: 'mdi:folder-star',
    emoji: '📂',
    lightBg: '#eff6ff',
    textColor: '#1d4ed8',
  },
  '笔记': {
    color: '#ea580c',
    gradient: 'linear-gradient(135deg, #fb923c 0%, #ea580c 100%)',
    icon: 'mdi:notebook-outline',
    emoji: '📓',
    lightBg: '#fff7ed',
    textColor: '#c2410c',
  },
  // 默认
  '其他': {
    color: '#64748b',
    gradient: 'linear-gradient(135deg, #94a3b8 0%, #64748b 100%)',
    icon: 'mdi:bookshelf',
    emoji: '📁',
    lightBg: '#f1f5f9',
    textColor: '#475569',
  },
  '默认': {
    color: '#64748b',
    gradient: 'linear-gradient(135deg, #94a3b8 0%, #64748b 100%)',
    icon: 'mdi:book',
    emoji: '📖',
    lightBg: '#f1f5f9',
    textColor: '#475569',
  },
};

// 关键词到学科的映射
const keywordMapping = {
  // 数学
  'math': '数学', '数学': '数学', '代数': '数学', '几何': '数学', '微积分': '数学',
  '函数': '数学', '方程': '数学', '公式': '数学', '计算': '数学', '概率': '数学', '统计': '数学',
  // 物理
  'physics': '物理', '物理': '物理', '力学': '物理', '电学': '物理', '磁学': '物理',
  '光学': '物理', '热学': '物理', '能量': '物理', '牛顿': '物理', '运动': '物理',
  // 化学
  'chemistry': '化学', '化学': '化学', '元素': '化学', '分子': '化学', '原子': '化学',
  '反应': '化学', '酸碱': '化学', '有机': '化学',
  // 生物
  'biology': '生物', '生物': '生物', '细胞': '生物', '遗传': '生物', '基因': '生物',
  '生命': '生物', '生态': '生物', '动物': '生物', '植物': '生物',
  // 语文
  'chinese': '语文', '语文': '语文', '文言文': '语文', '古诗': '语文', '阅读理解': '语文',
  '写作': '写作', '作文': '语文', '文学': '语文', '诗词': '语文',
  // 英语
  'english': '英语', '英语': '英语', '单词': '英语', '语法': '英语', '词汇': '英语',
  '口语': '英语', '听力': '英语', '翻译': '英语',
  // 历史
  'history': '历史', '历史': '历史', '朝代': '历史', '战争': '历史', '革命': '历史',
  '古代': '历史', '近代': '历史', '历史事件': '历史',
  // 地理
  'geography': '地理', '地理': '地理', '气候': '地理', '地形': '地理', '区域': '地理',
  '城市': '地理', '自然': '地理', '环境': '地理',
  // 政治
  'politics': '政治', '政治': '政治', '政策': '政治', '制度': '政治', '法律': '政治',
  '哲学': '政治', '思想': '政治',
  // 编程
  'programming': '编程', '编程': '编程', '代码': '编程', 'python': '编程', 'java': '编程',
  'javascript': '编程', '程序': '编程', '开发': '编程', '算法': '编程',
  // 计算机
  'computer': '计算机', '计算机': '计算机', '软件': '计算机', '硬件': '计算机',
  '网络': '计算机', '系统': '计算机', '数据库': '计算机',
  // 艺术
  'art': '艺术', '艺术': '艺术', '绘画': '艺术', '美术': '艺术', '设计': '艺术',
  '色彩': '艺术', '创作': '艺术',
  // 音乐
  'music': '音乐', '音乐': '音乐', '歌曲': '音乐', '乐器': '音乐', '旋律': '音乐',
  '节奏': '音乐', '音符': '音乐',
  // 体育
  'sports': '体育', '体育': '体育', '锻炼': '体育', '健身': '体育', '比赛': '体育', '训练': '体育',
  // 通识
  '学习方法': '学习方法', '学习技巧': '学习方法', '记忆': '学习方法', '复习': '学习方法',
  '笔记': '笔记', '思维导图': '学习方法', '效率': '学习方法',
  '考试': '考试技巧', 'exam': '考试技巧', '测验': '考试技巧', '答题': '考试技巧', '解题': '考试技巧',
  '阅读': '阅读', 'reading': '阅读', '书籍': '阅读', '文章': '阅读',
  '思维': '思维训练', '逻辑': '思维训练', '推理': '思维训练', '思考': '思维训练', '分析': '思维训练',
  '项目': '项目',
};

/**
 * 根据分类获取显示配置
 * @param {string} category - 分类名称
 * @returns {Object} 显示配置对象
 */
export function getSubjectConfig(category) {
  // 直接匹配
  if (subjectConfigs[category]) {
    return subjectConfigs[category];
  }
  
  // 模糊匹配
  const lowerCat = (category || '').toLowerCase();
  for (const [keyword, subject] of Object.entries(keywordMapping)) {
    if (lowerCat.includes(keyword.toLowerCase())) {
      return subjectConfigs[subject] || subjectConfigs['默认'];
    }
  }
  
  return subjectConfigs['默认'];
}

/**
 * 获取分类的颜色
 * @param {string} category - 分类名称
 * @returns {string} 颜色值
 */
export function getCategoryColor(category) {
  return getSubjectConfig(category).color;
}

/**
 * 获取分类的渐变色
 * @param {string} category - 分类名称
 * @returns {string} CSS渐变值
 */
export function getCategoryGradient(category) {
  return getSubjectConfig(category).gradient;
}

/**
 * 获取分类的图标
 * @param {string} category - 分类名称
 * @returns {string} Iconify图标名
 */
export function getCategoryIcon(category) {
  return getSubjectConfig(category).icon;
}

/**
 * 获取分类的Emoji
 * @param {string} category - 分类名称
 * @returns {string} Emoji字符
 */
export function getCategoryEmoji(category) {
  return getSubjectConfig(category).emoji;
}

/**
 * 获取分类的浅色背景
 * @param {string} category - 分类名称
 * @returns {string} 背景色值
 */
export function getCategoryLightBg(category) {
  return getSubjectConfig(category).lightBg;
}

/**
 * 获取所有可用的分类列表
 * @returns {Array<string>} 分类名称数组
 */
export function getAllCategories() {
  return Object.keys(subjectConfigs).filter(k => k !== '默认');
}

/**
 * ECharts饼图的渐变色配置生成
 * @param {Array} data - 分布数据
 * @returns {Array} 带渐变色的饼图数据
 */
export function generatePieChartData(data) {
  const gradientColors = [
    ['#3b82f6', '#1d4ed8'], // 蓝色
    ['#f59e0b', '#d97706'], // 橙色
    ['#ec4899', '#db2777'], // 粉色
    ['#8b5cf6', '#6d28d9'], // 紫色
    ['#10b981', '#059669'], // 绿色
    ['#06b6d4', '#0891b2'], // 青色
    ['#f97316', '#ea580c'], // 橘色
    ['#84cc16', '#65a30d'], // 酸橙
  ];

  return data.map((item, index) => {
    const config = getSubjectConfig(item.category);
    const colorPair = gradientColors[index % gradientColors.length];
    
    return {
      value: item.count || 0,
      name: item.category,
      itemStyle: {
        color: {
          type: 'linear',
          x: 0, y: 0, x2: 1, y2: 1,
          colorStops: [
            { offset: 0, color: config.color || colorPair[0] },
            { offset: 1, color: colorPair[1] || config.color }
          ]
        }
      }
    };
  });
}

export default {
  subjectConfigs,
  getSubjectConfig,
  getCategoryColor,
  getCategoryGradient,
  getCategoryIcon,
  getCategoryEmoji,
  getCategoryLightBg,
  getAllCategories,
  generatePieChartData,
};
