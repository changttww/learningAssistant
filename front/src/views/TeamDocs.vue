<template>
  <div class="w-full h-full flex flex-col bg-gray-50">
    <!-- 顶部导航 -->
    <div class="bg-white border-b border-gray-200 px-6 py-4 flex items-center justify-between shadow-sm z-10">
      <div class="flex items-center gap-4">
        <button
          @click="$router.back()"
          class="flex items-center px-3 py-2 text-sm font-medium text-gray-600 bg-gray-50 hover:bg-blue-50 hover:text-blue-600 rounded-lg transition-all duration-200"
        >
          <iconify-icon icon="mdi:arrow-left" class="mr-1"></iconify-icon>
          返回
        </button>
        <h1 class="text-xl font-bold text-gray-800 flex items-center gap-2">
          <iconify-icon icon="mdi:file-document-multiple" class="text-orange-500"></iconify-icon>
          团队协作文档
        </h1>
      </div>
      <div class="flex items-center gap-3">
        <span class="text-sm text-gray-500" v-if="saving">
          <iconify-icon icon="mdi:loading" class="animate-spin inline mr-1"></iconify-icon>
          保存中...
        </span>
        <span class="text-sm text-green-600" v-else-if="lastSaved">
          <iconify-icon icon="mdi:check" class="inline mr-1"></iconify-icon>
          已保存 {{ formatTime(lastSaved) }}
        </span>
        <button
          @click="createNewDoc"
          class="flex items-center px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors shadow-sm"
        >
          <iconify-icon icon="mdi:plus" class="mr-1"></iconify-icon>
          新建文档
        </button>
      </div>
    </div>

    <div class="flex-1 flex overflow-hidden">
      <!-- 左侧文档列表 -->
      <div class="w-64 bg-white border-r border-gray-200 flex flex-col">
        <div class="flex border-b border-gray-100">
          <button 
            @click="sidebarTab = 'docs'"
            class="flex-1 py-3 text-sm font-medium transition-colors relative"
            :class="sidebarTab === 'docs' ? 'text-blue-600' : 'text-gray-500 hover:text-gray-700'"
          >
            文档列表
            <div v-if="sidebarTab === 'docs'" class="absolute bottom-0 left-0 w-full h-0.5 bg-blue-600"></div>
          </button>
          <button 
            @click="sidebarTab = 'outline'"
            class="flex-1 py-3 text-sm font-medium transition-colors relative"
            :class="sidebarTab === 'outline' ? 'text-blue-600' : 'text-gray-500 hover:text-gray-700'"
          >
            大纲索引
            <div v-if="sidebarTab === 'outline'" class="absolute bottom-0 left-0 w-full h-0.5 bg-blue-600"></div>
          </button>
        </div>

        <div v-if="sidebarTab === 'docs'" class="flex-1 flex flex-col overflow-hidden">
          <div class="p-4 border-b border-gray-100">
            <div class="relative">
              <iconify-icon icon="mdi:magnify" class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400"></iconify-icon>
              <input
                v-model="searchQuery"
                type="text"
                placeholder="搜索文档..."
                class="w-full pl-9 pr-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:bg-white transition-all"
              />
            </div>
          </div>
          <div class="flex-1 overflow-y-auto p-2 space-y-1">
            <div
              v-for="doc in filteredDocs"
              :key="doc.id"
              @click="selectDoc(doc)"
              class="group flex items-center justify-between px-3 py-2.5 rounded-lg cursor-pointer transition-all"
              :class="currentDoc?.id === doc.id ? 'bg-blue-50 text-blue-700' : 'hover:bg-gray-50 text-gray-700'"
            >
              <div class="flex items-center gap-3 overflow-hidden">
                <iconify-icon
                  :icon="currentDoc?.id === doc.id ? 'mdi:file-document-edit' : 'mdi:file-document-outline'"
                  class="flex-shrink-0"
                  :class="currentDoc?.id === doc.id ? 'text-blue-600' : 'text-gray-400'"
                ></iconify-icon>
                <div class="flex flex-col overflow-hidden">
                  <span class="truncate font-medium text-sm">{{ doc.title || '无标题文档' }}</span>
                  <span class="text-xs text-gray-400 truncate">{{ formatDate(doc.updatedAt) }}</span>
                </div>
              </div>
              <button
                @click.stop="deleteDoc(doc)"
                class="opacity-0 group-hover:opacity-100 p-1 hover:bg-red-50 hover:text-red-600 rounded transition-all"
                title="删除"
              >
                <iconify-icon icon="mdi:trash-can-outline"></iconify-icon>
              </button>
            </div>
            
            <div v-if="filteredDocs.length === 0" class="text-center py-8 text-gray-400 text-sm">
              <p>暂无文档</p>
            </div>
          </div>
        </div>

        <div v-else class="flex-1 overflow-y-auto p-4">
          <div v-if="!currentDoc" class="text-center text-gray-400 text-sm py-8">
            请先选择文档
          </div>
          <div v-else-if="outline.length === 0" class="text-center text-gray-400 text-sm py-8">
            暂无标题
          </div>
          <div v-else class="space-y-1">
            <div
              v-for="(item, index) in outline"
              :key="index"
              @click="scrollToHeading(item.pos)"
              class="cursor-pointer hover:bg-gray-50 hover:text-blue-600 py-1.5 rounded text-sm text-gray-600 transition-colors truncate"
              :class="{
                'pl-2 font-medium': item.level === 1,
                'pl-6': item.level === 2,
                'pl-10 text-xs': item.level === 3
              }"
            >
              {{ item.text }}
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧编辑器 -->
      <div class="flex-1 flex flex-col bg-white overflow-hidden relative" v-if="currentDoc">
        <!-- 标题输入 -->
        <div class="px-8 pt-8 pb-4">
          <input
            v-model="currentDoc.title"
            @input="autoSave"
            @keydown.enter.prevent="focusEditor"
            type="text"
            placeholder="请输入文档标题"
            class="w-full text-3xl font-bold text-gray-800 placeholder-gray-300 border-none outline-none bg-transparent"
          />
        </div>
        
        <!-- 工具栏 -->
        <div class="px-8 py-2 border-y border-gray-100 flex items-center gap-2 flex-wrap bg-gray-50/50" v-if="editor">
          <button @click="editor.chain().focus().setParagraph().run()" :class="{ 'is-active': editor.isActive('paragraph') }" class="editor-btn" title="正文">
            <iconify-icon icon="mdi:format-paragraph"></iconify-icon>
          </button>
          <div class="w-px h-4 bg-gray-300 mx-1"></div>
          <button @click="editor.chain().focus().toggleHeading({ level: 1 }).run()" :class="{ 'is-active': editor.isActive('heading', { level: 1 }) }" class="editor-btn" title="一级标题">
            H1
          </button>
          <button @click="editor.chain().focus().toggleHeading({ level: 2 }).run()" :class="{ 'is-active': editor.isActive('heading', { level: 2 }) }" class="editor-btn" title="二级标题">
            H2
          </button>
          <button @click="editor.chain().focus().toggleHeading({ level: 3 }).run()" :class="{ 'is-active': editor.isActive('heading', { level: 3 }) }" class="editor-btn" title="三级标题">
            H3
          </button>
          <div class="w-px h-4 bg-gray-300 mx-1"></div>
          <button @click="editor.chain().focus().toggleBold().run()" :class="{ 'is-active': editor.isActive('bold') }" class="editor-btn" title="加粗">
            <iconify-icon icon="mdi:format-bold"></iconify-icon>
          </button>
          <button @click="editor.chain().focus().toggleItalic().run()" :class="{ 'is-active': editor.isActive('italic') }" class="editor-btn" title="斜体">
            <iconify-icon icon="mdi:format-italic"></iconify-icon>
          </button>
          <button @click="editor.chain().focus().toggleStrike().run()" :class="{ 'is-active': editor.isActive('strike') }" class="editor-btn" title="删除线">
            <iconify-icon icon="mdi:format-strikethrough"></iconify-icon>
          </button>
          <div class="w-px h-4 bg-gray-300 mx-1"></div>
          <button @click="editor.chain().focus().toggleBulletList().run()" :class="{ 'is-active': editor.isActive('bulletList') }" class="editor-btn" title="无序列表">
            <iconify-icon icon="mdi:format-list-bulleted"></iconify-icon>
          </button>
          <button @click="editor.chain().focus().toggleOrderedList().run()" :class="{ 'is-active': editor.isActive('orderedList') }" class="editor-btn" title="有序列表">
            <iconify-icon icon="mdi:format-list-numbered"></iconify-icon>
          </button>
          <div class="w-px h-4 bg-gray-300 mx-1"></div>
          <button @click="editor.chain().focus().toggleBlockquote().run()" :class="{ 'is-active': editor.isActive('blockquote') }" class="editor-btn" title="引用">
            <iconify-icon icon="mdi:format-quote-close"></iconify-icon>
          </button>
          <button @click="editor.chain().focus().setHorizontalRule().run()" class="editor-btn" title="分割线">
            <iconify-icon icon="mdi:minus"></iconify-icon>
          </button>
          <button @click="editor.chain().focus().setHardBreak().run()" class="editor-btn" title="换行">
            <iconify-icon icon="mdi:keyboard-return"></iconify-icon>
          </button>
        </div>

        <!-- 编辑区域 -->
        <div class="flex-1 overflow-y-auto">
          <editor-content :editor="editor" />
        </div>
      </div>

      <!-- 空状态 -->
      <div v-else class="flex-1 flex flex-col items-center justify-center bg-gray-50 text-gray-400">
        <iconify-icon icon="mdi:file-document-outline" width="64" height="64" class="mb-4 opacity-50"></iconify-icon>
        <p class="text-lg font-medium">选择或创建一个文档开始协作</p>
      </div>
    </div>
  </div>
