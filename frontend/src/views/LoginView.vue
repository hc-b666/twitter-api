<template>
  <main class="auth-page">
    <div class="auth-page-content">
      <h2>Login Form</h2>
      <form @submit.prevent="handleLogin" class="auth-page-content-form">
        <FloatLabel variant="in">
          <InputText v-model="email" id="email" type="text" />
          <label for="email">Email</label>
        </FloatLabel>
        <FloatLabel variant="in">
          <InputText v-model="password" id="password" type="password" />
          <label for="password">Password</label>
        </FloatLabel>
        <p>Don't you have an account? <router-link to="/register">Register here.</router-link></p>
        <p v-if="error" class="error">{{ error }}</p>
        <Button type="submit" label="Login" />
      </form>
    </div>
  </main>
</template>

<script setup>
import { ref } from 'vue';
import FloatLabel from 'primevue/floatlabel';
import InputText from 'primevue/inputtext';
import Button from 'primevue/button';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/store/auth';
import { useToast } from 'primevue';

const router = useRouter();
const authStore = useAuthStore();
const toast = useToast();

const email = ref('');
const password = ref('');
const error = ref('');

async function handleLogin() {
  const user = {
    email: email.value,
    password: password.value,
  };
  try {
    const success = await authStore.login(user);
    if (success) {
      router.push('/profile');
      toast.add({
        severity: 'success',
        summary: 'Logged in',
        detail: 'Successfully logged in!',
        life: 3000,
      });
    }
    else error.value = err.message || 'An error occured';

  } catch (err) {
    error.value = err.message || 'An error occured';
  }
}
</script>
