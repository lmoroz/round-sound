<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  cover: string;
}>()

const coverStyle = computed(() => {
  if (!props.cover) return {}
  return {
    backgroundImage: `url(${props.cover})`,
  }
})
</script>

<template>
  <div
    class="album-cover"
    :style="coverStyle"
  >
    <!-- Fallback icon when no cover -->
    <div
      v-if="!cover"
      class="cover-placeholder"
    >
      <svg
        fill="none"
        stroke="currentColor"
        stroke-width="1.5"
        viewBox="0 0 24 24"
      >
        <circle
          cx="12"
          cy="12"
          r="10"
        />
        <circle
          cx="12"
          cy="12"
          r="3"
        />
      </svg>
    </div>
  </div>
</template>

<style scoped>
.album-cover {
  position: absolute;
  inset: 0;
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  background-color: rgba(30, 30, 40, 1);
}

.cover-placeholder {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(
    135deg,
    rgba(0, 212, 170, 0.1) 0%,
    rgba(0, 184, 212, 0.1) 100%
  );
}

.cover-placeholder svg {
  width: 80px;
  height: 80px;
  color: var(--color-primary);
  opacity: 0.4;
}
</style>
