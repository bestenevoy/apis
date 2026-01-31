<template>
  <div>
    <header class="admin-header">
      <div>
        <p class="eyebrow">NAV CONSOLE</p>
        <h1>后台管理</h1>
        <p>维护分类、链接与访问控制。</p>
      </div>
      <div class="admin-actions">
        <a class="admin-link" href="/">返回导航</a>
        <button class="ghost" type="button" @click="logout">退出登录</button>
      </div>
    </header>

    <main class="admin-main">
      <section class="panel">
        <h2>类别管理</h2>
        <form @submit.prevent="saveCategory">
          <div>
            <label>类别名称</label>
            <input type="text" v-model.trim="categoryForm.name" required />
          </div>
          <div>
            <label>排序</label>
            <input type="number" v-model.number="categoryForm.order" />
          </div>
          <div class="operator">
            <div>
              <label>&nbsp;</label>
              <button type="submit">
                {{ editingCategoryId == null ? '添加类别' : '保存类别' }}
              </button>
            </div>
            <div>
              <label>&nbsp;</label>
              <button type="button" class="ghost" @click="resetCategoryForm">取消</button>
            </div>
          </div>
        </form>
        <table>
          <thead>
            <tr>
              <th>名称</th>
              <th>排序</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="cat in categories" :key="cat.id">
              <td>{{ cat.name }}</td>
              <td>{{ cat.order }}</td>
              <td>
                <div class="actions">
                  <button type="button" class="secondary" @click="startEditCategory(cat)">
                    编辑
                  </button>
                  <button type="button" class="ghost" @click="deleteCategory(cat.id)">删除</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </section>

      <section class="panel">
        <h2>
          链接管理 <span class="pill">{{ itemCount }} items</span>
        </h2>
        <form @submit.prevent="saveItem">
          <div>
            <label>名称</label>
            <input type="text" v-model.trim="itemForm.name" required />
          </div>
          <div>
            <label>链接</label>
            <input type="url" v-model.trim="itemForm.url" placeholder="https://" required />
          </div>
          <div>
            <label>头像</label>
            <input type="url" v-model.trim="itemForm.avatar_url" placeholder="https://" />
          </div>
          <div>
            <label>简介</label>
            <input type="text" v-model.trim="itemForm.summary" placeholder="一句话简介" />
          </div>
          <div>
            <label>所属类别</label>
            <select v-model="itemForm.category_id">
              <option :value="null">未分类</option>
              <option v-for="cat in categories" :key="cat.id" :value="cat.id">
                {{ cat.name }}
              </option>
            </select>
          </div>
          <div>
            <label>排序</label>
            <input type="number" v-model.number="itemForm.order" />
          </div>
          <div class="operator">
            <div>
              <label>&nbsp;</label>
              <button type="submit">{{ editingId == null ? '添加' : '保存' }}</button>
            </div>
            <div>
              <label>&nbsp;</label>
              <button type="button" class="ghost" @click="resetItemForm">取消</button>
            </div>
          </div>
        </form>

        <div class="toolbar">
          <button type="button" class="secondary" @click="refresh">刷新</button>
          <button type="button" @click="backupData">备份数据</button>
          <label class="ghost" style="padding: 10px 14px; border-radius: 10px; cursor: pointer">
            恢复数据
            <input
              type="file"
              accept="application/json"
              style="display: none"
              @change="restoreData"
            />
          </label>
        </div>

        <table>
          <thead>
            <tr>
              <th>名称</th>
              <th>链接</th>
              <th>类别</th>
              <th>排序</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in items" :key="item.id">
              <td>{{ item.name }}</td>
              <td>
                <a :href="item.url" target="_blank">{{ item.url }}</a>
              </td>
              <td>{{ categoryName(item.category_id) }}</td>
              <td>{{ item.order }}</td>
              <td>
                <div class="actions">
                  <button type="button" class="secondary" @click="startEditItem(item)">编辑</button>
                  <button type="button" class="ghost" @click="deleteItem(item.id)">删除</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </section>

      <section class="panel">
        <h2>修改密码</h2>
        <form class="password-form" @submit.prevent="changePassword">
          <div>
            <label>旧密码</label>
            <input type="password" v-model.trim="passwordForm.old" required />
          </div>
          <div>
            <label>新密码</label>
            <input type="password" v-model.trim="passwordForm.new" required />
          </div>
          <div class="operator">
            <label>&nbsp;</label>
            <button type="submit" class="secondary">修改密码</button>
          </div>
        </form>
        <p class="muted">{{ passwordMsg }}</p>
      </section>
    </main>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'

const props = defineProps({
  data: { type: Object, required: true },
  categories: { type: Array, default: () => [] },
  items: { type: Array, default: () => [] },
  itemCount: { type: Number, default: 0 },
  categoryName: { type: Function, required: true },
  refresh: { type: Function, required: true },
})

const editingId = ref(null)
const editingCategoryId = ref(null)
const passwordMsg = ref('')

const itemForm = reactive({
  name: '',
  url: '',
  avatar_url: '',
  summary: '',
  category_id: null,
  order: 0,
})

const categoryForm = reactive({
  name: '',
  order: 0,
})

const passwordForm = reactive({
  old: '',
  new: '',
})

const resetItemForm = () => {
  editingId.value = null
  Object.assign(itemForm, {
    name: '',
    url: '',
    avatar_url: '',
    summary: '',
    category_id: null,
    order: 0,
  })
}

