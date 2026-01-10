<template>
  <div class="knowledge-chat-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">💬 知识问答助手</h1>
        <p class="page-subtitle">基于您的知识库进行智能问答，回答带引用溯源</p>
      </div>
    </div>

    <!-- 主体区域 -->
    <div class="chat-wrapper">
      <!-- 聊天记录 -->
      <div class="chat-messages" ref="messagesContainer">
        <!-- 欢迎消息 -->
        <div v-if="messages.length === 0" class="welcome-message">
          <div class="welcome-icon">🤖</div>
          <h2>欢迎使用知识问答助手</h2>
          <p>我可以基于您的知识库回答问题，并展示参考来源。</p>
          <div class="quick-questions">
            <p class="quick-title">试试这些问题：</p>
            <button 
              v-for="q in quickQuestions" 
              :key="q" 
              class="quick-btn"
              @click="sendMessage(q)"
            >
              {{ q }}
            </button>
          </div>
        </div>

        <!-- 消息列表 -->
        <div 
          v-for="(msg, index) in messages" 
          :key="index" 
          :class="['message', msg.role]"
        >
          <div class="message-avatar">
            {{ msg.role === 'user' ? '👤' : '🤖' }}
          </div>
          <div class="message-content">
            <div class="message-text" v-html="formatMessage(msg.content)"></div>
            
            <!-- 引用来源 -->
            <div v-if="msg.citations && msg.citations.length > 0" class="citations">
              <div class="citations-title">📚 参考来源：</div>
              <div class="citation-list">
                <div 
                  v-for="cite in msg.citations" 
                  :key="cite.id" 
                  class="citation-item"
                  @click="showCitationDetail(cite)"
                >
                  <div class="citation-main">
                    <span class="citation-title">{{ cite.title }}</span>
                    <span class="citation-category" v-if="cite.category">{{ cite.category }}</span>
                  </div>
                  <div class="citation-similarity">
                    相似度: {{ (cite.similarity * 100).toFixed(0) }}%
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 加载中 -->
        <div v-if="loading" class="message assistant loading">
          <div class="message-avatar">🤖</div>
          <div class="message-content">
            <div class="typing-indicator">
              <span></span>
              <span></span>
              <span></span>
            </div>
          </div>
        </div>
      </div>

      <!-- 输入区域 -->
      <div class="chat-input-area">
        <div class="input-wrapper">
          <textarea
            v-model="inputText"
            placeholder="输入您的问题..."
            @keydown.enter.exact.prevent="handleSend"
            @keydown.enter.shift.exact="handleNewline"
            rows="1"
            ref="inputArea"
          ></textarea>
          <button 
            class="send-btn" 
            @click="handleSend"
            :disabled="!inputText.trim() || loading"
          >
            <span v-if="loading">⏳</span>
            <span v-else>发送</span>
          </button>
        </div>
        <p class="input-hint">按 Enter 发送，Shift + Enter 换行</p>
      </div>
    </div>

    <!-- 引用详情弹窗 -->
    <div v-if="selectedCitation" class="citation-modal" @click.self="selectedCitation = null">
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ selectedCitation.title }}</h3>
          <button class="close-btn" @click="selectedCitation = null">×</button>
        </div>
        <div class="modal-body">
          <div class="detail-row">
            <span class="label">分类：</span>
            <span class="value">{{ selectedCitation.category || '未分类' }}</span>
          </div>
          <div class="detail-row">
            <span class="label">相似度：</span>
            <span class="value">{{ (selectedCitation.similarity * 100).toFixed(1) }}%</span>
          </div>
          <div class="detail-row" v-if="selectedCitation.summary">
            <span class="label">摘要：</span>
            <span class="value summary">{{ selectedCitation.summary }}</span>
          </div>
        </div>
        <div class="modal-footer">
          <router-link 
            :to="`/knowledge-base?highlight=${selectedCitation.id}`" 
            class="btn-primary"
          >
            查看完整内容
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ragChat } from '@/api/modules/knowledge';
import { marked } from 'marked';

// 配置 marked
marked.setOptions({
  breaks: true,        // 支持换行
  gfm: true,           // 支持 GitHub 风格 Markdown
});

