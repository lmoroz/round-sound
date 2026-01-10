<script setup lang="ts">
defineProps<{
  title: string;
  artist: string;
}>()
</script>

<template>
  <div class="track-info">
    <h2 class="track-title text-truncate">
      {{ title || 'No Track' }}
    </h2>
    <p class="track-artist text-truncate">
      {{ artist || 'Unknown Artist' }}
    </p>
  </div>
</template>

<style scoped>
.track-info {
  text-align: center;
  max-width: 100%;
  padding: 0 10px;
  margin-bottom: 15px;
}

.track-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: 4px;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.5);
  letter-spacing: 0.5px;
}

.track-artist {
  font-size: 13px;
  font-weight: 400;
  color: var(--color-text-secondary);
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.5);
}
</style>
