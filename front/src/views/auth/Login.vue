<template>
  <div class="login-container">
    <div class="login-card">
      <!-- 头部 -->
      <div class="login-header">
        <div class="logo">
          <iconify-icon
            icon="mdi:school"
            width="40"
            height="40"
            class="text-blue-600"
          ></iconify-icon>
        </div>
        <h1 class="title">学习助手</h1>
        <p class="subtitle">登录您的账户，开始高效学习之旅</p>
      </div>

      <!-- 登录表单 -->
      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label for="username" class="form-label">用户名或邮箱</label>
          <div class="input-wrapper">
            <iconify-icon
              icon="mdi:account"
              width="20"
              height="20"
              class="input-icon"
            ></iconify-icon>
            <input
              id="username"
              v-model="loginForm.username"
              type="text"
              class="form-input"
              placeholder="请输入用户名或邮箱"
              required
            />
          </div>
        </div>

        <div class="form-group">
          <label for="password" class="form-label">密码</label>
          <div class="input-wrapper">
            <iconify-icon
              icon="mdi:lock"
              width="20"
              height="20"
              class="input-icon"
            ></iconify-icon>
            <input
              id="password"
              v-model="loginForm.password"
              :type="showPassword ? 'text' : 'password'"
              class="form-input"
              placeholder="请输入密码"
              required
            />
            <button
              type="button"
              @click="showPassword = !showPassword"
              class="password-toggle"
            >
              <iconify-icon
                :icon="showPassword ? 'mdi:eye-off' : 'mdi:eye'"
                width="20"
                height="20"
              ></iconify-icon>
            </button>
          </div>
        </div>

        <div class="form-options">
          <label class="checkbox-wrapper">
            <input
              v-model="loginForm.remember"
              type="checkbox"
              class="checkbox"
            />
            <span class="checkbox-label">记住我</span>
          </label>
          <a href="#" class="forgot-link">忘记密码？</a>
        </div>

        <button type="submit" class="login-btn" :disabled="loading">
          <iconify-icon
            v-if="loading"
            icon="mdi:loading"
            width="20"
            height="20"
            class="animate-spin"
          ></iconify-icon>
          <span>{{ loading ? "登录中..." : "登录" }}</span>
        </button>
      </form>

      <!-- 分割线 -->
      <div class="divider">
        <span>或</span>
      </div>

      <!-- 第三方登录 -->
      <div class="social-login">
        <button class="social-btn github">
          <iconify-icon icon="mdi:github" width="20" height="20"></iconify-icon>
          <span>GitHub登录</span>
        </button>
        <button class="social-btn google">
          <iconify-icon icon="mdi:google" width="20" height="20"></iconify-icon>
          <span>Google登录</span>
        </button>
      </div>

      <!-- 注册链接 -->
      <div class="register-link">
        <span>还没有账户？</span>
        <a href="#" @click="goToRegister">立即注册</a>
      </div>
    </div>

    <!-- 背景装饰 -->
    <div class="bg-decoration">
      <div class="decoration-circle circle-1"></div>
      <div class="decoration-circle circle-2"></div>
      <div class="decoration-circle circle-3"></div>
    </div>
  </div>
</template>

<script>
  export default {
    name: "Login",
    data() {
      return {
        loginForm: {
          username: "",
          password: "",
          remember: false,
        },
        showPassword: false,
        loading: false,
      };
    },
    methods: {
      async handleLogin() {
        this.loading = true;

        try {
          // 模拟登录请求
          await new Promise((resolve) => setTimeout(resolve, 1500));

          // 这里应该调用实际的登录API
          // const response = await authApi.login(this.loginForm)

          // 模拟登录成功
          console.log("登录成功", this.loginForm);

          // 跳转到首页或之前访问的页面
          const redirect = this.$route.query.redirect || "/";
          this.$router.push(redirect);
        } catch (error) {
          console.error("登录失败:", error);
          // 这里应该显示错误提示
          alert("登录失败，请检查用户名和密码");
        } finally {
          this.loading = false;
        }
      },

      goToRegister() {
        // 跳转到注册页面
        console.log("跳转到注册页面");
      },
    },

    mounted() {
      document.title = "登录 | 学习助手";
    },
  };
</script>

