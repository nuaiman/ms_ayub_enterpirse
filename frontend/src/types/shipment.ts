export type ShipmentType = 'pickup' | 'delivery'
export type ShipmentStatus = 'pending' | 'in_transit' | 'completed' | 'cancelled'

export interface Shipment {
  id: number
  user_id?: number | null
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
  status: ShipmentStatus
  notes?: string | null
  shipped_at?: string | null
  received_at?: string | null
  created_at: string
}
