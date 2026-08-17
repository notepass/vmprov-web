package libvirt

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

// ErrKeyPassphraseProtected is returned when the key file is a
// passphrase-protected private key, which is not supported.
var ErrKeyPassphraseProtected = errors.New("passphrase-protected SSH keys are not supported")

// ValidateSSHKey checks that the key file exists, is readable, and parses
// as a non-passphrase-protected SSH private key.
func ValidateSSHKey(path string) error {
	if path == "" {
		return fmt.Errorf("ssh key path is empty")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read ssh key file %q: %w", path, err)
	}

	if _, err := ssh.ParsePrivateKey(data); err != nil {
		var passphraseErr *ssh.PassphraseMissingError
		if errors.As(err, &passphraseErr) {
			return fmt.Errorf("ssh key file %q: %w", path, ErrKeyPassphraseProtected)
		}
		return fmt.Errorf("ssh key file %q is not a valid private key: %w", path, err)
	}

	return nil
}
