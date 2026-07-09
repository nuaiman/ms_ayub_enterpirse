// src/types/item.ts

export type QuantityUnit =
  | "bag"
  | "bottle"
  | "box"
  | "can"
  | "carton"
  | "cup"
  | "dozen"
  | "gallon"
  | "pack"
  | "pair"
  | "pcs"
  | "roll"
  | "set"
  | "sheet"
  | "unit";

export type WeightUnit = "mg" | "g" | "oz" | "lb" | "kg" | "ton";
export type DimensionUnit = "mm" | "cm" | "in" | "ft" | "m" | "yd" | "km";
export type DurationType = "day" | "week" | "month" | "year";
export type ItemStatus = "active" | "complete";

export interface Item {
  id: number;
  user_id: number;
  notes?: string | null;

  // Customer Fields
  name: string;
  customer_phone: string;
  customer_email?: string | null;

  // Category
  category?: string | null;
  subcategory?: string | null;

  // Item Details
  quantity_unit: QuantityUnit;
  quantity: number;

  // Dimensions
  weight?: number | null;
  weight_unit?: WeightUnit | null;
  width?: number | null;
  height?: number | null;
  length?: number | null;
  dimension_unit?: DimensionUnit | null;

  // Storage Contract (Optional)
  duration_type?: DurationType | null;
  duration?: number | null;
  start_date?: string | null;
  amount: number;
  deposit: number;
  customer_paid: number;

  // Status
  status?: ItemStatus | null;
  image_url?: string | null;

  created_at: string;
  updated_at: string;
}

export interface CreateItemPayload {
  // Customer Fields (Required)
  name: string;
  customer_phone: string;
  customer_email?: string | null;

  // Category (Optional)
  category?: string | null;
  subcategory?: string | null;

  // Item Details (Optional with defaults)
  quantity_unit?: QuantityUnit; // defaults to 'pcs'
  quantity?: number; // defaults to 1

  // Dimensions (Optional)
  weight?: number | null;
  weight_unit?: WeightUnit | null;
  width?: number | null;
  height?: number | null;
  length?: number | null;
  dimension_unit?: DimensionUnit | null;

  // Storage Contract (Optional)
  duration_type?: DurationType | null;
  duration?: number | null;
  start_date?: string | null;
  amount?: number; // Required if creating contract
  deposit?: number;
  customer_paid?: number;

  // Notes and Image
  notes?: string | null;
  image_url?: string | null;
}

export interface UpdateItemPayload {
  // Customer Fields
  name?: string;
  customer_phone?: string;
  customer_email?: string | null;

  // Category
  category?: string | null;
  subcategory?: string | null;

  // Item Details
  quantity_unit?: QuantityUnit;
  quantity?: number;

  // Dimensions
  weight?: number | null;
  weight_unit?: WeightUnit | null;
  width?: number | null;
  height?: number | null;
  length?: number | null;
  dimension_unit?: DimensionUnit | null;

  // Storage Contract
  duration_type?: DurationType | null;
  duration?: number | null;
  start_date?: string | null;
  amount?: number;
  deposit?: number;
  customer_paid?: number;
  status?: ItemStatus | null;

  // Notes
  notes?: string | null;
}

// Helper function to check if item has a storage contract
export function hasContract(item: Item): boolean {
  return !!item.duration_type && !!item.duration && !!item.status;
}

// Helper function to check if item storage is active
export function isActiveContract(item: Item): boolean {
  return item.status === "active";
}