const startEditItem = (item) => {
  editingId.value = item.id
  Object.assign(itemForm, {
    name: item.name || '',
    url: item.url || '',
    avatar_url: item.avatar_url || '',
    summary: item.summary || '',
    category_id: item.category_id == null ? null : item.category_id,
    order: Number(item.order || 0),
  })
}

const saveItem = async () => {
  const payload = {
    name: itemForm.name.trim(),
    url: itemForm.url.trim(),
    avatar_url: itemForm.avatar_url.trim(),
    summary: itemForm.summary.trim(),
    category_id: itemForm.category_id == null ? null : Number(itemForm.category_id),
    order: Number(itemForm.order || 0),
  }
  if (!payload.name || !payload.url) return
  if (editingId.value == null) {
    await fetch('/api/item', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    })
  } else {
    await fetch(`/api/item/${editingId.value}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    })
  }
  resetItemForm()
  await props.refresh()
}

const deleteItem = async (id) => {
  if (!confirm('确认删除该链接?')) return
  await fetch(`/api/item/${id}`, { method: 'DELETE' })
  await props.refresh()
}

const resetCategoryForm = () => {
  editingCategoryId.value = null
  Object.assign(categoryForm, { name: '', order: 0 })
}

const startEditCategory = (cat) => {
  editingCategoryId.value = cat.id
  Object.assign(categoryForm, { name: cat.name, order: Number(cat.order || 0) })
}

const saveCategory = async () => {
  const payload = {
    name: categoryForm.name.trim(),
    order: Number(categoryForm.order || 0),
  }
  if (!payload.name) return
  if (editingCategoryId.value == null) {
    await fetch('/api/category', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    })
  } else {
    await fetch(`/api/category/${editingCategoryId.value}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    })
  }
  resetCategoryForm()
  await props.refresh()
}

const deleteCategory = async (id) => {
  if (!confirm('确认删除该类别?')) return
  await fetch(`/api/category/${id}`, { method: 'DELETE' })
  await props.refresh()
}

const backupData = async () => {
  const res = await fetch('/api/data')
  const text = await res.text()
  const blob = new Blob([text], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `nav-backup-${new Date().toISOString().slice(0, 10)}.json`
  a.click()
  URL.revokeObjectURL(url)
}

const restoreData = async (event) => {
  const file = event.target.files && event.target.files[0]
  if (!file) return
  const text = await file.text()
  await fetch('/api/data', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: text,
  })
  event.target.value = ''
  await props.refresh()
}

const logout = async () => {
  await fetch('/api/logout', { method: 'POST' })
  window.location.href = '/login'
}

const changePassword = async () => {
  const res = await fetch('/api/password', {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      old_password: passwordForm.old.trim(),
      new_password: passwordForm.new.trim(),
    }),
  })
  passwordMsg.value = res.ok ? '密码已更新' : '修改失败'
}
</script>

<style scoped>
.admin-header {
  padding: 32px 6vw 18px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 18px;
}

.admin-header h1 {
  margin: 0;
  font-family: 'Sora', sans-serif;
  font-size: clamp(26px, 4vw, 38px);
}

.admin-header p {
  margin: 6px 0 0;
  color: var(--muted);
}

.eyebrow {
  margin: 0 0 6px;
  font-size: 11px;
  letter-spacing: 0.3em;
  text-transform: uppercase;
  color: rgba(15, 23, 42, 0.6);
  font-weight: 600;
}

.admin-link {
  text-decoration: none;
  padding: 10px 16px;
  border-radius: 999px;
  border: 1px solid var(--border);
  background: var(--bg-elevated);
  font-weight: 600;
}

.admin-actions {
  display: flex;
  gap: 10px;
  align-items: center;
}

.admin-main {
  padding: 0 6vw 40px;
  display: grid;
  gap: 20px;
}

.panel {
  background: var(--bg-elevated);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 18px;
  box-shadow: var(--shadow);
}

.panel h2 {
  margin-top: 0;
  font-family: 'Sora', sans-serif;
  font-size: 20px;
}

form {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 10px;
  margin-bottom: 14px;
}

label {
  font-size: 12px;
  color: var(--muted);
}

input,
select {
  width: 100%;
  padding: 10px 12px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border);
  font-size: 14px;
  background: #fff;
}

button {
  border: none;
  background: var(--accent);
  color: white;
  padding: 10px 14px;
  border-radius: 12px;
  font-weight: 600;
  cursor: pointer;
  transition:
    transform 0.2s ease,
    box-shadow 0.2s ease;
}

button.secondary {
  background: var(--accent-2);
}

button.ghost {
  background: transparent;
  color: var(--muted);
  border: 1px dashed var(--border);
}

button:hover {
  transform: translateY(-1px);
  box-shadow: 0 10px 18px rgba(15, 23, 42, 0.14);
}

table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
}

th,
td {
  text-align: left;
  padding: 8px 6px;
  border-bottom: 1px solid var(--border);
}

.actions {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.toolbar {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  margin-bottom: 12px;
}

.pill {
  display: inline-block;
  padding: 4px 10px;
  border-radius: 999px;
  background: rgba(47, 107, 255, 0.12);
  color: var(--accent);
  font-size: 12px;
}

.muted {
  color: var(--muted);
  font-size: 13px;
}

.operator {
  display: flex;
  align-content: center;
  flex-wrap: nowrap;
  align-items: flex-end;
}

@media (max-width: 900px) {
  .admin-header {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
