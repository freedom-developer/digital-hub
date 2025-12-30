<template>
    <div class="register-container">
      <el-card class="register-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <el-icon :size="24"><User /></el-icon>
            <span>用户注册</span>
          </div>
        </template>
  
        <el-form
          ref="registerFormRef"
          :model="registerForm"
          :rules="rules"
          label-width="80px"
          size="large"
        >
          <el-form-item label="用户名" prop="username">
            <el-input
              v-model="registerForm.username"
              placeholder="请输入用户名（3-20个字符）"
              clearable
            >
              <template #prefix>
                <el-icon><User /></el-icon>
              </template>
            </el-input>
          </el-form-item>
  
          <el-form-item label="邮箱" prop="email">
            <el-input
              v-model="registerForm.email"
              type="email"
              placeholder="请输入邮箱"
              clearable
            >
              <template #prefix>
                <el-icon><Message /></el-icon>
              </template>
            </el-input>
          </el-form-item>
  
          <el-form-item label="密码" prop="password">
            <el-input
              v-model="registerForm.password"
              type="password"
              placeholder="请输入密码（至少3个字符）"
              show-password
              clearable
            >
              <template #prefix>
                <el-icon><Lock /></el-icon>
              </template>
            </el-input>
          </el-form-item>
  
          <el-form-item label="确认密码" prop="confirmPassword">
            <el-input
              v-model="registerForm.confirmPassword"
              type="password"
              placeholder="请再次输入密码"
              show-password
              clearable
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
              @click="handleRegister"
              style="width: 100%"
            >
              注册
            </el-button>
          </el-form-item>
  
          <el-form-item>
            <el-link type="primary" @click="goToLogin">已有账号？立即登录</el-link>
          </el-form-item>
        </el-form>
      </el-card>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref, reactive } from 'vue'
  import { useRouter } from 'vue-router'
  import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
  import { User, Message, Lock } from '@element-plus/icons-vue'
  import { userApi } from '@/api/user'
  
  const router = useRouter()
  const registerFormRef = ref<FormInstance>()
  const loading = ref(false)
  
  const registerForm = reactive({
    username: '',
    email: '',
    password: '',
    confirmPassword: ''
  })
  
  // 验证确认密码
  const validateConfirmPassword = (rule: any, value: any, callback: any) => {
    if (value === '') {
      callback(new Error('请再次输入密码'))
    } else if (value !== registerForm.password) {
      callback(new Error('两次输入的密码不一致'))
    } else {
      callback()
    }
  }
  
  const rules = reactive<FormRules>({
    username: [
      { required: true, message: '请输入用户名', trigger: 'blur' },
      { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
    ],
    email: [
      { required: false, message: '请输入邮箱', trigger: 'blur' },
      { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
    ],
    password: [
      { required: true, message: '请输入密码', trigger: 'blur' },
      { min: 3, message: '密码长度至少 3 个字符', trigger: 'blur' }
    ],
    confirmPassword: [
      { required: true, validator: validateConfirmPassword, trigger: 'blur' }
    ]
  })
  
  async function handleRegister() {
    if (!registerFormRef.value) return
    
    await registerFormRef.value.validate(async (valid) => {
      if (valid) {
        try {
          loading.value = true
          
          await userApi.register({
            username: registerForm.username,
            email: registerForm.email,
            password: registerForm.password
          })
          
          ElMessage.success('注册成功！请登录')
          
          // 跳转到登录页
          setTimeout(() => {
            router.push('/user/login')
          }, 1000)
        } catch (err: any) {
          ElMessage.error(err.response?.data?.message || '注册失败')
        } finally {
          loading.value = false
        }
      }
    })
  }
  
  function goToLogin() {
    router.push('/user/login')
  }
  </script>
  
  <style scoped>
  .register-container {
    min-height: 100vh;
    display: flex;
    justify-content: center;
    align-items: center;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    padding: 20px;
  }
  
  .register-card {
    width: 100%;
    max-width: 500px;
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
  