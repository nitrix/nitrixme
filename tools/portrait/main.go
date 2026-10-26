// Command portrait rebuilds the frames of the portrait's orbit from the clip, with `go generate` (see main.go in the project folder):
// static/nitrixme-orbit.mp4 -> static/nitrixme-orbit.bin. It needs ffmpeg and runs from the project folder.
//
// The .bin is one file the page loads in a single request, with the frames already cut into JPEGs. Layout, little endian:
// "AGRING1\n", frame count, width, height (uint32 each), the length of every JPEG (uint32 each), then the JPEGs one after the other.
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
)

const (
	size    = 448 // the side of a frame: twice the 220 px the portrait is shown at, about what the still is
	quality = "4" // ffmpeg's -q:v for the frames
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "portrait:", err)
		os.Exit(1)
	}
}

func ffmpeg(args ...string) error {
	if out, err := exec.Command("ffmpeg", append([]string{"-y", "-loglevel", "error"}, args...)...).CombinedOutput(); err != nil {
		return fmt.Errorf("ffmpeg: %v: %s", err, out)
	}
	return nil
}

func run() error {
	tmp, err := os.MkdirTemp("", "portrait-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	// -fps_mode passthrough: every frame of the clip once, none repeated or dropped to fit a rate
	if err := ffmpeg("-i", "static/nitrixme-orbit.mp4", "-an", "-fps_mode", "passthrough",
		"-vf", fmt.Sprintf("scale=%d:%d:flags=lanczos", size, size), "-q:v", quality, filepath.Join(tmp, "f%04d.jpg")); err != nil {
		return err
	}
	names, err := filepath.Glob(filepath.Join(tmp, "f*.jpg"))
	if err != nil || len(names) == 0 {
		return fmt.Errorf("no frames in the clip")
	}
	sort.Strings(names)
	frames := make([][]byte, len(names))
	for i, n := range names {
		if frames[i], err = os.ReadFile(n); err != nil {
			return err
		}
	}

	var b bytes.Buffer
	b.WriteString("AGRING1\n")
	for _, v := range []int{len(frames), size, size} {
		binary.Write(&b, binary.LittleEndian, uint32(v))
	}
	for _, f := range frames {
		binary.Write(&b, binary.LittleEndian, uint32(len(f)))
	}
	for _, f := range frames {
		b.Write(f)
	}
	if err := os.WriteFile("static/nitrixme-orbit.bin", b.Bytes(), 0o644); err != nil {
		return err
	}
	fmt.Printf("portrait: %d frames, %d bytes\n", len(frames), b.Len())
	return nil
}
