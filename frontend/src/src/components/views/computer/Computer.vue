<template>
  <div class="page">
    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <el-icon><Monitor /></el-icon>
          <span>计算机</span>
        </div>
      </template>
      <p>欢迎来到计算机页面!</p>
    </el-card>
  </div>
</template>
  
<script lang="ts">
import { defineComponent } from 'vue'
import { Monitor } from '@element-plus/icons-vue'

export default defineComponent({
  name: 'Computer',
  components: {
    Monitor
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