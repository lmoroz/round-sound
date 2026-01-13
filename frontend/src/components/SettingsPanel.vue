<script setup lang="ts">
import {
  computed,
  nextTick,
  onMounted,
  onUnmounted,
  ref,
} from 'vue'
import {
  AlertTriangle,
  ExternalLink,
  Settings,
  X,
} from 'lucide-vue-next'
import { useSettings } from '@/composables/useSettings'
import { FFT_SIZE_OPTIONS } from '@/types/settings'
import {
  ChangeWNPPort,
  GetWNPPort,
  IsAutorunEnabled,
  IsWNPConnected,
  SetAutorun,
} from '../../wailsjs/go/app/App'
import { EventsOff, EventsOn } from '../../wailsjs/runtime/runtime'

const { audioSettings, colorScheme, wnpSettings, updateAudioSettings, updatePrimaryColor, updateWNPSettings, resetToDefaults } = useSettings()

const isOpen = ref(false)
const primaryColorInput = ref(colorScheme.value.primary)
const autorunEnabled = ref(false)
const wnpPortInput = ref(wnpSettings.value.port)
const wnpConnected = ref(false)
const wnpPortError = ref('')
const showCustomAdapterHint = ref(false)
const wnpSectionRef = ref<HTMLElement | null>(null)

const fftSizeLabel = computed(() => {
  const size = audioSettings.value.fftSize
  if (size < 2048) return 'Быстро (низкая точность)'
  if (size === 2048) return 'Сбалансировано'
  if (size === 4096) return 'Точно'
  return 'Очень точно (медленнее)'
})

onMounted(async () => {
  try {
    autorunEnabled.value = await IsAutorunEnabled()
    wnpConnected.value = await IsWNPConnected()
    wnpPortInput.value = await GetWNPPort()
  }
  catch (error) {
    console.error('[Settings] Failed to load initial state:', error)
  }

  // Listen for WNP port busy event — auto-open settings
  EventsOn('wnp:port_busy', (data: { port: number; message: string }) => {
    console.log('[Settings] WNP port busy:', data)
    showCustomAdapterHint.value = true
    wnpPortError.value = data.message
    wnpConnected.value = false
    isOpen.value = true

    // Scroll to WNP section after modal opens
    nextTick(() => {
      setTimeout(() => {
        wnpSectionRef.value?.scrollIntoView({ behavior: 'smooth', block: 'start' })
      }, 300)
    })
  })

  EventsOn('wnp:port_changed', (port: number) => {
    console.log('[Settings] WNP port changed:', port)
    wnpPortInput.value = port
    wnpConnected.value = true
    wnpPortError.value = ''
    showCustomAdapterHint.value = false
    updateWNPSettings({ port, showCustomAdapterHint: false })
  })

  EventsOn('wnp:port_error', (data: { port: number; message: string }) => {
    console.log('[Settings] WNP port error:', data)
    wnpPortError.value = data.message
    wnpConnected.value = false
  })
})

onUnmounted(() => {
  EventsOff('wnp:port_busy')
  EventsOff('wnp:port_changed')
  EventsOff('wnp:port_error')
})

function toggleModal() {
  isOpen.value = !isOpen.value
}

function handlePrimaryColorChange() {
  updatePrimaryColor(primaryColorInput.value)
}

async function handleAutorunToggle() {
  try {
    await SetAutorun(autorunEnabled.value)
  }
  catch (error) {
    console.error('[Settings] Failed to set autorun:', error)
    autorunEnabled.value = !autorunEnabled.value
  }
}

async function handleWNPPortChange() {
  const port = wnpPortInput.value
  if (port < 1024 || port > 65535) {
    wnpPortError.value = 'Порт должен быть от 1024 до 65535'
    return
  }

  wnpPortError.value = ''

  try {
    await ChangeWNPPort(port)
    updateWNPSettings({ port })
  }
  catch (error) {
    console.error('[Settings] Failed to change WNP port:', error)
    wnpPortError.value = 'Не удалось сменить порт'
  }
}

function handleReset() {
  if (confirm('Сбросить все настройки к значениям по умолчанию?')) {
    resetToDefaults()
    primaryColorInput.value = colorScheme.value.primary
    wnpPortInput.value = wnpSettings.value.port
  }
}

function dismissCustomAdapterHint() {
  showCustomAdapterHint.value = false
  updateWNPSettings({ showCustomAdapterHint: false })
}
</script>

