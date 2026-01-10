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
  isPlaying: boolean;
}>()

const { colorScheme } = useSettings()

const canvasRef = ref<HTMLCanvasElement | null>(null)
let animationId: number | null = null
let currentLevels: number[] = []

const size = 580
const innerRadius = 160
const maxRayLength = 130
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

    // Create gradient using color scheme from settings
    const gradient = generateRayGradient(ctx, startX, startY, endX, endY, colorScheme.value, props.isPlaying)

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
