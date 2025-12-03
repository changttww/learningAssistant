<template>
  <div class="flex items-center justify-between mb-5">
    <h2 class="text-2xl font-bold text-gray-700">动态学习看板</h2>
    <div class="inline-flex bg-gray-100 p-1 rounded-xl">
      <button
        class="px-4 py-1 rounded-lg font-medium transition-colors"
        :class="buttonClass('week')"
        @click="updateFilter('week')"
      >
        周
      </button>
      <button
        class="px-4 py-1 rounded-lg font-medium transition-colors"
        :class="buttonClass('month')"
        @click="updateFilter('month')"
      >
        月
      </button>
    </div>
  </div>
</template>

<script>
  export default {
    name: "LearningBoardHeader",
    props: {
      activeFilter: {
        type: String,
        required: true,
      },
    },
    emits: ["update:activeFilter"],
    methods: {
      updateFilter(filter) {
        this.$emit("update:activeFilter", filter);
      },
      buttonClass(filter) {
        return this.activeFilter === filter
          ? "bg-white shadow-sm text-blue-600"
          : "text-gray-600 hover:text-gray-800";
      },
    },
  };
</script>
