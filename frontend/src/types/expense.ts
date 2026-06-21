export type ExpenseCategory = 'rent' | 'salary' | 'bill' | 'other'

export interface Expense {
  id: number
  user_id?: number | null
  title: string
  category: ExpenseCategory
  amount: number
  expense_date: string
  notes?: string | null
  image_url?: string | null
  created_at: string
}
