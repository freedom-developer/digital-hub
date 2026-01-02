import Login from '@/components/views/user/Login.vue'
import request from '@/utils/request'

export interface User {
  id: string
  username: string
  email: string
  role: number
  active: boolean
  created_at: string
}

export interface RegisterReq {
  username: string
  email: string
  password: string
}

export type RegisterRsp = User

export interface LoginReq {
  username: string
  password: string
}

export interface LoginRsp {
  token: string
  user: User
}


export const userApi = {
  // 注册用户
  register(req: RegisterReq): Promise<User> {
    return request.post('/user/register', req)
  },

  // 用户登录
  login(req: LoginReq):  Promise<LoginRsp>  {
    return request.post('/user/login', req)
  },

  getCurrentUser(): Promise<User> {
    return request.get('/user/me')
  },
  

}
