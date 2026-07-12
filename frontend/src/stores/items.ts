/* eslint-disable @typescript-eslint/no-explicit-any */
import { defineStore } from "pinia";
import { ref, computed } from "vue";
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
      console.log("📤 Creating item with payload:", payload);

      const res = await api.post<ApiResponse<Item>>("/items", payload);

      console.log("📥 Create item response:", res.data);

      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }
      items.value.unshift(res.data.data);
      push.success(res.data.message || "Item created successfully");
      return res.data.data;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      console.error("❌ Create item error:", err);
      push.error(err.response?.data?.message || "Failed to create item");
      return null;
    } finally {
      destroyLoader();
    }
  };

  const updateItem = async (id: number, payload: UpdateItemPayload) => {
    displayLoader();
    try {
      const cleanPayload: any = { ...payload };
      Object.keys(cleanPayload).forEach((key) => {
        if (cleanPayload[key] === undefined) {
          delete cleanPayload[key];
        }
      });

      console.log("📤 Updating item:", id, cleanPayload);

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
      console.error("❌ Update item error:", err);
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

  // Computed helpers
  const getItemsByCategory = (category: string) => {
    return items.value.filter((item) => item.category === category);
  };

  const getCategories = computed(() => {
    const categories = new Set<string>();
    items.value.forEach((item) => {
      if (item.category) categories.add(item.category);
    });
    return Array.from(categories);
  });

  const searchItems = (query: string) => {
    const searchTerm = query.toLowerCase();
    return items.value.filter(
      (item) =>
        item.product_name.toLowerCase().includes(searchTerm) ||
        (item.customer_phone && item.customer_phone.includes(searchTerm)) ||
        (item.customer_email && item.customer_email.toLowerCase().includes(searchTerm)) ||
        (item.category && item.category.toLowerCase().includes(searchTerm)) ||
        (item.subcategory && item.subcategory.toLowerCase().includes(searchTerm)) ||
        (item.lot_number && item.lot_number.toLowerCase().includes(searchTerm)) ||
        (item.storage_name && item.storage_name.toLowerCase().includes(searchTerm)) ||
        (item.account_name && item.account_name.toLowerCase().includes(searchTerm))
    );
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
    getItemsByCategory,
    getCategories,
    searchItems,
  };
});