package utils

import (
	"encoding/json"
	"golang.org/x/xerrors"
	"os"
)

func WriteFile(filepath string, data any) error {
	d, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return xerrors.Errorf("failed to marshal:%w", err)
	}
	if err = os.WriteFile(filepath, d, os.ModePerm); err != nil {
		return xerrors.Errorf("failed to write to %s:%w", filepath, err)
	}

	return nil
}

func ReadFile(filepath string, data any) error {
	d, err := os.ReadFile(filepath)
	if err != nil {
		return xerrors.Errorf("failed to read %s:%w", filepath, err)
	}

	if err = json.Unmarshal(d, data); err != nil {
		return xerrors.Errorf("failed to unmarshal %s:%w", filepath, err)
	}

	return nil
}

func Mkdir(dir string) error {
	_, err := os.Stat(dir)
	if err == nil {
		return nil
	}
	if !os.IsNotExist(err) {
		return err
	}

	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	return nil
}
