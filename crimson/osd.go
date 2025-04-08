package crimson

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func GetOSDID(socketFilename string) (string, error) {
	filename := filepath.Base(socketFilename)
	split := strings.Split(filename, ".")
	if len(split) != 3 {
		return "", errors.New("unexpected filename format for osd socket")
	}
	_, err := strconv.Atoi(split[1])
	if err != nil {
		return "", err
	}
	return split[1], nil
}

func GetOSDSocketFiles(basepath string) (files []string, err error) {
	sockFiles, err := os.ReadDir(basepath)
	if err != nil {
		return nil, err
	}

	for _, f := range sockFiles {
		// We should not proceed reading any file that is not named according
		// to the ceph socket naming conventions.
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".asok") {
			continue
		}

		if !strings.Contains(f.Name(), "osd") {
			continue
		}

		files = append(files, f.Name())
	}

	return
}

func RunAsokDumpMetrics(f string) ([]byte, error) {
	args := []string{
		"--admin-daemon",
		f,
		"dump_metrics",
	}
	out, err := exec.Command("ceph", args...).Output()
	if err != nil {
		return nil, err
	}
	return out, nil
}
