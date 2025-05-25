<template>
  <div class="home-layout">

    <ScrollPanel class="home-layout__posts">
      <form @submit.prevent="handleSubmit" class="home-layout__posts__form">
        <h3>Create Post</h3>
        <FloatLabel variant="on">
          <Textarea v-model="content" id="content" type="text" autoResize />
          <label for="content">Content</label>
        </FloatLabel>
        <div style="display: grid; grid-template-columns: repeat(2, 1fr); gap: 2rem;">
          <FileUpload ref="fileupload" mode="basic" :maxFileSize="100000000" @select="onSelect" chooseLabel="Browse"
            style="width: 100%;" />
          <Button type="submit" label="Post" style="margin-top: auto;" />
        </div>
      </form>

      <div v-if="loadPosts" style="display: flex; flex-direction: column; gap: 1rem;">
        <Skeleton width="100%" height="16rem" />
        <Skeleton width="100%" height="16rem" />
        <Skeleton width="100%" height="16rem" />
      </div>

      <div v-for="post in posts" :key="post.id" @click="router.push(`/post/${post.id}`)" class="home-layout__posts__post">
        <div class="home-layout__posts__post-header">
          <router-link to="/profile">
            <Avatar :label="post.email[0].toUpperCase()" class="mr-2" size="normal" style="background-color: #6ee7b7; 
            color: #000" shape="circle" />
            <span>{{ post.email }}</span>
          </router-link>

          <span>{{ formatDateAndHour(post.created_at) }}</span>
        </div>

        <div class="home-layout__posts__post-content">
          <p>{{ post.content }}</p>
          <img v-if="isImage(post)" :src="post.file_url" :alt="post.content.slice(0, 20)" />
          <a :href="post.file_url" target="_blank" v-else-if="post.file_url"
            class="home-layout__posts__post-content-attachment">
            <icon-file />
            <span>View the attachment</span>
          </a>
        </div>
      </div>

      <ScrollTop target="parent" :threshold="100" :buttonProps="{ raised: true, rounded: true }">
        <icon-chevron-up />
      </ScrollTop>
    </ScrollPanel>

  </div>
</template>

<script setup>
import { usePostsStore } from "@/store/posts";
import { formatDateAndHour, isImage } from "@/utils/utils";
import { FileUpload, FloatLabel, Textarea, Button, useToast, ScrollPanel, ScrollTop } from 'primevue';
import { Avatar } from "primevue";
import { computed, onMounted, ref } from "vue";
import IconFile from "@/icons/IconFile.vue";
import IconChevronUp from "@/icons/IconChevronUp.vue";
import { Skeleton } from "primevue";
import { useRouter } from "vue-router";

const postsStore = usePostsStore();
const toast = useToast();
const router = useRouter();

const posts = ref([]);
const content = ref('');
const fileupload = ref();
const file = ref(null);

const loadPosts = computed(() => postsStore.isLoading);

onMounted(async () => {
  try {
    posts.value = await postsStore.getAllPosts();
  } catch (err) {
    console.error("err getting posts:", err);
  }
});

function onSelect(e) {
  file.value = e.files[0];
}

async function handleSubmit() {
  const formData = new FormData();
  formData.append('content', content.value);
  if (file.value) {
    formData.append('file', file.value);
  }

  try {
    const res = await postsStore.createPost(formData);
    if (res) {
      toast.add({ severity: 'success', summary: 'Post Created', detail: 'Your post has been created successfully.', life: 3000 });
      content.value = '';
      file.value = null;
      fileupload.value.clear();
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to create post.', life: 3000 });
    }

  } catch (err) {
    console.error('Error creating post:', err);
  }
}

</script>

<style lang="scss" scoped>
.home-layout {
  &__posts {
    padding: 1rem;
    height: 100vh;

    flex: 2;

    overflow-y: auto;

    &__form {
      width: 100%;

      margin-bottom: 1rem;

      display: flex;
      flex-direction: column;
      gap: 0.5rem;

      textarea {
        width: 100%;
      }
    }

    &__post {
      color: #fff;
      text-decoration: none;

      padding: 0.75rem;

      border: 1px solid #27272a;

      display: flex;
      flex-direction: column;
      gap: 0.5rem;

      &-header {
        display: flex;
        align-items: center;
        justify-content: space-between;

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
</style>
