// 用户类型
export interface User {
  id: number
  name: string
  age: number
}

// 菜单项类型
export interface MenuItem {
  id: string
  name: string
}

// 音乐类型
export interface Music {
  id: number
  name: string
  file_path: string
}

// API 响应类型
export interface ApiResponse<T> {
  code: number
  data: T
  message: string
}
