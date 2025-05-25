<template>
  <main class="profile-layout">
    <div>
      <h1>Profile</h1>
      <div v-if="loading">Loading profile...</div>
      <div v-else-if="user">
        <p>Email: {{ user.email }}</p>
        <p>Role: {{ user.role }}</p>
      </div>
      <Button @click="handleLogout" label="Logout" />
    </div>

    <div class="profile-layout__posts">
      <div v-for="(post, idx) in posts" :key="post.id" class="profile-layout__posts__post">
        <div class="profile-layout__posts__post-header">
          <span>{{ formatDateAndHour(post.created_at) }}</span>
          <Button type="button" @click="toggle($event, idx)" aria-haspopup="true" :aria-controls="`overlay_menu_${idx}`"
            class="menu-button">
            <icon-ellipsis style="color: #fff" />
          </Button>
          <Menu :ref="el => setMenuRef(el, idx)" :id="`overlay_menu_${idx}`" :model="getMenuItems(post, idx)"
            :popup="true" />
        </div>

        <div class="profile-layout__posts__post-content">
          <p>{{ post.content }}</p>
          <img v-if="isImage(post)" :src="post.file_url" :alt="post.content.slice(0, 20)" />
          <a :href="post.file_url" target="_blank" v-else-if="post.file_url"
            class="profile-layout__posts__post-content-attachment">
            <icon-file />
            <span>View the attachment</span>
          </a>
        </div>
      </div>
    </div>
  </main>
</template>

<script setup>
import Button from "primevue/button";
import Menu from "primevue/menu";
import { useRouter } from 'vue-router';
import { computed, onMounted, ref } from 'vue';
import { useAuthStore } from '@/store/auth';
import { useToast } from "primevue/usetoast";
import { usePostsStore } from "@/store/posts";
import { formatDateAndHour, isImage } from "@/utils/utils";
import IconFile from "@/icons/IconFile.vue";
import IconEllipsis from "@/icons/IconEllipsis.vue";

const router = useRouter();
const authStore = useAuthStore();
const toast = useToast();
const postsStore = usePostsStore();

const posts = ref([]);
const user = computed(() => authStore.getUser);
const loading = computed(() => authStore.loading);

const menuRefs = ref({});

onMounted(async () => {
  if (!user.value) {
    await authStore.fetchProfile();
  }

  if (user.value) {
    posts.value = await postsStore.getPostsByUserId(user.value.id);
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
}

function setMenuRef(el, idx) {
  if (el) {
    menuRefs.value[idx] = el;
  }
}

function getMenuItems(post, idx) {
  return [
    {
      label: "Options",
      items: [
        {
          label: 'Edit',
          icon: 'pi pi-pencil',
          command: () => {
            console.log('edit clicked for post:', post.id);
          }
        },
        {
          label: 'Delete',
          icon: 'pi pi-trash',
          command: () => {
            console.log('delete clicked for post:', post.id);
          },
        }
      ],
    },
  ];
}

function toggle(event, idx) {
  const menu = menuRefs.value[idx];
  if (menu) {
    menu.toggle(event);
  } else {
    console.error(`Menu ref not found for index ${idx}`);
  }
}
</script>

<style lang="scss" scoped>
.profile-layout {
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 2rem;
  padding-bottom: 4rem;

  &__posts {
    height: 100vh;
    overflow-y: auto;

    &__post {
      color: #fff;
      text-decoration: none;
      padding: 0.75rem;
      border: 1px solid #27272a;
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
      margin-bottom: 1rem;

      &-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        position: relative;

        a {
          color: #fff;
          font-size: 14px;
          text-decoration: none;
          width: auto;
          padding: 7px 10.5px;
          border-radius: 8px;
          display: flex;
          align-items: center;
          gap: 0.5rem;

          &:hover {
            background-color: oklch(27.4% 0.006 286.033);
          }
        }

        .menu-button {
          width: 32px;
          height: 32px;
          background: transparent;
          border: none;
          border-radius: 8px;
          outline: none;
          cursor: pointer;

          &:hover {
            background-color: oklch(27.4% 0.006 286.033);
          }

          svg {
            width: 20px;
            height: 20px;
          }
        }

        &>span {
          font-size: 12px;
        }
      }

      &-content {
        img {
          width: 100%;
        }

        &-attachment {
          margin-top: 0.75rem;
          text-decoration: none;
          display: flex;
          align-items: center;
          gap: 0.5rem;

          svg {
            width: 20px;
            height: 20px;
          }

          span {
            font-size: 14px;
          }
        }
      }
    }
  }
}

:deep(.p-menu) {
  z-index: 1000;
}
</style>