<template>
  <div v-if="visible" class="modal-container">
    <div class="modal">
      <div class="modal-header">
        <h2 class="modal-title">创建学习房间</h2>
        <div class="close-btn" @click="closeModal">×</div>
      </div>

      <div class="modal-body">
        <div class="form-group">
          <label class="form-label">AI 房间灵感</label>
          <div class="ai-row">
            <input
              v-model="aiPrompt"
              type="text"
              placeholder="例如：考研冲刺 / 前端刷题 / 论文阅读"
            />
            <button class="ai-btn" type="button" :disabled="aiLoading" @click="generateRoomIdea">
              {{ aiLoading ? '生成中...' : 'AI 生成' }}
            </button>
          </div>
          <div class="ai-hint">让 AI 帮你起名字、写简介并推荐标签。</div>
        </div>

        <!-- 房间名称 -->
        <div class="form-group">
          <label class="form-label">房间名称</label>
          <input 
            v-model="formData.roomName" 
            type="text" 
            placeholder="请输入房间名称（2-20字）" 
          />
          <div class="error-msg">{{ errors.roomName }}</div>
        </div>

        <!-- 房间类型 -->
        <div class="form-group">
          <label class="form-label">房间类型</label>
          <select v-model="formData.roomType">
            <option value="public">公开房间（任何人可见）</option>
            <option value="private">私密房间（需密码加入）</option>
          </select>
          <div class="error-msg">{{ errors.roomType }}</div>
        </div>

        <!-- 房间描述 -->
        <div class="form-group">
          <label class="form-label">房间描述</label>
          <textarea
            v-model="formData.description"
            class="textarea"
            placeholder="请简要描述房间用途（0-100字）"
          ></textarea>
          <div class="error-msg">{{ errors.description }}</div>
        </div>

        <!-- 密码设置 -->
        <div v-if="formData.roomType === 'private'" class="form-group">
          <label class="form-label password-label">
            房间密码 <span>（6-12位数字/字母）</span>
          </label>
          <div class="password-input-wrapper">
            <input 
              v-model="formData.password"
              :type="showPassword ? 'text' : 'password'" 
              placeholder="请输入房间密码" 
            />
            <span class="password-toggle" @click="togglePassword">
              <iconify-icon 
                :icon="showPassword ? 'mdi:eye-off-outline' : 'mdi:eye-outline'" 
                width="18"
              ></iconify-icon>
            </span>
          </div>
          <div class="error-msg">{{ errors.password }}</div>
        </div>

        <!-- 最大人数 -->
        <div class="form-group">
          <label class="form-label">最大人数</label>
          <select v-model="formData.maxUsers">
            <option value="5">5人</option>
            <option value="10">10人</option>
            <option value="20">20人</option>
            <option value="50">50人</option>
            <option value="unlimited">不限</option>
          </select>
          <div class="error-msg">{{ errors.maxUsers }}</div>
        </div>

        <!-- 学习主题 -->
        <div class="form-group">
          <label class="form-label">学习主题</label>
          <div class="tags-container">
            <div 
              v-for="tag in availableTags" 
              :key="tag"
              :class="['tag', { selected: formData.tags.includes(tag) }]"
              @click="toggleTag(tag)"
            >
              {{ tag }}
            </div>
          </div>
          <div class="error-msg">{{ errors.tags }}</div>
        </div>
      </div>

      <div class="action-buttons">
        <button class="cancel-btn" @click="closeModal">取消</button>
        <button 
          class="create-btn" 
          :disabled="isCreating"
          @click="createRoom"
        >
          <div v-if="isCreating" class="loading-dots">
            <span class="loading-dot"></span>
            <span class="loading-dot"></span>
            <span class="loading-dot"></span>
          </div>
          <span v-else>创建房间</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { ElMessage } from 'element-plus'
import { createStudyRoom } from '@/api/modules/study'
import { generateRoomIdea } from '@/api/modules/ai'