export default {
  name: 'KnowledgeChat',
  data() {
    return {
      inputText: '',
      messages: [],
      loading: false,
      selectedCitation: null,
      quickQuestions: [
        '帮我总结一下我学过的内容',
        '我在哪些方面需要加强学习？',
        '给我一些学习建议'
      ]
    };
  },
  mounted() {
    this.adjustTextareaHeight();
  },
  methods: {
    async handleSend() {
      const text = this.inputText.trim();
      if (!text || this.loading) return;
      await this.sendMessage(text);
    },

    async sendMessage(text) {
      // 添加用户消息
      this.messages.push({
        role: 'user',
        content: text
      });
      this.inputText = '';
      this.scrollToBottom();
      
      // 发送请求
      this.loading = true;
      try {
        const res = await ragChat(text, 5);
        const data = res.data || res;
        const result = data.data || data;
        
        // 添加助手回复
        this.messages.push({
          role: 'assistant',
          content: result.answer || '抱歉，我无法回答这个问题。',
          citations: result.citations || []
        });
      } catch (error) {
        console.error('问答失败:', error);
        this.messages.push({
          role: 'assistant',
          content: '抱歉，处理您的问题时出现错误，请稍后重试。',
          citations: []
        });
      } finally {
        this.loading = false;
        this.scrollToBottom();
      }
    },

    handleNewline() {
      // Shift+Enter 正常换行
    },

    formatMessage(content) {
      if (!content) return '';
      // 使用 marked 渲染 Markdown
      return marked.parse(content);
    },

    showCitationDetail(cite) {
      this.selectedCitation = cite;
    },

    scrollToBottom() {
      this.$nextTick(() => {
        const container = this.$refs.messagesContainer;
        if (container) {
          container.scrollTop = container.scrollHeight;
        }
      });
    },

    adjustTextareaHeight() {
      const textarea = this.$refs.inputArea;
      if (textarea) {
        textarea.style.height = 'auto';
        textarea.style.height = Math.min(textarea.scrollHeight, 120) + 'px';
      }
    }
  },
  watch: {
    inputText() {
      this.adjustTextareaHeight();
    }
  }
};
</script>

<style scoped>
.knowledge-chat-container {
  padding: 24px;
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: linear-gradient(135deg, #f0f4ff 0%, #fef3f2 100%);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  flex-shrink: 0;
}

.header-left .page-title {
  font-size: 24px;
  font-weight: 700;
  color: #1e293b;
  margin: 0;
}

.header-left .page-subtitle {
  font-size: 14px;
  color: #64748b;
  margin: 4px 0 0 0;
}

.chat-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  overflow: hidden;
  min-height: 0;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

/* 欢迎消息 */
.welcome-message {
  text-align: center;
  padding: 60px 20px;
}

.welcome-icon {
  font-size: 64px;
  margin-bottom: 16px;
}

.welcome-message h2 {
  font-size: 22px;
  color: #1e293b;
  margin: 0 0 8px 0;
}

.welcome-message p {
  color: #64748b;
  margin: 0 0 24px 0;
}

.quick-questions {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
}

.quick-title {
  font-size: 14px;
  color: #94a3b8;
  margin: 0;
}

.quick-btn {
  padding: 10px 20px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 20px;
  color: #475569;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.quick-btn:hover {
  background: #2D5BFF;
  color: white;
  border-color: #2D5BFF;
}

/* 消息样式 */
.message {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
}

.message.user {
  flex-direction: row-reverse;
}

.message-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #f1f5f9;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  flex-shrink: 0;
}

.message.user .message-avatar {
  background: #dbeafe;
}

.message-content {
  max-width: 70%;
}

.message.user .message-content {
  text-align: right;
}

.message-text {
  padding: 12px 16px;
  border-radius: 16px;
  font-size: 15px;
  line-height: 1.6;
}

.message.user .message-text {
  background: #2D5BFF;
  color: white;
  border-bottom-right-radius: 4px;
}

.message.assistant .message-text {
  background: #f8fafc;
  color: #1e293b;
  border-bottom-left-radius: 4px;
}

.message-text :deep(p) {
  margin: 0 0 8px 0;
}

.message-text :deep(p:last-child) {
  margin-bottom: 0;
}

.message-text :deep(ul), .message-text :deep(ol) {
  margin: 8px 0;
  padding-left: 20px;
}

.message-text :deep(li) {
  margin: 4px 0;
}

.message-text :deep(strong) {
  font-weight: 600;
}

/* Markdown 标题样式 */
.message-text :deep(h1),
.message-text :deep(h2),
.message-text :deep(h3),
.message-text :deep(h4) {
  margin: 16px 0 8px 0;
  font-weight: 600;
  line-height: 1.4;
}

.message-text :deep(h1) {
  font-size: 1.4em;
}

.message-text :deep(h2) {
  font-size: 1.25em;
}

.message-text :deep(h3) {
  font-size: 1.1em;
}

.message-text :deep(h4) {
  font-size: 1em;
}

/* 代码块样式 */
.message-text :deep(pre) {
  background: #1e293b;
  color: #e2e8f0;
  padding: 12px 16px;
  border-radius: 8px;
  overflow-x: auto;
  margin: 12px 0;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.5;
}

.message-text :deep(code) {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 0.9em;
}

.message-text :deep(:not(pre) > code) {
  background: #e2e8f0;
  color: #be185d;
  padding: 2px 6px;
  border-radius: 4px;
}

