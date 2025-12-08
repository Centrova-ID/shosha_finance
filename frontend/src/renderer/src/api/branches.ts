import { apiClient, APIResponse } from './client'

export interface Branch {
  id: string
  code: string
  name: string
  description: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface BranchRequest {
  code: string
  name: string
  description?: string
}

export const branchAPI = {
  getAll: async (): Promise<APIResponse<Branch[]>> => {
    const response = await apiClient.get('/branches')
    return response.data
  },

  getActive: async (): Promise<APIResponse<Branch[]>> => {
    const response = await apiClient.get('/branches/active')
    return response.data
  },

  getById: async (id: string): Promise<APIResponse<Branch>> => {
    const response = await apiClient.get(`/branches/${id}`)
    return response.data
  },

  create: async (data: BranchRequest): Promise<APIResponse<Branch>> => {
    const response = await apiClient.post('/branches', data)
    return response.data
  },

  update: async (id: string, data: BranchRequest): Promise<APIResponse<Branch>> => {
    const response = await apiClient.put(`/branches/${id}`, data)
    return response.data
  },

  delete: async (id: string): Promise<APIResponse<null>> => {
    const response = await apiClient.delete(`/branches/${id}`)
    return response.data
  }
}