export default {
  name: 'CreateRoom',
  props: {
    visible: {
      type: Boolean,
      default: false
    }
  },
  emits: ['close', 'created'],
  data() {
    return {
      showPassword: false,
      isCreating: false,
      aiPrompt: '',
      aiLoading: false,
      formData: {
        roomName: '',
        roomType: 'public',
        description: '',
        password: '',
        maxUsers: '20',
        tags: ['后端开发']
      },
      errors: {},
      availableTags: [
        '前端开发',
        '后端开发',
        'Python',
        'Java',
        '英语学习',
        '考研备考',
        '职业技能',
        '其他'
      ]
    }
  },
  methods: {
    closeModal() {
      this.$emit('close')
      this.resetForm()
    },
    
    togglePassword() {
      this.showPassword = !this.showPassword
    },
    
    toggleTag(tag) {
      const index = this.formData.tags.indexOf(tag)
      if (index > -1) {
        this.formData.tags.splice(index, 1)
      } else {
        if (this.formData.tags.length < 3) {
          this.formData.tags.push(tag)
        }
      }
    },
    async generateRoomIdea() {
      if (this.aiLoading) return
      this.aiLoading = true
      try {
        const response = await generateRoomIdea({ prompt: this.aiPrompt })
        const data = response?.data || response
        const name = (data?.name || '').trim()
        const description = (data?.description || '').trim()
        const tags = Array.isArray(data?.tags) ? data.tags : []

        if (name) this.formData.roomName = name
        if (description) this.formData.description = description

        if (tags.length) {
          const normalized = tags
            .map((tag) => String(tag).trim())
            .filter(Boolean)
            .slice(0, 3)
          normalized.forEach((tag) => {
            if (!this.availableTags.includes(tag)) {
              this.availableTags.push(tag)
            }
          })
          this.formData.tags = normalized
        }
      } catch (error) {
        console.error('生成房间创意失败:', error)
        ElMessage.error(error?.message || '生成房间创意失败，请稍后重试')
      } finally {
        this.aiLoading = false
      }
    },
    
    validateForm() {
      this.errors = {}
      let isValid = true
      
      // 验证房间名称
      if (!this.formData.roomName.trim()) {
        this.errors.roomName = '请输入房间名称'
        isValid = false
      } else if (this.formData.roomName.length < 2 || this.formData.roomName.length > 20) {
        this.errors.roomName = '房间名称长度需为2-20字'
        isValid = false
      }
      
      // 验证私密房间密码
      if (this.formData.roomType === 'private') {
        if (!this.formData.password.trim()) {
          this.errors.password = '请输入房间密码'
          isValid = false
        } else if (this.formData.password.length < 6 || this.formData.password.length > 12) {
          this.errors.password = '密码长度需为6-12位'
          isValid = false
        }
      }
      
      // 验证描述长度
      if (this.formData.description.length > 100) {
        this.errors.description = '描述不能超过100字'
        isValid = false
      }
      
      return isValid
    },
    
    async createRoom() {
      if (!this.validateForm()) {
        return
      }

      this.isCreating = true

      const payload = {
        name: this.formData.roomName.trim(),
        room_type: this.formData.roomType,
        description: this.formData.description.trim(),
        password: this.formData.roomType === 'private' ? this.formData.password : '',
        max_members: this.formData.maxUsers === 'unlimited' ? 0 : Number(this.formData.maxUsers),
        tags: [...this.formData.tags]
      }

      try {
        const response = await createStudyRoom(payload)
        const room = response?.data?.room || response?.data
        this.$emit('created', room)
        ElMessage.success('房间创建成功')
        this.closeModal()
      } catch (error) {
        console.error('创建房间失败:', error)
        ElMessage.error(error?.message || '创建房间失败，请稍后重试')
      } finally {
        this.isCreating = false
      }
    },
    
    resetForm() {
      this.formData = {
        roomName: '',
        roomType: 'public',
        description: '',
        password: '',
        maxUsers: '20',
        tags: ['后端开发']
      }
      this.errors = {}
      this.showPassword = false
      this.aiPrompt = ''
      this.aiLoading = false
    }
  }
}
</script>

