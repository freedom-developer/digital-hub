import request from '@/utils/request'

export interface Music {
  id: number
  name: string
  file_path: string
}

export const musicApi = {
  // 获取音乐列表
  getMusicList(): Promise<Music[]>  {
    return request.get('/music')
  },

  // 播放音乐
  playMusic: async(musicId: number): Promise<string> => {
    const rsp = await request.get<{ url: string }>(`/music/play/${musicId}`)
    return rsp.url || ''
  },

  // 添加收藏
  addFavorite(musicId: number) {
    return request.post('/music/favorite', { music_id: musicId })
  },

  // 取消收藏
  removeFavorite(musicId: number) {
    return request.delete(`/music/favorite/${musicId}`)
  },

  // 获取收藏列表
  getFavoriteMusic(): Promise<Music[]> {
    return request.get('/music/favorite')
  },

  // 获取收藏的音乐ID列表
  getFavoriteMusicIds(): Promise<number[]> {
    return request.get('/music/favorite/ids')
  },

  // 检查是否已收藏
  checkFavorite(musicId: number) {
    return request.get<{ is_favorite: boolean }>(`/music/favorite/check/${musicId}`)
  }
}
