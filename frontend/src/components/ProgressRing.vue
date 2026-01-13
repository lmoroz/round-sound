<script setup lang="ts">
import {
  computed,
  ref,
} from 'vue'

const props = defineProps<{
  progress: number; // 0 to 1
  duration: number; // seconds
}>()

const emit = defineEmits<{
  (e: 'seek', position: number): void;
}>()

const size = 310
const thumbPadding = 20 // Extra space for thumb + glow
const viewBoxSize = size + thumbPadding * 2
const strokeWidth = 3
const radius = (size - strokeWidth) / 2
const circumference = 2 * Math.PI * radius
const center = viewBoxSize / 2

const isHovered = ref(false)
const isDragging = ref(false)
const dragProgress = ref(0)
const svgRef = ref<SVGSVGElement | null>(null)

const strokeDashoffset = computed(() => {
  return circumference * (1 - props.progress)
})

const currentProgress = computed(() => {
  return isDragging.value ? dragProgress.value : props.progress
})

// Thumb position (clockwise starting from top)
// Since we rotate SVG by -90deg, we need to calculate position accordingly
const thumbPosition = computed(() => {
  // Progress 0 = top (12 o'clock), Progress 1 = full circle back to top
  // In the rotated SVG space, 0 progress = right side (3 o'clock in unrotated)
  const angle = currentProgress.value * 2 * Math.PI
  return {
    x: center + radius * Math.cos(angle),
    y: center + radius * Math.sin(angle),
  }
})

// Calculate progress from mouse/touch position relative to SVG center
function getProgressFromClientPosition(clientX: number, clientY: number): number {
  const svg = svgRef.value
  if (!svg) return 0

  const rect = svg.getBoundingClientRect()
  const svgCenterX = rect.left + rect.width / 2
  const svgCenterY = rect.top + rect.height / 2

  // Vector from center to mouse in screen coordinates
  const dx = clientX - svgCenterX
  const dy = clientY - svgCenterY

  // The SVG is rotated -90deg, so we need to account for that
  // In screen space: up is negative Y, right is positive X
  // After -90deg rotation: what was "right" in SVG is now "up" on screen
  // So we need to rotate the mouse vector by +90deg to get SVG coordinates

  // atan2(dy, dx) gives angle from positive X axis (right)
  // We want angle from positive Y axis (down in screen = right in rotated SVG)
  // After rotation: screen up (-Y) corresponds to SVG right (+X at 0 progress)

  let angle = Math.atan2(dy, dx)
  // Rotate by +90 degrees to account for SVG rotation
  angle += Math.PI / 2
  // Normalize to 0-2PI
  if (angle < 0) angle += 2 * Math.PI

  return Math.min(1, Math.max(0, angle / (2 * Math.PI)))
}

function onMouseDown(event: MouseEvent | TouchEvent) {
  event.preventDefault()
  isDragging.value = true

  const clientX = 'touches' in event ? event.touches[0]?.clientX ?? 0 : event.clientX
  const clientY = 'touches' in event ? event.touches[0]?.clientY ?? 0 : event.clientY
  dragProgress.value = getProgressFromClientPosition(clientX, clientY)

  const onMove = (e: MouseEvent | TouchEvent) => {
    if (isDragging.value) {
      e.preventDefault()
      const cx = 'touches' in e ? e.touches[0]?.clientX ?? 0 : e.clientX
      const cy = 'touches' in e ? e.touches[0]?.clientY ?? 0 : e.clientY
      dragProgress.value = getProgressFromClientPosition(cx, cy)
    }
  }

  const onUp = () => {
    if (isDragging.value) {
      isDragging.value = false
      const seekPosition = dragProgress.value * props.duration
      emit('seek', seekPosition)
    }
    document.removeEventListener('mousemove', onMove)
    document.removeEventListener('mouseup', onUp)
    document.removeEventListener('touchmove', onMove)
    document.removeEventListener('touchend', onUp)
  }

  document.addEventListener('mousemove', onMove)
  document.addEventListener('mouseup', onUp)
  document.addEventListener('touchmove', onMove, { passive: false })
  document.addEventListener('touchend', onUp)
}

