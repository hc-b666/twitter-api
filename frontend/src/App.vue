<template>
  <router-link to="/">home</router-link>
  <router-link to="/register">register</router-link>
  <router-link to="/login">login</router-link>
  <router-link to="/profile">profile</router-link>
  <router-view />
</template>

<script setup>
const name = "world";
</script>

<style lang="scss" scoped></style>
