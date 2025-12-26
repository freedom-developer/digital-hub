import axios from 'axios'
import type { Music, ApiResponse } from '@/types'

const apiClient = axios.create({
  baseURL: '/api',
  timeout: 10000
})

export const musicApi = {
  // 获取音乐列表
  getMusicList: async (): Promise<Music[]> => {
    const response = await apiClient.get<ApiResponse<Music[]>>('/music')
    return response.data.data
  },

  // 播放音乐
  playMusic: async (musicId: number): Promise<string> => {
    const response = await apiClient.get<ApiResponse<{ url: string }>>(`/music/play/${musicId}`)
    return response.data.data?.url || response.data.message
  }
}
