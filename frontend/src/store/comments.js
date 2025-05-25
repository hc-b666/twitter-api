import axios from "axios";
import { defineStore } from "pinia";

export const useCommentsStore = defineStore("comments", {
  state: () => ({
    loading: false,
  }),
  getters: {
    isLoading: (state) => state.loading,
  },
  actions: {
    async createComment(postId, body) {
      this.loading = true;

      try {
        const res = await axios.post(`/comments/${postId}`, body);
        return res.data;
      } catch (err) {
        console.error("create comment", err);
        return false;
      } finally {
        this.loading = false;
      }
    },

    async getCommentsByPostId(postId) {
      this.loading = true;

      try {
        const res = await axios.get(`/comments/${postId}`);
        return res.data;
      } catch (err) {
        console.error("get comments by post id", err);
        return false;
      } finally {
        this.loading = false;
      }
    },

    async softDeleteComment(commentId) {
      this.loading = true;

      try {
        const res = await axios.post(`/comments/delete/${commentId}`);
        return res.data;
      } catch (err) {
        console.error("delete comment", err);
        return false;
      } finally {
        this.loading = false;
      }
    },

    async updateComment(commentId, body) {
      this.loading = true;

      try {
        const res = await axios.put(`/comments/${commentId}`, body);
        return res.data;
      } catch (err) {
        console.error("update comment", err);
        return false;
      } finally {
        this.loading = false;
      }
    },
  },
});
