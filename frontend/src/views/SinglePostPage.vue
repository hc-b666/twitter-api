<template>
  <div class="post">
    <Button @click="$router.push('/')" class="mb-4">
      <icon-chevron-left />
      Go Back
    </Button>

    <div class="post-header">
      <router-link to="/profile">
        <Avatar :label="post?.email[0].toUpperCase()" class="mr-2" size="normal" style="background-color: #6ee7b7; 
            color: #000" shape="circle" />
        <span>{{ post?.email }}</span>
      </router-link>

      <span>{{ formatDateAndHour(post?.created_at) }}</span>
    </div>

    <div class="post-content">
      <p>{{ post?.content }}</p>
      <img v-if="isImage(post)" :src="post?.file_url" :alt="post?.content.slice(0, 20)" />
      <a :href="post?.file_url" target="_blank" v-else-if="post?.file_url" class="post-content-attachment">
        <icon-file />
        <span>View the attachment</span>
      </a>
    </div>

    <form @submit.prevent="handleSubmit" class="post-comment-form">
      <h3>Create Comment</h3>
      <FloatLabel variant="on">
        <Textarea v-model="commentContent" id="content" type="text" autoResize />
        <label for="content">Content</label>
      </FloatLabel>
      <Button type="submit" label="Comment" style="margin-top: 1rem;" />
    </form>

    <div class="post-comments">
      <div v-for="(comment, idx) in comments" :key="comment.id" class="post-comments__comment">
        <div class="post-comment-header">
          <Avatar :label="comment.email[0].toUpperCase()" class="mr-2" size="normal" style="background-color: #6ee7b7; 
            color: #000" shape="circle" />
          <span>{{ comment.email }}</span>
          <span>{{ formatDateAndHour(comment.created_at) }}</span>

          <Button v-if="user?.id === comment.user_id" type="button" @click="toggle($event, idx)" aria-haspopup="true"
            :aria-controls="`overlay_menu_${idx}`" class="menu-button">
            <icon-ellipsis style="color: #fff" />
          </Button>
          <Menu :ref="el => setMenuRef(el, idx)" :id="`overlay_menu_${idx}`" :model="getMenuItems(post, idx)"
            :popup="true" />
        </div>
        <p>{{ comment.content }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { usePostsStore } from '@/store/posts';
import { formatDateAndHour, isImage } from '@/utils/utils';
import { Avatar, FloatLabel, Textarea, useToast, Button, Menu } from 'primevue';
import { computed, onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';
import IconFile from "@/icons/IconFile.vue";
import IconChevronLeft from '@/icons/IconChevronLeft.vue';
import { useCommentsStore } from '@/store/comments';
import IconEllipsis from '@/icons/IconEllipsis.vue';
import { useAuthStore } from '@/store/auth';

const authStore = useAuthStore();
const postsStore = usePostsStore();
const commentsStore = useCommentsStore();
const toast = useToast();
const route = useRoute();
const postId = route.params.postId;

const post = ref(null);
const comments = ref([]);
const commentContent = ref('');
const menuRefs = ref({});

const user = computed(() => authStore.getUser);

onMounted(async () => {
  post.value = await postsStore.getPostById(postId);
  comments.value = await commentsStore.getCommentsByPostId(postId);
});

async function handleSubmit() {
  if (!commentContent.value.trim()) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Comment content cannot be empty.',
      life: 3000
    });
    return;
  }

  try {
    const res = await commentsStore.createComment(postId, { content: commentContent.value });

    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: res.message,
      life: 3000
    });
    commentContent.value = '';

    comments.value = await commentsStore.getCommentsByPostId(postId);
  } catch (error) {
    console.error('Error creating comment:', error);
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to create comment.',
      life: 3000
    });
  }
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
          command: async () => {
            try {
              const res = await commentsStore.softDeleteComment(comments.value[idx].id);
              if (res) {
                toast.add({
                  severity: 'success',
                  summary: 'Success',
                  detail: 'Comment deleted successfully.',
                  life: 3000
                });
                comments.value = await commentsStore.getCommentsByPostId(postId);
              } else {
                toast.add({
                  severity: 'error',
                  summary: 'Error',
                  detail: 'Failed to delete comment.',
                  life: 3000
                });
              }
            } catch (err) {
              console.error('Error deleting comment:', err);
              toast.add({
                severity: 'error',
                summary: 'Error',
                detail: 'Failed to delete comment.',
                life: 3000
              });
            }
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
.post {
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

  &-comment-form {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;

    textarea {
      width: 100%;
    }
  }

  &-comments {
    display: flex;
    flex-direction: column;
  }
}
</style>
