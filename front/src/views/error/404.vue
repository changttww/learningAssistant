<template>
  <div class="error-container">
    <div class="error-content">
      <!-- 404图标 -->
      <div class="error-icon">
        <iconify-icon
          icon="mdi:file-question-outline"
          width="120"
          height="120"
          class="text-gray-400"
        ></iconify-icon>
      </div>

      <!-- 错误信息 -->
      <div class="error-info">
        <h1 class="error-code">404</h1>
        <h2 class="error-title">页面不存在</h2>
        <p class="error-description">
          抱歉，您访问的页面不存在或已被移除。
          <br />
          请检查网址是否正确，或返回首页继续浏览。
        </p>
      </div>

      <!-- 操作按钮 -->
      <div class="error-actions">
        <button @click="goBack" class="btn btn-secondary">
          <iconify-icon
            icon="mdi:arrow-left"
            width="16"
            height="16"
          ></iconify-icon>
          返回上页
        </button>
        <button @click="goHome" class="btn btn-primary">
          <iconify-icon icon="mdi:home" width="16" height="16"></iconify-icon>
          回到首页
        </button>
      </div>

      <!-- 建议链接 -->
      <div class="suggested-links">
        <h3>您可能想要访问：</h3>
        <div class="link-grid">
          <router-link to="/" class="link-item">
            <iconify-icon icon="mdi:home" width="20" height="20"></iconify-icon>
            <span>首页</span>
          </router-link>
          <router-link to="/personal-tasks" class="link-item">
            <iconify-icon
              icon="mdi:clipboard-list"
              width="20"
              height="20"
            ></iconify-icon>
            <span>个人任务</span>
          </router-link>
          <router-link to="/team-tasks" class="link-item">
            <iconify-icon
              icon="mdi:account-group"
              width="20"
              height="20"
            ></iconify-icon>
            <span>团队任务</span>
          </router-link>
          <router-link to="/study-room" class="link-item">
            <iconify-icon
              icon="mdi:school"
              width="20"
              height="20"
            ></iconify-icon>
            <span>在线自习室</span>
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  export default {
    name: "NotFound",
    methods: {
      goBack() {
        // 如果有历史记录就返回上一页，否则回到首页
        if (window.history.length > 1) {
          this.$router.go(-1);
        } else {
          this.$router.push("/");
        }
      },
      goHome() {
        this.$router.push("/");
      },
    },
    mounted() {
      // 设置页面标题
      document.title = "404 - 页面不存在 | 学习助手";
    },
  };
</script>

<style scoped>
  .error-container {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    padding: 20px;
  }

  .error-content {
    background: white;
    border-radius: 16px;
    padding: 48px 32px;
    text-align: center;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.1);
    max-width: 600px;
    width: 100%;
  }

  .error-icon {
    margin-bottom: 24px;
  }

  .error-code {
    font-size: 72px;
    font-weight: 800;
    color: #667eea;
    margin: 0 0 16px 0;
    line-height: 1;
  }

  .error-title {
    font-size: 32px;
    font-weight: 600;
    color: #1f2937;
    margin: 0 0 16px 0;
  }

  .error-description {
    font-size: 16px;
    color: #6b7280;
    line-height: 1.6;
    margin: 0 0 32px 0;
  }

  .error-actions {
    display: flex;
    gap: 16px;
    justify-content: center;
    margin-bottom: 40px;
    flex-wrap: wrap;
  }

  .btn {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    padding: 12px 24px;
    border-radius: 8px;
    font-weight: 500;
    text-decoration: none;
    border: none;
    cursor: pointer;
    transition: all 0.2s ease;
    font-size: 14px;
  }

  .btn-primary {
    background: #667eea;
    color: white;
  }

  .btn-primary:hover {
    background: #5a6fd8;
    transform: translateY(-1px);
  }

  .btn-secondary {
    background: #f3f4f6;
    color: #374151;
    border: 1px solid #d1d5db;
  }

  .btn-secondary:hover {
    background: #e5e7eb;
    transform: translateY(-1px);
  }

  .suggested-links {
    border-top: 1px solid #e5e7eb;
    padding-top: 32px;
  }

  .suggested-links h3 {
    font-size: 18px;
    font-weight: 600;
    color: #1f2937;
    margin: 0 0 20px 0;
  }

  .link-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
    gap: 12px;
  }

  .link-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
    padding: 16px 12px;
    background: #f9fafb;
    border-radius: 8px;
    text-decoration: none;
    color: #374151;
    transition: all 0.2s ease;
    border: 1px solid transparent;
  }

  .link-item:hover {
    background: #667eea;
    color: white;
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
  }

  .link-item span {
    font-size: 14px;
    font-weight: 500;
  }

  /* 响应式设计 */
  @media (max-width: 640px) {
    .error-content {
      padding: 32px 24px;
    }

    .error-code {
      font-size: 56px;
    }

    .error-title {
      font-size: 24px;
    }

    .error-actions {
      flex-direction: column;
      align-items: center;
    }

    .btn {
      width: 100%;
      max-width: 200px;
    }

    .link-grid {
      grid-template-columns: repeat(2, 1fr);
    }
  }
</style>
