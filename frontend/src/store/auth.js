import axios from "axios";
import { defineStore } from "pinia";

export const useAuthStore = defineStore("auth", {
  state: () => ({
    accessToken: localStorage.getItem("accessToken") || null,
    refreshToken: localStorage.getItem("refreshToken") || null,
    user: JSON.parse(localStorage.getItem("user")) || null,
    loading: false,
  }),

  getters: {
    isAuthenticated: (state) => !!state.accessToken,
    isAdmin: (state) => (state.user ? state.user.role === "admin" : false),
    getUser: (state) => state.user,
    getAccessToken: (state) => state.accessToken,
  },

  actions: {
    setTokens(accessToken, refreshToken) {
      this.accessToken = accessToken;
      this.refreshToken = refreshToken;
      localStorage.setItem("accessToken", accessToken);
      localStorage.setItem("refreshToken", refreshToken);

      axios.defaults.headers.common["Authorization"] = `Bearer ${accessToken}`;
    },

    logout() {
      this.accessToken = null;
      (this.refreshToken = null), (this.user = null);
      localStorage.removeItem("accessToken");
      localStorage.removeItem("refreshToken");
      localStorage.removeItem("user");

      delete axios.defaults.headers.common["Authorization"];
    },

    async login(credentials) {
      this.loading = true;

      try {
        const res = await axios.post("/auth/login", credentials);
        this.setTokens(res.data.access_token, res.data.refresh_token);
        await this.fetchProfile();
        return true;
      } catch (err) {
        console.error("login failed", err);
        return false;
      } finally {
        this.loading = false;
      }
    },

    async register(userDto) {
      this.loading = true;

      try {
        const res = await axios.post("/auth/register", userDto);
        // res.data.message add toast
        return true;
      } catch (err) {
        console.error("register failed", err);
        return false;
      } finally {
        this.loading = false;
      }
    },

    async refresh() {
      try {
        const res = await axios.post("/auth/refresh", {
          token: this.refreshToken,
        });
        console.log(res);
        this.accessToken = res.data.access_token;
        localStorage.setItem("accessToken", res.data.access_token);
        return true;
      } catch (err) {
        console.error("refresh failed", err);
        this.logout();
        return false;
      }
    },

    async fetchProfile() {
      try {
        const res = await axios.get("/user/profile");
        this.user = res.data;
        localStorage.setItem("user", JSON.stringify(res.data));
        return res.data;
      } catch (err) {
        console.error("fetch profile failed", err);
        const refreshed = await this.refresh();
        if (refreshed) {
          return this.fetchProfile();
        } else {
          this.logout();
        }
      }

      return null;
    },
  },
});
