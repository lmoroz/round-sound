<script setup lang="ts">
import {
  onMounted,
  onUnmounted,
  ref,
  watch,
} from 'vue'
import { useSettings } from '@/composables/useSettings'
import { generateRayGradient } from '@/utils/colors'

const props = defineProps<{
  levels: number[];
}>()

const { colorScheme } = useSettings()

const canvasRef = ref<HTMLCanvasElement | null>(null)
let animationId: number | null = null
let currentLevels: number[] = []
let currentDpr = 1

const size = 580
const innerRadius = 160
const maxRayLength = 130
const rayCount = 64

// Smoothing factor for level transitions
const smoothingFactor = 0.15

function initializeCanvas() {
  const canvas = canvasRef.value
  if (!canvas) return

  const dpr = window.devicePixelRatio || 1
  currentDpr = dpr

  canvas.width = size * dpr
  canvas.height = size * dpr
  canvas.style.width = `${size}px`
  canvas.style.height = `${size}px`

  const ctx = canvas.getContext('2d')
  if (ctx) ctx.scale(dpr, dpr)
}

function checkDprChange() {
  const dpr = window.devicePixelRatio || 1
  if (dpr !== currentDpr) {
    console.warn(`[AudioLevelsRays] DPR changed: ${currentDpr} → ${dpr}, reinitializing canvas`)
    initializeCanvas()
  }
}

function draw() {
  const canvas = canvasRef.value
  if (!canvas) return

  // Check for DPR changes (e.g., after display scaling or system recovery)
  checkDprChange()

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

  // Determine if there's any sound (system or player)
  const soundThreshold = 0.02 // Минимальный порог для определения звука
  const hasSound = currentLevels.some(level => level > soundThreshold)

  // Draw rays
  for (let i = 0; i < rayCount; i++) {
    const angle = (i / rayCount) * Math.PI * 2 - Math.PI / 2
    const levelIndex = Math.floor((i / rayCount) * currentLevels.length)
    const level = currentLevels[levelIndex] || 0

    // Calculate ray length based on level (always react to system audio)
    const rayLength = Math.max(level * maxRayLength, maxRayLength * 0.05)

    // Start and end points
    const startX = centerX + Math.cos(angle) * innerRadius
    const startY = centerY + Math.sin(angle) * innerRadius
    const endX = centerX + Math.cos(angle) * (innerRadius + rayLength)
    const endY = centerY + Math.sin(angle) * (innerRadius + rayLength)

    // Create gradient - colored when there's any sound (system or player)
    const gradient = generateRayGradient(ctx, startX, startY, endX, endY, colorScheme.value, hasSound)

    // Draw ray
    ctx.beginPath()
    ctx.moveTo(startX, startY)
    ctx.lineTo(endX, endY)
    ctx.strokeStyle = gradient
    ctx.lineWidth = 6
    ctx.lineCap = 'round'
    ctx.stroke()
  }

  animationId = requestAnimationFrame(draw)
}

onMounted(() => {
  initializeCanvas()
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
