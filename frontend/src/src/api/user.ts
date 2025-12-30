import axios from 'axios'
import type { User } from '@/types'

// API 基础配置
const apiClient = axios.create({
  baseURL: '/api',
  timeout: 10000
})

export interface RegisterData {
  username: string
  email: string
  password: string
}

export interface LoginData {
  username: string
  password: string
}

export interface UserInfo {
  id: string
  username: string
  email: string
  role: number
  active: boolean
  created_at: string
}

export const userApi = {
  // 注册用户
  register(data: RegisterData): Promise<UserInfo> {
    return apiClient.post('/users/register', data)
  },

  // 用户登录
  login(data: LoginData): Promise<UserInfo> {
    return apiClient.post('/users/login', data)
  }

}
