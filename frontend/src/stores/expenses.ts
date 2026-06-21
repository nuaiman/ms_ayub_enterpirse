import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/utils/axios'
import type { ApiResponse } from '@/types/auth'
import type { Expense, ExpenseCategory } from '@/types/expense'
import { push } from 'notivue'
import { useGlobalLoader } from 'vue-global-loader'
import type { AxiosError } from 'axios'

export interface CreateExpensePayload {
  title: string
  category: ExpenseCategory
  amount: number
  expense_date: string
  notes?: string | null
  image_url?: string | null
}

export interface UpdateExpensePayload {
  title: string
  category: ExpenseCategory
  amount: number
  expense_date: string
  notes?: string | null
}

export const useExpensesStore = defineStore('expenses', () => {
  const { displayLoader, destroyLoader } = useGlobalLoader()
  const expenses = ref<Expense[]>([])

  const fetchExpenses = async () => {
    displayLoader()
    try {
      const res = await api.get<ApiResponse<Expense[]>>('/expenses')
      if (!res.data.success) {
        push.error(res.data.message)
        return []
      }
      expenses.value = res.data.data
      return expenses.value
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to fetch expenses')
      return []
    } finally {
      destroyLoader()
    }
  }

  const createExpense = async (payload: CreateExpensePayload) => {
    displayLoader()
    try {
      const res = await api.post<ApiResponse<Expense>>('/expenses', payload)
      if (!res.data.success) {
        push.error(res.data.message)
        return null
      }
      expenses.value.push(res.data.data)
      push.success(res.data.message)
      return res.data.data
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to create expense')
      return null
    } finally {
      destroyLoader()
    }
  }

  const updateExpense = async (id: number, payload: UpdateExpensePayload) => {
    displayLoader()
    try {
      const res = await api.patch<ApiResponse<Expense>>(`/expenses/${id}`, payload)
      if (!res.data.success) {
        push.error(res.data.message)
        return false
      }
      const index = expenses.value.findIndex((e) => e.id === id)
      if (index !== -1) expenses.value[index] = res.data.data
      push.success(res.data.message || 'Expense updated')
      return true
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to update expense')
      return false
    } finally {
      destroyLoader()
    }
  }

  const deleteExpense = async (id: number) => {
    displayLoader()
    try {
      const res = await api.delete<ApiResponse<null>>(`/expenses/${id}`)
      if (!res.data.success) {
        push.error(res.data.message)
        return false
      }
      expenses.value = expenses.value.filter((e) => e.id !== id)
      push.success(res.data.message || 'Expense deleted')
      return true
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to delete expense')
      return false
    } finally {
      destroyLoader()
    }
  }

  return { expenses, fetchExpenses, createExpense, updateExpense, deleteExpense }
})
