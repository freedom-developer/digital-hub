import axios from 'axios'
import type { User } from '@/types'

// API 基础配置
const apiClient = axios.create({
  baseURL: '/api',
  timeout: 10000
})

// 用户相关 API
export const userApi = {
  // 获取用户信息
  getUser: async (): Promise<User> => {
    const response = await apiClient.get<User>('/user')
    return response.data
  }
}