<template>
  <div class="settings-wrapper">
    <!-- Settings button -->
    <button
      class="settings-button"
      :class="{ active: isOpen }"
      title="Настройки"
      @click="toggleModal"
    >
      <Settings :size="20" />
    </button>

    <!-- Modal backdrop and panel -->
    <Transition name="modal-fade">
      <div
        v-if="isOpen"
        class="modal-backdrop"
        @click.self="toggleModal"
      >
        <div class="settings-modal">
          <div class="modal-header">
            <h2>Настройки</h2>
            <button
              class="close-button"
              @click="toggleModal"
            >
              <X :size="24" />
            </button>
          </div>

          <div class="modal-content">
            <!-- Audio Settings -->
            <section class="settings-section">
              <h3>Аудио анализ</h3>

              <div class="setting-item">
                <label for="fft-size">
                  Размер FFT
                  <span class="setting-hint">{{ fftSizeLabel }}</span>
                </label>
                <select
                  id="fft-size"
                  :value="audioSettings.fftSize"
                  @change="e => updateAudioSettings({ fftSize: Number((e.target as HTMLSelectElement).value) })"
                >
                  <option
                    v-for="size in FFT_SIZE_OPTIONS"
                    :key="size"
                    :value="size"
                  >
                    {{ size }}
                  </option>
                </select>
              </div>

              <div class="setting-item">
                <label for="freq-min">
                  Минимальная частота (Hz)
                </label>
                <input
                  id="freq-min"
                  max="1000"
                  min="10"
                  step="10"
                  type="number"
                  :value="audioSettings.freqMin"
                  @input="e => updateAudioSettings({ freqMin: Number((e.target as HTMLInputElement).value) })"
                >
              </div>

              <div class="setting-item">
                <label for="freq-max">
                  Максимальная частота (Hz)
                </label>
                <input
                  id="freq-max"
                  max="24000"
                  min="1000"
                  step="1000"
                  type="number"
                  :value="audioSettings.freqMax"
                  @input="e => updateAudioSettings({ freqMax: Number((e.target as HTMLInputElement).value) })"
                >
              </div>
            </section>

            <!-- Color Settings -->
            <section class="settings-section">
              <h3>Цветовая схема</h3>

              <div class="setting-item">
                <label for="primary-color">
                  Основной цвет
                  <span class="setting-hint">Определяет всю цветовую палитру</span>
                </label>
                <div class="color-picker-wrapper">
                  <input
                    id="primary-color"
                    v-model="primaryColorInput"
                    class="color-input"
                    type="color"
                    @input="handlePrimaryColorChange"
                  >
                  <input
                    v-model="primaryColorInput"
                    class="color-text"
                    maxlength="7"
                    type="text"
                    @change="handlePrimaryColorChange"
                  >
                </div>
              </div>

              <div class="color-preview">
                <div
                  class="color-swatch"
                  :style="{ backgroundColor: colorScheme.primary }"
                >
                  <span>Primary</span>
                </div>
                <div
                  class="color-swatch"
                  :style="{ backgroundColor: colorScheme.secondary }"
                >
                  <span>Secondary</span>
                </div>
                <div
                  class="color-swatch"
                  :style="{ backgroundColor: colorScheme.accent }"
                >
                  <span>Accent</span>
                </div>
              </div>
            </section>

            <!-- WebNowPlaying Settings -->
            <section
              ref="wnpSectionRef"
              class="settings-section"
            >
              <h3>WebNowPlaying</h3>

              <!-- Custom Adapter Hint (показывается при конфликте портов) -->
              <div
                v-if="showCustomAdapterHint"
                class="alert-box warning"
              >
                <div class="alert-header">
                  <AlertTriangle :size="18" />
                  <span>Порт занят</span>
                  <button
                    class="alert-close"
                    @click="dismissCustomAdapterHint"
                  >
                    <X :size="16" />
                  </button>
                </div>
                <p>
                  Стандартный порт 8974 занят (скорее всего, Rainmeter WebNowPlaying.dll).
                </p>
                <p>
                  <strong>Решение:</strong> Укажите другой порт (например, 9000) и добавьте его как
                  <strong>Custom Adapter</strong> в настройках расширения WebNowPlaying в браузере.
                </p>
                <a
                  class="hint-link"
                  href="https://wnp.keifufu.dev"
                  target="_blank"
                >
                  <ExternalLink :size="14" />
                  Документация WebNowPlaying
                </a>
              </div>

              <div class="setting-item">
                <label for="wnp-port">
                  Порт подключения
                  <span class="setting-hint">
                    Round Sound слушает подключения от расширения WebNowPlaying на этом порту
                  </span>
                </label>
                <div class="port-input-wrapper">
                  <input
                    id="wnp-port"
                    v-model.number="wnpPortInput"
                    max="65535"
                    min="1024"
                    type="number"
                    @keydown.enter="handleWNPPortChange"
                  >
                  <button
                    class="port-apply-button"
                    @click="handleWNPPortChange"
                  >
                    Применить
                  </button>
                </div>
                <div
                  v-if="wnpPortError"
                  class="setting-error"
                >
                  {{ wnpPortError }}
                </div>
              </div>

              <div class="setting-item">
                <div class="status-row">
                  <span class="status-label">Статус:</span>
                  <span
                    class="status-badge"
                    :class="{ connected: wnpConnected, disconnected: !wnpConnected }"
                  >
                    {{ wnpConnected ? 'Сервер запущен' : 'Сервер не запущен' }}
                  </span>
                </div>
              </div>

              <div class="info-box">
                <p>
                  <strong>Как подключить:</strong>
                </p>
                <ol>
                  <li>
                    Установите расширение <a
                      href="https://chrome.google.com/webstore/detail/webnowplaying/jfakgfcdgpghbbefmdfjkbdlibjgnbli"
                      target="_blank"
                    >WebNowPlaying</a> в браузере
                  </li>
                  <li>Откройте настройки расширения</li>
                  <li>Нажмите "Add custom adapter"</li>
                  <li>Укажите порт: <strong>{{ wnpPortInput }}</strong></li>
                  <li>Включите адаптер</li>
                </ol>
              </div>
            </section>

            <!-- System Settings -->
            <section class="settings-section">
              <h3>Система</h3>

              <div class="setting-item">
                <label class="checkbox-label">
                  <input
                    v-model="autorunEnabled"
                    type="checkbox"
                    @change="handleAutorunToggle"
                  >
                  <span>Запускать при старте Windows</span>
                </label>
              </div>
            </section>
          </div>

          <div class="modal-footer">
            <button
              class="reset-button"
              @click="handleReset"
            >
              Сбросить к дефолту
            </button>
            <button
              class="save-button"
              @click="toggleModal"
            >
              Закрыть
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.settings-wrapper {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 1000;
}

