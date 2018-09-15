/*
 * Copyright 2018 Johannes Donath <johannesd@torchmind.com>
 * and other copyright owners as documented in the project's IP log.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package sstv

import (
  "github.com/go-audio/audio"
  "image"
)

// represents an arbitray SSTV encoder
type Encoder interface {
  // retrieves the vis which is to be encoded within the handshake
  Vis() uint8
  // retrieves the standard resolution for this transmission format
  // while this width and height is typically not required for successful encoding, it is
  // recommended to stick to them as most decoders will expect the standard sizes
  Resolution() image.Rectangle
  // encodes a given image into an SSTV audio signal represented by an array of raw PCM samples
  Encode(image image.Image) *audio.FloatBuffer
}
