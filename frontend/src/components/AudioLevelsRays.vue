<script setup lang="ts">
import {
  onMounted,
  onUnmounted,
  ref,
  watch,
} from 'vue'

const props = defineProps<{
  levels: number[];
  isPlaying: boolean;
}>()

const canvasRef = ref<HTMLCanvasElement | null>(null)
let animationId: number | null = null
let currentLevels: number[] = []

const size = 380
const innerRadius = 160
const maxRayLength = 30
const rayCount = 64

// Smoothing factor for level transitions
const smoothingFactor = 0.15

function draw() {
  const canvas = canvasRef.value
  if (!canvas) return

  const ctx = canvas.getContext('2d')
  if (!ctx) return

  // Clear canvas
  ctx.clearRect(0, 0, size, size)

  const centerX = size / 2
  const centerY = size / 2

  // Smooth level transitions
  if (currentLevels.length !== props.levels.length) {
    currentLevels = [...props.levels]
  }
  else {
    for (let i = 0; i < props.levels.length; i++) {
      const current = currentLevels[i] ?? 0
      const target = props.levels[i] ?? 0
      currentLevels[i] = current + (target - current) * smoothingFactor
    }
  }

  // Draw rays
  for (let i = 0; i < rayCount; i++) {
    const angle = (i / rayCount) * Math.PI * 2 - Math.PI / 2
    const levelIndex = Math.floor((i / rayCount) * currentLevels.length)
    const level = currentLevels[levelIndex] || 0

    // Calculate ray length based on level
    const rayLength = props.isPlaying ? level * maxRayLength : maxRayLength * 0.1

    // Start and end points
    const startX = centerX + Math.cos(angle) * innerRadius
    const startY = centerY + Math.sin(angle) * innerRadius
    const endX = centerX + Math.cos(angle) * (innerRadius + rayLength)
    const endY = centerY + Math.sin(angle) * (innerRadius + rayLength)

    // Create gradient for each ray
    const gradient = ctx.createLinearGradient(startX, startY, endX, endY)

    if (props.isPlaying) {
      gradient.addColorStop(0, 'rgba(0, 212, 170, 0.8)')
      gradient.addColorStop(0.5, 'rgba(0, 229, 255, 0.6)')
      gradient.addColorStop(1, 'rgba(0, 184, 212, 0.2)')
    }
    else {
      gradient.addColorStop(0, 'rgba(100, 100, 100, 0.3)')
      gradient.addColorStop(1, 'rgba(100, 100, 100, 0.1)')
    }

    // Draw ray
    ctx.beginPath()
    ctx.moveTo(startX, startY)
    ctx.lineTo(endX, endY)
    ctx.strokeStyle = gradient
    ctx.lineWidth = 2
    ctx.lineCap = 'round'
    ctx.stroke()
  }

  animationId = requestAnimationFrame(draw)
}

onMounted(() => {
  // Set canvas size for high DPI displays
  const canvas = canvasRef.value
  if (canvas) {
    const dpr = window.devicePixelRatio || 1
    canvas.width = size * dpr
    canvas.height = size * dpr
    canvas.style.width = `${size}px`
    canvas.style.height = `${size}px`

    const ctx = canvas.getContext('2d')
    if (ctx) ctx.scale(dpr, dpr)
  }

  draw()
})

onUnmounted(() => {
  if (animationId) {
    cancelAnimationFrame(animationId)
  }
})

// Restart animation when levels change significantly
watch(() => props.levels, () => {
  // Animation loop handles updates automatically
}, { deep: true })
</script>

<template>
  <canvas
    ref="canvasRef"
    class="audio-rays"
  />
</template>

<style scoped>
.audio-rays {
  position: absolute;
  pointer-events: none;
}
</style>
