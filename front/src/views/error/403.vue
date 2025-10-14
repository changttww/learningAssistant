<template>
  <div class="error-container">
    <div class="error-content">
      <!-- 403图标 -->
      <div class="error-icon">
        <iconify-icon
          icon="mdi:shield-alert-outline"
          width="120"
          height="120"
          class="text-red-400"
        ></iconify-icon>
      </div>

      <!-- 错误信息 -->
      <div class="error-info">
        <h1 class="error-code">403</h1>
        <h2 class="error-title">访问被拒绝</h2>
        <p class="error-description">
          抱歉，您没有权限访问此页面。
          <br />
          请联系管理员获取相应权限，或返回您有权限访问的页面。
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

      <!-- 权限说明 -->
      <div class="permission-info">
        <h3>权限说明</h3>
        <div class="info-grid">
          <div class="info-item">
            <iconify-icon
              icon="mdi:account"
              width="24"
              height="24"
              class="text-blue-500"
            ></iconify-icon>
            <div>
              <h4>当前用户</h4>
              <p>普通用户</p>
            </div>
          </div>
          <div class="info-item">
            <iconify-icon
              icon="mdi:shield-check"
              width="24"
              height="24"
              class="text-green-500"
            ></iconify-icon>
            <div>
              <h4>所需权限</h4>
              <p>管理员权限</p>
            </div>
          </div>
          <div class="info-item">
            <iconify-icon
              icon="mdi:email"
              width="24"
              height="24"
              class="text-purple-500"
            ></iconify-icon>
            <div>
              <h4>申请权限</h4>
              <p>联系管理员</p>
            </div>
          </div>
        </div>
      </div>

      <!-- 可访问的页面 -->
      <div class="accessible-links">
        <h3>您可以访问的页面：</h3>
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
          <router-link to="/study-room" class="link-item">
            <iconify-icon
              icon="mdi:school"
              width="20"
              height="20"
            ></iconify-icon>
            <span>在线自习室</span>
          </router-link>
          <router-link to="/profile" class="link-item">
            <iconify-icon
              icon="mdi:account-circle"
              width="20"
              height="20"
            ></iconify-icon>
            <span>个人资料</span>
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  export default {
    name: "Forbidden",
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
      document.title = "403 - 访问被拒绝 | 学习助手";
    },
  };
</script>

<style scoped>
  .error-container {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #ff6b6b 0%, #ee5a24 100%);
    padding: 20px;
  }

  .error-content {
    background: white;
    border-radius: 16px;
    padding: 48px 32px;
    text-align: center;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.1);
    max-width: 700px;
    width: 100%;
  }

  .error-icon {
    margin-bottom: 24px;
  }

  .error-code {
    font-size: 72px;
    font-weight: 800;
    color: #ff6b6b;
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
    background: #ff6b6b;
    color: white;
  }

  .btn-primary:hover {
    background: #ff5252;
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

  .permission-info {
    border-top: 1px solid #e5e7eb;
    border-bottom: 1px solid #e5e7eb;
    padding: 32px 0;
    margin-bottom: 32px;
  }

  .permission-info h3 {
    font-size: 18px;
    font-weight: 600;
    color: #1f2937;
    margin: 0 0 20px 0;
  }

  .info-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 20px;
  }

  .info-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 16px;
    background: #f9fafb;
    border-radius: 8px;
    text-align: left;
  }

  .info-item h4 {
    font-size: 14px;
    font-weight: 600;
    color: #1f2937;
    margin: 0 0 4px 0;
  }

  .info-item p {
    font-size: 12px;
    color: #6b7280;
    margin: 0;
  }

  .accessible-links h3 {
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
    background: #ff6b6b;
    color: white;
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(255, 107, 107, 0.3);
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

    .info-grid {
      grid-template-columns: 1fr;
    }

    .link-grid {
      grid-template-columns: repeat(2, 1fr);
    }
  }
</style>
