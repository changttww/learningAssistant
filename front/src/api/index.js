/**
 * API服务统一入口
 * 导出所有业务模块的API接口
 */

// 认证相关API
export * from "./modules/auth";

// 用户相关API
export * from "./modules/user";

// 任务相关API
export * from "./modules/task";

// 团队相关API
export * from "./modules/team";

// 学习相关API
export * from "./modules/study";

// 文件上传相关API
export * from "./modules/upload";

// 统计相关API
export * from "./modules/statistics";

// 通知相关API
export * from "./modules/notification";

// 学习效率分析相关API
export * from "./modules/analysis";