.settings-button {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: var(--color-button);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: var(--color-text);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  backdrop-filter: blur(10px);
}

.settings-button:hover {
  background: var(--color-button-hover);
  transform: scale(1.05);
}

.settings-button.active {
  background: var(--color-primary);
  color: white;
}

/* Modal */
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: transparent;
  backdrop-filter: blur(12px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
  padding: 0 10px;
  border-radius: 20px;
  --wails-draggable: drag;
}

.settings-modal {
  width: 100%;
  background: var(--color-background);
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.15);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.6);
  display: flex;
  flex-direction: column;
  max-height: 95vh;
  overflow: hidden;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  --wails-draggable: drag;
}

.modal-header h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: var(--color-text);
}

.close-button {
  width: 30px;
  height: 30px;
  border-radius: 100%;
  background: transparent;
  border: none;
  color: var(--color-text);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  --wails-draggable: no-drag;
}

.close-button:hover {
  background: rgba(255, 255, 255, 0.1);
}

.modal-content {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  --wails-draggable: no-drag;
}

.settings-section {
  margin-bottom: 20px;
}

.settings-section:last-child {
  margin-bottom: 0;
}

.settings-section h3 {
  margin: 0 0 16px 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--color-accent);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.setting-item {
  margin-bottom: 20px;
}

.setting-item:last-child {
  margin-bottom: 0;
}

.setting-item label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text);
}

.setting-hint {
  display: block;
  font-size: 12px;
  font-weight: 400;
  color: var(--color-text-secondary);
  margin-top: 4px;
}

.setting-item input[type="number"],
.setting-item input[type="text"],
.setting-item select {
  width: 100%;
  padding: 10px 14px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  color: var(--color-text);
  font-size: 14px;
  outline: none;
  transition: all 0.2s ease;
}

.setting-item select option {
  background: rgba(20, 20, 30, 0.98);
  color: var(--color-text);
  padding: 8px;
}

.setting-item input[type="number"]:hover,
.setting-item input[type="text"]:hover,
.setting-item select:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.2);
}

.setting-item input[type="number"]:focus,
.setting-item input[type="text"]:focus,
.setting-item select:focus {
  background: rgba(255, 255, 255, 0.1);
  border-color: var(--color-primary);
}

.color-picker-wrapper {
  display: flex;
  gap: 12px;
  align-items: center;
}

.color-input {
  width: 60px;
  height: 44px;
  padding: 4px;
  background: transparent;
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  cursor: pointer;
}

.color-input::-webkit-color-swatch-wrapper {
  padding: 4px;
}

.color-input::-webkit-color-swatch {
  border-radius: 4px;
  border: none;
}

.color-text {
  flex: 1;
  font-family: 'Consolas', 'Monaco', monospace;
  text-transform: uppercase;
}

.color-preview {
  display: flex;
  gap: 12px;
  margin-top: 16px;
}

.color-swatch {
  flex: 1;
  height: 60px;
  border-radius: 8px;
  display: flex;
  align-items: flex-end;
  justify-content: center;
  padding: 10px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.color-swatch span {
  font-size: 11px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.9);
  text-shadow: 0 1px 3px rgba(0, 0, 0, 0.5);
}

