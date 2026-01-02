<template>
    <div class="login-container">
      <el-card class="login-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <el-icon :size="24"><UserFilled /></el-icon>
            <span>用户登录</span>
          </div>
        </template>
  
        <el-form
          ref="loginFormRef"
          :model="loginForm"
          :rules="rules"
          label-width="80px"
          size="large"
        >
          <el-form-item label="用户名" prop="username">
            <el-input
              v-model="loginForm.username"
              placeholder="请输入用户名"
              clearable
            >
              <template #prefix>
                <el-icon><User /></el-icon>
              </template>
            </el-input>
          </el-form-item>
  
          <el-form-item label="密码" prop="password">
            <el-input
              v-model="loginForm.password"
              type="password"
              placeholder="请输入密码"
              show-password
              clearable
              @keyup.enter="handleLogin"
            >
              <template #prefix>
                <el-icon><Lock /></el-icon>
              </template>
            </el-input>
          </el-form-item>
  
          <el-form-item>
            <el-button
              type="primary"
              :loading="loading"
              @click="handleLogin"
              style="width: 100%"
            >
              登录
            </el-button>
          </el-form-item>
  
          <el-form-item>
            <el-link type="primary" @click="goToRegister">没有账号？立即注册</el-link>
          </el-form-item>
        </el-form>
      </el-card>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref, reactive } from 'vue'
  import { useRouter } from 'vue-router'
  import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
  import { User, Lock, UserFilled } from '@element-plus/icons-vue'
  import { userApi } from '@/api/user'
  
  const router = useRouter()
  const loginFormRef = ref<FormInstance>()
  const loading = ref(false)
  
  const loginForm = reactive({
    username: '',
    password: ''
  })
  
  const rules = reactive<FormRules>({
    username: [
      { required: true, message: '请输入用户名', trigger: 'blur' }
    ],
    password: [
      { required: true, message: '请输入密码', trigger: 'blur' }
    ]
  })
  
  async function handleLogin() {
    if (!loginFormRef.value) return
    
    await loginFormRef.value.validate(async (valid) => {
      if (valid) {
        try {
          loading.value = true
          
          const loginRsp = await userApi.login({
            username: loginForm.username,
            password: loginForm.password
          })

          localStorage.setItem('token', loginRsp.token)
          localStorage.setItem('user', JSON.stringify(loginRsp.user))
          
          ElMessage.success('登录成功！')
          
          // 跳转到首页
          setTimeout(() => {
            router.push('/music')
          }, 1000)
        } catch (err: any) {
          ElMessage.error(err.response?.data?.message || '登录失败')
        } finally {
          loading.value = false
        }
      }
    })
  }
  
  function goToRegister() {
    router.push('/register')
  }
  </script>
  
  <style scoped>
  .login-container {
    min-height: 100vh;
    display: flex;
    justify-content: center;
    align-items: center;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    padding: 20px;
  }
  
  .login-card {
    width: 100%;
    max-width: 450px;
  }
  
  .card-header {
    display: flex;
    align-items: center;
    gap: 12px;
    font-size: 20px;
    font-weight: bold;
  }
  
  :deep(.el-form-item__label) {
    font-weight: 500;
  }
  </style>
  