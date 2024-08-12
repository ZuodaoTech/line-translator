package tasker

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os/exec"

	"github.com/zuodaotech/line-translator/core"
)

func (w *Worker) ProcessTaskActionFetchAudioAndTranscript(ctx context.Context, task *core.Task) (core.JSONMap, error) {
	messageID := task.Params.GetString("message_id")
	if messageID == "" {
		return nil, fmt.Errorf("messageID is empty")
	}

	groupID := task.Params.GetString("group_id")
	if groupID == "" {
		return nil, fmt.Errorf("group_id is empty")
	}

	cli, err := w.GetLineClient(groupID)
	if err != nil {
		return nil, err
	}

	buf, err := cli.GetContent(messageID)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("buf: %v\n", len(buf))

	wavBuf, err := convertMP4ToMP3(buf)
	if err != nil {
		slog.Error("failed to convert mp3 to wav", "error", err)
		return nil, err
	}

	// fmt.Printf("wavBuf: %v\n", len(wavBuf))

	lang, err := w.speechCli.ToText("zh-CN", wavBuf)
	if err != nil {
		return nil, err
	}

	fmt.Printf("lang: %v\n", lang)

	return nil, nil
}

func convertMP4ToMP3(inputData []byte) ([]byte, error) {
	// Prepare the ffmpeg command to read from stdin and write to stdout
	cmd := exec.Command("ffmpeg", "-i", "pipe:0", "-ar", "16000", "-acodec", "pcm_s16le", "-f", "wav", "pipe:1")

	// Set up pipes to send input data and capture output data
	cmd.Stdin = bytes.NewReader(inputData)
	var outBuffer bytes.Buffer
	cmd.Stdout = &outBuffer

	// Capture any errors from stderr
	var errBuffer bytes.Buffer
	cmd.Stderr = &errBuffer

	// Run the command
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg command failed: %v, %s", err, errBuffer.String())
	}

	// Return the converted MP3 data
	return outBuffer.Bytes(), nil
}
