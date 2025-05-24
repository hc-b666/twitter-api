<template>
  <div>
    <DataTable :value="users" size="small" showGridlines>
      <Column field="id" header="ID"></Column>
      <Column field="email" header="Email"></Column>
      <Column field="role" header="Role"></Column>
      <Column field="created_at" header="Created At">
        <template #body="slotProps">
          {{ formatDate(slotProps.data.created_at) }}
        </template>
      </Column>
      <Column field="updated_at" header="Updated At">
        <template #body="slotProps">
          {{ formatDate(slotProps.data.updated_at) }}
        </template>
      </Column>
    </DataTable>
  </div>
</template>

<script setup>
import { useAdminStore } from '@/store/admin';
import { onMounted, ref } from 'vue';
import { DataTable, Column } from 'primevue';
import { formatDate } from '@/utils/utils';

const adminStore = useAdminStore();
const users = ref(null);

onMounted(async () => {
  users.value = await adminStore.getAllUsers();
});

</script>

<style lang="scss" scoped></style>
