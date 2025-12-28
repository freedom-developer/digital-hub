<template>
  <div class="page">
    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <el-icon><VideoCamera /></el-icon>
          <span>电影</span>
        </div>
      </template>
      <p>欢迎来到电影页面!</p>
    </el-card>
  </div>
</template>

<script lang="ts">
  import { defineComponent } from 'vue';
  import { VideoCamera } from '@element-plus/icons-vue';

  export default defineComponent({
    name: 'Movie',
    components: {
      VideoCamera
    }
  })
</script>

<style scoped>
  .page {
    max-width: 800px;
    margin: 0 auto;
  }

  .card-header {
    display: flex;
    align-items: center;
    gap: 10px;
    font-size: 18px;
    font-weight: bold;
  }
</style>