<style scoped>
.modal-container {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal {
  width: 500px;
  max-height: 80vh;
  min-height: 400px;
  background-color: #ffffff;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  padding: 24px;
  display: flex;
  flex-direction: column;
}

.modal-header {
  height: 40px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.modal-title {
  font-size: 18px;
  font-weight: bold;
  color: #333333;
}

.close-btn {
  font-size: 16px;
  color: #999999;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: color 0.2s;
  border-radius: 4px;
}

.close-btn:hover {
  color: #333333;
  background-color: #f5f5f5;
}

.modal-body {
  flex: 1;
  overflow-y: auto;
  margin-bottom: 24px;
}

.form-group {
  margin-bottom: 20px;
}

.ai-row {
  display: flex;
  gap: 10px;
  align-items: center;
}

.ai-row input {
  flex: 1;
}

.ai-btn {
  height: 40px;
  padding: 0 14px;
  border-radius: 6px;
  border: 1px solid #1e88e5;
  background: #e6f4ff;
  color: #1e88e5;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.ai-btn:hover:not(:disabled) {
  background: #1e88e5;
  color: #fff;
}

.ai-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.ai-hint {
  font-size: 12px;
  color: #7b8794;
  margin-top: 6px;
}

.form-label {
  display: block;
  font-size: 14px;
  color: #333333;
  margin-bottom: 8px;
}

.password-label span {
  font-size: 12px;
  color: #999999;
  margin-left: 4px;
}

input,
select,
textarea {
  width: 100%;
  height: 40px;
  border: 1px solid #e5e6eb;
  border-radius: 4px;
  padding: 0 12px;
  font-size: 14px;
  font-family: "Microsoft YaHei", sans-serif;
  outline: none;
}

input:focus,
select:focus,
textarea:focus {
  border-color: #1e88e5;
}

input::placeholder,
textarea::placeholder {
  color: #999999;
}

.textarea {
  height: 80px;
  padding: 12px;
  resize: none;
}

.password-input-wrapper {
  position: relative;
}

.password-toggle {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: #999999;
  cursor: pointer;
  font-size: 18px;
}

.tags-container {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  row-gap: 12px;
}

.tag {
  background-color: white;
  border: 1px solid #e5e6eb;
  border-radius: 16px;
  padding: 6px 12px;
  font-size: 14px;
  color: #666666;
  cursor: pointer;
  transition: all 0.2s;
}

.tag.selected {
  background-color: #e6f7ff;
  border-color: #91d5ff;
  color: #1890ff;
}

.error-msg {
  height: 16px;
  font-size: 12px;
  color: #f5222d;
  margin-top: 4px;
}

.action-buttons {
  display: flex;
  justify-content: flex-end;
  gap: 16px;
  height: 50px;
}

.cancel-btn,
.create-btn {
  height: 40px;
  font-size: 16px;
  font-weight: bold;
  border-radius: 4px;
  font-family: "Microsoft YaHei", sans-serif;
  cursor: pointer;
  border: 1px solid #e5e6eb;
}

.cancel-btn {
  width: 100px;
  background-color: white;
  color: #666666;
}

.cancel-btn:hover {
  background-color: #f5f5f5;
}

.create-btn {
  width: 120px;
  background-color: #1890ff;
  color: white;
  border: none;
}

.create-btn:hover:not(:disabled) {
  background-color: #096dd9;
}

.create-btn:disabled {
  background-color: #b3d8ff;
  cursor: not-allowed;
}

.loading-dots {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 4px;
}

.loading-dot {
  width: 6px;
  height: 6px;
  background-color: white;
  border-radius: 50%;
  display: inline-block;
  animation: dot-flashing 1s infinite;
  opacity: 0.3;
}

.loading-dot:nth-child(1) {
  animation-delay: 0;
}

.loading-dot:nth-child(2) {
  animation-delay: 0.2s;
}

.loading-dot:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes dot-flashing {
  0% {
    opacity: 0.3;
  }
  50% {
    opacity: 1;
  }
  100% {
    opacity: 0.3;
  }
}

/* 滚动条样式 */
.modal-body::-webkit-scrollbar {
  width: 6px;
}

.modal-body::-webkit-scrollbar-track {
  background: #f5f5f5;
  border-radius: 3px;
}

.modal-body::-webkit-scrollbar-thumb {
  background: #cccccc;
  border-radius: 3px;
}

/* 隐藏滚动条 */
.hide-scrollbar::-webkit-scrollbar {
  display: none;
}
.hide-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
