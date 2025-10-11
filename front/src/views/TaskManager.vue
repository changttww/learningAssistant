<template>
  <div class="bg-gray-50 min-h-screen py-8">
    <div class="w-[1440px] mx-auto">
      <div class="flex gap-6">
        <!-- 左侧：学习档案卡 -->
        <aside class="w-80 flex-shrink-0 sidebar-shadow rounded-2xl bg-white p-6">
          <!-- 统一的学习档案卡 -->
          <div class="card mb-6 p-6 text-center bg-gradient-to-br from-blue-50 to-indigo-50">
            <!-- 个人信息区 -->
            <div class="flex justify-center mb-4">
              <div class="w-24 h-24 rounded-full bg-blue-100 flex items-center justify-center overflow-hidden">
                <div class="w-24 h-24 bg-gray-300 rounded-full"></div>
              </div>
            </div>
            <h2 class="text-xl font-bold mb-2">李明</h2>
            <div class="flex justify-center items-center mb-4">
              <iconify-icon icon="mdi:medal" class="text-orange-500 text-xl"></iconify-icon>
              <span class="ml-1 text-orange-500 font-medium">黄金会员</span>
            </div>
            
            <!-- 学习统计 -->
            <div class="grid grid-cols-3 gap-3 mb-6">
              <div class="bg-white rounded-lg p-3">
                <p class="text-gray-500 text-xs mb-1">连续签到</p>
                <p class="text-lg font-bold text-blue-600">28<span class="text-xs">天</span></p>
              </div>
              <div class="bg-white rounded-lg p-3">
                <p class="text-gray-500 text-xs mb-1">正在学习</p>
                <p class="text-lg font-bold text-blue-600">3<span class="text-xs">门课</span></p>
              </div>
              <div class="bg-white rounded-lg p-3">
                <p class="text-gray-500 text-xs mb-1">总任务</p>
                <p class="text-lg font-bold text-blue-600">102<span class="text-xs">项</span></p>
              </div>
            </div>
            
            <!-- 积分展示 -->
            <div class="mb-6">
              <div class="flex justify-between items-center mb-3">
                <h3 class="font-bold text-gray-700 flex items-center">
                  <iconify-icon icon="mdi:crystal-ball" class="mr-2 text-blue-600"></iconify-icon> 
                  我的积分
                </h3>
                <span class="badge inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                  TOP 15%
                </span>
              </div>
              
              <div class="flex justify-center mb-4">
                <div class="relative">
                  <div class="w-20 h-20 rounded-full relative flex items-center justify-center">
                    <div class="absolute inset-0 bg-blue-600 bg-opacity-10 rounded-full"></div>
                    <span class="text-2xl font-bold text-blue-600 relative z-10">3860</span>
                  </div>
                  <div class="absolute bottom-0 right-0 transform translate-y-1">
                    <iconify-icon icon="mdi:star-four-points" class="text-yellow-400 text-lg"></iconify-icon>
                  </div>
                </div>
              </div>
            </div>
            
            <!-- 成长等级进度条 -->
            <div class="mb-4">
              <div class="flex justify-between items-center mb-2">
                <span class="text-sm font-medium text-gray-700">学霸 Lv.4</span>
                <span class="text-xs text-gray-500">距离下一级</span>
              </div>
              <div class="w-full bg-gray-200 rounded-full h-3 mb-2">
                <div class="bg-gradient-to-r from-blue-500 to-purple-500 h-3 rounded-full" style="width: 75%"></div>
              </div>
              <div class="text-xs text-gray-600">还需 250 积分升级到 Lv.5</div>
            </div>
            
            <!-- 操作按钮 -->
            <div class="flex gap-2">
              <button class="flex-1 px-3 py-2 bg-white text-blue-600 border border-blue-200 rounded-lg font-medium text-sm hover:bg-blue-50">
                查看明细
              </button>
              <button class="flex-1 px-3 py-2 bg-blue-600 text-white rounded-lg font-medium text-sm hover:bg-blue-700 flex items-center justify-center">
                <iconify-icon icon="mdi:gift-outline" class="mr-1"></iconify-icon>
                兑换奖品
              </button>
            </div>
          </div>
          
          <!-- 快速导航 -->
          <div class="space-y-3">
            <router-link to="/personal-tasks" class="flex items-center p-3 text-gray-700 hover:bg-gray-50 rounded-lg transition">
              <iconify-icon icon="mdi:calendar-check" class="text-xl text-green-500 mr-3"></iconify-icon>
              <span>今日任务</span>
            </router-link>
            <button @click="showAchievements = true" class="w-full flex items-center p-3 text-gray-700 hover:bg-gray-50 rounded-lg transition">
              <iconify-icon icon="mdi:trophy" class="text-xl text-yellow-500 mr-3"></iconify-icon>
              <span>我的成就</span>
            </button>
            <button @click="showSettings = true" class="w-full flex items-center p-3 text-gray-700 hover:bg-gray-50 rounded-lg transition">
              <iconify-icon icon="mdi:cog" class="text-xl text-gray-500 mr-3"></iconify-icon>
              <span>系统设置</span>
            </button>
          </div>
        </aside>
        
        <!-- 中间：动态学习看板 -->
        <main class="flex-1">
          <!-- 看板标题和时间切换 -->
          <div class="flex items-center justify-between mb-5">
            <h2 class="text-2xl font-bold text-gray-700">动态学习看板</h2>
            <div class="inline-flex bg-gray-100 p-1 rounded-xl">
              <button 
                class="px-4 py-1 rounded-lg font-medium transition-colors"
                :class="activeTimeFilter === 'week' ? 'bg-white shadow-sm text-blue-600' : 'text-gray-600 hover:text-gray-800'"
                @click="setTimeFilter('week')"
              >
                周
              </button>
              <button 
                class="px-4 py-1 rounded-lg font-medium transition-colors"
                :class="activeTimeFilter === 'month' ? 'bg-white shadow-sm text-blue-600' : 'text-gray-600 hover:text-gray-800'"
                @click="setTimeFilter('month')"
              >
                月
              </button>
              <button 
                class="px-4 py-1 rounded-lg font-medium transition-colors"
                :class="activeTimeFilter === 'quarter' ? 'bg-white shadow-sm text-blue-600' : 'text-gray-600 hover:text-gray-800'"
                @click="setTimeFilter('quarter')"
              >
                季度
              </button>
            </div>
          </div>
          
          <!-- 任务进度图表区 -->
          <div class="grid grid-cols-3 gap-4 mb-6">
            <!-- 环形进度图 -->
            <div class="card col-span-1 p-5 cursor-pointer hover:shadow-lg transition-shadow" @click="showTaskDetails">
              <div class="flex justify-between mb-3">
                <div>
                  <p class="text-gray-500 mb-2">整体完成率</p>
                  <h3 class="text-3xl font-bold text-blue-600">{{ currentTimeData.completionRate }}%</h3>
                </div>
                <div class="h-16 w-16" ref="ringProgress"></div>
              </div>
              <div class="mt-4">
                <div class="flex justify-between text-sm text-gray-500 mb-1">
                  <span>已完成</span>
                  <span>{{ currentTimeData.completedTasks }}/{{ currentTimeData.totalTasks }}</span>
                </div>
                <div class="w-full h-3 bg-gray-200 rounded-full">
                  <div class="h-full rounded-full bg-green-500" :style="`width: ${currentTimeData.completionRate}%`"></div>
                </div>
              </div>
            </div>
            
            <!-- 柱状图容器 -->
            <div class="card col-span-2 p-5">
              <div class="h-64" ref="taskProgressChart"></div>
            </div>
          </div>
          
          <!-- 任务标签页 -->
          <div class="mb-4">
            <div class="flex space-x-1 bg-gray-100 p-1 rounded-lg w-fit">
              <button 
                class="px-4 py-2 rounded-md font-medium transition-colors"
                :class="activeTab === 'inProgress' ? 'bg-white shadow-sm text-blue-600' : 'text-gray-600 hover:text-gray-800'"
                @click="activeTab = 'inProgress'"
              >
                进行中 (3)
              </button>
              <button 
                class="px-4 py-2 rounded-md font-medium transition-colors"
                :class="activeTab === 'pending' ? 'bg-white shadow-sm text-orange-600' : 'text-gray-600 hover:text-gray-800'"
                @click="activeTab = 'pending'"
              >
                待开始 (5)
              </button>
              <button 
                class="px-4 py-2 rounded-md font-medium transition-colors"
                :class="activeTab === 'completed' ? 'bg-white shadow-sm text-green-600' : 'text-gray-600 hover:text-gray-800'"
                @click="activeTab = 'completed'"
              >
                已完成 (86)
              </button>
            </div>
          </div>
          
          <!-- 任务列表 -->
          <div class="space-y-4">
            <!-- 进行中的任务 -->
            <div v-show="activeTab === 'inProgress'">
              <!-- 任务卡片1 -->
              <div class="card p-5 flex items-center hover:shadow-md transition-shadow">
                <div class="w-12 h-12 rounded-xl bg-blue-50 flex items-center justify-center mr-4">
                  <iconify-icon icon="mdi:code-tags" class="text-2xl text-blue-600"></iconify-icon>
                </div>
                <div class="flex-1">
                  <h4 class="font-bold text-gray-800">前端开发课程学习</h4>
                  <div class="flex items-center mt-1">
                    <div class="text-xs text-gray-500">最后学习：今天 12:30</div>
                    <span class="ml-2 px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded-full">高优先级</span>
                  </div>
                </div>
                <div class="w-1/3">
                  <div class="flex justify-between text-sm text-gray-500 mb-1">
                    <span>学习进度</span>
                    <span>64%</span>
                  </div>
                  <div class="h-2 bg-gray-200 rounded-full">
                    <div class="h-full rounded-full bg-blue-600" style="width: 64%"></div>
                  </div>
                </div>
                <button class="ml-6 bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-lg">继续</button>
              </div>
              
              <!-- 任务卡片2 -->
              <div class="card p-5 flex items-center hover:shadow-md transition-shadow">
                <div class="w-12 h-12 rounded-xl bg-orange-50 flex items-center justify-center mr-4">
                  <iconify-icon icon="mdi:translate" class="text-2xl text-orange-600"></iconify-icon>
                </div>
                <div class="flex-1">
                  <h4 class="font-bold text-gray-800">英语四级备考测试</h4>
                  <div class="flex items-center mt-1">
                    <div class="text-xs text-gray-500">最后测试：昨天 16:45</div>
                    <span class="ml-2 px-2 py-1 bg-orange-100 text-orange-800 text-xs rounded-full">中优先级</span>
                  </div>
                </div>
                <div class="w-1/3">
                  <div class="flex justify-between text-sm text-gray-500 mb-1">
                    <span>完成进度</span>
                    <span>48%</span>
                  </div>
                  <div class="h-2 bg-gray-200 rounded-full">
                    <div class="h-full rounded-full bg-orange-500" style="width: 48%"></div>
                  </div>
                </div>
                <button class="ml-6 bg-orange-500 hover:bg-orange-600 text-white font-medium py-2 px-4 rounded-lg">开始</button>
              </div>
              
              <!-- 任务卡片3 -->
              <div class="card p-5 flex items-center hover:shadow-md transition-shadow">
                <div class="w-12 h-12 rounded-xl bg-green-50 flex items-center justify-center mr-4">
                  <iconify-icon icon="mdi:book-open-page-variant" class="text-2xl text-green-600"></iconify-icon>
                </div>
                <div class="flex-1">
                  <h4 class="font-bold text-gray-800">每周阅读计划 - 《Web开发实战》</h4>
                  <div class="flex items-center mt-1">
                    <div class="text-xs text-gray-500">本周任务：3/5章节</div>
                    <span class="ml-2 px-2 py-1 bg-green-100 text-green-800 text-xs rounded-full">低优先级</span>
                  </div>
                </div>
                <div class="w-1/3">
                  <div class="flex justify-between text-sm text-gray-500 mb-1">
                    <span>进度</span>
                    <span>82%</span>
                  </div>
                  <div class="h-2 bg-gray-200 rounded-full">
                    <div class="h-full rounded-full bg-green-500" style="width: 82%"></div>
                  </div>
                </div>
                <button class="ml-6 bg-gray-800 hover:bg-black text-white font-medium py-2 px-4 rounded-lg">阅读</button>
              </div>
            </div>
            
            <!-- 待开始的任务 -->
            <div v-show="activeTab === 'pending'">
              <div class="card p-5 flex items-center hover:shadow-md transition-shadow">
                <div class="w-12 h-12 rounded-xl bg-purple-50 flex items-center justify-center mr-4">
                  <iconify-icon icon="mdi:database" class="text-2xl text-purple-600"></iconify-icon>
                </div>
                <div class="flex-1">
                  <h4 class="font-bold text-gray-800">数据库设计与优化</h4>
                  <div class="flex items-center mt-1">
                    <div class="text-xs text-gray-500">计划开始：明天 09:00</div>
                    <span class="ml-2 px-2 py-1 bg-purple-100 text-purple-800 text-xs rounded-full">待开始</span>
                  </div>
                </div>
                <div class="w-1/3">
                  <div class="flex justify-between text-sm text-gray-500 mb-1">
                    <span>预计时长</span>
                    <span>40小时</span>
                  </div>
                  <div class="h-2 bg-gray-200 rounded-full">
                    <div class="h-full rounded-full bg-gray-300" style="width: 0%"></div>
                  </div>
                </div>
                <button class="ml-6 bg-purple-600 hover:bg-purple-700 text-white font-medium py-2 px-4 rounded-lg">开始学习</button>
              </div>
            </div>
            
            <!-- 已完成的任务 -->
            <div v-show="activeTab === 'completed'">
              <div class="card p-5 flex items-center hover:shadow-md transition-shadow">
                <div class="w-12 h-12 rounded-xl bg-green-50 flex items-center justify-center mr-4">
                  <iconify-icon icon="mdi:check-circle" class="text-2xl text-green-600"></iconify-icon>
                </div>
                <div class="flex-1">
                  <h4 class="font-bold text-gray-800">JavaScript基础课程</h4>
                  <div class="flex items-center mt-1">
                    <div class="text-xs text-gray-500">完成时间：2024-01-15 18:30</div>
                    <span class="ml-2 px-2 py-1 bg-green-100 text-green-800 text-xs rounded-full">已完成</span>
                  </div>
                </div>
                <div class="w-1/3">
                  <div class="flex justify-between text-sm text-gray-500 mb-1">
                    <span>最终得分</span>
                    <span>95分</span>
                  </div>
                  <div class="h-2 bg-gray-200 rounded-full">
                    <div class="h-full rounded-full bg-green-500" style="width: 100%"></div>
                  </div>
                </div>
                <button class="ml-6 bg-gray-500 text-white font-medium py-2 px-4 rounded-lg">查看证书</button>
              </div>
            </div>
          </div>
        </main>
        
        <!-- 右侧：学习互动面板 -->
        <aside class="w-80 flex-shrink-0 sidebar-shadow rounded-2xl bg-white p-6">
          <!-- 互动面板标题 -->
          <div class="flex justify-between items-center mb-5">
            <h3 class="text-xl font-bold text-gray-700">学习互动</h3>
            <div class="flex items-center gap-2">
              <div class="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
              <span class="text-sm text-green-600">3人在线</span>
            </div>
          </div>
          
          <!-- 最近互动 -->
          <div class="mb-6">
            <div class="flex justify-between items-center mb-3">
              <h4 class="font-medium text-gray-700">最近互动</h4>
              <span class="text-xs text-blue-600 cursor-pointer hover:underline">查看全部</span>
            </div>
            <div class="space-y-2">
              <div class="flex items-center p-2 bg-blue-50 rounded-lg">
                <div class="w-8 h-8 rounded-full bg-gray-300 mr-2"></div>
                <div class="flex-1">
                  <div class="text-sm font-medium">张伟</div>
                  <div class="text-xs text-gray-500">点赞了你的学习笔记</div>
                </div>
                <div class="text-xs text-gray-400">2分钟前</div>
              </div>
              <div class="flex items-center p-2 bg-green-50 rounded-lg">
                <div class="w-8 h-8 rounded-full bg-gray-300 mr-2"></div>
                <div class="flex-1">
                  <div class="text-sm font-medium">刘燕</div>
                  <div class="text-xs text-gray-500">评论了你的项目</div>
                </div>
                <div class="text-xs text-gray-400">5分钟前</div>
              </div>
            </div>
          </div>
          
          <!-- 点赞排行 -->
          <div class="mb-6">
            <div class="flex justify-between items-center mb-3">
              <h4 class="font-medium text-gray-700">本周点赞排行</h4>
              <iconify-icon icon="mdi:trophy" class="text-yellow-500"></iconify-icon>
            </div>
            <div class="space-y-2">
              <div class="flex items-center justify-between p-2 bg-yellow-50 rounded-lg">
                <div class="flex items-center">
                  <div class="w-6 h-6 bg-yellow-500 text-white rounded-full flex items-center justify-center text-xs font-bold mr-2">1</div>
                  <div class="w-8 h-8 rounded-full bg-gray-300 mr-2"></div>
                  <span class="text-sm font-medium">王浩</span>
                </div>
                <div class="flex items-center text-yellow-600">
                  <iconify-icon icon="mdi:thumb-up" class="mr-1"></iconify-icon>
                  <span class="text-sm font-bold">128</span>
                </div>
              </div>
              <div class="flex items-center justify-between p-2 bg-gray-50 rounded-lg">
                <div class="flex items-center">
                  <div class="w-6 h-6 bg-gray-400 text-white rounded-full flex items-center justify-center text-xs font-bold mr-2">2</div>
                  <div class="w-8 h-8 rounded-full bg-gray-300 mr-2"></div>
                  <span class="text-sm font-medium">李明</span>
                </div>
                <div class="flex items-center text-gray-600">
                  <iconify-icon icon="mdi:thumb-up" class="mr-1"></iconify-icon>
                  <span class="text-sm font-bold">95</span>
                </div>
              </div>
            </div>
          </div>
          
          <!-- 学习伙伴 -->
          <div class="mb-6">
            <div class="flex justify-between items-center mb-3">
              <h4 class="font-medium text-gray-700">学习伙伴</h4>
              <button class="text-blue-600 hover:text-blue-800 flex items-center text-sm">
                <iconify-icon icon="mdi:plus-circle" class="mr-1"></iconify-icon>
                添加
              </button>
            </div>
            
            <!-- 好友列表 -->
            <div class="space-y-3">
              <!-- 好友1 -->
              <div class="friend-card card p-3 flex items-center cursor-pointer hover:bg-blue-50 rounded-xl">
                <div class="relative">
                  <div class="w-10 h-10 rounded-full bg-gray-300"></div>
                  <div class="absolute bottom-0 right-0 w-3 h-3 rounded-full bg-green-500 border-2 border-white"></div>
                </div>
                <div class="ml-3 flex-1">
                  <div class="font-medium text-sm">张伟</div>
                  <div class="text-xs text-gray-500">前端开发课程学习中</div>
                </div>
                <div class="flex gap-1">
                  <button class="w-7 h-7 flex items-center justify-center bg-gray-100 text-gray-600 rounded-full hover:bg-blue-100 text-sm">
                    <iconify-icon icon="mdi:thumb-up"></iconify-icon>
                  </button>
                  <button class="w-7 h-7 flex items-center justify-center bg-blue-100 text-blue-600 rounded-full hover:bg-blue-200 text-sm">
                    <iconify-icon icon="mdi:message"></iconify-icon>
                  </button>
                </div>
              </div>
              
              <!-- 好友2 -->
              <div class="friend-card card p-3 flex items-center cursor-pointer hover:bg-blue-50 rounded-xl">
                <div class="relative">
                  <div class="w-10 h-10 rounded-full bg-gray-300"></div>
                  <div class="absolute bottom-0 right-0 w-3 h-3 rounded-full bg-gray-300 border-2 border-white"></div>
                </div>
                <div class="ml-3 flex-1">
                  <div class="font-medium text-sm">刘燕</div>
                  <div class="text-xs text-gray-500">英语学习进行中</div>
                </div>
                <div class="flex gap-1">
                  <button class="w-7 h-7 flex items-center justify-center bg-gray-100 text-gray-600 rounded-full hover:bg-blue-100 text-sm">
                    <iconify-icon icon="mdi:thumb-up"></iconify-icon>
                  </button>
                  <button class="w-7 h-7 flex items-center justify-center bg-blue-100 text-blue-600 rounded-full hover:bg-blue-200 text-sm">
                    <iconify-icon icon="mdi:message"></iconify-icon>
                  </button>
                </div>
              </div>
            </div>
          </div>
          
          <!-- 评论提醒 -->
          <div class="mb-4">
            <div class="flex justify-between items-center mb-3">
              <h4 class="font-medium text-gray-700">评论提醒</h4>
              <div class="w-5 h-5 bg-red-500 text-white rounded-full flex items-center justify-center text-xs">3</div>
            </div>
            <div class="space-y-2">
              <div class="p-2 bg-red-50 border-l-4 border-red-400 rounded">
                <div class="text-sm font-medium text-red-800">新评论</div>
                <div class="text-xs text-red-600">王浩评论了你的Vue项目</div>
              </div>
            </div>
          </div>
          
          <!-- 互动消息 -->
          <div>
            <div class="flex justify-between items-center mb-3">
              <h4 class="font-medium text-gray-700">互动消息</h4>
              <span class="text-blue-600 hover:text-blue-800 text-xs cursor-pointer">查看全部</span>
            </div>
            
            <!-- 消息1 -->
            <div class="flex mb-3">
              <div class="mr-2">
                <div class="w-8 h-8 rounded-full bg-gray-300"></div>
              </div>
              <div class="flex-1">
                <div class="bubble-left bg-gray-100 p-2 rounded-xl rounded-tl-none relative max-w-[80%]">
                  <div class="font-medium text-sm">陈敏</div>
                  <p class="text-xs">你的前端项目太棒了！😄</p>
                </div>
                <div class="text-xs text-gray-500 mt-1">1小时前</div>
              </div>
            </div>
            
            <!-- 消息输入框 -->
            <div class="mt-3">
              <div class="flex items-center">
                <input type="text" placeholder="回复消息..." class="flex-1 bg-gray-100 border-0 focus:ring-0 rounded-full py-2 px-3 text-sm">
                <button class="w-8 h-8 flex items-center justify-center bg-blue-600 text-white rounded-full ml-2">
                  <iconify-icon icon="mdi:send" class="text-sm"></iconify-icon>
                </button>
              </div>
            </div>
          </div>
        </aside>
      </div>
    </div>
    
    <!-- 我的成就弹窗 -->
    <div v-if="showAchievements" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click="showAchievements = false">
      <div class="bg-white rounded-2xl p-6 w-[600px] max-h-[80vh] overflow-y-auto" @click.stop>
        <div class="flex justify-between items-center mb-6">
          <h2 class="text-2xl font-bold text-gray-800 flex items-center">
            <iconify-icon icon="mdi:trophy" class="text-yellow-500 mr-2"></iconify-icon>
            我的成就
          </h2>
          <button @click="showAchievements = false" class="text-gray-400 hover:text-gray-600">
            <iconify-icon icon="mdi:close" class="text-2xl"></iconify-icon>
          </button>
        </div>
        
        <!-- 成就统计 -->
        <div class="grid grid-cols-3 gap-4 mb-6">
          <div class="bg-gradient-to-br from-yellow-50 to-orange-50 p-4 rounded-xl text-center">
            <iconify-icon icon="mdi:medal" class="text-3xl text-yellow-500 mb-2"></iconify-icon>
            <div class="text-2xl font-bold text-gray-800">12</div>
            <div class="text-sm text-gray-600">已获得成就</div>
          </div>
          <div class="bg-gradient-to-br from-blue-50 to-indigo-50 p-4 rounded-xl text-center">
            <iconify-icon icon="mdi:star" class="text-3xl text-blue-500 mb-2"></iconify-icon>
            <div class="text-2xl font-bold text-gray-800">3860</div>
            <div class="text-sm text-gray-600">成就积分</div>
          </div>
          <div class="bg-gradient-to-br from-green-50 to-emerald-50 p-4 rounded-xl text-center">
            <iconify-icon icon="mdi:target" class="text-3xl text-green-500 mb-2"></iconify-icon>
            <div class="text-2xl font-bold text-gray-800">85%</div>
            <div class="text-sm text-gray-600">完成度</div>
          </div>
        </div>
        
        <!-- 成就列表 -->
        <div class="space-y-4">
          <h3 class="text-lg font-bold text-gray-700 mb-3">最近获得</h3>
          
          <!-- 成就项目 -->
          <div class="flex items-center p-4 bg-yellow-50 border border-yellow-200 rounded-xl">
            <div class="w-12 h-12 bg-yellow-500 rounded-full flex items-center justify-center mr-4">
              <iconify-icon icon="mdi:school" class="text-white text-xl"></iconify-icon>
            </div>
            <div class="flex-1">
              <h4 class="font-bold text-gray-800">学习达人</h4>
              <p class="text-sm text-gray-600">连续学习30天</p>
            </div>
            <div class="text-right">
              <div class="text-yellow-600 font-bold">+500积分</div>
              <div class="text-xs text-gray-500">2024-01-15</div>
            </div>
          </div>
          
          <div class="flex items-center p-4 bg-blue-50 border border-blue-200 rounded-xl">
            <div class="w-12 h-12 bg-blue-500 rounded-full flex items-center justify-center mr-4">
              <iconify-icon icon="mdi:code-tags" class="text-white text-xl"></iconify-icon>
            </div>
            <div class="flex-1">
              <h4 class="font-bold text-gray-800">编程高手</h4>
              <p class="text-sm text-gray-600">完成10个编程项目</p>
            </div>
            <div class="text-right">
              <div class="text-blue-600 font-bold">+300积分</div>
              <div class="text-xs text-gray-500">2024-01-10</div>
            </div>
          </div>
          
          <div class="flex items-center p-4 bg-green-50 border border-green-200 rounded-xl">
            <div class="w-12 h-12 bg-green-500 rounded-full flex items-center justify-center mr-4">
              <iconify-icon icon="mdi:account-group" class="text-white text-xl"></iconify-icon>
            </div>
            <div class="flex-1">
              <h4 class="font-bold text-gray-800">团队协作者</h4>
              <p class="text-sm text-gray-600">参与5个团队项目</p>
            </div>
            <div class="text-right">
              <div class="text-green-600 font-bold">+200积分</div>
              <div class="text-xs text-gray-500">2024-01-05</div>
            </div>
          </div>
          
          <!-- 未获得成就 -->
          <h3 class="text-lg font-bold text-gray-700 mb-3 mt-6">待解锁</h3>
          
          <div class="flex items-center p-4 bg-gray-50 border border-gray-200 rounded-xl opacity-60">
            <div class="w-12 h-12 bg-gray-400 rounded-full flex items-center justify-center mr-4">
              <iconify-icon icon="mdi:lightning-bolt" class="text-white text-xl"></iconify-icon>
            </div>
            <div class="flex-1">
              <h4 class="font-bold text-gray-800">速度之王</h4>
              <p class="text-sm text-gray-600">单日完成5个任务</p>
              <div class="w-full bg-gray-200 rounded-full h-2 mt-2">
                <div class="bg-blue-500 h-2 rounded-full" style="width: 60%"></div>
              </div>
              <div class="text-xs text-gray-500 mt-1">进度: 3/5</div>
            </div>
            <div class="text-right">
              <div class="text-gray-500 font-bold">+800积分</div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 系统设置弹窗 -->
    <div v-if="showSettings" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click="showSettings = false">
      <div class="bg-white rounded-2xl p-6 w-[500px] max-h-[80vh] overflow-y-auto" @click.stop>
        <div class="flex justify-between items-center mb-6">
          <h2 class="text-2xl font-bold text-gray-800 flex items-center">
            <iconify-icon icon="mdi:cog" class="text-gray-500 mr-2"></iconify-icon>
            系统设置
          </h2>
          <button @click="showSettings = false" class="text-gray-400 hover:text-gray-600">
            <iconify-icon icon="mdi:close" class="text-2xl"></iconify-icon>
          </button>
        </div>
        
        <!-- 设置选项 -->
        <div class="space-y-6">
          <!-- 通知设置 -->
          <div>
            <h3 class="text-lg font-bold text-gray-700 mb-3 flex items-center">
              <iconify-icon icon="mdi:bell" class="text-blue-500 mr-2"></iconify-icon>
              通知设置
            </h3>
            <div class="space-y-3">
              <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                <div>
                  <div class="font-medium text-gray-800">学习提醒</div>
                  <div class="text-sm text-gray-600">每日学习时间提醒</div>
                </div>
                <label class="relative inline-flex items-center cursor-pointer">
                  <input type="checkbox" class="sr-only peer" checked>
                  <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
                </label>
              </div>
              
              <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                <div>
                  <div class="font-medium text-gray-800">任务截止提醒</div>
                  <div class="text-sm text-gray-600">任务即将到期时提醒</div>
                </div>
                <label class="relative inline-flex items-center cursor-pointer">
                  <input type="checkbox" class="sr-only peer" checked>
                  <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
                </label>
              </div>
              
              <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                <div>
                  <div class="font-medium text-gray-800">社交互动</div>
                  <div class="text-sm text-gray-600">好友动态和消息通知</div>
                </div>
                <label class="relative inline-flex items-center cursor-pointer">
                  <input type="checkbox" class="sr-only peer">
                  <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
                </label>
              </div>
            </div>
          </div>
          
          <!-- 学习偏好 -->
          <div>
            <h3 class="text-lg font-bold text-gray-700 mb-3 flex items-center">
              <iconify-icon icon="mdi:account-cog" class="text-green-500 mr-2"></iconify-icon>
              学习偏好
            </h3>
            <div class="space-y-3">
              <div class="p-3 bg-gray-50 rounded-lg">
                <label class="block text-sm font-medium text-gray-700 mb-2">每日学习目标</label>
                <select class="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent">
                  <option>1小时</option>
                  <option selected>2小时</option>
                  <option>3小时</option>
                  <option>4小时</option>
                </select>
              </div>
              
              <div class="p-3 bg-gray-50 rounded-lg">
                <label class="block text-sm font-medium text-gray-700 mb-2">学习模式</label>
                <div class="space-y-2">
                  <label class="flex items-center">
                    <input type="radio" name="studyMode" class="text-blue-600" checked>
                    <span class="ml-2 text-sm text-gray-700">专注模式</span>
                  </label>
                  <label class="flex items-center">
                    <input type="radio" name="studyMode" class="text-blue-600">
                    <span class="ml-2 text-sm text-gray-700">轻松模式</span>
                  </label>
                  <label class="flex items-center">
                    <input type="radio" name="studyMode" class="text-blue-600">
                    <span class="ml-2 text-sm text-gray-700">挑战模式</span>
                  </label>
                </div>
              </div>
            </div>
          </div>
          
          <!-- 隐私设置 -->
          <div>
            <h3 class="text-lg font-bold text-gray-700 mb-3 flex items-center">
              <iconify-icon icon="mdi:shield-account" class="text-purple-500 mr-2"></iconify-icon>
              隐私设置
            </h3>
            <div class="space-y-3">
              <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                <div>
                  <div class="font-medium text-gray-800">学习进度公开</div>
                  <div class="text-sm text-gray-600">允许好友查看学习进度</div>
                </div>
                <label class="relative inline-flex items-center cursor-pointer">
                  <input type="checkbox" class="sr-only peer" checked>
                  <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
                </label>
              </div>
              
              <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                <div>
                  <div class="font-medium text-gray-800">在线状态显示</div>
                  <div class="text-sm text-gray-600">显示在线状态给其他用户</div>
                </div>
                <label class="relative inline-flex items-center cursor-pointer">
                  <input type="checkbox" class="sr-only peer" checked>
                  <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
                </label>
              </div>
            </div>
          </div>
          
          <!-- 操作按钮 -->
          <div class="flex gap-3 pt-4">
            <button @click="showSettings = false" class="flex-1 px-4 py-2 bg-gray-200 text-gray-700 rounded-lg font-medium hover:bg-gray-300">
              取消
            </button>
            <button @click="saveSettings" class="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg font-medium hover:bg-blue-700">
              保存设置
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import * as echarts from 'echarts'

