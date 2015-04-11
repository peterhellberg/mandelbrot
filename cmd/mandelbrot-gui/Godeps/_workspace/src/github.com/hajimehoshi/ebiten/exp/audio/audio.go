// Copyright 2015 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package audio

import (
	"github.com/hajimehoshi/ebiten/exp/audio/internal"
)

// SampleRate returns the sampling frequency (e.g. 44100).
func SampleRate() int {
	return internal.SampleRate
}

// MaxChannel is a max number of channels.
var MaxChannel = internal.MaxChannel

// Play appends the given data to the given channel.
//
// channel must be -1 or a channel index. If channel is -1, an empty channel is automatically selected.
// If the channel is not empty, this function does nothing and returns false. This returns true otherwise.
//
// This function is useful to play SE or a note of PCM synthesis immediately.
func Play(channel int, l []int16, r []int16) bool {
	return internal.Play(channel, l, r)
}

// Queue queues the given data to the given channel.
// The given data is queued to the end of the buffer and not played immediately.
//
// channel must be a channel index. You can't give -1 to channel.
//
// This function is useful to play streaming data.
func Queue(channel int, l []int16, r []int16) {
	internal.Queue(channel, l, r)
}

// IsPlaying returns a boolean value which indicates if the channel buffer has data to play.
func IsPlaying(channel int) bool {
	return internal.IsPlaying(channel)
}
