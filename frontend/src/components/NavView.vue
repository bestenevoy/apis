<template>
  <div class="nav-page">
    <header class="nav-hero">
      <div class="nav-hero__inner">
        <div class="nav-brand">
          <div class="nav-logo">N</div>
          <div>
            <div class="nav-brand__name">OpenNav</div>
            <div class="nav-brand__sub">Navigation Hub</div>
          </div>
        </div>
        <nav class="nav-menu">
          <a href="/" class="active">导航</a>
          <a href="/admin">管理</a>
          <a href="/login">登录</a>
        </nav>
      </div>
    </header>
    <section class="nav-intro">
      <div>
        <p class="nav-intro__kicker">Choose your region and platform to get started.</p>
        <h1>资源导航平台</h1>
      </div>
      <!-- <a class="admin-link" href="/admin">后台管理</a> -->
    </section>
    <main class="nav-main">
      <section class="nav-section" v-for="group in groupedItems" :key="group.name">
        <div class="nav-section__head">
          <div>
            <p class="nav-section__eyebrow">Category</p>
            <h2>{{ group.name }}</h2>
            <p>已收录 {{ group.items.length }} 条链接。</p>
          </div>
          <div class="nav-pill">{{ group.items.length }} 项</div>
        </div>
        <div class="nav-grid">
          <button
            class="nav-card"
            type="button"
            v-for="item in group.items"
            :key="item.id"
            @click="emit('open-item', item)"
          >
            <div class="nav-card__top">
              <span class="nav-tag">RECOMMENDED</span>
              <!-- <span class="nav-open">Deploy ↗</span> -->
            </div>
            <div class="nav-meta">
              <img v-if="item.avatar_url" class="nav-avatar" :src="item.avatar_url" :alt="item.name" />
              <div v-else class="nav-fallback">{{ item.name ? item.name.slice(0, 1) : '?' }}</div>
              <div>
                <h3>{{ item.name }}</h3>
                <p>{{ item.summary || '暂无简介' }}</p>
              </div>
            </div>
            <div class="nav-chip-row">
              <span class="nav-chip">{{ group.name }}</span>
              <span class="nav-chip">快速访问</span>
            </div>
          </button>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup>
import { computed, toRefs } from 'vue'

const props = defineProps({
  items: { type: Array, default: () => [] },
  categoryName: { type: Function, required: true }
});

const emit = defineEmits(['open-item'])
const { items, categoryName } = toRefs(props)

const groupedItems = computed(() => {
  const groups = []
  const seen = new Map()
  ;(items.value || []).forEach((item) => {
    const name = categoryName.value(item.category_id)
    if (!seen.has(name)) {
      const group = { name, items: [] }
      seen.set(name, group)
      groups.push(group)
    }
    seen.get(name).items.push(item)
  })
  return groups
})
</script>

<style scoped>
.nav-page {
  min-height: 100vh;
}

.nav-hero {
  padding: 18px 6vw 0;
  border-bottom: 1px solid var(--border);
  background: #fff;
}

.nav-hero__inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding-bottom: 14px;
}

.nav-brand {
  display: flex;
  align-items: center;
  gap: 12px;
  font-family: 'Sora', sans-serif;
}

.nav-logo {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: var(--accent);
  color: #fff;
  display: grid;
  place-items: center;
  font-weight: 700;
}

.nav-brand__name {
  font-weight: 700;
}

.nav-brand__sub {
  font-size: 12px;
  color: var(--muted);
}

.nav-menu {
  display: flex;
  gap: 18px;
  font-size: 14px;
}

.nav-menu a {
  text-decoration: none;
  color: var(--muted);
}

.nav-menu a.active {
  color: var(--accent);
  font-weight: 600;
}

.nav-intro {
  padding: 24px 6vw 10px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.nav-intro__kicker {
  margin: 0 0 8px;
  color: var(--muted);
  font-size: 13px;
}

.nav-intro h1 {
  margin: 0;
  font-family: 'Sora', sans-serif;
  font-size: clamp(26px, 4vw, 34px);
}

.admin-link {
  text-decoration: none;
  padding: 10px 16px;
  border-radius: 999px;
  border: 1px solid var(--border);
  color: var(--ink);
  font-weight: 600;
  background: #fff;
  transition: all 0.2s ease;
}

.admin-link:hover {
  transform: translateY(-1px);
  box-shadow: 0 10px 18px rgba(15, 23, 42, 0.12);
}

.nav-main {
  padding: 0 6vw 40px;
}

.nav-section {
  margin-top: 18px;
}

.nav-section__head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  margin-bottom: 18px;
}

.nav-section__head h2 {
  margin: 0;
  font-family: 'Sora', sans-serif;
  font-size: 18px;
}

.nav-section__head p {
  margin: 6px 0 0;
  color: var(--muted);
}

.nav-section__eyebrow {
  margin: 0 0 6px;
  font-size: 11px;
  letter-spacing: 0.3em;
  text-transform: uppercase;
  color: var(--accent);
  font-weight: 600;
}

.nav-pill {
  padding: 6px 14px;
  border-radius: 999px;
  background: var(--accent-soft);
  color: var(--accent);
  font-weight: 600;
  font-size: 12px;
}

.nav-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 16px;
}

.nav-card {
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  background: #ffffff;
  text-align: left;
  cursor: pointer;
  transition:
    transform 0.2s ease,
    box-shadow 0.2s ease,
    border-color 0.2s ease;
}

.nav-card:hover {
  transform: translateY(-2px);
  border-color: rgba(232, 77, 61, 0.4);
  box-shadow: var(--shadow);
}

.nav-card h3 {
  margin: 0;
  font-size: 16px;
}

.nav-card p {
  margin: 4px 0 8px;
  color: var(--muted);
  font-size: 13px;
}

.nav-card__top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 12px;
}

.nav-tag {
  padding: 4px 10px;
  border-radius: 999px;
  background: var(--accent-soft);
  color: var(--accent);
  font-weight: 600;
}

.nav-open {
  color: var(--accent);
  font-weight: 600;
}

.nav-meta {
  display: flex;
  align-items: center;
  gap: 12px;
}

.nav-avatar {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  object-fit: cover;
  border: 1px solid var(--border);
}

.nav-fallback {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: var(--accent-soft);
  display: grid;
  place-items: center;
  color: var(--accent);
  font-weight: 700;
}

.nav-chip-row {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.nav-chip {
  font-size: 11px;
  padding: 4px 8px;
  border-radius: 8px;
  border: 1px solid var(--border);
  color: var(--muted);
}

@media (max-width: 900px) {
  .nav-hero__inner {
    flex-direction: column;
    align-items: flex-start;
  }

  .nav-section__head {
    flex-direction: column;
    align-items: flex-start;
  }

  .nav-intro {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
