<template>
  <div class="w-full h-full overflow-auto px-4">
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- 左侧主要内容 -->
      <div class="lg:col-span-2 space-y-6">
        <div class="card">
          <div class="flex items-center justify-between mb-6">
            <h1 class="text-2xl font-bold">团队任务</h1>
            <button class="bg-[#2D5BFF] text-white font-medium py-2 px-4 rounded-lg text-sm hover:bg-opacity-90 transition-colors flex items-center gap-2">
              <iconify-icon icon="mdi:plus" width="16" height="16"></iconify-icon>
              创建任务
            </button>
          </div>

          <!-- 团队任务进度图表 -->
          <div class="mb-6">
            <h2 class="font-bold text-xl mb-4">团队进度概览</h2>
            <div class="chart-container" ref="teamProgressChart"></div>
          </div>

          <!-- 任务列表 -->
          <div class="space-y-4">
            <h3 class="font-bold text-lg">当前任务</h3>
            
            <div class="p-4 border border-gray-200 rounded-lg">
              <div class="flex items-start justify-between">
                <div class="flex-1">
                  <h4 class="font-medium">登录功能开发</h4>
                  <p class="text-sm text-gray-600 mt-1">实现用户登录、注册和密码重置功能</p>
                  <div class="flex items-center gap-4 mt-3">
                    <div class="flex items-center gap-2">
                      <iconify-icon icon="mdi:account" width="16" height="16" class="text-gray-500"></iconify-icon>
                      <span class="text-sm">王同学</span>
                    </div>
                    <div class="flex items-center gap-2">
                      <iconify-icon icon="mdi:calendar" width="16" height="16" class="text-gray-500"></iconify-icon>
                      <span class="text-sm">2024-01-18</span>
                    </div>
                  </div>
                  <div class="mt-3">
                    <div class="flex items-center justify-between mb-1">
                      <span class="text-sm">进度</span>
                      <span class="text-sm font-medium">85%</span>
                    </div>
                    <div class="progress-bar">
                      <div class="progress-fill" style="width: 85%"></div>
                    </div>
                  </div>
                </div>
                <span class="bg-green-100 text-green-800 px-2 py-1 rounded text-xs">进行中</span>
              </div>
            </div>

            <div class="p-4 border border-gray-200 rounded-lg">
              <div class="flex items-start justify-between">
                <div class="flex-1">
                  <h4 class="font-medium">支付模块设计</h4>
                  <p class="text-sm text-gray-600 mt-1">设计支付流程和UI界面</p>
                  <div class="flex items-center gap-4 mt-3">
                    <div class="flex items-center gap-2">
                      <iconify-icon icon="mdi:account" width="16" height="16" class="text-gray-500"></iconify-icon>
                      <span class="text-sm">陈同学</span>
                    </div>
                    <div class="flex items-center gap-2">
                      <iconify-icon icon="mdi:calendar" width="16" height="16" class="text-gray-500"></iconify-icon>
                      <span class="text-sm">2024-01-22</span>
                    </div>
                  </div>
                  <div class="mt-3">
                    <div class="flex items-center justify-between mb-1">
                      <span class="text-sm">进度</span>
                      <span class="text-sm font-medium">30%</span>
                    </div>
                    <div class="progress-bar">
                      <div class="progress-fill" style="width: 30%"></div>
                    </div>
                  </div>
                </div>
                <span class="bg-blue-100 text-blue-800 px-2 py-1 rounded text-xs">设计中</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧栏 -->
      <div class="space-y-5">
        <div class="card">
          <div class="flex justify-between items-center">
            <h3 class="font-bold text-lg">团队动态</h3>
            <span class="text-xs text-blue-600">更多</span>
          </div>
          
          <div class="mt-4 space-y-3">
            <div class="flex items-start">
              <div class="relative">
                <div class="w-10 h-10 rounded-full bg-gray-300"></div>
                <div class="online-indicator"></div>
              </div>
              <div class="ml-3 flex-1">
                <div>
                  <span class="font-medium">王同学</span> 
                  <span class="text-sm text-gray-500">完成了 登录功能开发</span>
                </div>
                <div class="text-xs text-gray-500 mt-1">20分钟前</div>
                <div class="flex items-center mt-1">
                  <iconify-icon icon="mdi:heart-outline" width="16" height="16" class="interaction-btn text-gray-500"></iconify-icon>
                  <span class="text-xs text-gray-500 ml-1 mr-3">8</span>
                  <iconify-icon icon="mdi:comment-outline" width="16" height="16" class="interaction-btn text-gray-500"></iconify-icon>
                  <span class="text-xs text-gray-500 ml-1">2</span>
                </div>
              </div>
            </div>
            
            <div class="flex items-start">
              <div class="relative">
                <div class="w-10 h-10 rounded-full bg-gray-300"></div>
              </div>
              <div class="ml-3 flex-1">
                <div>
                  <span class="font-medium">钱同学</span> 
                  <span class="text-sm text-gray-500">创建了新任务: 支付模块设计</span>
                </div>
                <div class="text-xs text-gray-500 mt-1">1小时前</div>
              </div>
            </div>
          </div>
        </div>

        <div class="card">
          <div class="flex justify-between items-center">
            <h3 class="font-bold text-lg">团队积分排名</h3>
            <span class="text-xs text-blue-600">详情</span>
          </div>
          <div class="mt-4 space-y-3">
            <div class="flex items-center p-2 hover:bg-gray-50 rounded-lg transition-colors">
              <div class="w-6 text-center font-bold text-[#FF6B35]">1</div>
              <div class="w-8 h-8 rounded-full bg-gray-300 ml-2"></div>
              <div class="ml-3 flex-1">
                <div class="font-medium">王同学</div>
                <div class="text-xs text-gray-500">前端开发</div>
              </div>
              <div class="text-[#FF6B35] font-semibold">1,580</div>
            </div>
            <div class="flex items-center p-2 hover:bg-gray-50 rounded-lg transition-colors">
              <div class="w-6 text-center font-bold text-[#FF9500]">2</div>
              <div class="w-8 h-8 rounded-full bg-gray-300 ml-2"></div>
              <div class="ml-3 flex-1">
                <div class="font-medium">李同学</div>
                <div class="text-xs text-gray-500">后端开发</div>
              </div>
              <div class="text-[#FF9500] font-semibold">1,420</div>
            </div>
            <div class="flex items-center p-2 hover:bg-gray-50 rounded-lg transition-colors">
              <div class="w-6 text-center font-bold text-[#FFC107]">3</div>
              <div class="w-8 h-8 rounded-full bg-gray-300 ml-2"></div>
              <div class="ml-3 flex-1">
                <div class="font-medium">陈同学</div>
                <div class="text-xs text-gray-500">UI设计</div>
              </div>
              <div class="text-[#FFC107] font-semibold">1,350</div>
            </div>
          </div>
        </div>

        <div class="card">
          <h3 class="font-bold text-lg mb-4">协作工具</h3>
          <div class="grid grid-cols-2 gap-3">
            <button class="py-3 bg-blue-50 hover:bg-blue-100 rounded-lg flex flex-col items-center justify-center transition-colors">
              <iconify-icon icon="mdi:video" width="24" height="24" style="color: #2D5BFF;"></iconify-icon>
              <span class="mt-2 text-sm text-gray-700">视频会议</span>
            </button>
            <button class="py-3 bg-orange-50 hover:bg-orange-100 rounded-lg flex flex-col items-center justify-center transition-colors">
              <iconify-icon icon="mdi:file-document" width="24" height="24" style="color: #FF9500;"></iconify-icon>
              <span class="mt-2 text-sm text-gray-700">协作文档</span>
            </button>
            <button class="py-3 bg-green-50 hover:bg-green-100 rounded-lg flex flex-col items-center justify-center transition-colors">
              <iconify-icon icon="mdi:calendar" width="24" height="24" style="color: #4CAF50;"></iconify-icon>
              <span class="mt-2 text-sm text-gray-700">团队日历</span>
            </button>
            <button class="py-3 bg-purple-50 hover:bg-purple-100 rounded-lg flex flex-col items-center justify-center transition-colors">
              <iconify-icon icon="mdi:chart-pie" width="24" height="24" style="color: #9C27B0;"></iconify-icon>
              <span class="mt-2 text-sm text-gray-700">数据报告</span>
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
  name: 'TeamTasks',
  mounted() {
    this.initChart()
  },
  methods: {
    initChart() {
      const chartDom = this.$refs.teamProgressChart
      const chart = echarts.init(chartDom)
      
      const option = {
        tooltip: {
          trigger: 'axis',
          formatter: function(params) {
            return `${params[0].name}<br/>进度: ${params[0].value}%<br/>任务数: ${Math.floor(params[0].value / 10)}`
          }
        },
        legend: {
          data: ['项目进度', '任务完成率'],
          bottom: 0,
        },
        grid: {
          left: '3%',
          right: '4%',
          bottom: '15%',
          top: '10%',
          containLabel: true
        },
        xAxis: {
          type: 'category',
          data: ['第一周', '第二周', '第三周', '第四周', '当前周']
        },
        yAxis: {
          type: 'value',
          axisLabel: {
            formatter: '{value}%'
          },
          max: 100
        },
        series: [
          {
            name: '项目进度',
            type: 'bar',
            data: [20, 35, 45, 60, 78],
            itemStyle: {
              color: '#2D5BFF'
            },
            barWidth: '40%'
          },
          {
            name: '任务完成率',
            type: 'line',
            data: [25, 40, 55, 70, 86],
            smooth: true,
            lineStyle: {
              width: 3,
              color: '#4CAF50'
            },
            symbol: 'circle',
            symbolSize: 8
          }
        ]
      }
      
      chart.setOption(option)
      
      window.addEventListener('resize', () => {
        chart.resize()
      })
    }
  }
}
</script>

<style scoped>
.container {
  max-width: 1440px;
  margin: 0 auto;
  padding: 20px;
}
</style>
