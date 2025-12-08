import axios from 'axios'

const API_BASE_URL = 'http://localhost:8080/api/v1'
const TOKEN_KEY = 'shosha_token'

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json'
  }
})

apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem(TOKEN_KEY)
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem(TOKEN_KEY)
      localStorage.removeItem('shosha_user')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export const setAuthToken = (token: string) => {
  localStorage.setItem(TOKEN_KEY, token)
}

export const removeAuthToken = () => {
  localStorage.removeItem(TOKEN_KEY)
}

export const getAuthToken = () => {
  return localStorage.getItem(TOKEN_KEY)
}

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
