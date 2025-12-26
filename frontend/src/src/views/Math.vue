<template>
  <div class="page">
    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <el-icon><Edit /></el-icon>
          <span>数学</span>
        </div>
      </template>
      <p>欢迎来到数学页面！</p>
    </el-card>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { Edit } from '@element-plus/icons-vue'

export default defineComponent({
  name: 'Math',
  components: {
    Edit
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
