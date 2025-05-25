import axios from "axios";
import { defineStore } from "pinia";

export const usePostsStore = defineStore("posts", {
  state: () => ({
    loading: false,
  }),
  getters: {
    isLoading: (state) => state.loading,
  },
  actions: {
    async createPost(body) {
      this.loading = true;

      try {
        const res = await axios.post("/posts", body);
        return res.data;
      } catch (err) {
        console.error("create post", err);
        return false;
      } finally {
        this.loading = false;
      }
    },

    async getPostsByUserId(userId) {
      this.loading = true;

      try {
        const res = await axios.get(`/posts/u/${userId}`);
        return res.data;
      } catch (err) {
        console.error("get posts by user id", err);
        return false;
      } finally {
        this.loading = false;
      }
    },

    async getAllPosts() {
      this.loading = true;

      try {
        const res = await axios.get("/posts");
        return res.data;
      } catch (err) {
        console.error("get all posts", err);
        return false;
      } finally {
        this.loading = false;
      }
    },

    async getPostById(postId) {
      this.loading = true;

      try {
        const res = await axios.get(`/posts/${postId}`);
        return res.data;
      } catch (err) {
        console.error("get post by id", err);
        return false;
      } finally {
        this.loading = false;
      }
    },

    async deletePost(postId) {
      this.loading = true;

      try {
        const res = await axios.post(`/posts/${postId}`);
        return res.data;
      } catch (err) {
        console.error("delete post", err);
        return false;
      } finally {
        this.loading = false;
      }
    },

    async updatePost(postId, body) {
      this.loading = true;

      try {
        const res = await axios.put(`/posts/${postId}`, body);
        return res.data;
      } catch (err) {
        console.error("update post", err);
        return false;
      } finally {
        this.loading = false;
      }
    },
  },
});
