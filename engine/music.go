package engine

import (
	"errors"
	"github.com/banthar/Go-SDL/mixer"
	"path"
)

var music *mixer.Music

func initMusic() error {
	//defaults for now
	status := mixer.OpenAudio(mixer.DEFAULT_FREQUENCY, mixer.DEFAULT_FORMAT,
		mixer.DEFAULT_CHANNELS, 4096)

	if status != 0 {
		return errors.New("Error initializing music mixer from SDL.")
	}

	return nil
}

//PlayMusicFile plays the passed in file.  Filetype support is
// determined by SDL_mixer.  If file isn't an absolute path to a file,
// it'll look for the file in the engine's data directory.
// fadeIn is number of miliseconds to spend fading in
// 0 is no fade
func MusicPlayFile(file string, loop bool, fadeIn int) {
	//free the last music played
	// we'll see if this causes uncessary lag before play or
	// a waste of resources while no music is playing
	// If so we'll look to clean up music occasionally in the main loop
	if music != nil {
		music.Free()
	}
	if !path.IsAbs(file) {
		file = path.Join(dataDir, file)
	}
	music = mixer.LoadMUS(file)

	var loops = 1

	if loop {
		loops = -1
	}

	if fadeIn == 0 {
		music.PlayMusic(loops)
	} else {
		music.FadeInMusic(loops, fadeIn)
	}
}

func MusicFaceOut(ms int) {
	mixer.FadeOutMusic(ms)
}
func MusicSetVolume(volume int) {
	mixer.VolumeMusic(volume)
}

func MusicPause() {
	mixer.PauseMusic()
}

func ResumeMusic() {
	mixer.ResumeMusic()
}

func StopMusic() {
	mixer.HaltMusic()
}
