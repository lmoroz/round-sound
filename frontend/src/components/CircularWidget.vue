<script setup lang="ts">
import { computed, ref } from 'vue'

import { Volume, Volume1, Volume2, VolumeX } from 'lucide-vue-next'

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
  seek,
  setVolume,
} = useMediaPlayer()

const { levels } = useAudioLevels(64)

const progress = computed(() => {
  if (!player.value.duration) return 0
  return player.value.position / player.value.duration
})

const isPlaying = computed(() => player.value.state === StateMode.Playing)

// -- Volume Control --
const isVolumeVisible = ref(false)
let volumeTimer: ReturnType<typeof setTimeout> | null = null

const volumeIcon = computed(() => {
  const v = player.value.volume
  if (v === 0) return VolumeX
  if (v < 30) return Volume
  if (v < 70) return Volume1
  return Volume2
})

function handleWheel(e: WheelEvent) {
  console.log('isConnected.value = ', isConnected.value)
  // Only handle if connected
  if (!isConnected.value) return

  // Prevent default scroll behavior
  e.preventDefault()

  // Determine direction (scroll up = increase volume)
  // deltaY is negative when scrolling up
  const delta = e.deltaY > 0 ? -5 : 5

  const newVolume = Math.min(100, Math.max(0, player.value.volume + delta))

  if (newVolume !== player.value.volume) {
    setVolume(newVolume)
    showVolumeOverlay()
  }
}

function showVolumeOverlay() {
  isVolumeVisible.value = true
  if (volumeTimer) clearTimeout(volumeTimer)
  volumeTimer = setTimeout(() => {
    isVolumeVisible.value = false
  }, 1500)
}
</script>

<template>
  <div
    class="widget-container"
    @wheel="handleWheel"
  >
    <!-- Audio visualization rays (outermost layer) -->
    <AudioLevelsRays
      :is-playing="isPlaying"
      :levels="levels"
    />

    <!-- Progress ring -->
    <ProgressRing
      :duration="player.duration"
      :progress="progress"
      @seek="seek"
      @wheel.stop
    />

    <!-- Main circular content -->
    <div class="widget-main">
      <!-- Album cover background -->
      <AlbumCover :cover="player.cover" />

      <!-- Content overlay -->
      <div class="widget-content">
        <!-- Volume Overlay -->
        <Transition name="fade">
          <div
            v-if="isVolumeVisible"
            class="volume-overlay"
          >
            <component
              :is="volumeIcon"
              class="volume-icon"
              :size="48"
            />
            <span class="volume-text">{{ player.volume }}%</span>
          </div>
        </Transition>

        <!-- Normal Content (hide when volume is changing) -->
        <Transition name="fade">
          <div
            v-if="!isVolumeVisible"
            class="main-content-wrapper"
          >
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

            <div
              class="volume-handler-wrapper"
              @wheel="handleWheel"
            >
              <component
                :is="volumeIcon"
                class="volume-icon hover"
                :size="40"
              />
            </div>
          </div>
        </Transition>
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

.main-content-wrapper {
  padding-top: 43px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 100%;
}

.volume-overlay {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 10px;
  color: var(--color-primary);
  filter: drop-shadow(0 0 10px rgba(0, 0, 0, 0.5));
}

.volume-handler-wrapper {
  padding-top: 12px;
  display: flex;
  justify-content: center;
  align-items: center;
  position: relative;
  z-index: 100;
  cursor: pointer;
  opacity: 0;
}
.widget-container:hover .volume-handler-wrapper {
  opacity: 0.7;
}

.volume-icon {
  margin-bottom: 5px;
}

.volume-text {
  font-size: 24px;
  font-weight: bold;
  font-feature-settings: "tnum";
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

/* Transitions */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
