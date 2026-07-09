// src/types/checkout.ts

export type CheckoutType = "pickup" | "delivery";
export type CheckoutStatus = "pending" | "in_transit" | "complete";

export interface Checkout {
  id: number;
  user_id: number;
  item_id: number;
  quantity: number;
  checkout_date: string;
  receiver_name?: string | null;
  receiver_phone?: string | null;
  type: CheckoutType;

  // Delivery Details (only for delivery type)
  vehicle_number?: string | null;
  driver_name?: string | null;
  driver_phone?: string | null;
  from_location?: string | null;
  to_location?: string | null;

  // Financial (only for delivery)
  delivery_charge?: number | null;
  delivery_cost?: number | null;
  customer_paid?: number | null;

  // Status
  status: CheckoutStatus;
  notes?: string | null;
  image_url?: string | null;

  created_at: string;
  updated_at: string;
}

export interface CreateCheckoutPayload {
  type: CheckoutType;
  item_id: number;
  quantity: number;
  receiver_name?: string | null;
  receiver_phone?: string | null;

  // Delivery Details
  vehicle_number?: string | null;
  driver_name?: string | null;
  driver_phone?: string | null;
  from_location?: string | null;
  to_location?: string | null;

  // Financial
  delivery_charge?: number | null;
  delivery_cost?: number | null;
  customer_paid?: number | null;

  notes?: string | null;
  image_url?: string | null;
}

export interface UpdateCheckoutPayload {
  quantity?: number;
  receiver_name?: string | null;
  receiver_phone?: string | null;

  // Delivery Details
  vehicle_number?: string | null;
  driver_name?: string | null;
  driver_phone?: string | null;
  from_location?: string | null;
  to_location?: string | null;

  // Financial
  delivery_charge?: number | null;
  delivery_cost?: number | null;
  customer_paid?: number | null;

  notes?: string | null;
}
