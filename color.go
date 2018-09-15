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
  "image/color"
)

func convertRGB(color color.Color) (float64, float64, float64) {
  r, g, b, _ := color.RGBA()

  return float64(r>>8) / float64(256),
    float64(g>>8) / float64(256),
    float64(b>>8) / float64(256)
}

func convertYUV(color color.Color) (byte, byte, byte) {
  r, g, b, _ := color.RGBA()

  rf := float64(int(r) >> 8)
  gf := float64(int(g) >> 8)
  bf := float64(int(b) >> 8)

  y := clamp(16.0 + (.003906 * ((65.738 * rf) + (129.057 * gf) + (25.064 * bf))))
  u := clamp(128.0 + (.003906 * ((-37.945 * rf) + (-74.494 * gf) + (112.439 * bf))))
  v := clamp(128.0 + (.003906 * ((112.439 * rf) + (-94.154 * gf) + (-18.285 * bf))))

  return byte(y), byte(u), byte(v)
}

func clamp(input float64) float64 {
  if input < 0 {
    return 0
  }
  if input > 255 {
    return 255
  }
  return input
}