/* Checkbox */
.checkbox-label {
  display: flex!important;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  user-select: none;
}

.checkbox-label input[type="checkbox"] {
  appearance: none;
  width: 48px;
  height: 26px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 13px;
  position: relative;
  cursor: pointer;
  transition: all 0.3s ease;
  flex-shrink: 0;
}

.checkbox-label input[type="checkbox"]::before {
  content: '';
  position: absolute;
  width: 20px;
  height: 20px;
  background: white;
  border-radius: 50%;
  top: 2px;
  left: 2px;
  transition: all 0.3s ease;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.checkbox-label input[type="checkbox"]:checked {
  background: var(--color-primary);
  border-color: var(--color-primary);
}

.checkbox-label input[type="checkbox"]:checked::before {
  left: 24px;
}

.checkbox-label input[type="checkbox"]:hover {
  background: rgba(255, 255, 255, 0.15);
}

.checkbox-label input[type="checkbox"]:checked:hover {
  background: var(--color-secondary);
}

.checkbox-label span {
  font-size: 14px;
  color: var(--color-text);
}

.modal-footer {
  display: flex;
  gap: 12px;
  padding: 12px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.reset-button,
.save-button {
  flex: 1;
  padding: 16px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s ease;
}

.reset-button {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: var(--color-text);
}

.reset-button:hover {
  background: rgba(255, 100, 100, 0.2);
  border-color: rgba(255, 100, 100, 0.5);
  color: #ff6464;
}

.save-button {
  background: var(--color-primary);
  border: 1px solid var(--color-primary);
  color: white;
}

.save-button:hover {
  background: var(--color-secondary);
  border-color: var(--color-secondary);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px var(--color-primary-glow);
}

/* Transitions */
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: all 0.3s ease;
}

.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}

.modal-fade-enter-from .settings-modal,
.modal-fade-leave-to .settings-modal {
  transform: scale(0.9);
  opacity: 0;
}

/* Scrollbar */
.modal-content::-webkit-scrollbar {
  width: 8px;
}

.modal-content::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 4px;
}

.modal-content::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 4px;
}

.modal-content::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.3);
}

/* Alert Box */
.alert-box {
  padding: 16px;
  border-radius: 10px;
  margin-bottom: 20px;
}

.alert-box.warning {
  background: rgba(255, 180, 50, 0.15);
  border: 1px solid rgba(255, 180, 50, 0.4);
}

.alert-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
  color: #ffb432;
  font-weight: 600;
}

.alert-close {
  margin-left: auto;
  background: transparent;
  border: none;
  color: rgba(255, 255, 255, 0.5);
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.alert-close:hover {
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.8);
}

.alert-box p {
  margin: 8px 0;
  font-size: 13px;
  color: var(--color-text-secondary);
  line-height: 1.5;
}

.hint-link {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: var(--color-accent);
  font-size: 13px;
  text-decoration: none;
  margin-top: 8px;
}

.hint-link:hover {
  text-decoration: underline;
}

/* Port Input */
.port-input-wrapper {
  display: flex;
  gap: 10px;
}

.port-input-wrapper input {
  flex: 1;
  padding: 10px 14px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  color: var(--color-text);
  font-size: 14px;
  outline: none;
  transition: all 0.2s ease;
}

.port-input-wrapper input:focus {
  background: rgba(255, 255, 255, 0.1);
  border-color: var(--color-primary);
}

.port-apply-button {
  padding: 10px 18px;
  background: var(--color-primary);
  border: 1px solid var(--color-primary);
  border-radius: 8px;
  color: white;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.port-apply-button:hover {
  background: var(--color-secondary);
  border-color: var(--color-secondary);
}

/* Setting Error */
.setting-error {
  color: #ff6464;
  font-size: 12px;
  margin-top: 6px;
}

/* Status Row */
.status-row {
  display: flex;
  align-items: center;
  gap: 10px;
}

.status-label {
  font-size: 14px;
  color: var(--color-text-secondary);
}

.status-badge {
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
}

.status-badge.connected {
  background: rgba(80, 200, 120, 0.2);
  color: #50c878;
}

.status-badge.disconnected {
  background: rgba(255, 100, 100, 0.2);
  color: #ff6464;
}

/* Info Box */
.info-box {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 10px;
  padding: 16px;
  margin-top: 16px;
}

.info-box p {
  margin: 0 0 10px 0;
  font-size: 13px;
  color: var(--color-text);
}

.info-box ol {
  margin: 0;
  padding-left: 20px;
}

.info-box li {
  font-size: 13px;
  color: var(--color-text-secondary);
  margin-bottom: 6px;
  line-height: 1.5;
}

.info-box a {
  color: var(--color-accent);
  text-decoration: none;
}

.info-box a:hover {
  text-decoration: underline;
}
</style>
