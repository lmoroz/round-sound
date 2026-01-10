<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  progress: number; // 0 to 1
}>()

const size = 310
const strokeWidth = 3
const radius = (size - strokeWidth) / 2
const circumference = 2 * Math.PI * radius

const strokeDashoffset = computed(() => {
  return circumference * (1 - props.progress)
})
</script>

<template>
  <svg
    class="progress-ring"
    :height="size"
    :viewBox="`0 0 ${size} ${size}`"
    :width="size"
  >
    <!-- Background track -->
    <circle
      class="progress-track"
      :cx="size / 2"
      :cy="size / 2"
      fill="none"
      :r="radius"
      :stroke-width="strokeWidth"
    />

    <!-- Progress arc -->
    <circle
      class="progress-bar"
      :cx="size / 2"
      :cy="size / 2"
      fill="none"
      :r="radius"
      :stroke-dasharray="circumference"
      :stroke-dashoffset="strokeDashoffset"
      :stroke-width="strokeWidth"
    />
  </svg>
</template>

<style scoped>
.progress-ring {
  position: absolute;
  transform: rotate(-90deg);
  pointer-events: none;
}

.progress-track {
  stroke: rgba(255, 255, 255, 0.1);
}

.progress-bar {
  stroke: var(--color-primary);
  stroke-linecap: round;
  transition: stroke-dashoffset 0.3s ease-out;
  filter: drop-shadow(0 0 4px var(--color-primary-glow));
}
</style>
