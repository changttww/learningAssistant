/**
 * 应用配置管理
 * 统一管理环境变量和应用配置
 */

// 获取环境变量的辅助函数
const getEnvValue = (key, defaultValue = '') => {
  return import.meta.env[key] || defaultValue
}

// API 配置
export const apiConfig = {
  baseURL: getEnvValue('VITE_API_BASE_URL', 'http://localhost:3000/api'),
  timeout: parseInt(getEnvValue('VITE_API_TIMEOUT', '10000')),
  enableMock: getEnvValue('VITE_ENABLE_MOCK', 'false') === 'true'
}

// 应用配置
export const appConfig = {
  title: getEnvValue('VITE_APP_TITLE', '学习助手'),
  version: getEnvValue('VITE_APP_VERSION', '1.0.0'),
  enableDevtools: getEnvValue('VITE_ENABLE_DEVTOOLS', 'true') === 'true'
}

// 认证配置
export const authConfig = {
  tokenKey: getEnvValue('VITE_TOKEN_KEY', 'learning_assistant_token'),
  refreshTokenKey: getEnvValue('VITE_REFRESH_TOKEN_KEY', 'learning_assistant_refresh_token'),
  tokenExpireTime: parseInt(getEnvValue('VITE_TOKEN_EXPIRE_TIME', '7200000'))
}

// 上传配置
export const uploadConfig = {
  maxSize: parseInt(getEnvValue('VITE_UPLOAD_MAX_SIZE', '10485760')), // 10MB
  allowedTypes: getEnvValue('VITE_UPLOAD_ALLOWED_TYPES', 'image/jpeg,image/png,image/gif,application/pdf').split(',')
}

// 导出所有配置
export default {
  api: apiConfig,
  app: appConfig,
  auth: authConfig,
  upload: uploadConfig
}