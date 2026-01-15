<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'

const props = defineProps<{
  x: number;
  y: number;
  show: boolean;
}>()

const emit = defineEmits<{
  close: [];
  quit: [];
}>()

const menuRef = ref<HTMLDivElement | null>(null)

// Adjust menu position if it would go off-screen
const menuStyle = computed(() => {
  const adjustedX = Math.min(props.x, window.innerWidth - 150) // 150px - approximate menu width
  const adjustedY = Math.min(props.y, window.innerHeight - 50) // 50px - approximate menu height

  return {
    left: `${adjustedX}px`,
    top: `${adjustedY}px`,
  }
})

function handleClickOutside(e: MouseEvent) {
  if (menuRef.value && !menuRef.value.contains(e.target as Node)) emit('close')
}

function handleQuit() {
  emit('quit')
  emit('close')
}

onMounted(() => {
  // Add click listener to close menu when clicking outside
  document.addEventListener('click', handleClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<template>
  <Teleport to="body">
    <Transition name="context-menu">
      <div
        v-if="show"
        ref="menuRef"
        class="context-menu"
        :style="menuStyle"
      >
        <div
          class="menu-item"
          @click="handleQuit"
        >
          Выход
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.context-menu {
  position: fixed;
  z-index: 9999;
  background: rgba(20, 20, 20, 0.95);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  box-shadow:
    0 4px 20px rgba(0, 0, 0, 0.5),
    0 0 0 1px rgba(255, 255, 255, 0.05);
  overflow: hidden;
  min-width: 140px;
}

.menu-item {
  padding: 10px 16px;
  color: #e0e0e0;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s ease;
  user-select: none;
}

.menu-item:hover {
  background: rgba(255, 255, 255, 0.1);
  color: var(--color-primary);
}

/* Transition animations */
.context-menu-enter-active,
.context-menu-leave-active {
  transition: all 0.15s ease;
}

.context-menu-enter-from {
  opacity: 0;
  transform: scale(0.95) translateY(-5px);
}

.context-menu-leave-to {
  opacity: 0;
  transform: scale(0.95);
}
</style>
