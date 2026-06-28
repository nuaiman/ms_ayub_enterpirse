// src/stores/shipments.ts

import { defineStore } from "pinia";
import { ref } from "vue";
import api from "@/utils/axios";
import type { ApiResponse } from "@/types/auth";
import type {
  Shipment,
  CreateShipmentPayload,
  UpdateShipmentPayload,
  ShipmentStatus,
} from "@/types/shipment";
import { push } from "notivue";
import { useGlobalLoader } from "vue-global-loader";
import type { AxiosError } from "axios";

export const useShipmentsStore = defineStore("shipments", () => {
  const { displayLoader, destroyLoader } = useGlobalLoader();
  const shipments = ref<Shipment[]>([]);
  const currentShipment = ref<Shipment | null>(null);

  const fetchShipments = async () => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Shipment[]>>("/shipments");
      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }
      shipments.value = res.data.data;
      return shipments.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch shipments");
      return [];
    } finally {
      destroyLoader();
    }
  };

  const fetchShipment = async (id: number) => {
    displayLoader();
    try {
      const res = await api.get<ApiResponse<Shipment>>(`/shipments/${id}`);
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }
      currentShipment.value = res.data.data;
      return currentShipment.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch shipment");
      return null;
    } finally {
      destroyLoader();
    }
  };

  const createShipment = async (payload: CreateShipmentPayload) => {
    displayLoader();
    try {
      const res = await api.post<ApiResponse<Shipment>>("/shipments", payload);
      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }
      shipments.value.unshift(res.data.data);
      push.success(res.data.message || "Shipment created successfully");
      return res.data.data;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to create shipment");
      return null;
    } finally {
      destroyLoader();
    }
  };

  const updateShipment = async (id: number, payload: UpdateShipmentPayload) => {
    displayLoader();
    try {
      const res = await api.patch<ApiResponse<Shipment>>(`/shipments/${id}`, payload);
      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }
      const index = shipments.value.findIndex((s) => s.id === id);
      if (index !== -1) shipments.value[index] = res.data.data;
      if (currentShipment.value?.id === id) currentShipment.value = res.data.data;
      push.success(res.data.message || "Shipment updated successfully");
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to update shipment");
      return false;
    } finally {
      destroyLoader();
    }
  };

  const changeShipmentStatus = async (id: number, status: ShipmentStatus) => {
    displayLoader();
    try {
      const res = await api.patch<ApiResponse<Shipment>>(`/shipments/${id}/status`, { status });
      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }
      const index = shipments.value.findIndex((s) => s.id === id);
      if (index !== -1) shipments.value[index] = res.data.data;
      if (currentShipment.value?.id === id) currentShipment.value = res.data.data;
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

  const deleteShipment = async (id: number) => {
    displayLoader();
    try {
      const res = await api.delete<ApiResponse<null>>(`/shipments/${id}`);
      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }
      shipments.value = shipments.value.filter((s) => s.id !== id);
      if (currentShipment.value?.id === id) currentShipment.value = null;
      push.success(res.data.message || "Shipment deleted successfully");
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to delete shipment");
      return false;
    } finally {
      destroyLoader();
    }
  };

  return {
    shipments,
    currentShipment,
    fetchShipments,
    fetchShipment,
    createShipment,
    updateShipment,
    changeShipmentStatus,
    deleteShipment,
  };
});
