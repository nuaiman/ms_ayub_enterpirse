import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/utils/axios'
import type { ApiResponse } from '@/types/auth'
import type { Shipment, ShipmentType, ShipmentStatus } from '@/types/shipment'
import { push } from 'notivue'
import { useGlobalLoader } from 'vue-global-loader'
import type { AxiosError } from 'axios'

export interface CreateShipmentPayload {
  shipment_type: ShipmentType
  item_id: number
  quantity: number
  vehicle_number?: string | null
  driver_name?: string | null
  driver_phone?: string | null
  from_location?: string | null
  company_name?: string | null
  to_location?: string | null
  receiver_name?: string | null
  receiver_phone?: string | null
  customer_charge?: number | null
  company_cost?: number | null
  company_paid?: number | null
  customer_paid?: number | null
  notes?: string | null
}

export interface UpdateShipmentPayload {
  vehicle_number?: string | null
  driver_name?: string | null
  driver_phone?: string | null
  from_location?: string | null
  company_name?: string | null
  to_location?: string | null
  receiver_name?: string | null
  receiver_phone?: string | null
  customer_charge?: number | null
  company_cost?: number | null
  company_paid?: number | null
  customer_paid?: number | null
  status: ShipmentStatus
  notes?: string | null
}

export const useShipmentsStore = defineStore('shipments', () => {
  const { displayLoader, destroyLoader } = useGlobalLoader()
  const shipments = ref<Shipment[]>([])

  const fetchShipments = async () => {
    displayLoader()
    try {
      const res = await api.get<ApiResponse<Shipment[]>>('/shipments')
      if (!res.data.success) {
        push.error(res.data.message)
        return []
      }
      shipments.value = res.data.data
      return shipments.value
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to fetch shipments')
      return []
    } finally {
      destroyLoader()
    }
  }

  const createShipment = async (payload: CreateShipmentPayload) => {
    displayLoader()
    try {
      const res = await api.post<ApiResponse<Shipment>>('/shipments', payload)
      if (!res.data.success) {
        push.error(res.data.message)
        return null
      }
      shipments.value.push(res.data.data)
      push.success(res.data.message)
      return res.data.data
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to create shipment')
      return null
    } finally {
      destroyLoader()
    }
  }

  const updateShipment = async (id: number, payload: UpdateShipmentPayload) => {
    displayLoader()
    try {
      const res = await api.patch<ApiResponse<Shipment>>(`/shipments/${id}`, payload)
      if (!res.data.success) {
        push.error(res.data.message)
        return false
      }
      const index = shipments.value.findIndex((s) => s.id === id)
      if (index !== -1) shipments.value[index] = res.data.data
      push.success(res.data.message || 'Shipment updated')
      return true
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to update shipment')
      return false
    } finally {
      destroyLoader()
    }
  }

  const changeShipmentStatus = async (id: number, status: ShipmentStatus) => {
    displayLoader()
    try {
      const res = await api.patch<ApiResponse<Shipment>>(`/shipments/${id}/status`, { status })
      if (!res.data.success) {
        push.error(res.data.message)
        return false
      }
      const index = shipments.value.findIndex((s) => s.id === id)
      if (index !== -1) shipments.value[index] = res.data.data
      push.success(res.data.message || `Shipment ${status}`)
      return true
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to update status')
      return false
    } finally {
      destroyLoader()
    }
  }

  const deleteShipment = async (id: number) => {
    displayLoader()
    try {
      const res = await api.delete<ApiResponse<null>>(`/shipments/${id}`)
      if (!res.data.success) {
        push.error(res.data.message)
        return false
      }
      shipments.value = shipments.value.filter((s) => s.id !== id)
      push.success(res.data.message || 'Shipment deleted')
      return true
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to delete shipment')
      return false
    } finally {
      destroyLoader()
    }
  }

  return {
    shipments,
    fetchShipments,
    createShipment,
    updateShipment,
    changeShipmentStatus,
    deleteShipment,
  }
})
