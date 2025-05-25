import { createRouter, createWebHistory } from "vue-router";
import { useAuthStore } from "@/store/auth";

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
          component: () => import("@/views/HomeView.vue"),
          meta: {
            requiresAuth: true,
          },
        },
        {
          path: "profile",
          name: "profile",
          component: () => import("@/views/ProfileView.vue"),
          meta: {
            requiresAuth: true,
          },
        },
        {
          path: "post/:postId",
          name: "post",
          component: () => import("@/views/SinglePostPage.vue"),
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
          component: () => import("@/views/RegisterView.vue"),
          meta: {
            requiresGuest: true,
          },
        },
        {
          path: "login",
          name: "login",
          component: () => import("@/views/LoginView.vue"),
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