<style scoped>
  .login-container {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    padding: 20px;
    position: relative;
    overflow: hidden;
  }

  .login-card {
    background: white;
    border-radius: 16px;
    padding: 40px;
    width: 100%;
    max-width: 420px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.1);
    position: relative;
    z-index: 10;
  }

  .login-header {
    text-align: center;
    margin-bottom: 32px;
  }

  .logo {
    margin-bottom: 16px;
  }

  .title {
    font-size: 28px;
    font-weight: 700;
    color: #1f2937;
    margin: 0 0 8px 0;
  }

  .subtitle {
    font-size: 14px;
    color: #6b7280;
    margin: 0;
  }

  .login-form {
    margin-bottom: 24px;
  }

  .form-group {
    margin-bottom: 20px;
  }

  .form-label {
    display: block;
    font-size: 14px;
    font-weight: 500;
    color: #374151;
    margin-bottom: 8px;
  }

  .input-wrapper {
    position: relative;
  }

  .input-icon {
    position: absolute;
    left: 12px;
    top: 50%;
    transform: translateY(-50%);
    color: #9ca3af;
  }

  .form-input {
    width: 100%;
    padding: 12px 12px 12px 44px;
    border: 1px solid #d1d5db;
    border-radius: 8px;
    font-size: 14px;
    transition: all 0.2s ease;
    box-sizing: border-box;
  }

  .form-input:focus {
    outline: none;
    border-color: #667eea;
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
  }

  .password-toggle {
    position: absolute;
    right: 12px;
    top: 50%;
    transform: translateY(-50%);
    background: none;
    border: none;
    color: #9ca3af;
    cursor: pointer;
    padding: 0;
  }

  .password-toggle:hover {
    color: #667eea;
  }

  .form-options {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
  }

  .checkbox-wrapper {
    display: flex;
    align-items: center;
    gap: 8px;
    cursor: pointer;
  }

  .checkbox {
    width: 16px;
    height: 16px;
    accent-color: #667eea;
  }

  .checkbox-label {
    font-size: 14px;
    color: #374151;
  }

  .forgot-link {
    font-size: 14px;
    color: #667eea;
    text-decoration: none;
  }

  .forgot-link:hover {
    text-decoration: underline;
  }

  .login-btn {
    width: 100%;
    background: #667eea;
    color: white;
    border: none;
    border-radius: 8px;
    padding: 12px;
    font-size: 16px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
  }

  .login-btn:hover:not(:disabled) {
    background: #5a6fd8;
    transform: translateY(-1px);
  }

  .login-btn:disabled {
    opacity: 0.7;
    cursor: not-allowed;
  }

  .divider {
    text-align: center;
    margin: 24px 0;
    position: relative;
  }

  .divider::before {
    content: "";
    position: absolute;
    top: 50%;
    left: 0;
    right: 0;
    height: 1px;
    background: #e5e7eb;
  }

  .divider span {
    background: white;
    padding: 0 16px;
    color: #9ca3af;
    font-size: 14px;
  }

  .social-login {
    display: flex;
    gap: 12px;
    margin-bottom: 24px;
  }

  .social-btn {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    padding: 10px;
    border: 1px solid #d1d5db;
    border-radius: 8px;
    background: white;
    color: #374151;
    font-size: 14px;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .social-btn:hover {
    background: #f9fafb;
    transform: translateY(-1px);
  }

  .social-btn.github:hover {
    border-color: #24292e;
    color: #24292e;
  }

  .social-btn.google:hover {
    border-color: #ea4335;
    color: #ea4335;
  }

  .register-link {
    text-align: center;
    font-size: 14px;
    color: #6b7280;
  }

  .register-link a {
    color: #667eea;
    text-decoration: none;
    font-weight: 500;
    margin-left: 4px;
  }

  .register-link a:hover {
    text-decoration: underline;
  }

  /* 背景装饰 */
  .bg-decoration {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    pointer-events: none;
  }

  .decoration-circle {
    position: absolute;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.1);
  }

  .circle-1 {
    width: 200px;
    height: 200px;
    top: -100px;
    right: -100px;
  }

  .circle-2 {
    width: 150px;
    height: 150px;
    bottom: -75px;
    left: -75px;
  }

  .circle-3 {
    width: 100px;
    height: 100px;
    top: 50%;
    left: -50px;
    transform: translateY(-50%);
  }

  /* 动画 */
  .animate-spin {
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(360deg);
    }
  }

  /* 响应式设计 */
  @media (max-width: 480px) {
    .login-card {
      padding: 32px 24px;
    }

    .title {
      font-size: 24px;
    }

    .social-login {
      flex-direction: column;
    }
  }
</style>