</template>

<script>
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import { debounce } from 'lodash'

export default {
  name: 'TeamDocs',
  components: {
    EditorContent,
  },
  data() {
    return {
      teamId: null,
      docs: [],
      currentDoc: null,
      searchQuery: '',
      editor: null,
      saving: false,
      lastSaved: null,
      debouncedSave: null,
      sidebarTab: 'docs', // 'docs' or 'outline'
      outline: [],
    }
  },
  computed: {
    filteredDocs() {
      if (!this.searchQuery) return this.docs;
      const query = this.searchQuery.toLowerCase();
      return this.docs.filter(doc => 
        (doc.title || '').toLowerCase().includes(query) || 
        (doc.content || '').toLowerCase().includes(query)
      );
    }
  },
  created() {
    this.teamId = this.$route.params.teamId || this.$route.query.teamId;
    if (!this.teamId) {
      const storedTeamId = sessionStorage.getItem("currentTeamId");
      if (storedTeamId) {
        this.teamId = storedTeamId;
      } else {
        alert("未找到团队信息");
        this.$router.push({ name: "TeamTasks" });
        return;
      }
    }
    
    this.loadDocs();
    
    // 防抖保存
    this.debouncedSave = debounce(this.saveDocs, 1000);
  },
  mounted() {
    this.editor = useEditor({
      content: '',
      extensions: [
        StarterKit.configure({
          heading: {
            levels: [1, 2, 3],
          },
        }),
      ],
      editorProps: {
        attributes: {
          class: 'prose prose-blue max-w-none focus:outline-none min-h-[500px] px-8 py-6',
        },
      },
      onUpdate: ({ editor }) => {
        if (this.currentDoc) {
          this.currentDoc.content = editor.getHTML();
          this.currentDoc.updatedAt = Date.now();
          this.autoSave();
          this.updateOutline();
        }
      },
      onSelectionUpdate: ({ editor }) => {
        // 可以在这里高亮当前所在的大纲位置
      }
    })
  },
  beforeUnmount() {
    if (this.editor) {
      this.editor.destroy()
    }
  },
  methods: {
    focusEditor() {
      if (this.editor) {
        this.editor.chain().focus().run();
      }
    },
    updateOutline() {
      if (!this.editor) return;
      const headings = [];
      this.editor.state.doc.descendants((node, pos) => {
        if (node.type.name === 'heading') {
          headings.push({
            level: node.attrs.level,
            text: node.textContent,
            pos
          });
        }
      });
      this.outline = headings;
    },
    scrollToHeading(pos) {
      if (!this.editor) return;
      this.editor.chain().focus().setTextSelection(pos).run();
      const dom = this.editor.view.domAtPos(pos).node;
      if (dom && dom.scrollIntoView) {
        dom.scrollIntoView({ behavior: 'smooth', block: 'start' });
      }
    },
    loadDocs() {
      try {
        const key = `team_docs_${this.teamId}`;
        const saved = localStorage.getItem(key);
        if (saved) {
          this.docs = JSON.parse(saved);
          // 按更新时间排序
          this.docs.sort((a, b) => b.updatedAt - a.updatedAt);
        } else {
          // 默认创建一个欢迎文档
          this.docs = [{
            id: Date.now(),
            title: '欢迎使用团队文档',
            content: '<h2>👋 欢迎使用团队协作文档</h2><p>在这里，您可以：</p><ul><li>记录团队会议纪要</li><li>共享项目资料</li><li>协同编辑创意方案</li></ul><p>所有数据将自动保存在本地。</p>',
            createdAt: Date.now(),
            updatedAt: Date.now()
          }];
          this.saveDocs();
        }
      } catch (e) {
        console.error("加载文档失败", e);
      }
    },
    saveDocs() {
      this.saving = true;
      try {
        const key = `team_docs_${this.teamId}`;
        localStorage.setItem(key, JSON.stringify(this.docs));
        this.lastSaved = Date.now();
      } catch (e) {
        console.error("保存文档失败", e);
      } finally {
        setTimeout(() => {
          this.saving = false;
        }, 500);
      }
    },
    autoSave() {
      this.saving = true;
      this.debouncedSave();
    },
    createNewDoc() {
      const newDoc = {
        id: Date.now(),
        title: '',
        content: '',
        createdAt: Date.now(),
        updatedAt: Date.now()
      };
      this.docs.unshift(newDoc);
      this.selectDoc(newDoc);
      this.saveDocs();
    },
    selectDoc(doc) {
      this.currentDoc = doc;
      if (this.editor) {
        this.editor.commands.setContent(doc.content || '');
        this.updateOutline();
      }
    },
    deleteDoc(doc) {
      if (!confirm(`确定要删除文档 "${doc.title || '无标题'}" 吗？`)) return;
      
      const index = this.docs.findIndex(d => d.id === doc.id);
      if (index > -1) {
        this.docs.splice(index, 1);
        if (this.currentDoc?.id === doc.id) {
          this.currentDoc = null;
          this.editor.commands.setContent('');
        }
        this.saveDocs();
      }
    },
    formatDate(timestamp) {
      if (!timestamp) return '';
      return new Date(timestamp).toLocaleDateString('zh-CN', {
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      });
    },
    formatTime(timestamp) {
      if (!timestamp) return '';
      return new Date(timestamp).toLocaleTimeString('zh-CN', {
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
      });
    }
  }
}
</script>

<style scoped>
.editor-btn {
  @apply p-1.5 rounded text-gray-600 hover:bg-gray-200 hover:text-gray-900 transition-colors text-sm font-medium min-w-[32px] flex items-center justify-center;
}
.editor-btn.is-active {
  @apply bg-blue-100 text-blue-600;
}

/* Prose Mirror 样式覆盖 */
:deep(.ProseMirror) {
  min-height: 300px;
}
:deep(.ProseMirror p) {
  margin-bottom: 1em;
}
:deep(.ProseMirror p.is-editor-empty:first-child::before) {
  color: #adb5bd;
  content: attr(data-placeholder);
  float: left;
  height: 0;
  pointer-events: none;
}
</style>
