// Player state modes
export enum StateMode {
  Playing = 0,
  Paused = 1,
  Stopped = 2,
}

// Repeat modes
export enum RepeatMode {
  None = 1,
  All = 2,
  One = 4,
}

// Rating system types
export enum RatingSystem {
  None = 0,
  Like = 1,
  LikeDislike = 2,
  Scale = 3,
}

// Player state from WebNowPlaying
export interface Player {
  id: number;
  name: string;
  title: string;
  artist: string;
  album: string;
  cover: string;
  state: StateMode;
  position: number;  // seconds
  duration: number;  // seconds
  volume: number;    // 0-100
  rating: number;    // 0=none, 1=dislike, 5=like
  repeat: RepeatMode;
  shuffle: boolean;
  ratingSystem: RatingSystem;
  availableRepeat: number;
  canSetState: boolean;
  canSkipPrevious: boolean;
  canSkipNext: boolean;
  canSetPosition: boolean;
  canSetVolume: boolean;
  canSetRating: boolean;
  canSetRepeat: boolean;
  canSetShuffle: boolean;
  createdAt: number;
  updatedAt: number;
  activeAt: number;
}

// Default empty player
export const defaultPlayer: Player = {
  id: 0,
  name: '',
  title: '',
  artist: '',
  album: '',
  cover: '',
  state: StateMode.Stopped,
  position: 0,
  duration: 0,
  volume: 100,
  rating: 0,
  repeat: RepeatMode.None,
  shuffle: false,
  ratingSystem: RatingSystem.None,
  availableRepeat: RepeatMode.None,
  canSetState: false,
  canSkipPrevious: false,
  canSkipNext: false,
  canSetPosition: false,
  canSetVolume: false,
  canSetRating: false,
  canSetRepeat: false,
  canSetShuffle: false,
  createdAt: 0,
  updatedAt: 0,
  activeAt: 0,
}
