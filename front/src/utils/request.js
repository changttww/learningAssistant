/**
 * HTTP请求工具
 * 基于axios封装，提供统一的请求拦截、响应处理、错误处理
 */

import axios from "axios";
import { apiConfig, authConfig } from "@/config";
import { getToken, removeToken } from "@/utils/auth";
import { ElMessage } from "element-plus";

// 创建axios实例
const service = axios.create({
  baseURL: apiConfig.baseURL,
  timeout: apiConfig.timeout,
  headers: {
    "Content-Type": "application/json;charset=UTF-8",
  },
});

// 请求拦截器
service.interceptors.request.use(
  (config) => {
    // 添加认证token
    const token = getToken();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    // 添加请求时间戳（防止缓存）
    if (config.method === "get") {
      config.params = {
        ...config.params,
        _t: Date.now(),
      };
    }

    // 打印请求信息（开发环境）
    if (import.meta.env.DEV) {
      console.log(`🚀 [${config.method?.toUpperCase()}] ${config.url}`, config);
    }

    return config;
  },
  (error) => {
    console.error("请求拦截器错误:", error);
    return Promise.reject(error);
  }
);

// 响应拦截器
service.interceptors.response.use(
  (response) => {
    const { data, status } = response;

    // 打印响应信息（开发环境）
    if (import.meta.env.DEV) {
      console.log(`✅ [${status}] ${response.config.url}`, data);
    }

    // 处理业务状态码
    if (data.code !== undefined) {
      switch (data.code) {
        case 0:
        case 200:
        case 201: // ✅ 建议加上 201，防止后端在 body 里也返回 201
          return data;
        case 401:
          ElMessage.error("登录已过期，请重新登录");
          removeToken();
          window.location.href = "/login";
          return Promise.reject(new Error("登录已过期"));
        case 403:
          ElMessage.error("没有权限访问该资源");
          return Promise.reject(new Error("权限不足"));
        case 404:
          ElMessage.error("请求的资源不存在");
          return Promise.reject(new Error("资源不存在"));
        case 500:
          ElMessage.error("服务器内部错误");
          return Promise.reject(new Error("服务器错误"));
        default:
          // ✅ 修复点：兼容 data.msg
          const errorMsg = data.msg || data.message || "请求失败";
          ElMessage.error(errorMsg);
          return Promise.reject(new Error(errorMsg));
      }
    }

    return data;
  },
  (error) => {
    console.error("响应拦截器错误:", error);
    if (error.response) {
      console.error("错误状态码:", error.response.status);
      console.error("错误数据:", error.response.data);
    }

    // 网络错误处理
    if (!error.response) {
      ElMessage.error("网络连接失败，请检查网络设置");
      return Promise.reject(error);
    }

    const { status, data } = error.response;

    // HTTP状态码处理
    switch (status) {
      case 400:
        ElMessage.error(data?.message || "请求参数错误");
        break;
      case 401:
        ElMessage.error("未授权，请重新登录");
        removeToken();
        window.location.href = "/login";
        break;
      case 403:
        ElMessage.error("拒绝访问");
        break;
      case 404:
        ElMessage.error("请求地址不存在");
        break;
      case 408:
        ElMessage.error("请求超时");
        break;
      case 500:
        ElMessage.error("服务器内部错误");
        break;
      case 501:
        ElMessage.error("服务未实现");
        break;
      case 502:
        ElMessage.error("网关错误");
        break;
      case 503:
        ElMessage.error("服务不可用");
        break;
      case 504:
        ElMessage.error("网关超时");
        break;
      case 505:
        ElMessage.error("HTTP版本不受支持");
        break;
      default:
        ElMessage.error(data?.message || `连接错误${status}`);
    }

    return Promise.reject(error);
  }
);

/**
 * 通用请求方法
 */
export const request = {
  // GET请求
  get(url, params = {}, config = {}) {
    return service({
      method: "get",
      url,
      params,
      ...config,
    });
  },

  // POST请求
  post(url, data = {}, config = {}) {
    return service({
      method: "post",
      url,
      data,
      ...config,
    });
  },

  // PUT请求
  put(url, data = {}, config = {}) {
    return service({
      method: "put",
      url,
      data,
      ...config,
    });
  },

  // DELETE请求
  delete(url, params = {}, config = {}) {
    return service({
      method: "delete",
      url,
      params,
      ...config,
    });
  },

  // PATCH请求
  patch(url, data = {}, config = {}) {
    return service({
      method: "patch",
      url,
      data,
      ...config,
    });
  },

  // 文件上传
  upload(url, formData, config = {}) {
    return service({
      method: "post",
      url,
      data: formData,
      headers: {
        "Content-Type": "multipart/form-data",
      },
      ...config,
    });
  },

  // 文件下载
  download(url, params = {}, config = {}) {
    return service({
      method: "get",
      url,
      params,
      responseType: "blob",
      ...config,
    });
  },
};

export default service;
