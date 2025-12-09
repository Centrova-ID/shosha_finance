import { createContext, useContext, useState, useEffect, ReactNode } from 'react'
import { authAPI, User } from '@/api/auth'
import { setAuthToken, removeAuthToken, getAuthToken } from '@/api/client'

export type UserRole = 'admin' | 'manager' | 'staff'

export type { User }

interface AuthContextType {
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (identifier: string, password: string) => Promise<boolean>
  logout: () => void
}

const AuthContext = createContext<AuthContextType | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    const initAuth = async () => {
      const token = getAuthToken()
      if (token) {
        try {
          const response = await authAPI.me()
          if (response.success) {
            setUser(response.data)
            localStorage.setItem('shosha_user', JSON.stringify(response.data))
          } else {
            removeAuthToken()
            localStorage.removeItem('shosha_user')
          }
        } catch {
          removeAuthToken()
          localStorage.removeItem('shosha_user')
        }
      }
      setIsLoading(false)
    }

    initAuth()
  }, [])

  const login = async (identifier: string, password: string): Promise<boolean> => {
    try {
      const response = await authAPI.login({ identifier, password })
      if (response.success) {
        setUser(response.data.user)
        setAuthToken(response.data.token)
        localStorage.setItem('shosha_user', JSON.stringify(response.data.user))
        return true
      }
      return false
    } catch {
      return false
    }
  }

  const logout = async () => {
    try {
      await authAPI.logout()
    } catch {
      // Ignore logout errors
    }
    setUser(null)
    removeAuthToken()
    localStorage.removeItem('shosha_user')
  }

  return (
    <AuthContext.Provider value={{ user, isAuthenticated: !!user, isLoading, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
