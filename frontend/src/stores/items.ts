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
      // Clean up payload - remove undefined fields, but keep null values
      const cleanPayload: any = { ...payload };
      Object.keys(cleanPayload).forEach((key) => {
        if (cleanPayload[key] === undefined) {
          delete cleanPayload[key];
        }
        // Keep null values - they are valid for clearing fields
      });

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
        // Keep null values - they are valid for clearing fields
      });

      console.log("Sending update payload:", cleanPayload);

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

  const getItemsWithContracts = () => {
    return items.value.filter((item) => item.duration_type && item.status);
  };

  const getItemsWithoutContracts = () => {
    return items.value.filter((item) => !item.duration_type || !item.status);
  };

  const getItemsByCustomer = (customerId: number) => {
    return items.value.filter((item) => item.customer_id === customerId);
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
    getItemsWithContracts,
    getItemsWithoutContracts,
    getItemsByCustomer,
  };
});
