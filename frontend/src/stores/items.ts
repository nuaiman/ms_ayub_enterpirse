import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/utils/axios'
import type { ApiResponse } from '@/types/auth'
import type { Item, QuantityUnit } from '@/types/item'

import { push } from 'notivue'
import { useGlobalLoader } from 'vue-global-loader'
import type { AxiosError } from 'axios'

export interface CreateItemPayload {
  contract_id?: number | null
  name: string
  quantity_unit?: QuantityUnit
  quantity?: number
  weight?: number | null
  width?: number | null
  height?: number | null
  length?: number | null
  notes?: string | null
  image_url?: string | null
}

export interface UpdateItemPayload {
  contract_id?: number | null
  name: string
  quantity_unit?: QuantityUnit
  quantity?: number
  weight?: number | null
  width?: number | null
  height?: number | null
  length?: number | null
  notes?: string | null
}

export const useItemsStore = defineStore('items', () => {
  const { displayLoader, destroyLoader } = useGlobalLoader()
  const items = ref<Item[]>([])

  const fetchItems = async () => {
    displayLoader()
    try {
      const res = await api.get<ApiResponse<Item[]>>('/items')
      if (!res.data.success) {
        push.error(res.data.message)
        return []
      }
      items.value = res.data.data
      return items.value
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to fetch items')
      return []
    } finally {
      destroyLoader()
    }
  }

  const createItem = async (payload: CreateItemPayload) => {
    displayLoader()
    try {
      const res = await api.post<ApiResponse<Item>>('/items', payload)
      if (!res.data.success) {
        push.error(res.data.message)
        return null
      }
      items.value.push(res.data.data)
      push.success(res.data.message)
      return res.data.data
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to create item')
      return null
    } finally {
      destroyLoader()
    }
  }

  const updateItem = async (id: number, payload: UpdateItemPayload) => {
    displayLoader()
    try {
      const res = await api.patch<ApiResponse<Item>>(`/items/${id}`, payload)
      if (!res.data.success) {
        push.error(res.data.message)
        return false
      }
      const index = items.value.findIndex((i) => i.id === id)
      if (index !== -1) items.value[index] = res.data.data
      push.success(res.data.message || 'Item updated')
      return true
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to update item')
      return false
    } finally {
      destroyLoader()
    }
  }

  const deleteItem = async (id: number) => {
    displayLoader()
    try {
      const res = await api.delete<ApiResponse<null>>(`/items/${id}`)
      if (!res.data.success) {
        push.error(res.data.message)
        return false
      }
      items.value = items.value.filter((i) => i.id !== id)
      push.success(res.data.message || 'Item deleted')
      return true
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to delete item')
      return false
    } finally {
      destroyLoader()
    }
  }

  return { items, fetchItems, createItem, updateItem, deleteItem }
})
