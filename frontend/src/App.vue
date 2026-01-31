<template>
  <NavView
    v-if="view === 'nav'"
    :items="itemsSorted"
    :category-name="categoryName"
    @open-item="openItem"
  />
  <AdminView
    v-else-if="view === 'admin'"
    :data="data"
    :categories="categoriesSorted"
    :items="itemsSorted"
    :item-count="itemCount"
    :category-name="categoryName"
    :refresh="refresh"
  />
  <LoginView v-else />
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import NavView from './components/NavView.vue'
import AdminView from './components/AdminView.vue'
import LoginView from './components/LoginView.vue'

const view = ref('nav')
const data = reactive({ next_id: 1, items: [], categories: [] })

const categoriesSorted = computed(() =>
  [...(data.categories || [])].sort((a, b) => a.order - b.order || a.name.localeCompare(b.name)),
)

const itemsSorted = computed(() =>
  [...(data.items || [])].sort((a, b) => a.order - b.order || a.name.localeCompare(b.name)),
)

const itemCount = computed(() => (data.items || []).length)

const detectView = () => {
  const path = window.location.pathname || '/'
  if (path.startsWith('/admin')) return 'admin'
  if (path.startsWith('/login')) return 'login'
  return 'nav'
}

const categoryName = (id) => {
  if (id == null) return '未分类'
  const cat = (data.categories || []).find((entry) => entry.id === id)
  return cat ? cat.name : '未知'
}

const loadData = async () => {
  const res = await fetch('/api/data')
  if (!res.ok) {
    if (view.value === 'admin' && res.status === 401) {
      window.location.href = '/login'
      return
    }
    throw new Error('load failed')
  }
  const payload = await res.json()
  console.log('api/data payload', payload, payload.items?.length)
  if (!payload.categories) payload.categories = []
  if (!payload.items) payload.items = []
  Object.assign(data, payload)
}

const refresh = async () => {
  try {
    await loadData()
  } catch (err) {
    if (view.value !== 'admin') return
    throw err
  }
}

const openItem = (item) => {
  if (item.url) {
    window.open(item.url, '_blank')
  }
}

onMounted(async () => {
  view.value = detectView()
  document.body.className = `${view.value}-page`
  await refresh()
})
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Manrope:wght@400;500;600;700&family=Sora:wght@500;700&display=swap');

:root {
  --bg: #f6f9fc;
  --bg-elevated: #ffffff;
  --ink: #111827;
  --muted: #6b7280;
  --accent: #e84d3d;
  --accent-soft: #fee7e4;
  --border: #e6edf5;
  --shadow: 0 18px 40px rgba(15, 23, 42, 0.08);
  --radius-lg: 16px;
  --radius-md: 12px;
  --radius-sm: 10px;
}

* {
  box-sizing: border-box;
}

body {
  margin: 0;
  font-family: 'Manrope', 'Segoe UI', sans-serif;
  color: var(--ink);
  background: var(--bg);
}

a {
  color: inherit;
}
</style>
