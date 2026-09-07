package ui

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
)

func writeClipboard(text string) error {
	osc := fmt.Sprintf("\x1b]52;c;%s\x07", base64.StdEncoding.EncodeToString([]byte(text)))
	_, _ = os.Stdout.WriteString(osc)

	if os.Getenv("WAYLAND_DISPLAY") != "" {
		if path, err := exec.LookPath("wl-copy"); err == nil {
			cmd := exec.Command(path)
			in, err := cmd.StdinPipe()
			if err == nil {
				if err := cmd.Start(); err == nil {
					_, _ = in.Write([]byte(text))
					_ = in.Close()
					_ = cmd.Wait()
					return nil
				}
			}
		}
	}

	if os.Getenv("DISPLAY") != "" {
		if path, err := exec.LookPath("xclip"); err == nil {
			cmd := exec.Command(path, "-selection", "clipboard")
			in, err := cmd.StdinPipe()
			if err == nil {
				if err := cmd.Start(); err == nil {
					_, _ = in.Write([]byte(text))
					_ = in.Close()
					_ = cmd.Wait()
					return nil
				}
			}
		}
		if path, err := exec.LookPath("xsel"); err == nil {
			cmd := exec.Command(path, "--clipboard", "--input")
			in, err := cmd.StdinPipe()
			if err == nil {
				if err := cmd.Start(); err == nil {
					_, _ = in.Write([]byte(text))
					_ = in.Close()
					_ = cmd.Wait()
					return nil
				}
			}
		}
	}

	if path, err := exec.LookPath("pbcopy"); err == nil {
		cmd := exec.Command(path)
		in, err := cmd.StdinPipe()
		if err == nil {
			if err := cmd.Start(); err == nil {
				_, _ = in.Write([]byte(text))
				_ = in.Close()
				_ = cmd.Wait()
				return nil
			}
		}
	}

	if path, err := exec.LookPath("clip.exe"); err == nil {
		cmd := exec.Command(path)
		in, err := cmd.StdinPipe()
		if err == nil {
			if err := cmd.Start(); err == nil {
				_, _ = in.Write([]byte(text))
				_ = in.Close()
				_ = cmd.Wait()
				return nil
			}
		}
	}

	return nil
}
