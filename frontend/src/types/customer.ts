import type { IDType } from './auth'

export interface Customer {
  id: number
  user_id?: number | null
  name: string
  phone: string
  email?: string | null
  company?: string | null
  address?: string | null
  id_type?: IDType | null
  id_number?: string | null
  notes?: string | null
  created_at: string
  updated_at: string
}
