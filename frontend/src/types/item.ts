export type QuantityUnit = "bag" | "bottle" | "box" | "can" | "carton" | "cup" | "dozen" | "gallon" | "pack" | "pair" | "pcs" | "roll" | "set" | "sheet" | "unit";
export type WeightUnit = "mg" | "g" | "oz" | "lb" | "kg" | "ton";

export interface Item {
  id: number; user_id: number; notes?: string | null;
  product_name: string; storage_name?: string | null; account_name?: string | null; lot_number?: string | null;
  customer_phone?: string | null; customer_email?: string | null;
  category?: string | null; subcategory?: string | null;
  quantity_unit: QuantityUnit; quantity: number;
  weight?: number | null; weight_unit?: WeightUnit | null;
  amount: number; deposit: number; customer_paid: number;
  image_url?: string | null; created_at: string; updated_at: string;
}

export interface CreateItemPayload {
  product_name: string; storage_name?: string | null; account_name?: string | null; lot_number?: string | null;
  customer_phone?: string | null; customer_email?: string | null;
  category?: string | null; subcategory?: string | null;
  quantity_unit?: QuantityUnit; quantity?: number;
  weight?: number | null; weight_unit?: WeightUnit | null;
  amount?: number; deposit?: number; customer_paid?: number;
  notes?: string | null; image_url?: string | null;
}

export interface UpdateItemPayload {
  product_name?: string; storage_name?: string | null; account_name?: string | null; lot_number?: string | null;
  customer_phone?: string | null; customer_email?: string | null;
  category?: string | null; subcategory?: string | null;
  quantity_unit?: QuantityUnit; quantity?: number;
  weight?: number | null; weight_unit?: WeightUnit | null;
  amount?: number; deposit?: number; customer_paid?: number;
  notes?: string | null;
}