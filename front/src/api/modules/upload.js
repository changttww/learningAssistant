/**
 * 文件上传相关API
 */

import { request } from '@/utils/request'

/**
 * 上传单个文件
 */
export function uploadFile(formData, config = {}) {
  return request.upload('/upload/file', formData, {
    onUploadProgress: (progressEvent) => {
      const progress = Math.round((progressEvent.loaded * 100) / progressEvent.total)
      config.onProgress && config.onProgress(progress)
    },
    ...config
  })
}

/**
 * 上传多个文件
 */
export function uploadMultipleFiles(formData, config = {}) {
  return request.upload('/upload/files', formData, {
    onUploadProgress: (progressEvent) => {
      const progress = Math.round((progressEvent.loaded * 100) / progressEvent.total)
      config.onProgress && config.onProgress(progress)
    },
    ...config
  })
}

/**
 * 上传图片
 */
export function uploadImage(formData, config = {}) {
  return request.upload('/upload/image', formData, {
    onUploadProgress: (progressEvent) => {
      const progress = Math.round((progressEvent.loaded * 100) / progressEvent.total)
      config.onProgress && config.onProgress(progress)
    },
    ...config
  })
}

/**
 * 上传头像
 */
export function uploadAvatar(formData, config = {}) {
  return request.upload('/upload/avatar', formData, {
    onUploadProgress: (progressEvent) => {
      const progress = Math.round((progressEvent.loaded * 100) / progressEvent.total)
      config.onProgress && config.onProgress(progress)
    },
    ...config
  })
}

/**
 * 上传文档
 */
export function uploadDocument(formData, config = {}) {
  return request.upload('/upload/document', formData, {
    onUploadProgress: (progressEvent) => {
      const progress = Math.round((progressEvent.loaded * 100) / progressEvent.total)
      config.onProgress && config.onProgress(progress)
    },
    ...config
  })
}

/**
 * 获取文件信息
 */
export function getFileInfo(fileId) {
  return request.get(`/upload/files/${fileId}`)
}

/**
 * 删除文件
 */
export function deleteFile(fileId) {
  return request.delete(`/upload/files/${fileId}`)
}

/**
 * 获取文件列表
 */
export function getFileList(params = {}) {
  return request.get('/upload/files', params)
}

/**
 * 下载文件
 */
export function downloadFile(fileId, fileName) {
  return request.download(`/upload/files/${fileId}/download`, {}, {
    responseType: 'blob'
  }).then(response => {
    // 创建下载链接
    const url = window.URL.createObjectURL(new Blob([response]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', fileName)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  })
}

/**
 * 获取文件预览URL
 */
export function getFilePreviewUrl(fileId) {
  return request.get(`/upload/files/${fileId}/preview-url`)
}

/**
 * 检查文件是否存在
 */
export function checkFileExists(fileHash) {
  return request.get('/upload/check-exists', { hash: fileHash })
}

/**
 * 分片上传初始化
 */
export function initChunkUpload(data) {
  return request.post('/upload/chunk/init', data)
}

/**
 * 上传文件分片
 */
export function uploadChunk(formData, config = {}) {
  return request.upload('/upload/chunk/upload', formData, config)
}

/**
 * 合并文件分片
 */
export function mergeChunks(data) {
  return request.post('/upload/chunk/merge', data)
}