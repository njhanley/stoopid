package bot

import (
	dg "github.com/bwmarrin/discordgo"
	"layeh.com/gopus"
)

const (
	channels     = 2
	frameSize    = 960
	sampleRate   = 48000
	pcmFrameSize = frameSize * channels
)

func Send(v *dg.VoiceConnection, pcm <-chan []int16) error {
	encoder, err := gopus.NewEncoder(sampleRate, channels, gopus.Audio)
	if err != nil {
		return err
	}

	v.Speaking(true)
	defer v.Speaking(false)

	for frame := range pcm {
		opus, err := encoder.Encode(frame, frameSize, pcmFrameSize*2)
		if err != nil {
			return err
		}
		v.OpusSend <- opus
	}

	return nil
}
