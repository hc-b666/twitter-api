import axios from "axios";
import { defineStore } from "pinia";

export const useAdminStore = defineStore("admin", {
  state: () => ({
    loading: false,
  }),
  getters: {
    isLoading: (state) => state.loading,
  },
  actions: {
    async getAllUsers() {
      this.loading = true;

      try {
        const res = await axios.get("/admin/users");
        return res.data;
      } catch (err) {
        console.error("get all users", err);
        return false;
      } finally {
        this.loading = false;
      }
    },
  },
});
