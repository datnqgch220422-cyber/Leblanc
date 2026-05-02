<script setup>
import { useRouter } from 'vue-router'
import { clearSessionUser } from '@/composables/useSessionAuth'

const router = useRouter()

const handleLogout = () => {
  clearSessionUser()
  router.push('/login')
}
</script>

<template>
  <div class="admin-layout">
    <aside class="admin-sidebar">
      <h2>LeBlanc Admin</h2>
      <nav>
        <router-link to="/admin/users">Quản lý Users</router-link>
        <router-link to="/admin/bookings">Quản lý Đặt bàn</router-link>
        <router-link to="/admin/drinks">Quản lý Đồ uống</router-link>
      </nav>
      <button @click="handleLogout" class="btn-logout">Đăng xuất</button>
    </aside>

    <main class="admin-content">
      <router-view />
    </main>
  </div>
</template>

<style scoped>
.admin-layout {
  display: flex;
  min-height: 100vh;
}
.admin-sidebar {
  width: 250px;
  background: var(--paper);
  border-right: 1px solid rgba(0,0,0,0.1);
  padding: 20px;
  display: flex;
  flex-direction: column;
}
.admin-sidebar nav {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-top: 20px;
  flex-grow: 1;
}
.admin-content {
  flex-grow: 1;
  padding: 24px;
  background: var(--bg);
}
</style>