export default {
  name: 'TaskManager',
  data() {
    return {
      showAchievements: false,
      showSettings: false,
      activeTab: 'inProgress',
      activeTimeFilter: 'week', // 修正数据属性名称
      chartInstance: null,
      progressChartInstance: null,
      // 不同时间段的数据
      timeFilterData: {
        week: {
          chartData: [45, 52, 68, 73, 64, 42, 30],
          chartLabels: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'],
          completionRate: 72,
          completedTasks: 86,
          totalTasks: 120
        },
        month: {
          chartData: [65, 72, 58, 83, 76, 69, 74, 81, 67, 79, 85, 78, 72, 88, 91, 69, 75, 82, 77, 84, 73, 86, 79, 81, 75, 83, 78, 80, 76, 89],
          chartLabels: ['1日', '2日', '3日', '4日', '5日', '6日', '7日', '8日', '9日', '10日', '11日', '12日', '13日', '14日', '15日', '16日', '17日', '18日', '19日', '20日', '21日', '22日', '23日', '24日', '25日', '26日', '27日', '28日', '29日', '30日'],
          completionRate: 78,
          completedTasks: 234,
          totalTasks: 300
        },
        quarter: {
          chartData: [68, 74, 82],
          chartLabels: ['1月', '2月', '3月'],
          completionRate: 81,
          completedTasks: 486,
          totalTasks: 600
        }
      },
      tasks: {
        inProgress: [
          {
            id: 1,
            title: 'Vue.js 组件开发',
            description: '学习Vue组件的高级用法',
            progress: 75,
            priority: 'high',
            dueDate: '2024-01-20',
            tags: ['前端', 'Vue']
          },
          {
            id: 2,
            title: 'JavaScript ES6+',
            description: '掌握现代JavaScript语法',
            progress: 60,
            priority: 'medium',
            dueDate: '2024-01-25',
            tags: ['JavaScript', '基础']
          },
          {
            id: 3,
            title: 'CSS Grid 布局',
            description: '学习CSS Grid布局系统',
            progress: 40,
            priority: 'low',
            dueDate: '2024-01-30',
            tags: ['CSS', '布局']
          }
        ],
        toStart: [
          {
            id: 4,
            title: 'React Hooks',
            description: '学习React Hooks的使用',
            progress: 0,
            priority: 'medium',
            dueDate: '2024-02-01',
            tags: ['React', '前端']
          },
          {
            id: 5,
            title: 'Node.js 后端开发',
            description: '构建RESTful API',
            progress: 0,
            priority: 'high',
            dueDate: '2024-02-05',
            tags: ['Node.js', '后端']
          }
        ],
        completed: [
          {
            id: 6,
            title: 'HTML5 基础',
            description: '掌握HTML5新特性',
            progress: 100,
            priority: 'low',
            dueDate: '2024-01-10',
            tags: ['HTML', '基础']
          },
          {
            id: 7,
            title: 'Git 版本控制',
            description: '学习Git基本操作',
            progress: 100,
            priority: 'medium',
            dueDate: '2024-01-15',
            tags: ['Git', '工具']
          }
        ]
      }
    }
  },
  computed: {
    // 当前时间筛选器对应的数据
    currentTimeData() {
      return this.timeFilterData[this.activeTimeFilter];
    }
  },
  mounted() {
    this.initCharts()
  },
  methods: {
    // 设置时间筛选器
    setTimeFilter(filter) {
      this.activeTimeFilter = filter;
      this.updateCharts();
    },
    
    // 更新图表数据
    updateCharts() {
      this.updateRingChart();
      this.updateBarChart();
    },
    
    // 更新环形进度图
    updateRingChart() {
      if (this.chartInstance) {
        const completionRate = this.currentTimeData.completionRate / 100;
        this.chartInstance.setOption({
          tooltip: { show: false },
          series: [{
            type: 'gauge',
            startAngle: 180,
            endAngle: 0,
            radius: '100%',
            min: 0,
            max: 100,
            splitNumber: 10,
            pointer: { show: false },
            axisLine: {
              lineStyle: {
                width: 15,
                color: [[completionRate, '#2D5BFF'], [1, '#F5F7FA']]
              }
            },
            axisLabel: { show: false },
            axisTick: { show: false },
            splitLine: { show: false },
            detail: { show: false }
          }]
        })
        
        // 柱状图
        this.progressChartInstance = echarts.init(this.$refs.taskProgressChart)
        this.updateBarChart()
        
        // 窗口大小变化时重绘图表
        window.addEventListener('resize', () => {
          if (this.chartInstance) this.chartInstance.resize()
          if (this.progressChartInstance) this.progressChartInstance.resize()
        })
      }
    },
    
    // 更新柱状图
    updateBarChart() {
      if (this.progressChartInstance) {
        const isLineChart = this.activeTimeFilter === 'month' || this.activeTimeFilter === 'quarter';
        
        this.progressChartInstance.setOption({
          xAxis: {
            data: this.currentTimeData.chartLabels
          },
          series: [{
            data: this.currentTimeData.chartData,
            type: isLineChart ? 'line' : 'bar',
            barWidth: isLineChart ? undefined : 24,
            smooth: isLineChart ? true : undefined,
            symbol: isLineChart ? 'circle' : undefined,
            symbolSize: isLineChart ? 6 : undefined,
            lineStyle: isLineChart ? {
              width: 3,
              color: '#2D5BFF'
            } : undefined,
            itemStyle: isLineChart ? {
              color: '#2D5BFF',
              borderColor: '#fff',
              borderWidth: 2
            } : {
              color: {
                type: 'linear',
                x: 0,
                y: 0,
                x2: 0,
                y2: 1,
                colorStops: [{
                  offset: 0, color: '#2D5BFF'
                }, {
                  offset: 1, color: '#5D8AFE'
                }]
              },
              borderRadius: [8, 8, 0, 0]
            },
            emphasis: {
              itemStyle: {
                shadowColor: 'rgba(45,91,255,0.5)',
                shadowBlur: 8
              }
            }
          }]
        });
      }
    },
    initCharts() {
      // 环形进度图
      this.chartInstance = echarts.init(this.$refs.ringProgress)
      const completionRate = this.currentTimeData.completionRate / 100;
      this.chartInstance.setOption({
        tooltip: { show: false },
        series: [{
          type: 'gauge',
          startAngle: 180,
          endAngle: 0,
          radius: '100%',
          min: 0,
          max: 100,
          splitNumber: 10,
          pointer: { show: false },
          axisLine: {
            lineStyle: {
              width: 15,
              color: [[completionRate, '#2D5BFF'], [1, '#F5F7FA']]
            }
          },
          axisLabel: { show: false },
          axisTick: { show: false },
          splitLine: { show: false },
          detail: { show: false }
        }]
      })
      
      // 柱状图/折线图
      this.progressChartInstance = echarts.init(this.$refs.taskProgressChart)
      const isLineChart = this.activeTimeFilter === 'month' || this.activeTimeFilter === 'quarter';
      
      this.progressChartInstance.setOption({
        tooltip: {
          trigger: 'axis',
          formatter: '{b}<br/>{c}% 完成'
        },
        grid: {
          left: '3%',
          right: '4%',
          bottom: '3%',
          top: '5%',
          containLabel: true
        },
        xAxis: {
          type: 'category',
          data: this.currentTimeData.chartLabels,
          axisLine: { lineStyle: { color: '#E5E7EB' } },
          axisTick: { show: false }
        },
        yAxis: {
          type: 'value',
          max: 100,
          axisLine: { show: false },
          axisTick: { show: false },
          splitLine: {
            lineStyle: { color: '#F0F2F5' }
          },
          axisLabel: { formatter: '{value}%' }
        },
        series: [{
          data: this.currentTimeData.chartData,
          type: isLineChart ? 'line' : 'bar',
          barWidth: isLineChart ? undefined : 24,
          smooth: isLineChart ? true : undefined,
          symbol: isLineChart ? 'circle' : undefined,
          symbolSize: isLineChart ? 6 : undefined,
          lineStyle: isLineChart ? {
            width: 3,
            color: '#2D5BFF'
          } : undefined,
          itemStyle: isLineChart ? {
            color: '#2D5BFF',
            borderColor: '#fff',
            borderWidth: 2
          } : {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [{
                offset: 0, color: '#2D5BFF'
              }, {
                offset: 1, color: '#5D8AFE'
              }]
            },
            borderRadius: [8, 8, 0, 0]
          },
          emphasis: {
            itemStyle: {
              shadowColor: 'rgba(45,91,255,0.5)',
              shadowBlur: 8
            }
          }
        }]
      })
      
      // 窗口大小变化时重绘图表
      window.addEventListener('resize', () => {
        if (this.chartInstance) this.chartInstance.resize()
        if (this.progressChartInstance) this.progressChartInstance.resize()
      })
    },
    showTaskDetails() {
      // 点击环形图显示任务详情的联动功能
      console.log('显示任务详情')
    },
    saveSettings() {
      // 保存设置逻辑
      console.log('设置已保存');
      this.showSettings = false;
      // 这里可以添加实际的保存逻辑，比如调用API
    }
  }
}
</script>

<style scoped>
.sidebar-shadow {
  box-shadow: 0px 0px 20px rgba(0, 0, 0, 0.05);
}

.badge {
  animation: pulse 1.5s infinite;
}

@keyframes pulse {
  0% { box-shadow: 0 0 0 0 rgba(45, 91, 255, 0.4); }
  70% { box-shadow: 0 0 0 10px rgba(45, 91, 255, 0); }
  100% { box-shadow: 0 0 0 0 rgba(45, 91, 255, 0); }
}

.bubble-left:before {
  content: '';
  width: 0;
  height: 0;
  position: absolute;
  left: -10px;
  top: 10px;
  border-top: 10px solid transparent;
  border-bottom: 10px solid transparent;
  border-right: 10px solid #E5E7EB;
}

.bubble-right:before {
  content: '';
  width: 0;
  height: 0;
  position: absolute;
  right: -10px;
  top: 10px;
  border-top: 10px solid transparent;
  border-bottom: 10px solid transparent;
  border-left: 10px solid #2D5BFF;
}

.friend-card:hover {
  background-color: rgba(45, 91, 255, 0.05);
}
</style>