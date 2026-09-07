package ui

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"
)

var errNoClipboardBackend = errors.New("no supported clipboard utility available")

func runClipboardCommand(name string, args []string, text string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	if _, err := in.Write([]byte(text)); err != nil {
		_ = in.Close()
		_ = cmd.Wait()
		return err
	}
	if err := in.Close(); err != nil {
		_ = cmd.Wait()
		return err
	}

	return cmd.Wait()
}

func writeClipboard(text string) error {
	osc := fmt.Sprintf("\x1b]52;c;%s\x07", base64.StdEncoding.EncodeToString([]byte(text)))
	_, _ = os.Stdout.WriteString(osc)

	if os.Getenv("WAYLAND_DISPLAY") != "" {
		if path, err := exec.LookPath("wl-copy"); err == nil {
			if err := runClipboardCommand(path, nil, text); err == nil {
				return nil
			}
		}
	}

	if os.Getenv("DISPLAY") != "" {
		if path, err := exec.LookPath("xclip"); err == nil {
			if err := runClipboardCommand(path, []string{"-selection", "clipboard"}, text); err == nil {
				return nil
			}
		}
		if path, err := exec.LookPath("xsel"); err == nil {
			if err := runClipboardCommand(path, []string{"--clipboard", "--input"}, text); err == nil {
				return nil
			}
		}
	}

	if path, err := exec.LookPath("pbcopy"); err == nil {
		if err := runClipboardCommand(path, nil, text); err == nil {
			return nil
		}
	}

	if path, err := exec.LookPath("clip.exe"); err == nil {
		if err := runClipboardCommand(path, nil, text); err == nil {
			return nil
		}
	}

	if os.Getenv("TERM") != "" {
		return nil
	}

	return errNoClipboardBackend
}
