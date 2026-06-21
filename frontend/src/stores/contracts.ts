import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/utils/axios'
import type { ApiResponse } from '@/types/auth'
import type { Contract, DurationType, ContractStatus } from '@/types/contract'

import { push } from 'notivue'
import { useGlobalLoader } from 'vue-global-loader'
import type { AxiosError } from 'axios'

export interface CreateContractPayload {
  customer_id: number
  duration_type: DurationType
  duration?: number | null
  price?: number | null
  security_deposit?: number | null
  estimated_value?: number | null
  notes?: string | null
}

export interface UpdateContractPayload {
  customer_id: number
  duration_type: DurationType
  duration?: number | null
  price?: number | null
  security_deposit?: number | null
  estimated_value?: number | null
  status: ContractStatus
  notes?: string | null
}

export const useContractsStore = defineStore('contracts', () => {
  const { displayLoader, destroyLoader } = useGlobalLoader()

  const contracts = ref<Contract[]>([])

  const fetchContracts = async () => {
    displayLoader()
    try {
      const res = await api.get<ApiResponse<Contract[]>>('/contracts')
      if (!res.data.success) {
        push.error(res.data.message)
        return []
      }
      contracts.value = res.data.data
      return contracts.value
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to fetch contracts')
      return []
    } finally {
      destroyLoader()
    }
  }

  const createContract = async (payload: CreateContractPayload) => {
    displayLoader()
    try {
      const res = await api.post<ApiResponse<Contract>>('/contracts', payload)
      if (!res.data.success) {
        push.error(res.data.message)
        return null
      }
      contracts.value.push(res.data.data)
      push.success(res.data.message)
      return res.data.data
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to create contract')
      return null
    } finally {
      destroyLoader()
    }
  }

  const updateContract = async (id: number, payload: UpdateContractPayload) => {
    displayLoader()
    try {
      const res = await api.patch<ApiResponse<Contract>>(`/contracts/${id}`, payload)
      if (!res.data.success) {
        push.error(res.data.message)
        return false
      }
      const index = contracts.value.findIndex((c) => c.id === id)
      if (index !== -1) {
        contracts.value[index] = res.data.data
      }
      push.success(res.data.message || 'Contract updated')
      return true
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to update contract')
      return false
    } finally {
      destroyLoader()
    }
  }

  const changeContractStatus = async (id: number, status: ContractStatus) => {
    displayLoader()
    try {
      const res = await api.patch<ApiResponse<Contract>>(`/contracts/${id}/status`, { status })
      if (!res.data.success) {
        push.error(res.data.message)
        return false
      }
      const index = contracts.value.findIndex((c) => c.id === id)
      if (index !== -1) {
        contracts.value[index] = res.data.data
      }
      push.success(res.data.message || `Contract ${status}`)
      return true
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to update contract status')
      return false
    } finally {
      destroyLoader()
    }
  }

  const deleteContract = async (id: number) => {
    displayLoader()
    try {
      const res = await api.delete<ApiResponse<null>>(`/contracts/${id}`)
      if (!res.data.success) {
        push.error(res.data.message)
        return false
      }
      contracts.value = contracts.value.filter((c) => c.id !== id)
      push.success(res.data.message || 'Contract deleted')
      return true
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to delete contract')
      return false
    } finally {
      destroyLoader()
    }
  }

  return {
    contracts,
    fetchContracts,
    createContract,
    updateContract,
    changeContractStatus,
    deleteContract,
  }
})