.message.user .message-text :deep(:not(pre) > code) {
  background: rgba(255, 255, 255, 0.2);
  color: #fce7f3;
}

/* 引用块样式 */
.message-text :deep(blockquote) {
  border-left: 4px solid #3b82f6;
  margin: 12px 0;
  padding: 8px 16px;
  background: #f1f5f9;
  color: #475569;
  border-radius: 0 8px 8px 0;
}

.message.user .message-text :deep(blockquote) {
  border-left-color: rgba(255, 255, 255, 0.5);
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.9);
}

/* 表格样式 */
.message-text :deep(table) {
  width: 100%;
  border-collapse: collapse;
  margin: 12px 0;
  font-size: 14px;
}

.message-text :deep(th),
.message-text :deep(td) {
  border: 1px solid #e2e8f0;
  padding: 8px 12px;
  text-align: left;
}

.message-text :deep(th) {
  background: #f1f5f9;
  font-weight: 600;
}

/* 分隔线 */
.message-text :deep(hr) {
  border: none;
  border-top: 1px solid #e2e8f0;
  margin: 16px 0;
}

/* 引用来源 */
.citations {
  margin-top: 12px;
  padding: 12px;
  background: #fefce8;
  border-radius: 12px;
  border: 1px solid #fef08a;
}

.citations-title {
  font-size: 13px;
  font-weight: 600;
  color: #854d0e;
  margin-bottom: 8px;
}

.citation-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.citation-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background: white;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.citation-item:hover {
  background: #fef9c3;
}

.citation-main {
  display: flex;
  align-items: center;
  gap: 8px;
}

.citation-title {
  font-size: 13px;
  color: #1e293b;
  font-weight: 500;
}

.citation-category {
  font-size: 11px;
  padding: 2px 6px;
  background: #e0f2fe;
  color: #0369a1;
  border-radius: 4px;
}

.citation-similarity {
  font-size: 12px;
  color: #65a30d;
  font-weight: 500;
}

/* 加载动画 */
.typing-indicator {
  display: flex;
  gap: 4px;
  padding: 12px 16px;
  background: #f8fafc;
  border-radius: 16px;
  border-bottom-left-radius: 4px;
}

.typing-indicator span {
  width: 8px;
  height: 8px;
  background: #94a3b8;
  border-radius: 50%;
  animation: typing 1.4s infinite ease-in-out;
}

.typing-indicator span:nth-child(2) {
  animation-delay: 0.2s;
}

.typing-indicator span:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes typing {
  0%, 60%, 100% { transform: translateY(0); }
  30% { transform: translateY(-8px); }
}

/* 输入区域 */
.chat-input-area {
  padding: 16px 24px 20px;
  border-top: 1px solid #f1f5f9;
  background: #fafafa;
}

.input-wrapper {
  display: flex;
  gap: 12px;
  background: white;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 8px 12px;
}

.input-wrapper textarea {
  flex: 1;
  border: none;
  outline: none;
  resize: none;
  font-size: 15px;
  line-height: 1.5;
  min-height: 24px;
  max-height: 120px;
  padding: 4px 0;
}

.input-wrapper textarea::placeholder {
  color: #94a3b8;
}

.send-btn {
  padding: 8px 20px;
  background: #2D5BFF;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  align-self: flex-end;
}

.send-btn:hover:not(:disabled) {
  background: #1e40af;
}

.send-btn:disabled {
  background: #94a3b8;
  cursor: not-allowed;
}

.input-hint {
  font-size: 12px;
  color: #94a3b8;
  margin: 8px 0 0 0;
  text-align: center;
}

/* 引用详情弹窗 */
.citation-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.citation-modal .modal-content {
  background: white;
  border-radius: 16px;
  width: 480px;
  max-width: 90vw;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #f1f5f9;
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
  color: #1e293b;
}

.close-btn {
  background: none;
  border: none;
  font-size: 24px;
  color: #94a3b8;
  cursor: pointer;
  padding: 0;
  line-height: 1;
}

.close-btn:hover {
  color: #475569;
}

.modal-body {
  padding: 20px;
}

.detail-row {
  display: flex;
  margin-bottom: 12px;
}

.detail-row:last-child {
  margin-bottom: 0;
}

.detail-row .label {
  width: 70px;
  color: #64748b;
  font-size: 14px;
  flex-shrink: 0;
}

.detail-row .value {
  font-size: 14px;
  color: #1e293b;
}

.detail-row .value.summary {
  line-height: 1.6;
}

.modal-footer {
  padding: 16px 20px;
  border-top: 1px solid #f1f5f9;
  display: flex;
  justify-content: flex-end;
}

.btn-primary {
  padding: 10px 24px;
  background: #2D5BFF;
  color: white;
  border-radius: 8px;
  text-decoration: none;
  font-size: 14px;
  transition: background 0.2s;
}

.btn-primary:hover {
  background: #1e40af;
}
</style>
