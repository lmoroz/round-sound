export interface AudioSettings {
  fftSize: number;
  freqMin: number;
  freqMax: number;
}

export interface ColorScheme {
  primary: string;
  primaryGlow: string;
  secondary: string;
  accent: string;
}

export interface AppSettings {
  audio: AudioSettings;
  colors: ColorScheme;
}

export const DEFAULT_SETTINGS: AppSettings = {
  audio: {
    fftSize: 2048,
    freqMin: 20,
    freqMax: 20000,
  },
  colors: {
    primary: '#ff8c42',
    primaryGlow: 'rgba(255, 140, 66, 0.6)',
    secondary: '#ff6b35',
    accent: '#ffaa66',
  },
}

export const FFT_SIZE_OPTIONS = [1024, 2048, 4096, 8192] as const
export type FFTSize = typeof FFT_SIZE_OPTIONS[number]
