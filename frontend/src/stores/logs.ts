import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/utils/axios'
import type { ApiResponse } from '@/types/auth'
import type { Log } from '@/types/log'
import { push } from 'notivue'
import { useGlobalLoader } from 'vue-global-loader'
import type { AxiosError } from 'axios'

export const useLogsStore = defineStore('logs', () => {
  const { displayLoader, destroyLoader } = useGlobalLoader()
  const logs = ref<Log[]>([])

  const fetchLogs = async () => {
    displayLoader()
    try {
      const res = await api.get<ApiResponse<Log[]>>('/logs')
      if (!res.data.success) {
        push.error(res.data.message)
        return []
      }
      logs.value = res.data.data
      return logs.value
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to fetch logs')
      return []
    } finally {
      destroyLoader()
    }
  }

  return { logs, fetchLogs }
})
