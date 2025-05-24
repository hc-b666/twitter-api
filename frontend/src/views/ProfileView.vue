<template>
  <main class="profile-page">
    <h1>Profile</h1>
    <div v-if="loading">Loading profile...</div>
    <div v-else-if="user">
      <p>{{ user.email }}</p>
      <p>{{ user.role }}</p>
    </div>
    <Button @click="handleLogout" label="Logout" />
  </main>
</template>

<script setup>
import Button from "primevue/button";
import { useRouter } from 'vue-router';
import { computed, onMounted } from 'vue';
import { useAuthStore } from '@/store/auth';
import { useToast } from "primevue";

const router = useRouter();
const authStore = useAuthStore();
const toast = useToast();

const user = computed(() => authStore.getUser);
const loading = computed(() => authStore.loading);

onMounted(async () => {
  if (!user.value) {
    await authStore.fetchProfile();
  }
});

function handleLogout() {
  authStore.logout();
  router.push("/login");
  toast.add({
    severity: "success",
    summary: "Logged out",
    detail: "Successfully logged out!",
    life: 3000,
  });
};
</script>
