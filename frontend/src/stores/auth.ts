import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { useRouter } from 'vue-router'

import api, { setAccessToken } from '@/utils/axios'
import type { ApiResponse, User, IDType } from '@/types/auth'

import { push } from 'notivue'
import { useGlobalLoader } from 'vue-global-loader'
import type { AxiosError } from 'axios'

export interface UpdateProfilePayload {
  name: string
  email?: string | null
  phone?: string | null
  address?: string | null
  id_type?: IDType | null
  id_number?: string | null
  image_url?: string | null
}

export const useAuthStore = defineStore('auth', () => {
  const router = useRouter()
  const { displayLoader, destroyLoader } = useGlobalLoader()

  // STATE
  const user = ref<User | null>(null)
  const accessToken = ref<string | null>(localStorage.getItem('access_token'))

  const isAuthenticated = computed(() => !!user.value && !!accessToken.value)

  // HELPERS
  const setSession = (token: string, u: User) => {
    accessToken.value = token
    user.value = u
    setAccessToken(token)
    localStorage.setItem('access_token', token)
  }

  const clearSession = () => {
    accessToken.value = null
    user.value = null
    setAccessToken(null)
    localStorage.removeItem('access_token')
  }

  const redirect = (_role: User['role']) => {
    router.push('/home')
  }

  // LOGIN
  const login = async (username: string, password: string) => {
    displayLoader()

    try {
      const res = await api.post<ApiResponse<{ access_token: string }>>('/users/login', {
        username,
        password,
      })

      if (!res.data.success) {
        push.error(res.data.message || 'Login failed')
        return false
      }

      const token = res.data.data.access_token

      const me = await api.get<ApiResponse<User>>('/users/current-user', {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })

      setSession(token, me.data.data)

      push.success(res.data.message)

      redirect(me.data.data.role)

      return true
    } catch (error) {
      const err = error as AxiosError<ApiResponse>

      push.error(err.response?.data?.message || 'Login failed')

      clearSession()

      return false
    } finally {
      destroyLoader()
    }
  }

  // REFRESH SESSION
  const refreshSession = async () => {
    try {
      const res = await api.post<ApiResponse<{ access_token: string }>>('/users/refresh-session')

      if (!res.data.success) {
        clearSession()
        return false
      }

      const token = res.data.data.access_token

      const me = await api.get<ApiResponse<User>>('/users/current-user', {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })

      setSession(token, me.data.data)

      return true
    } catch {
      clearSession()
      return false
    }
  }

  // INIT AUTH (APP START)
  const initAuth = async () => {
    displayLoader()

    try {
      // 1. Try direct current user (works if access token still valid)
      const me = await api.get<ApiResponse<User>>('/users/current-user')

      if (me.data.success) {
        setSession(accessToken.value || '', me.data.data)
        redirect(me.data.data.role)
        return true
      }
    } catch {}

    // 2. If failed → try refresh session (cookie based)
    try {
      const res = await api.post<ApiResponse<{ access_token: string }>>('/users/refresh-session')

      if (!res.data.success) {
        clearSession()
        return false
      }

      const token = res.data.data.access_token
      setAccessToken(token)

      const me = await api.get<ApiResponse<User>>('/users/current-user')

      setSession(token, me.data.data)

      redirect(me.data.data.role)

      return true
    } catch {
      clearSession()
      return false
    } finally {
      destroyLoader()
    }
  }

  // LOGOUT
  const logout = async () => {
    displayLoader()

    try {
      await api.delete('/users/logout')
    } catch {}

    clearSession()

    push.info('Logged out')

    router.push('/auth')

    destroyLoader()
  }

  // CHANGE PASSWORD
  const changePassword = async (
    current_password: string,
    new_password: string,
  ): Promise<boolean> => {
    displayLoader()

    try {
      const res = await api.patch<ApiResponse<null>>('/users/change-password', {
        current_password,
        new_password,
      })

      if (!res.data.success) {
        push.error(res.data.message || 'Password change failed')
        return false
      }

      push.success(res.data.message || 'Password changed successfully')
      return true
    } catch (error) {
      const err = error as AxiosError<{ message: string }>

      push.error(err.response?.data?.message || 'Password change failed')
      return false
    } finally {
      destroyLoader()
    }
  }

  // UPDATE PROFILE
  const updateProfile = async (payload: UpdateProfilePayload): Promise<boolean> => {
    displayLoader()

    try {
      const res = await api.patch<ApiResponse<User>>('/users/profile', payload)

      if (!res.data.success) {
        push.error(res.data.message || 'Profile update failed')
        return false
      }

      // Update local user state with returned data
      user.value = res.data.data

      push.success(res.data.message || 'Profile updated successfully')
      return true
    } catch (error) {
      const err = error as AxiosError<ApiResponse>

      push.error(err.response?.data?.message || 'Profile update failed')
      return false
    } finally {
      destroyLoader()
    }
  }

  return {
    user,
    accessToken,
    isAuthenticated,
    login,
    logout,
    refreshSession,
    initAuth,
    changePassword,
    updateProfile,
    clearSession,
  }
})
