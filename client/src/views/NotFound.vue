<template>
  <main class="not-found-shell">
    <p class="not-found-code">404</p>
    <h1 class="not-found-title">{{ $t('route.notFound') }}</h1>
    <router-link class="not-found-link" to="/ai-chat">{{ $t('admin.backToChat') }}</router-link>
  </main>
</template>

<script setup>
</script>

<style scoped>
.not-found-shell {
  min-height: 100vh;
  display: grid;
  place-content: center;
  justify-items: center;
  gap: 20px;
  background:
    radial-gradient(circle at 50% 20%, rgba(var(--mint-rgb), 0.18), transparent 34%),
    var(--bg-1);
  color: var(--text-1);
}

.not-found-code {
  margin: 0;
  font-family: 'Orbitron', sans-serif;
  font-size: clamp(72px, 18vw, 180px);
  line-height: 0.9;
  color: var(--accent-strong);
}

.not-found-title {
  margin: 0;
  font-size: 24px;
}

.not-found-link {
  padding: 12px 20px;
  border-radius: 999px;
  background: var(--accent);
  color: var(--btn-text-dark);
  text-decoration: none;
  font-weight: 700;
}
</style>
