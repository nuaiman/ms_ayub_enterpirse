/* eslint-disable @typescript-eslint/no-explicit-any */
// src/stores/items.ts

import { defineStore } from "pinia";
import { ref } from "vue";
import api from "@/utils/axios";
import type { ApiResponse } from "@/types/auth";
import type { Item, CreateItemPayload, UpdateItemPayload } from "@/types/item";
import { push } from "notivue";
import { useGlobalLoader } from "vue-global-loader";
import type { AxiosError } from "axios";

export const useItemsStore = defineStore("items", () => {
  const { displayLoader, destroyLoader } = useGlobalLoader();
  const items = ref<Item[]>([]);
  const currentItem = ref<Item | null>(null);

  const fetchItems = async () => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Item[]>>("/items");
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      items.value = res.data.data;
      return items.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch items");
      return [];
    } finally {
      destroyLoader();
    }
  };

  const fetchItem = async (id: number) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Item>>(`/items/${id}`);
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }
      currentItem.value = res.data.data;
      return currentItem.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch item");
      return null;
    } finally {
      destroyLoader();
    }
  };

  const createItem = async (payload: CreateItemPayload) => {
    displayLoader();
    try {
      // Set defaults
      const cleanPayload: any = {
        ...payload,
        quantity_unit: payload.quantity_unit || "pcs",
        quantity: payload.quantity || 1,
        amount: payload.amount || 0,
        deposit: payload.deposit || 0,
      };

      // Remove undefined fields, but keep null values
      Object.keys(cleanPayload).forEach((key) => {
        if (cleanPayload[key] === undefined) {
          delete cleanPayload[key];
        }
      });

      // Validate required fields
      if (!cleanPayload.name || !cleanPayload.customer_phone) {
        push.error("Name and customer phone are required");
        return null;
      }

      // Validate contract fields
      if (cleanPayload.duration_type && !cleanPayload.amount) {
        push.error("Amount is required for storage contracts");
        return null;
      }

      const res = await api.post<ApiResponse<Item>>("/items", cleanPayload);
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }
      items.value.unshift(res.data.data);
      push.success(res.data.message || "Item created successfully");
      return res.data.data;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to create item");
      return null;
    } finally {
      destroyLoader();
    }
  };

  const updateItem = async (id: number, payload: UpdateItemPayload) => {
    displayLoader();
    try {
      // Clean up payload - remove undefined fields, but KEEP null values
      const cleanPayload: any = { ...payload };
      Object.keys(cleanPayload).forEach((key) => {
        if (cleanPayload[key] === undefined) {
          delete cleanPayload[key];
        }
      });

      const res = await api.patch<ApiResponse<Item>>(`/items/${id}`, cleanPayload);
      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }
      const index = items.value.findIndex((i) => i.id === id);
      if (index !== -1) items.value[index] = res.data.data;
      if (currentItem.value?.id === id) currentItem.value = res.data.data;
      push.success(res.data.message || "Item updated successfully");
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to update item");
      return false;
    } finally {
      destroyLoader();
    }
  };

  const updateItemStatus = async (id: number, status: string) => {
    displayLoader();
    try {
      const res = await api.patch<ApiResponse<Item>>(`/items/${id}/status`, { status });
      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }
      const index = items.value.findIndex((i) => i.id === id);
      if (index !== -1) items.value[index] = res.data.data;
      if (currentItem.value?.id === id) currentItem.value = res.data.data;
      push.success(res.data.message || "Status updated successfully");
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to update status");
      return false;
    } finally {
      destroyLoader();
    }
  };

  const deleteItem = async (id: number) => {
    displayLoader();
    try {
      const res = await api.delete<ApiResponse<null>>(`/items/${id}`);
      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }
      items.value = items.value.filter((i) => i.id !== id);
      if (currentItem.value?.id === id) currentItem.value = null;
      push.success(res.data.message || "Item deleted successfully");
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to delete item");
      return false;
    } finally {
      destroyLoader();
    }
  };

  // Get items with active storage contracts
  const getActiveContracts = () => {
    return items.value.filter((item) => item.duration_type && item.status === "active");
  };

  // Get items without contracts
  const getItemsWithoutContracts = () => {
    return items.value.filter((item) => !item.duration_type || !item.status);
  };

  // Get completed contracts
  const getCompletedContracts = () => {
    return items.value.filter((item) => item.status === "complete");
  };

  // Search items by customer name or phone
  const searchItems = (query: string) => {
    const searchTerm = query.toLowerCase();
    return items.value.filter(
      (item) =>
        item.name.toLowerCase().includes(searchTerm) ||
        item.customer_phone.includes(searchTerm) ||
        (item.customer_email && item.customer_email.toLowerCase().includes(searchTerm)),
    );
  };

  // Get items by category
  const getItemsByCategory = (category: string) => {
    return items.value.filter((item) => item.category === category);
  };

  // Get all unique categories
  const getCategories = () => {
    const categories = new Set<string>();
    items.value.forEach((item) => {
      if (item.category) categories.add(item.category);
    });
    return Array.from(categories);
  };

  return {
    items,
    currentItem,
    fetchItems,
    fetchItem,
    createItem,
    updateItem,
    updateItemStatus,
    deleteItem,
    getActiveContracts,
    getItemsWithoutContracts,
    getCompletedContracts,
    searchItems,
    getItemsByCategory,
    getCategories,
  };
});