function onTrackClick(event: MouseEvent) {
  if (isDragging.value) return
  const progress = getProgressFromClientPosition(event.clientX, event.clientY)
  const seekPosition = progress * props.duration
  emit('seek', seekPosition)
}
</script>

<template>
  <svg
    ref="svgRef"
    class="progress-ring"
    :class="{ 'is-interactive': isHovered || isDragging }"
    :height="viewBoxSize"
    :viewBox="`0 0 ${viewBoxSize} ${viewBoxSize}`"
    :width="viewBoxSize"
    @mouseenter="isHovered = true"
    @mouseleave="isHovered = false"
  >
    <!-- Invisible wider track for easier clicking/hovering -->
    <circle
      class="progress-hit-area"
      :cx="center"
      :cy="center"
      fill="none"
      :r="radius"
      stroke-width="32"
      @click="onTrackClick"
      @mousedown.prevent="onMouseDown"
      @touchstart.prevent="onMouseDown"
    />

    <!-- Background track -->
    <circle
      class="progress-track"
      :cx="center"
      :cy="center"
      fill="none"
      :r="radius"
      :stroke-width="strokeWidth"
    />

    <!-- Progress arc -->
    <circle
      class="progress-bar"
      :class="{ 'is-dragging': isDragging }"
      :cx="center"
      :cy="center"
      fill="none"
      :r="radius"
      :stroke-dasharray="circumference"
      :stroke-dashoffset="isDragging ? circumference * (1 - dragProgress) : strokeDashoffset"
      :stroke-width="strokeWidth"
    />

    <!-- Thumb button -->
    <g
      class="progress-thumb-group"
      :class="{ 'is-visible': isHovered || isDragging }"
    >
      <!-- Outer glow -->
      <circle
        class="progress-thumb-glow"
        :cx="thumbPosition.x"
        :cy="thumbPosition.y"
        r="16"
      />
      <!-- Main thumb -->
      <circle
        class="progress-thumb"
        :class="{ 'is-dragging': isDragging }"
        :cx="thumbPosition.x"
        :cy="thumbPosition.y"
        :r="isDragging ? 12 : 10"
      />
      <!-- Inner highlight -->
      <circle
        class="progress-thumb-inner"
        :cx="thumbPosition.x"
        :cy="thumbPosition.y"
        r="4"
      />
    </g>
  </svg>
</template>

<style scoped>
.progress-ring {
  position: absolute;
  transform: rotate(-90deg);
  cursor: pointer;
  /* Compensate for thumbPadding to keep ring centered */
  margin: -20px;
  /* Prevent window dragging when interacting with progress */
  --wails-draggable: no-drag;
}

.progress-ring.is-interactive {
  z-index: 10;
}

.progress-hit-area {
  stroke: transparent;
  cursor: pointer;
  pointer-events: stroke;
}

.progress-track {
  stroke: rgba(255, 255, 255, 0.1);
  pointer-events: none;
}

.progress-bar {
  stroke: var(--color-primary);
  stroke-linecap: round;
  transition: stroke-dashoffset 0.3s ease-out;
  filter: drop-shadow(0 0 4px var(--color-primary-glow));
  pointer-events: none;
}

.progress-bar.is-dragging {
  transition: none;
}

.progress-thumb-group {
  opacity: 0;
  transition: opacity 0.2s ease;
  pointer-events: none;
}

.progress-thumb-group.is-visible {
  opacity: 1;
}

.progress-thumb-glow {
  fill: var(--color-primary);
  opacity: 0.25;
  filter: blur(6px);
}

.progress-thumb {
  fill: var(--color-primary);
  filter: drop-shadow(0 0 8px var(--color-primary-glow));
  transition: r 0.15s ease;
  stroke: rgba(255, 255, 255, 0.3);
  stroke-width: 2;
}

.progress-thumb-inner {
  fill: rgba(255, 255, 255, 0.5);
  pointer-events: none;
}
</style>
