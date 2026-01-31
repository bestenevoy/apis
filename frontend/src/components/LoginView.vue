<template>
  <main class="login-main">
    <section class="panel login-panel">
      <p class="eyebrow">NAV ACCESS</p>
      <h2>后台登录</h2>
      <form class="login-form" @submit.prevent="login">
        <div>
          <label>账号</label>
          <input type="text" v-model.trim="loginForm.username" required />
        </div>
        <div>
          <label>密码</label>
          <input type="password" v-model.trim="loginForm.password" required />
        </div>
        <div>
          <label>&nbsp;</label>
          <button type="submit">登录</button>
        </div>
      </form>
      <p class="muted">{{ loginError }}</p>
    </section>
  </main>
</template>

<script setup>
import { reactive, ref } from 'vue'

const loginForm = reactive({
  username: '',
  password: '',
})

const loginError = ref('')

const login = async () => {
  const res = await fetch('/api/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      username: loginForm.username.trim(),
      password: loginForm.password.trim(),
    }),
  })
  if (res.ok) {
    window.location.href = '/admin'
  } else {
    loginError.value = '账号或密码错误'
  }
}
</script>

<style scoped>
.login-main {
  min-height: 100vh;
  display: grid;
  place-items: center;
  padding: 24px;
}

.panel {
  background: var(--bg-elevated);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 24px;
  box-shadow: var(--shadow);
}

.panel h2 {
  margin-top: 0;
  font-family: 'Sora', sans-serif;
  font-size: 22px;
}

.login-panel {
  max-width: 420px;
  width: 100%;
}

.login-form {
  grid-template-columns: 1fr;
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

input {
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

button:hover {
  transform: translateY(-1px);
  box-shadow: 0 10px 18px rgba(15, 23, 42, 0.14);
}

.eyebrow {
  margin: 0 0 6px;
  font-size: 11px;
  letter-spacing: 0.3em;
  text-transform: uppercase;
  color: rgba(15, 23, 42, 0.6);
  font-weight: 600;
}

.muted {
  color: var(--muted);
  font-size: 13px;
}
</style>
