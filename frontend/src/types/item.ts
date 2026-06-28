// src/types/item.ts

export type QuantityUnit = "pcs" | "g" | "kg" | "ton" | "ml" | "liter";
export type DurationType = "day" | "week" | "month" | "year";
export type ItemStatus = "active" | "completed" | "cancelled";

export interface Item {
  id: number;
  user_id?: number | null;
  notes?: string | null;

  // Customer Relationship (Optional)
  customer_id?: number | null;

  // Contract Fields (Optional - if present, item has a contract)
  duration_type?: DurationType | null;
  duration?: number | null;
  price?: number | null;
  security_deposit?: number | null;
  estimated_value?: number | null;
  status?: ItemStatus | null;
  ended_at?: string | null;

  // Item Fields (Required)
  name: string;
  quantity_unit: QuantityUnit;
  quantity: number;
  weight?: number | null;
  width?: number | null;
  height?: number | null;
  length?: number | null;
  image_url?: string | null;
  created_at: string;
  updated_at: string;
}

export interface CreateItemPayload {
  // Required
  name: string;

  // Optional with defaults
  quantity_unit?: QuantityUnit; // defaults to 'pcs'
  quantity?: number; // defaults to 1

  // Optional nullable fields
  customer_id?: number | null;
  duration_type?: DurationType | null;
  duration?: number | null;
  price?: number | null;
  security_deposit?: number | null;
  estimated_value?: number | null;
  weight?: number | null;
  width?: number | null;
  height?: number | null;
  length?: number | null;
  notes?: string | null;
  image_url?: string | null;
}

export interface UpdateItemPayload {
  // Required fields for update
  name?: string; // Optional but should be sent if updating

  // Optional fields that can be updated
  customer_id?: number | null;
  quantity_unit?: QuantityUnit;
  quantity?: number;
  weight?: number | null;
  width?: number | null;
  height?: number | null;
  length?: number | null;
  notes?: string | null;

  // Contract fields (optional, can be set to null to remove)
  duration_type?: DurationType | null;
  duration?: number | null;
  price?: number | null;
  security_deposit?: number | null;
  estimated_value?: number | null;
  status?: ItemStatus | null;
}

// Helper function to check if item has a contract
export function hasContract(item: Item): boolean {
  return !!item.duration_type && !!item.status;
}

// Helper function to check if item has a customer
export function hasCustomer(item: Item): boolean {
  return !!item.customer_id;
}
