<template>
  <div class="home">
    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <span>用户信息</span>
        </div>
      </template>
      
      <div v-loading="loading">
        <el-alert
          v-if="error"
          :title="error"
          type="error"
          :closable="false"
        />
        
        <el-descriptions v-else :column="1" border>
          <el-descriptions-item label="ID">{{ user.id }}</el-descriptions-item>
          <el-descriptions-item label="姓名">{{ user.name }}</el-descriptions-item>
          <el-descriptions-item label="年龄">{{ user.age }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </el-card>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import axios, { AxiosError } from 'axios'
import type { User } from '@/types'
import { userApi } from '@/api/user';

export default defineComponent({
  name: 'Home',
  data() {
    return {
      user: {} as User,
      loading: true,
      error: null as string | null
    }
  },
  mounted() {
    this.fetchUser()
  },
  methods: {
    async fetchUser(): Promise<void> {
      try {
        this.user = await userApi.getUser()
        this.loading = false
      } catch (err: any) {
        this.error = `获取用户信息失败：${err.message}`
        this.loading = false
      }
    }
  }
})
</script>

<style scoped>
.home {
  max-width: 800px;
  margin: 0 auto;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 18px;
  font-weight: bold;
}
</style>
