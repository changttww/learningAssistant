<template>
  <div class="register-container">
    <div class="register-card">
      <div class="register-header">
        <div class="logo">
          <iconify-icon
            icon="mdi:school-outline"
            width="40"
            height="40"
            class="text-blue-600"
          ></iconify-icon>
        </div>
        <h1 class="title">创建账号</h1>
        <p class="subtitle">加入学习助手，打造高效学习计划</p>
      </div>

      <form @submit.prevent="handleRegister" class="register-form">
        <div class="form-group">
          <label for="displayName" class="form-label">昵称</label>
          <div class="input-wrapper">
            <iconify-icon
              icon="mdi:account-circle"
              width="20"
              height="20"
              class="input-icon"
            ></iconify-icon>
            <input
              id="displayName"
              v-model="registerForm.displayName"
              type="text"
              class="form-input"
              placeholder="请输入昵称"
              required
            />
          </div>
        </div>

        <div class="form-group">
          <label for="username" class="form-label">用户名</label>
          <div class="input-wrapper">
            <iconify-icon
              icon="mdi:account"
              width="20"
              height="20"
              class="input-icon"
            ></iconify-icon>
            <input
              id="username"
              v-model="registerForm.username"
              type="text"
              class="form-input"
              placeholder="请输入用户名"
              required
            />
          </div>
        </div>

        <div class="form-group">
          <label for="email" class="form-label">邮箱</label>
          <div class="input-wrapper">
            <iconify-icon
              icon="mdi:email"
              width="20"
              height="20"
              class="input-icon"
            ></iconify-icon>
            <input
              id="email"
              v-model="registerForm.email"
              type="email"
              class="form-input"
              placeholder="请输入邮箱"
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
              v-model="registerForm.password"
              :type="showPassword ? 'text' : 'password'"
              class="form-input"
              placeholder="请输入密码"
              required
              minlength="6"
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

        <div class="form-group">
          <label for="confirmPassword" class="form-label">确认密码</label>
          <div class="input-wrapper">
            <iconify-icon
              icon="mdi:lock-check"
              width="20"
              height="20"
              class="input-icon"
            ></iconify-icon>
            <input
              id="confirmPassword"
              v-model="registerForm.confirmPassword"
              :type="showConfirmPassword ? 'text' : 'password'"
              class="form-input"
              placeholder="请再次输入密码"
              required
              minlength="6"
            />
            <button
              type="button"
              @click="showConfirmPassword = !showConfirmPassword"
              class="password-toggle"
            >
              <iconify-icon
                :icon="showConfirmPassword ? 'mdi:eye-off' : 'mdi:eye'"
                width="20"
                height="20"
              ></iconify-icon>
            </button>
          </div>
        </div>

        <button type="submit" class="register-btn" :disabled="loading">
          <iconify-icon
            v-if="loading"
            icon="mdi:loading"
            width="20"
            height="20"
            class="animate-spin"
          ></iconify-icon>
          <span>{{ loading ? "注册中..." : "注册并登录" }}</span>
        </button>
      </form>

      <div class="login-link">
        已有账号？
        <router-link :to="{ name: 'Login' }">立即登录</router-link>
      </div>
    </div>

    <div class="bg-decoration">
      <div class="decoration-circle circle-1"></div>
      <div class="decoration-circle circle-2"></div>
      <div class="decoration-circle circle-3"></div>
    </div>
  </div>
</template>

