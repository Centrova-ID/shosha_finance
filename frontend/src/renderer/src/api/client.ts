import axios from 'axios'

const API_BASE_URL = 'http://localhost:8080/api/v1'

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json'
  }
})

export interface APIResponse<T> {
  success: boolean
  message: string
  data: T
}

export interface PaginatedResponse<T> extends APIResponse<T> {
  meta: {
    page: number
    limit: number
    total: number
    total_pages: number
  }
}
