import type { ColorScheme } from '@/types/settings'

/**
 * Конвертирует HEX в RGB
 */
export function hexToRgb(hex: string): { r: number; g: number; b: number } | null {
  const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex)
  return result && result[1] && result[2] && result[3]
    ? {
        r: Number.parseInt(result[1], 16),
        g: Number.parseInt(result[2], 16),
        b: Number.parseInt(result[3], 16),
      }
    : null
}

/**
 * Конвертирует RGB в HEX
 */
export function rgbToHex(r: number, g: number, b: number): string {
  return `#${((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1)}`
}

/**
 * Adjust color lightness (HSL)
 */
export function adjustLightness(hex: string, percent: number): string {
  const rgb = hexToRgb(hex)
  if (!rgb) return hex

  const { r, g, b } = rgb
  const hsl = rgbToHsl(r, g, b)

  hsl.l = Math.max(0, Math.min(100, hsl.l + percent))

  const newRgb = hslToRgb(hsl.h, hsl.s, hsl.l)
  return rgbToHex(newRgb.r, newRgb.g, newRgb.b)
}

/**
 * Convert RGB to HSL
 */
function rgbToHsl(r: number, g: number, b: number): { h: number; s: number; l: number } {
  r /= 255
  g /= 255
  b /= 255

  const max = Math.max(r, g, b)
  const min = Math.min(r, g, b)
  let h = 0
  let s = 0
  const l = (max + min) / 2

  if (max !== min) {
    const d = max - min
    s = l > 0.5 ? d / (2 - max - min) : d / (max + min)

    switch (max) {
    case r:
      h = ((g - b) / d + (g < b ? 6 : 0)) / 6
      break
    case g:
      h = ((b - r) / d + 2) / 6
      break
    case b:
      h = ((r - g) / d + 4) / 6
      break
    }
  }

  return {
    h: Math.round(h * 360),
    s: Math.round(s * 100),
    l: Math.round(l * 100),
  }
}

/**
 * Convert HSL to RGB
 */
function hslToRgb(h: number, s: number, l: number): { r: number; g: number; b: number } {
  h /= 360
  s /= 100
  l /= 100

  let r: number
  let g: number
  let b: number

  if (s === 0) {
    r = g = b = l
  }
  else {
    const hue2rgb = (p: number, q: number, t: number) => {
      if (t < 0) t += 1
      if (t > 1) t -= 1
      if (t < 1 / 6) return p + (q - p) * 6 * t
      if (t < 1 / 2) return q
      if (t < 2 / 3) return p + (q - p) * (2 / 3 - t) * 6
      return p
    }

    const q = l < 0.5 ? l * (1 + s) : l + s - l * s
    const p = 2 * l - q

    r = hue2rgb(p, q, h + 1 / 3)
    g = hue2rgb(p, q, h)
    b = hue2rgb(p, q, h - 1 / 3)
  }

  return {
    r: Math.round(r * 255),
    g: Math.round(g * 255),
    b: Math.round(b * 255),
  }
}

/**
 * Генерирует цветовую схему из основного цвета
 */
export function generateColorScheme(primaryHex: string): ColorScheme {
  const rgb = hexToRgb(primaryHex)
  if (!rgb) {
    return {
      primary: primaryHex,
      primaryGlow: 'rgba(255, 140, 66, 0.6)',
      secondary: adjustLightness(primaryHex, -10),
      accent: adjustLightness(primaryHex, 10),
    }
  }

  const { r, g, b } = rgb

  return {
    primary: primaryHex,
    primaryGlow: `rgba(${r}, ${g}, ${b}, 0.6)`,
    secondary: adjustLightness(primaryHex, -10),
    accent: adjustLightness(primaryHex, 15),
  }
}

/**
 * Генерирует градиент для лучей из цветовой схемы
 * @param hasSound - если true, лучи будут цветными (реагируют на звук)
 */
export function generateRayGradient(
  ctx: CanvasRenderingContext2D,
  startX: number,
  startY: number,
  endX: number,
  endY: number,
  colors: ColorScheme,
  hasSound: boolean,
): CanvasGradient {
  const gradient = ctx.createLinearGradient(startX, startY, endX, endY)

  if (hasSound) {
    const rgb = hexToRgb(colors.primary)
    const rgbAccent = hexToRgb(colors.accent)
    const rgbSecondary = hexToRgb(colors.secondary)

    if (rgb && rgbAccent && rgbSecondary) {
      gradient.addColorStop(0, `rgba(${rgb.r}, ${rgb.g}, ${rgb.b}, 0.8)`)
      gradient.addColorStop(0.5, `rgba(${rgbAccent.r}, ${rgbAccent.g}, ${rgbAccent.b}, 0.6)`)
      gradient.addColorStop(1, `rgba(${rgbSecondary.r}, ${rgbSecondary.g}, ${rgbSecondary.b}, 0.2)`)
    }
    else {
      gradient.addColorStop(0, 'rgba(255, 140, 66, 0.8)')
      gradient.addColorStop(0.5, 'rgba(255, 170, 102, 0.6)')
      gradient.addColorStop(1, 'rgba(255, 107, 53, 0.2)')
    }
  }
  else {
    gradient.addColorStop(0, 'rgba(100, 100, 100, 0.3)')
    gradient.addColorStop(1, 'rgba(100, 100, 100, 0.1)')
  }

  return gradient
}
