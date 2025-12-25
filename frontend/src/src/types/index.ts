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
  
  // API 响应类型
  export interface ApiResponse<T> {
    data: T
    message?: string
    code?: number
  }
  
  // 错误类型
  export interface ApiError {
    message: string
    code?: number
  }
  