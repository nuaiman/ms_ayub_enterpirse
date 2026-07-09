// src/types/expense.ts

export interface Expense {
  id: number;
  user_id: number;

  // Salary Tracking
  is_salary: boolean;
  salary_user_id?: number | null;
  salary_month_year?: string | null;

  // Expense Details
  title: string;
  category?: string | null;
  amount: number;
  expense_date: string;
  notes?: string | null;
  image_url?: string | null;

  created_at: string;
  updated_at: string;
}

export interface CreateExpensePayload {
  // Salary Tracking
  is_salary?: boolean;
  salary_user_id?: number | null;
  salary_month_year?: string | null;

  // Expense Details
  title: string;
  category?: string | null;
  amount: number;
  expense_date: string;
  notes?: string | null;
  image_url?: string | null;
}

export interface UpdateExpensePayload {
  // Salary Tracking
  is_salary?: boolean;
  salary_user_id?: number | null;
  salary_month_year?: string | null;

  // Expense Details
  title?: string;
  category?: string | null;
  amount?: number;
  expense_date?: string;
  notes?: string | null;
}
