import { defineStore } from "pinia";
import { ref, computed } from "vue";
import api from "@/utils/axios";
import type { ApiResponse } from "@/types/auth";
import type {
  Checkout,
  CreateCheckoutPayload,
  UpdateCheckoutPayload,
} from "@/types/checkout";
import { push } from "notivue";
import { useGlobalLoader } from "vue-global-loader";
import type { AxiosError } from "axios";

export const useCheckoutsStore = defineStore("checkouts", () => {
  const { displayLoader, destroyLoader } = useGlobalLoader();
  const checkouts = ref<Checkout[]>([]);
  const currentCheckout = ref<Checkout | null>(null);

  const fetchCheckouts = async () => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Checkout[]>>("/checkouts");
      if (!res.data.success) { push.error(res.data.message); return []; }
      checkouts.value = res.data.data;
      return checkouts.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch checkouts");
      return [];
    } finally { destroyLoader(); }
  };

  const fetchCheckout = async (id: number) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Checkout>>(`/checkouts/${id}`);
      if (!res.data.success) { push.error(res.data.message); return null; }
      currentCheckout.value = res.data.data;
      return currentCheckout.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch checkout");
      return null;
    } finally { destroyLoader(); }
  };

  const createCheckout = async (payload: CreateCheckoutPayload) => {
    displayLoader();
    try {
      if (!payload.type || !payload.item_id || !payload.quantity) { push.error("Type, item, and quantity are required"); return null; }
      if (payload.quantity <= 0) { push.error("Quantity must be greater than 0"); return null; }
      if (payload.type === "delivery" && !payload.to_location) { push.error("Delivery location is required for delivery type"); return null; }
      const cleanPayload: any = { ...payload };
      Object.keys(cleanPayload).forEach((key) => { if (cleanPayload[key] === undefined) delete cleanPayload[key]; });
      const res = await api.post<ApiResponse<Checkout>>("/checkouts", cleanPayload);
      if (!res.data.success) { push.error(res.data.message); return null; }
      checkouts.value.unshift(res.data.data);
      push.success(res.data.message || "Checkout created successfully");
      return res.data.data;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to create checkout");
      return null;
    } finally { destroyLoader(); }
  };

  const updateCheckout = async (id: number, payload: UpdateCheckoutPayload) => {
    displayLoader();
    try {
      const cleanPayload: any = { ...payload };
      Object.keys(cleanPayload).forEach((key) => { if (cleanPayload[key] === undefined) delete cleanPayload[key]; });
      const res = await api.patch<ApiResponse<Checkout>>(`/checkouts/${id}`, cleanPayload);
      if (!res.data.success) { push.error(res.data.message); return false; }
      const index = checkouts.value.findIndex((c) => c.id === id);
      if (index !== -1) checkouts.value[index] = res.data.data;
      if (currentCheckout.value?.id === id) currentCheckout.value = res.data.data;
      push.success(res.data.message || "Checkout updated successfully");
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to update checkout");
      return false;
    } finally { destroyLoader(); }
  };

  const deleteCheckout = async (id: number, restoreQuantity: boolean = true) => {
    displayLoader();
    try {
      const res = await api.delete<ApiResponse<null>>(`/checkouts/${id}`, { params: { restore_quantity: restoreQuantity } });
      if (!res.data.success) { push.error(res.data.message); return false; }
      checkouts.value = checkouts.value.filter((c) => c.id !== id);
      if (currentCheckout.value?.id === id) currentCheckout.value = null;
      push.success(res.data.message || "Checkout deleted successfully");
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to delete checkout");
      return false;
    } finally { destroyLoader(); }
  };

  // Helper functions
  const getPickups = computed(() => checkouts.value.filter((c) => c.type === "pickup"));
  const getDeliveries = computed(() => checkouts.value.filter((c) => c.type === "delivery"));
  const getCheckoutsByItem = (itemId: number) => checkouts.value.filter((c) => c.item_id === itemId);

  const getTotalDeliveryCharges = computed(() => checkouts.value.filter((c) => c.type === "delivery").reduce((sum, c) => sum + (c.delivery_charge || 0), 0));
  const getTotalDeliveryCosts = computed(() => checkouts.value.filter((c) => c.type === "delivery").reduce((sum, c) => sum + (c.delivery_cost || 0), 0));
  const getDeliveryProfit = computed(() => getTotalDeliveryCharges.value - getTotalDeliveryCosts.value);

  return {
    checkouts, currentCheckout,
    fetchCheckouts, fetchCheckout, createCheckout, updateCheckout, deleteCheckout,
    getPickups, getDeliveries, getCheckoutsByItem,
    getTotalDeliveryCharges, getTotalDeliveryCosts, getDeliveryProfit,
  };
});