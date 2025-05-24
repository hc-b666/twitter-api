<template>
  <div class="main-layout">

    <!-- Left sidebar - Links -->
    <aside class="main-layout-sidebar">
      <div class="main-layout-sidebar-title">
        <router-link to="/">
          <h3>Twitter API</h3>
        </router-link>
      </div>

      <div class="main-layout-sidebar-links">
        <router-link v-if="isAdmin" to="/admin/users">
          <icon-users />
          <span>Users</span>
        </router-link>
      </div>

      <div class="main-layout-sidebar-user">
        <router-link to="/profile">
          <Avatar :label="user?.email[0].toUpperCase()" class="mr-2" size="normal" style="background-color: #6ee7b7; 
            color: #000" shape="circle" />
          <span>{{ user.email }}</span>
        </router-link>
      </div>
    </aside>

    <!-- Views -->
    <div class="main-layout-content">
      <router-view />
    </div>

    <!-- Right sidebar - Tags -->
    <div class="main-layout-tags">
      <div>
        <h4>Tags</h4>
        <p>Currently, on progress</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useAuthStore } from '@/store/auth';
import { Avatar } from 'primevue';
import { computed } from 'vue';
import IconUsers from '@/icons/IconUsers.vue';

const authStore = useAuthStore();

const user = computed(() => authStore.getUser);
const isAdmin = computed(() => authStore.isAdmin);

</script>

<style lang="scss" scoped>
.main-layout {
  margin: auto;

  max-width: 1240px;
  height: 100vh;

  display: flex;

  overflow: hidden;

  &-sidebar {
    padding: 0.5rem;

    width: 240px;

    border-left: 1px solid #27272a;
    border-right: 1px solid #27272a;

    display: flex;
    flex-direction: column;
    gap: 0.5rem;

    &-title {
      padding: 0.5rem;

      a {
        color: #fff;
        text-decoration: none;
      }
    }

    &-links {
      display: flex;
      flex-direction: column;

      a {
        color: #fff;
        text-decoration: none;

        padding: 7px 10.5px;

        border-radius: 8px;

        display: flex;
        align-items: center;
        gap: 0.5rem;

        &:hover {
          background-color: oklch(27.4% 0.006 286.033);
        }
      }
    }

    &-user {
      margin-top: auto;

      a {
        padding: 7px 10.5px;

        color: #fff;
        text-decoration: none;

        border-radius: 8px;

        display: flex;
        align-items: center;
        gap: 0.5rem;

        &:hover {
          background-color: oklch(27.4% 0.006 286.033);
        }

        & span {
          font-size: 14px;
        }
      }
    }
  }

  &-content {
    padding: 1rem;

    flex: 2;
  }

  &-tags {
    flex: 1;

    border-left: 1px solid #27272a;
    border-right: 1px solid #27272a;

    div {
      padding: 1rem;

      display: flex;
      flex-direction: column;
      gap: 0.5rem;

      p {
        font-size: 14px;
      }
    }
  }
}
</style>
