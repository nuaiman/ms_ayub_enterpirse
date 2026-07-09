import { defineStore } from "pinia";
import { ref } from "vue";
import api from "@/utils/axios";
import type { ApiResponse, User, IDType } from "@/types/auth";

import { push } from "notivue";
import { useGlobalLoader } from "vue-global-loader";
import type { AxiosError } from "axios";

export interface CreateUserPayload {
  name: string;
  username: string;
  password: string;
  role: string;
  email?: string | null;
  phone?: string | null;
  address?: string | null;
  id_type?: string | null;
  id_number?: string | null;
  monthly_salary?: number;
  image_url?: string | null;
}

export interface ChangeRolePayload {
  role: "admin" | "manager" | "accounts" | "staff";
}

export interface UpdateUserProfilePayload {
  name: string;
  email: string | null;
  phone: string | null;
  address: string | null;
  id_type: IDType | null;
  id_number: string | null;
  monthly_salary?: number | null; // ADDED: monthly_salary field
}

export const useUsersStore = defineStore("users", () => {
  const { displayLoader, destroyLoader } = useGlobalLoader();

  const users = ref<User[]>([]);

  // GET ALL USERS
  const fetchUsers = async () => {
    displayLoader();

    try {
      const res = await api.get<ApiResponse<User[]>>("/users");

      if (!res.data.success) {
        push.error(res.data.message);
        return [];
      }

      users.value = res.data.data;
      return users.value;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to fetch users");
      return [];
    } finally {
      destroyLoader();
    }
  };

  // CREATE USER
  const createUser = async (payload: CreateUserPayload) => {
    displayLoader();

    try {
      const res = await api.post<ApiResponse<User>>("/users", payload);

      if (!res.data.success) {
        push.error(res.data.message);
        return null;
      }

      users.value.push(res.data.data);

      push.success(res.data.message);
      return res.data.data;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to create user");
      return null;
    } finally {
      destroyLoader();
    }
  };

  // RESET ALL PASSWORDS
  const resetAllPasswords = async (new_password: string) => {
    displayLoader();

    try {
      const res = await api.patch<ApiResponse<null>>("/users/reset-all-passwords", {
        new_password,
      });

      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }

      push.success(res.data.message);
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to reset passwords");
      return false;
    } finally {
      destroyLoader();
    }
  };

  // CHANGE USER PASSWORD (ADMIN)
  const changeUserPassword = async (id: number, new_password: string) => {
    displayLoader();

    try {
      const res = await api.patch<ApiResponse<null>>(`/users/${id}/change-password`, {
        new_password,
      });

      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }

      push.success(res.data.message);
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to change password");
      return false;
    } finally {
      destroyLoader();
    }
  };

  // CHANGE ROLE
  const changeRole = async (id: number, role: ChangeRolePayload["role"]) => {
    displayLoader();

    try {
      const res = await api.patch<ApiResponse<null>>(`/users/${id}/change-role`, { role });

      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }

      const user = users.value.find((u) => u.id === id);

      if (user) {
        user.role = role;
      }

      push.success(res.data.message);
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to change role");
      return false;
    } finally {
      destroyLoader();
    }
  };

  // TOGGLE ACTIVE
  const toggleActive = async (id: number) => {
    displayLoader();

    try {
      const res = await api.patch<ApiResponse<null>>(`/users/${id}/toggle-active`);

      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }

      const user = users.value.find((u) => u.id === id);

      if (user) {
        user.is_active = !user.is_active;
      }

      push.success(res.data.message);
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to update user");
      return false;
    } finally {
      destroyLoader();
    }
  };

  // UPDATE USER PROFILE
  const updateUserProfile = async (id: number, payload: UpdateUserProfilePayload) => {
    displayLoader();

    try {
      const res = await api.patch<ApiResponse<User>>(`/users/${id}/profile`, payload);

      if (!res.data.success) {
        push.error(res.data.message);
        return false;
      }

      // Update the user in the local list
      const index = users.value.findIndex((u) => u.id === id);
      if (index !== -1) {
        users.value[index] = { ...users.value[index], ...res.data.data };
      }

      push.success(res.data.message || "Profile updated");
      return true;
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>;
      push.error(err.response?.data?.message || "Failed to update profile");
      return false;
    } finally {
      destroyLoader();
    }
  };

  return {
    users,
    fetchUsers,
    createUser,
    resetAllPasswords,
    changeUserPassword,
    changeRole,
    toggleActive,
    updateUserProfile,
  };
});
