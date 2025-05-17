import axios from "axios";
import { useAuthStore } from "@/store/auth";

axios.defaults.baseURL = "http://localhost:9999/api/v1";

const setupAxiosInterceptors = () => {
  axios.interceptors.request.use(
    (cfg) => {
      const authStore = useAuthStore();
      const token = authStore.getAccessToken;

      if (token) {
        cfg.headers["Authorization"] = `Bearer ${token}`;
      }

      return cfg;
    },
    (err) => {
      return Promise.reject(err);
    }
  );

  axios.interceptors.response.use(
    (res) => res,
    async (err) => {
      const orgReq = err.config;
      const authStore = useAuthStore();

      if (err.response?.status === 401 && !orgReq._retry) {
        orgReq._retry = true;

        try {
          const refreshed = await authStore.refresh();
          if (refreshed) {
            orgReq.headers[
              "Authorization"
            ] = `Bearer ${authStore.getAccessToken}`;
            return axios(orgReq);
          }
        } catch (refreshErr) {
          authStore.logout();
          return Promise.reject(refreshErr);
        }
      }

      return Promise.reject(err);
    }
  );
};

export default setupAxiosInterceptors;
