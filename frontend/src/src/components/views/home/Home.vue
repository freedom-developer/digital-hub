<!-- src/views/Home.vue -->
<template>
  <div class="home">
    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <span>用户信息</span>
          <el-button type="danger" size="small" @click="handleLogout">
            退出登录
          </el-button>
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
          <el-descriptions-item label="用户名">{{ user.username }}</el-descriptions-item>
          <el-descriptions-item label="邮箱">{{ user.email || '未设置' }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { userApi }  from '@/api/user'

const router = useRouter()
const user = ref<any>({})
const loading = ref(true)
const error = ref<string | null>(null)

onMounted(async () => {
  await checkLoginAndFetchUser()
})

// 检查登录状态并获取用户信息
async function checkLoginAndFetchUser() {
  const token = localStorage.getItem('token')
  
  if (!token) {
    // 未登录，跳转到登录页
    ElMessage.warning('请先登录')
    router.push('/user/login')
    return
  }
  
  try {
    loading.value = true
    const response = await userApi.getCurrentUser()
    user.value = response
  } catch (err: any) {
    if (err.response?.status === 401) {
      // token 过期或无效
      ElMessage.error('登录已过期，请重新登录')
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      router.push('/login')
    } else {
      error.value = `获取用户信息失败：${err.message}`
    }
  } finally {
    loading.value = false
  }
}

// 退出登录
async function handleLogout() {
  try {
    await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    // await authApi.logout()
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    ElMessage.success('已退出登录')
    router.push('/user/login')
  } catch (err) {
    // 用户取消操作
  }
}
</script>

<style scoped>
.home {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 18px;
  font-weight: bold;
}
</style>
