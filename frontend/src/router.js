import { createRouter, createWebHistory } from "vue-router";
import { useAuthStore } from "@/store/auth";

import HomeView from "@/views/HomeView.vue";
import RegisterView from "@/views/RegisterView.vue";
import LoginView from "@/views/LoginView.vue";
import ProfileView from "@/views/ProfileView.vue";
import MainLayout from "./layouts/MainLayout.vue";
import AuthLayout from "./layouts/AuthLayout.vue";

export const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      component: MainLayout,
      children: [
        {
          path: "",
          name: "home",
          component: HomeView,
          meta: {
            requiresAuth: true,
          },
        },
        {
          path: "profile",
          name: "profile",
          component: ProfileView,
          meta: {
            requiresAuth: true,
          },
        },
        {
          path: "admin/users",
          name: "admin-users",
          component: () => import("@/views/UsersView.vue"),
          meta: {
            requiresAuth: true,
            requiresAdmin: true,
          },
        },
      ],
    },
    {
      path: "/",
      component: AuthLayout,
      children: [
        {
          path: "register",
          name: "register",
          component: RegisterView,
          meta: {
            requiresGuest: true,
          },
        },
        {
          path: "login",
          name: "login",
          component: LoginView,
          meta: {
            requiresGuest: true,
          },
        },
      ],
    },
  ],
});

router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore();
  const isAuthenticated = authStore.isAuthenticated;

  if (to.meta.requiresAdmin && !authStore.isAdmin) {
    next("/");
  } else if (to.meta.requiresAuth && !isAuthenticated) {
    next("/login");
  } else if (to.meta.requiresGuest && isAuthenticated) {
    next("/");
  } else {
    next();
  }
});
