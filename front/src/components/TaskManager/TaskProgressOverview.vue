<template>
  <div class="grid grid-cols-3 gap-4 mb-6">
    <div
      class="card col-span-1 p-5 cursor-pointer hover:shadow-lg transition-shadow"
      @click="$emit('show-details')"
    >
      <div class="flex justify-between mb-3">
        <div>
          <p class="text-gray-500 mb-2">今日整体完成率</p>
          <h3 class="text-3xl font-bold text-blue-600">
            {{ dailyOverview.completionRate }}%
          </h3>
        </div>
        <div class="h-16 w-16" ref="ringProgress"></div>
      </div>
      <div class="mt-4">
        <div class="flex justify-between text-sm text-gray-500 mb-1">
          <span>已完成</span>
          <span>
            {{ dailyOverview.completedTasks }}/{{ dailyOverview.totalTasks }}
          </span>
        </div>
        <div class="w-full h-3 bg-gray-200 rounded-full">
          <div
            class="h-full rounded-full bg-green-500"
            :style="`width: ${dailyOverview.completionRate}%`"
          ></div>
        </div>
      </div>
    </div>

    <div class="card col-span-2 p-5">
      <div class="h-64" ref="taskProgressChart"></div>
    </div>
  </div>
</template>

<script>
  import * as echarts from "echarts";

  export default {
    name: "TaskProgressOverview",
    props: {
      dailyOverview: {
        type: Object,
        required: true,
      },
      currentTimeData: {
        type: Object,
        required: true,
      },
      activeTimeFilter: {
        type: String,
        required: true,
      },
    },
    emits: ["show-details"],
    data() {
      return {
        ringChart: null,
        progressChart: null,
      };
    },
    mounted() {
      this.initCharts();
      window.addEventListener("resize", this.handleResize);
    },
    beforeUnmount() {
      window.removeEventListener("resize", this.handleResize);
      if (this.ringChart) {
        this.ringChart.dispose();
      }
      if (this.progressChart) {
        this.progressChart.dispose();
      }
    },
    watch: {
      dailyOverview: {
        deep: true,
        handler() {
          this.renderRingChart();
        },
      },
      currentTimeData: {
        deep: true,
        handler() {
          this.renderProgressChart();
        },
      },
      activeTimeFilter() {
        this.renderProgressChart();
      },
    },
    methods: {
      initCharts() {
        this.$nextTick(() => {
          if (this.$refs.ringProgress) {
            this.ringChart = echarts.init(this.$refs.ringProgress);
          }
          if (this.$refs.taskProgressChart) {
            this.progressChart = echarts.init(this.$refs.taskProgressChart);
          }
          this.renderCharts();
        });
      },
      renderCharts() {
        this.renderRingChart();
        this.renderProgressChart();
      },
      renderRingChart() {
        if (!this.ringChart || !this.dailyOverview) return;
        const completionRate = (this.dailyOverview.completionRate ?? 0) / 100;
        this.ringChart.setOption({
          tooltip: { show: false },
          series: [
            {
              type: "gauge",
              startAngle: 180,
              endAngle: 0,
              radius: "100%",
              min: 0,
              max: 100,
              splitNumber: 10,
              pointer: { show: false },
              axisLine: {
                lineStyle: {
                  width: 15,
                  color: [
                    [completionRate, "#2D5BFF"],
                    [1, "#F5F7FA"],
                  ],
                },
              },
              axisLabel: { show: false },
              axisTick: { show: false },
              splitLine: { show: false },
              detail: { show: false },
            },
          ],
        });
      },
      renderProgressChart() {
        if (!this.progressChart || !this.currentTimeData) return;
        const isLineChart = this.activeTimeFilter === "month";

        this.progressChart.setOption({
          tooltip: {
            trigger: "axis",
            formatter: "{b}<br/>{c}% 完成",
          },
          grid: {
            left: "3%",
            right: "4%",
            bottom: "3%",
            top: "5%",
            containLabel: true,
          },
          xAxis: {
            type: "category",
            data: this.currentTimeData.chartLabels,
            axisLine: { lineStyle: { color: "#E5E7EB" } },
            axisTick: { show: false },
          },
          yAxis: {
            type: "value",
            max: 100,
            axisLine: { show: false },
            axisTick: { show: false },
            splitLine: {
              lineStyle: { color: "#F0F2F5" },
            },
            axisLabel: { formatter: "{value}%" },
          },
          series: [
            {
              data: this.currentTimeData.chartData,
              type: isLineChart ? "line" : "bar",
              barWidth: isLineChart ? undefined : 24,
              smooth: isLineChart ? true : undefined,
              symbol: isLineChart ? "circle" : undefined,
              symbolSize: isLineChart ? 6 : undefined,
              lineStyle: isLineChart
                ? {
                    width: 3,
                    color: "#2D5BFF",
                  }
                : undefined,
              itemStyle: isLineChart
                ? {
                    color: "#2D5BFF",
                    borderColor: "#fff",
                    borderWidth: 2,
                  }
                : {
                    color: {
                      type: "linear",
                      x: 0,
                      y: 0,
                      x2: 0,
                      y2: 1,
                      colorStops: [
                        { offset: 0, color: "#2D5BFF" },
                        { offset: 1, color: "#5D8AFE" },
                      ],
                    },
                    borderRadius: [8, 8, 0, 0],
                  },
              emphasis: {
                itemStyle: {
                  shadowColor: "rgba(45,91,255,0.5)",
                  shadowBlur: 8,
                },
              },
            },
          ],
        });
      },
      handleResize() {
        if (this.ringChart) {
          this.ringChart.resize();
        }
        if (this.progressChart) {
          this.progressChart.resize();
        }
      },
    },
  };
</script>
