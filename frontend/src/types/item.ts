export type QuantityUnit = 'pcs' | 'g' | 'kg' | 'ton' | 'ml' | 'liter'

export interface Item {
  id: number
  user_id?: number | null
  notes?: string | null
  contract_id?: number | null
  name: string
  quantity_unit: QuantityUnit
  quantity: number
  weight?: number | null
  width?: number | null
  height?: number | null
  length?: number | null
  image_url?: string | null
  created_at: string
  updated_at: string
}
