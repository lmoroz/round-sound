<script setup lang="ts">
import { computed } from 'vue'

import {
  Pause,
  Play,
  Repeat,
  Repeat1,
  Shuffle,
  SkipBack,
  SkipForward,
  ThumbsDown,
  Heart,
} from 'lucide-vue-next'

import { RepeatMode } from '@/types'

const props = defineProps<{
  isPlaying: boolean;
  shuffle: boolean;
  repeat: RepeatMode;
  rating: number;
  canSetRating: boolean;
  canSetShuffle: boolean;
  canSetRepeat: boolean;
}>()

const emit = defineEmits<{
  togglePlayPause: [];
  next: [];
  previous: [];
  toggleShuffle: [];
  toggleRepeat: [];
  setRating: [rating: number];
}>()

const isLiked = computed(() => props.rating === 5)
const isDisliked = computed(() => props.rating === 1)

const repeatIcon = computed(() => {
  if (props.repeat === RepeatMode.One) return Repeat1
  return Repeat
})

const isRepeatActive = computed(() =>
  props.repeat === RepeatMode.All || props.repeat === RepeatMode.One,
)

function handleLike() {
  emit('setRating', isLiked.value ? 0 : 5)
}

function handleDislike() {
  emit('setRating', isDisliked.value ? 0 : 1)
}
</script>

<template>
  <div class="media-controls">
    <!-- Like/Dislike Row -->
    <div
      v-if="canSetRating"
      class="rating-row"
    >
      <button
        class="control-button small"
        :class="{ 'active-heart': isLiked }"
        title="Like"
        @click="handleLike"
      >
        <Heart
          :fill="isLiked ? 'currentColor' : 'none'"
          :size="24"
        />
      </button>
      <button
        class="control-button small"
        :class="{ active: isDisliked }"
        title="Dislike"
        @click="handleDislike"
      >
        <ThumbsDown :size="16" />
      </button>
    </div>

    <!-- Main Controls Row -->
    <div class="main-controls">
      <button
        v-if="canSetShuffle"
        class="control-button"
        :class="{ active: shuffle }"
        title="Shuffle"
        @click="emit('toggleShuffle')"
      >
        <Shuffle :size="18" />
      </button>

      <button
        class="control-button"
        title="Previous"
        @click="emit('previous')"
      >
        <SkipBack :size="20" />
      </button>

      <button
        class="control-button play-button"
        :title="isPlaying ? 'Pause' : 'Play'"
        @click="emit('togglePlayPause')"
      >
        <Pause
          v-if="isPlaying"
          :size="24"
        />
        <Play
          v-else
          :size="24"
        />
      </button>

      <button
        class="control-button"
        title="Next"
        @click="emit('next')"
      >
        <SkipForward :size="20" />
      </button>

      <button
        v-if="canSetRepeat"
        class="control-button"
        :class="{ active: isRepeatActive }"
        title="Repeat"
        @click="emit('toggleRepeat')"
      >
        <component
          :is="repeatIcon"
          :size="18"
        />
      </button>
    </div>
  </div>
</template>

<style scoped>
.media-controls {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  width: calc(100% + 50px);
  background: var(--color-background);
  padding: 12px 0;
}

.rating-row {
  display: flex;
  gap: 15px;
  margin-bottom: 5px;
  position: relative;
  z-index: 100;
}

.main-controls {
  display: flex;
  align-items: center;
  gap: 8px;
  position: relative;
  z-index: 100;
}

.control-button {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: none;
  border-radius: 50%;
  background: transparent;
  color: var(--color-text-secondary);
  transition: color 0.2s ease,  transform 0.2s ease;
}

.control-button:hover {
  color: var(--color-text);
  transform: scale(1.2);
}

.control-button:active {
  transform: scale(0.95);
}

.control-button.active {
  color: var(--color-primary);
  background: rgba(0, 212, 170, 0.2);
}

.control-button.small {
  width: 28px;
  height: 28px;
}

.play-button {
  width: 44px;
  height: 44px;
  background: linear-gradient(
    135deg,
    var(--color-primary) 0%,
    var(--color-secondary) 100%
  );
  box-shadow: 0 4px 15px var(--color-primary-glow);
}

.play-button:hover {
  background: linear-gradient(
    135deg,
    var(--color-accent) 0%,
    var(--color-primary) 100%
  );
}

.control-button.active-heart {
  color: var(--color-primary);
  background: rgba(255, 255, 255, 0.2);
  box-shadow: 0 4px 15px var(--color-primary-glow);
}
</style>
