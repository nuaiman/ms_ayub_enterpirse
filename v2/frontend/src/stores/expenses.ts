import { defineStore } from "pinia";
import { ref, computed } from "vue";
import api from "@/utils/axios";
import type { ApiResponse } from "@/types/auth";
import type { Expense, CreateExpensePayload, UpdateExpensePayload } from "@/types/expense";
import { push } from "notivue";
import { useGlobalLoader } from "vue-global-loader";
import type { AxiosError } from "axios";

export const useExpensesStore = defineStore("expenses", () => {
  const { displayLoader, destroyLoader } = useGlobalLoader();
  const expenses = ref<Expense[]>([]);
  const currentExpense = ref<Expense | null>(null);

  const fetchExpenses = async () => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Expense[]>>("/expenses");
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      expenses.value = res.data.data;
      return expenses.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch expenses");
      return [];
    } finally {
      destroyLoader();
    }
  };

  const fetchExpense = async (id: number) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Expense>>(`/expenses/${id}`);
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }
      currentExpense.value = res.data.data;
      return currentExpense.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch expense");
      return null;
    } finally {
      destroyLoader();
    }
  };

  const createExpense = async (payload: CreateExpensePayload) => {
    displayLoader();
    try {
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const cleanPayload: any = { ...payload };

      // Set defaults
      if (cleanPayload.is_salary === undefined) {
        cleanPayload.is_salary = false;
      }

      // Remove undefined fields
      Object.keys(cleanPayload).forEach((key) => {
        if (cleanPayload[key] === undefined) {
          delete cleanPayload[key];
        }
      });

      const res = await api.post<ApiResponse<Expense>>("/expenses", cleanPayload);
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }
      expenses.value.unshift(res.data.data);
      push.success(res.data.message || "Expense created successfully");
      return res.data.data;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to create expense");
      return null;
    } finally {
      destroyLoader();
    }
  };

  const updateExpense = async (id: number, payload: UpdateExpensePayload) => {
    displayLoader();
    try {
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const cleanPayload: any = { ...payload };

      // Remove undefined fields, but keep null values
      Object.keys(cleanPayload).forEach((key) => {
        if (cleanPayload[key] === undefined) {
          delete cleanPayload[key];
        }
      });

      const res = await api.patch<ApiResponse<Expense>>(`/expenses/${id}`, cleanPayload);
      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }
      const index = expenses.value.findIndex((e) => e.id === id);
      if (index !== -1) expenses.value[index] = res.data.data;
      if (currentExpense.value?.id === id) currentExpense.value = res.data.data;
      push.success(res.data.message || "Expense updated");
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to update expense");
      return false;
    } finally {
      destroyLoader();
    }
  };

  const deleteExpense = async (id: number) => {
    displayLoader();
    try {
      const res = await api.delete<ApiResponse<null>>(`/expenses/${id}`);
      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }
      expenses.value = expenses.value.filter((e) => e.id !== id);
      if (currentExpense.value?.id === id) currentExpense.value = null;
      push.success(res.data.message || "Expense deleted");
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to delete expense");
      return false;
    } finally {
      destroyLoader();
    }
  };

  const getSalaries = computed(() => expenses.value.filter((e) => e.is_salary));
  const getRegularExpenses = computed(() => expenses.value.filter((e) => !e.is_salary));
  const getTotalAmount = computed(() => expenses.value.reduce((sum, e) => sum + e.amount, 0));
  const getTotalSalaries = computed(() =>
    expenses.value.filter((e) => e.is_salary).reduce((sum, e) => sum + e.amount, 0),
  );

  const getCategories = computed(() => {
    const categories = new Set<string>();
    expenses.value.forEach((e) => {
      if (e.category) categories.add(e.category);
    });
    return Array.from(categories);
  });

  return {
    expenses,
    currentExpense,
    fetchExpenses,
    fetchExpense,
    createExpense,
    updateExpense,
    deleteExpense,
    getSalaries,
    getRegularExpenses,
    getTotalAmount,
    getTotalSalaries,
    getCategories,
  };
});
