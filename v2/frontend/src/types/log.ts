export interface Log {
  id: number
  user_id: number
  action: string
  description: string
  entity_type: string
  entity_id: number
  old_data?: string | null
  new_data?: string | null
  ip_address?: string | null
  user_agent?: string | null
  created_at: string
}
