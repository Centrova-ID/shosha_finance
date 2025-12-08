import { apiClient, APIResponse } from './client'

export interface User {
  id: string
  username: string
  email?: string
  name: string
  role: 'admin' | 'manager' | 'staff'
  branch_id?: string
  is_active: boolean
}

export interface LoginRequest {
  identifier: string
  password: string
}

export interface LoginResponse {
  user: User
  token: string
}

export const authAPI = {
  login: async (data: LoginRequest): Promise<APIResponse<LoginResponse>> => {
    const response = await apiClient.post('/auth/login', data)
    return response.data
  },

  me: async (): Promise<APIResponse<User>> => {
    const response = await apiClient.get('/auth/me')
    return response.data
  },

  logout: async (): Promise<APIResponse<null>> => {
    const response = await apiClient.post('/auth/logout')
    return response.data
  }
}
