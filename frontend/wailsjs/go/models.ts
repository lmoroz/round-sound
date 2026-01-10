export namespace media {
	
	export class Player {
	    id: number;
	    name: string;
	    title: string;
	    artist: string;
	    album: string;
	    cover: string;
	    state: number;
	    position: number;
	    duration: number;
	    volume: number;
	    rating: number;
	    repeat: number;
	    shuffle: boolean;
	    ratingSystem: number;
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
	
	    static createFrom(source: any = {}) {
	        return new Player(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.title = source["title"];
	        this.artist = source["artist"];
	        this.album = source["album"];
	        this.cover = source["cover"];
	        this.state = source["state"];
	        this.position = source["position"];
	        this.duration = source["duration"];
	        this.volume = source["volume"];
	        this.rating = source["rating"];
	        this.repeat = source["repeat"];
	        this.shuffle = source["shuffle"];
	        this.ratingSystem = source["ratingSystem"];
	        this.availableRepeat = source["availableRepeat"];
	        this.canSetState = source["canSetState"];
	        this.canSkipPrevious = source["canSkipPrevious"];
	        this.canSkipNext = source["canSkipNext"];
	        this.canSetPosition = source["canSetPosition"];
	        this.canSetVolume = source["canSetVolume"];
	        this.canSetRating = source["canSetRating"];
	        this.canSetRepeat = source["canSetRepeat"];
	        this.canSetShuffle = source["canSetShuffle"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	        this.activeAt = source["activeAt"];
	    }
	}

}