<script>
  import { register, login } from "@/api/modules/auth";
  import {
    setToken,
    setRefreshToken,
    setUserInfo,
    setPermissions,
    setRoles,
  } from "@/utils/auth";
  import { useCurrentUser } from "@/composables/useCurrentUser";

  export default {
    name: "Register",
    setup() {
      const { setCurrentUser, loadCurrentUser, loadStudyStats } =
        useCurrentUser();
      return {
        setCurrentUser,
        loadCurrentUser,
        loadStudyStats,
      };
    },
    data() {
      return {
        registerForm: {
          displayName: "",
          username: "",
          email: "",
          password: "",
          confirmPassword: "",
        },
        showPassword: false,
        showConfirmPassword: false,
        loading: false,
      };
    },
    methods: {
      validateForm() {
        if (this.registerForm.password !== this.registerForm.confirmPassword) {
          alert("两次输入的密码不一致");
          return false;
        }
        if (this.registerForm.password.length < 6) {
          alert("密码长度至少为6位");
          return false;
        }
        return true;
      },
      async handleRegister() {
        if (!this.validateForm()) {
          return;
        }

        this.loading = true;

        try {
          const payload = {
            username: this.registerForm.username.trim(),
            email: this.registerForm.email.trim(),
            password: this.registerForm.password,
            display_name: this.registerForm.displayName.trim(),
          };

          const registerRes = await register(payload);
          const newUser = registerRes.data?.user;

          if (!newUser?.id) {
            throw new Error("注册失败，请稍后重试");
          }

          const loginRes = await login({
            identifier: payload.username,
            password: this.registerForm.password,
          });

          const loginData = loginRes.data || {};
          const {
            token,
            refresh_token: refreshToken,
            user,
            permissions = [],
            roles = [],
          } = loginData;

          if (!token || !user?.id) {
            throw new Error("自动登录失败，请手动登录");
          }

          setToken(token);
          if (refreshToken) {
            setRefreshToken(refreshToken);
          }
          setPermissions(permissions);
          setRoles(roles);
          setUserInfo(user);
          this.setCurrentUser(user);

          await Promise.all([
            this.loadCurrentUser(user.id, { force: true }),
            this.loadStudyStats(user.id, { force: true }),
          ]);

          alert("注册成功，已自动登录");
          this.$router.push("/");
        } catch (error) {
          console.error("注册失败:", error);
          alert(error?.message || "注册失败，请稍后再试");
        } finally {
          this.loading = false;
        }
      },
    },
    mounted() {
      document.title = "注册 | 学习助手";
    },
  };
</script>

<style scoped>
  .register-container {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #43cea2 0%, #185a9d 100%);
    padding: 20px;
    position: relative;
    overflow: hidden;
  }

  .register-card {
    background: white;
    border-radius: 16px;
    padding: 40px;
    width: 100%;
    max-width: 460px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.12);
    position: relative;
    z-index: 10;
  }

  .register-header {
    text-align: center;
    margin-bottom: 24px;
  }

  .logo {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 64px;
    height: 64px;
    background: rgba(59, 130, 246, 0.12);
    border-radius: 16px;
    margin-bottom: 16px;
  }

  .title {
    font-size: 28px;
    font-weight: 700;
    color: #1f2933;
  }

  .subtitle {
    font-size: 14px;
    color: #6b7280;
    margin-top: 8px;
  }

  .register-form {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .form-group {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .form-label {
    font-size: 14px;
    font-weight: 600;
    color: #374151;
  }

  .input-wrapper {
    position: relative;
    display: flex;
    align-items: center;
  }

  .input-icon {
    position: absolute;
    left: 12px;
    color: #9ca3af;
  }

  .form-input {
    width: 100%;
    padding: 12px 14px 12px 40px;
    border: 1px solid #d1d5db;
    border-radius: 10px;
    font-size: 14px;
    transition: border-color 0.2s ease, box-shadow 0.2s ease;
  }

  .form-input:focus {
    outline: none;
    border-color: #2563eb;
    box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.2);
  }

  .password-toggle {
    position: absolute;
    right: 12px;
    background: none;
    border: none;
    cursor: pointer;
    color: #9ca3af;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .register-btn {
    margin-top: 12px;
    background: linear-gradient(135deg, #2563eb 0%, #3b82f6 100%);
    border: none;
    color: white;
    padding: 12px;
    border-radius: 10px;
    font-size: 16px;
    font-weight: 600;
    cursor: pointer;
    transition: transform 0.2s ease, box-shadow 0.2s ease;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
  }

  .register-btn:hover {
    transform: translateY(-1px);
    box-shadow: 0 12px 20px rgba(37, 99, 235, 0.25);
  }

  .register-btn:disabled {
    opacity: 0.7;
    cursor: not-allowed;
    box-shadow: none;
  }

  .login-link {
    margin-top: 20px;
    text-align: center;
    color: #4b5563;
    font-size: 14px;
  }

  .login-link a {
    color: #2563eb;
    text-decoration: none;
    margin-left: 4px;
    font-weight: 600;
  }

  .login-link a:hover {
    text-decoration: underline;
  }

  .bg-decoration {
    position: absolute;
    inset: 0;
    overflow: hidden;
  }

  .decoration-circle {
    position: absolute;
    border-radius: 9999px;
    opacity: 0.25;
  }

  .circle-1 {
    width: 260px;
    height: 260px;
    background: rgba(59, 130, 246, 0.4);
    top: -80px;
    right: -80px;
  }

  .circle-2 {
    width: 180px;
    height: 180px;
    background: rgba(37, 99, 235, 0.35);
    bottom: -60px;
    left: -60px;
  }

  .circle-3 {
    width: 120px;
    height: 120px;
    background: rgba(96, 165, 250, 0.3);
    bottom: 60px;
    right: 90px;
  }

  @media (max-width: 768px) {
    .register-card {
      padding: 32px 24px;
    }
  }
</style>
