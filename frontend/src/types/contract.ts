export type DurationType = 'day' | 'week' | 'month' | 'year'
export type ContractStatus = 'active' | 'completed' | 'cancelled'

export interface Contract {
  id: number
  user_id?: number | null
  notes?: string | null
  customer_id: number
  duration_type: DurationType
  duration?: number | null
  price?: number | null
  security_deposit?: number | null
  estimated_value?: number | null
  status: ContractStatus
  ended_at?: string | null
  created_at: string
  updated_at: string
}
