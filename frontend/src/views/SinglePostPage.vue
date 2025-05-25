<template>
  <div class="post">
    <div>
      <Button @click="$router.go(-1)" class="mb-4 self-start">
        <icon-chevron-left />
        Go Back
      </Button>
    </div>

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
      <h4>Create Comment</h4>
      <FloatLabel variant="on">
        <Textarea v-model="commentContent" id="content" type="text" autoResize />
        <label for="content">Content</label>
      </FloatLabel>
      <Button type="submit" label="Comment" style="margin-top: 1rem;" />
    </form>

    <div class="post-comments">
      <div v-for="(comment, idx) in comments" :key="comment.id" class="post-comments__comment">
        <div class="post-comments__comment-header">
          <Avatar :label="comment.email[0].toUpperCase()" class="mr-2" size="normal" style="background-color: #6ee7b7; 
              color: #000" shape="circle" />
          <span>{{ comment.email }}</span>
          <span style="margin-left: auto; font-size: 12px;">{{ formatDateAndHour(comment.created_at) }}</span>

          <Button v-if="user?.id === comment.user_id" type="button" @click="toggle($event, idx)" aria-haspopup="true"
            :aria-controls="`overlay_menu_${idx}`" class="menu-button">
            <icon-ellipsis style="color: #fff" />
          </Button>
          <Menu :ref="el => setMenuRef(el, idx)" :id="`overlay_menu_${idx}`" :model="getMenuItems(comment, idx)"
            :popup="true" />
        </div>
        <p>{{ comment.content }}</p>
      </div>
    </div>

    <Dialog v-model:visible="editDialog" modal header="Edit comment" :style="{ width: '25rem' }">
      <span class="text-surface-500 dark:text-surface-400 block mb-8">Update your comment</span>
      <div class="flex items-center gap-4 mb-4">
        <label for="content" class="font-semibold w-24">Content</label>
        <InputText v-model="editCommentContent" id="content" class="flex-auto" autocomplete="off" />
      </div>

      <div style="margin-top: 8px;" class="flex justify-end gap-2">
        <Button type="button" label="Cancel" severity="secondary" @click="editDialog = false"></Button>
        <Button type="button" label="Save" @click="handleEdit"></Button>
      </div>
    </Dialog>

  </div>
</template>

<script setup>
import { usePostsStore } from '@/store/posts';
import { formatDateAndHour, isImage } from '@/utils/utils';
import { Avatar, FloatLabel, Textarea, useToast, Button, Menu, Dialog, InputText } from 'primevue';
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
const editCommentContent = ref('');
const menuRefs = ref({});
const editDialog = ref(false);
const selectedComment = ref(null);

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

async function handleEdit() {
  try {
    const res = await commentsStore.updateComment(selectedComment.value.id, { content: editCommentContent.value });
    if (res) {
      comments.value = await commentsStore.getCommentsByPostId(postId);
      editCommentContent.value = '';
      selectedComment.value = null;
      toast.add({
        severity: 'success',
        summary: 'Success',
        detail: 'Comment edited successfully.',
        life: 3000
      });
    } else {
      toast.add({
        severity: 'error',
        summary: 'Error',
        detail: 'Failed to edit comment.',
        life: 3000
      });
      return;
    }
  } catch (err) {
    console.error('Error editing comment:', err);
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to edit comment.',
      life: 3000
    });
    return;
  }

  editDialog.value = false;
}

function setMenuRef(el, idx) {
  if (el) {
    menuRefs.value[idx] = el;
  }
}

function getMenuItems(comment, idx) {
  return [
    {
      label: "Options",
      items: [
        {
          label: 'Edit',
          icon: 'pi pi-pencil',
          command: () => {
            editDialog.value = true;
            selectedComment.value = comment;
            editCommentContent.value = comment.content;
          }
        },
        {
          label: 'Delete',
          icon: 'pi pi-trash',
          command: async () => {
            try {
              const res = await commentsStore.softDeleteComment(comment.id);
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

  flex-grow: 1;

  overflow-y: auto;

  height: 100vh;

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
    padding: 0.5rem 0;

    border-top: 1px solid #27272a;
    border-bottom: 1px solid #27272a;

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

    &__comment {
      padding: 0.5rem;

      border: 1px solid #27272a;

      display: flex;
      flex-direction: column;
      gap: 0.5rem;

      &-header {
        display: flex;
        align-items: center;
        gap: 0.5rem;
      }
    }
  }
}
</style>
