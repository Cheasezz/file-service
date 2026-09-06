package repo

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func (fs *FileSystem) GetUserHashList(userID string) (map[string]string, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	return fs.readHashListFromDisk(userID)
}

func (fs *FileSystem) UpdateHash(userID, filename, hash string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	hl, err := fs.readHashListFromDisk(userID)
	if err != nil {
		return err
	}

	hl[filename] = hash

	return fs.writeHashListToDisk(userID, hl)
}

func (fs *FileSystem) readHashListFromDisk(userID string) (map[string]string, error) {
	var hl map[string]string

	data, err := os.ReadFile(filepath.Join(fs.path, userID, ".hash_list.json"))

	if os.IsNotExist(err) {
		return map[string]string{}, nil
	}
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &hl)
	if err != nil {
		return nil, err
	}

	return hl, nil
}

func (fs *FileSystem) writeHashListToDisk(userID string, hashList map[string]string) error {
	data, err := json.MarshalIndent(hashList, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(fs.path, userID, ".hash_list.json"), data, 0o644)
}
