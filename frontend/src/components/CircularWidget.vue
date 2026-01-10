<script setup lang="ts">
import { computed } from 'vue'

import { useAudioLevels } from '@/composables/useAudioLevels'
import { useMediaPlayer } from '@/composables/useMediaPlayer'
import { StateMode } from '@/types'

import AlbumCover from './AlbumCover.vue'
import AudioLevelsRays from './AudioLevelsRays.vue'
import MediaControls from './MediaControls.vue'
import ProgressRing from './ProgressRing.vue'
import TrackInfo from './TrackInfo.vue'

const {
  player,
  isConnected,
  togglePlayPause,
  next,
  previous,
  toggleShuffle,
  toggleRepeat,
  setRating,
} = useMediaPlayer()

const { levels } = useAudioLevels(64)

const progress = computed(() => {
  if (!player.value.duration) return 0
  return player.value.position / player.value.duration
})

const isPlaying = computed(() => player.value.state === StateMode.Playing)
</script>

<template>
  <div class="widget-container">
    <!-- Audio visualization rays (outermost layer) -->
    <AudioLevelsRays
      :is-playing="isPlaying"
      :levels="levels"
    />

    <!-- Progress ring -->
    <ProgressRing :progress="progress" />

    <!-- Main circular content -->
    <div class="widget-main">
      <!-- Album cover background -->
      <AlbumCover :cover="player.cover" />

      <!-- Content overlay -->
      <div class="widget-content">
        <!-- Track info -->
        <TrackInfo
          :artist="player.artist"
          :title="player.title"
        />

        <!-- Media controls -->
        <MediaControls
          :can-set-rating="player.canSetRating"
          :is-playing="isPlaying"
          :rating="player.rating"
          :repeat="player.repeat"
          :shuffle="player.shuffle"
          @next="next"
          @previous="previous"
          @set-rating="setRating"
          @toggle-play-pause="togglePlayPause"
          @toggle-repeat="toggleRepeat"
          @toggle-shuffle="toggleShuffle"
        />
      </div>
    </div>

    <!-- Connection indicator -->
    <div
      v-if="!isConnected"
      class="connection-indicator"
    >
      <span class="pulse" />
    </div>
  </div>
</template>

<style scoped>
.widget-container {
  position: relative;
  width: var(--rays-size);
  height: var(--rays-size);
  display: flex;
  align-items: center;
  justify-content: center;
}

.widget-main {
  position: absolute;
  width: var(--cover-size);
  height: var(--cover-size);
  border-radius: 50%;
  overflow: hidden;
  background: var(--color-background);
  box-shadow:
    0 0 30px rgba(0, 0, 0, 0.5),
    0 0 60px rgba(0, 212, 170, 0.15),
    inset 0 0 30px rgba(0, 0, 0, 0.3);
}

.widget-content {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 20px;
  background: linear-gradient(
    180deg,
    rgba(0, 0, 0, 0.3) 0%,
    rgba(0, 0, 0, 0.6) 50%,
    rgba(0, 0, 0, 0.8) 100%
  );
}

.connection-indicator {
  position: absolute;
  bottom: 10px;
  left: 50%;
  transform: translateX(-50%);
}

.pulse {
  display: block;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--color-primary);
  animation: pulse 1.5s ease-in-out infinite;
}
</style>
