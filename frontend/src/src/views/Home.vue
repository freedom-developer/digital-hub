<template>
  <div class="page">
    <h1>首页</h1>
    <div v-if="loading">加载中...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else class="user-card">
      <h2>用户信息</h2>
      <p><strong>ID:</strong> {{ user.id }}</p>
      <p><strong>姓名:</strong> {{ user.name }}</p>
      <p><strong>年龄:</strong> {{ user.age }}</p>
    </div>
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
.user-card {
  background: #f8f9fa;
  padding: 20px;
  border-radius: 8px;
  margin-top: 20px;
}

.user-card h2 {
  color: #667eea;
  font-size: 20px;
  margin-bottom: 15px;
}

.user-card p {
  margin: 10px 0;
}

.user-card strong {
  color: #333;
  display: inline-block;
  width: 80px;
}
</style>
