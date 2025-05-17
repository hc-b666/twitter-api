<template>
  <main class="auth-page">
    <div class="auth-page-content">
      <h2>Register Form</h2>
      <form @submit.prevent="handleRegister" class="auth-page-content-form">
        <FloatLabel variant="in">
          <InputText v-model="email" id="email" type="text" />
          <label for="email">Email</label>
        </FloatLabel>
        <FloatLabel variant="in">
          <InputText v-model="password" id="password" type="password" />
          <label for="password">Password</label>
        </FloatLabel>
        <p>Already have an account? <router-link to="/login">Login here.</router-link></p>
        <p v-if="error" class="error">{{ error }}</p>
        <Button type="submit" label="Register" />
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

const router = useRouter();
const authStore = useAuthStore();

const email = ref('');
const password = ref('');
const error = ref('');

async function handleRegister() {
  const user = {
    email: email.value,
    password: password.value,
  };

  try {
    const success = await authStore.register(user);
    if (success) router.push('/login');
    else error.value = err.message || 'An error occured';

  } catch (err) {
    error.value = err.message || 'An error occured';
  }
}
</script>
