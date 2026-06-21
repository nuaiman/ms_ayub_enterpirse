import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/utils/axios'
import type { ApiResponse } from '@/types/auth'
import type { Customer } from '@/types/customer'
import type { IDType } from '@/types/auth'

import { push } from 'notivue'
import { useGlobalLoader } from 'vue-global-loader'
import type { AxiosError } from 'axios'

export interface CreateCustomerPayload {
  name: string
  phone: string
  email?: string | null
  company?: string | null
  address?: string | null
  id_type?: IDType | null
  id_number?: string | null
  notes?: string | null
}

export interface UpdateCustomerPayload {
  name: string
  phone: string
  email?: string | null
  company?: string | null
  address?: string | null
  id_type?: IDType | null
  id_number?: string | null
  notes?: string | null
}

export const useCustomersStore = defineStore('customers', () => {
  const { displayLoader, destroyLoader } = useGlobalLoader()

  const customers = ref<Customer[]>([])

  // GET ALL CUSTOMERS
  const fetchCustomers = async () => {
    displayLoader()

    try {
      const res = await api.get<ApiResponse<Customer[]>>('/customers')

      if (!res.data.success) {
        push.error(res.data.message)
        return []
      }

      customers.value = res.data.data
      return customers.value
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to fetch customers')
      return []
    } finally {
      destroyLoader()
    }
  }

  // CREATE CUSTOMER
  const createCustomer = async (payload: CreateCustomerPayload) => {
    displayLoader()

    try {
      const res = await api.post<ApiResponse<Customer>>('/customers', payload)

      if (!res.data.success) {
        push.error(res.data.message)
        return null
      }

      customers.value.push(res.data.data)

      push.success(res.data.message)
      return res.data.data
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to create customer')
      return null
    } finally {
      destroyLoader()
    }
  }

  // UPDATE CUSTOMER
  const updateCustomer = async (id: number, payload: UpdateCustomerPayload) => {
    displayLoader()

    try {
      const res = await api.patch<ApiResponse<Customer>>(`/customers/${id}`, payload)

      if (!res.data.success) {
        push.error(res.data.message)
        return false
      }

      // Update the customer in the local list
      const index = customers.value.findIndex((c) => c.id === id)
      if (index !== -1) {
        customers.value[index] = res.data.data
      }

      push.success(res.data.message || 'Customer updated')
      return true
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to update customer')
      return false
    } finally {
      destroyLoader()
    }
  }

  // DELETE CUSTOMER
  const deleteCustomer = async (id: number) => {
    displayLoader()

    try {
      const res = await api.delete<ApiResponse<null>>(`/customers/${id}`)

      if (!res.data.success) {
        push.error(res.data.message)
        return false
      }

      // Remove from local list
      customers.value = customers.value.filter((c) => c.id !== id)

      push.success(res.data.message || 'Customer deleted')
      return true
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to delete customer')
      return false
    } finally {
      destroyLoader()
    }
  }

  return {
    customers,
    fetchCustomers,
    createCustomer,
    updateCustomer,
    deleteCustomer,
  }
})
