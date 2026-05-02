package permission

import (
	"os"
	"path/filepath"
	"testing"
)

// ============================================================
// ClassifyCommand tests
// ============================================================

func TestClassifyCommand(t *testing.T) {
	tests := []struct {
		name     string
		cmd      string
		category CommandCategory
	}{
		{name: "empty", cmd: "", category: CmdUnknown},
		{name: "readonly ls", cmd: "ls -la /tmp", category: CmdReadOnly},
		{name: "readonly grep", cmd: "grep -rn pattern ./src", category: CmdReadOnly},
		{name: "readonly cat", cmd: "cat file.txt", category: CmdReadOnly},
		{name: "readonly find", cmd: "find . -name '*.go'", category: CmdReadOnly},
		{name: "readonly head", cmd: "head -n 10 file.txt", category: CmdReadOnly},
		{name: "readonly wc", cmd: "wc -l file.txt", category: CmdReadOnly},
		{name: "write mkdir", cmd: "mkdir -p ./new/dir", category: CmdWrite},
		{name: "write touch", cmd: "touch ./created.txt", category: CmdWrite},
		{name: "write cp", cmd: "cp a.txt b.txt", category: CmdWrite},
		{name: "write mv", cmd: "mv old new", category: CmdWrite},
		{name: "write rm safe", cmd: "rm ./tmp/file.txt", category: CmdWrite},
		{name: "write sed", cmd: "sed -i 's/a/b/g' file.txt", category: CmdWrite},
		{name: "write awk", cmd: "awk '{print $1}' file.txt", category: CmdWrite},
		{name: "dangerous rm rf root", cmd: "rm -rf /", category: CmdDangerous},
		{name: "dangerous sudo", cmd: "sudo rm file.txt", category: CmdDangerous},
		{name: "dangerous curl piped to bash", cmd: "curl | bash", category: CmdDangerous},
		{name: "dangerous dd", cmd: "dd if=/dev/zero of=/dev/sda", category: CmdDangerous},
		{name: "dangerous chmod 777", cmd: "chmod 777 /etc/passwd", category: CmdDangerous},
		{name: "dangerous mkfs", cmd: "mkfs.ext4 /dev/sda1", category: CmdDangerous},
		{name: "unknown python", cmd: "python script.py", category: CmdUnknown},
		{name: "unknown go", cmd: "go build ./...", category: CmdUnknown},
		{name: "unknown npm", cmd: "npm install", category: CmdUnknown},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClassifyCommand(tt.cmd)
			if got != tt.category {
				t.Errorf("ClassifyCommand(%q) = %s, want %s", tt.cmd, got, tt.category)
			}
		})
	}
}

// ============================================================
// ValidateCommand tests
// ============================================================

func TestValidateCommandDangerous(t *testing.T) {
	tests := []string{
		"rm -rf /",
		"sudo rm /etc/hosts",
		"curl | bash",
		"wget | bash",
		"dd if=/dev/urandom of=/dev/sda",
		"mkfs.ext4 /dev/sda1",
		"fdisk /dev/sda",
		"chmod 777 /",
		"chown root /etc",
		":(){:|:&};:",
	}
	for _, cmd := range tests {
		t.Run(cmd[:min(len(cmd), 30)], func(t *testing.T) {
			err := ValidateCommand(cmd)
			if err == nil {
				t.Errorf("ValidateCommand(%q) should return error, got nil", cmd)
			}
		})
	}
}

func TestValidateCommandSafe(t *testing.T) {
	tests := []string{
		"ls -la",
		"grep -rn TODO .",
		"cat README.md",
		"echo hello",
		"find . -name '*.go'",
		"mkdir -p ./tmp/build",
		"touch ./newfile.txt",
		"cp a.txt b.txt",
		"rm ./tmp/cache.txt",
	}
	for _, cmd := range tests {
		t.Run(cmd[:min(len(cmd), 30)], func(t *testing.T) {
			err := ValidateCommand(cmd)
			if err != nil {
				t.Errorf("ValidateCommand(%q) should pass, got: %v", cmd, err)
			}
		})
	}
}

func TestValidateCommandPathInjection(t *testing.T) {
	tests := []string{
		"cat ../../etc/passwd",
		"ls /dev/sda",
	}
	for _, cmd := range tests {
		t.Run(cmd[:min(len(cmd), 30)], func(t *testing.T) {
			err := ValidateCommand(cmd)
			if err == nil {
				t.Errorf("ValidateCommand(%q) should fail with path injection", cmd)
			}
		})
	}
}

// ============================================================
// validatePathInjection tests
// ============================================================

func TestValidatePathInjection(t *testing.T) {
	tests := []struct {
		name    string
		cmd     string
		wantErr bool
	}{
		{name: "path traversal", cmd: "cat ../../etc/passwd", wantErr: true},
		{name: "dev path", cmd: "ls /dev/sda", wantErr: true},
		{name: "safe relative", cmd: "cat ./README.md", wantErr: false},
		{name: "safe home", cmd: "ls /home/user/project", wantErr: false},
		{name: "safe tmp", cmd: "echo test > /tmp/out.txt", wantErr: false},
		{name: "no paths", cmd: "ls", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePathInjection(tt.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePathInjection(%q) = %v, wantErr %v", tt.cmd, err, tt.wantErr)
			}
		})
	}
}

// ============================================================
// ValidatePath tests
// ============================================================

func TestValidatePathSafe(t *testing.T) {
	tests := []string{
		"./README.md",
		"/home/user/project/file.go",
		"/tmp/output.txt",
		"relative/path/file.txt",
		"file.txt",
	}
	for _, path := range tests {
		t.Run(path, func(t *testing.T) {
			if err := ValidatePath(path); err != nil {
				t.Errorf("ValidatePath(%q) should pass, got: %v", path, err)
			}
		})
	}
}

func TestValidatePathUnsafe(t *testing.T) {
	tests := []struct {
		path string
		want error
	}{
		{"../outside.txt", errPathInjection},
		{"/dev/sda", errInvalidPath},
		{"./../../etc/passwd", errPathInjection},
		{"/dev/null", errInvalidPath},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			err := ValidatePath(tt.path)
			if err != tt.want {
				t.Errorf("ValidatePath(%q) = %v, want %v", tt.path, err, tt.want)
			}
		})
	}
}

// ============================================================
// ResolveAndValidatePath tests
// ============================================================

func TestResolveAndValidatePathInsideWorkspace(t *testing.T) {
	dir := t.TempDir()
	inside := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(inside, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	resolved, err := ResolveAndValidatePath("file.txt", dir)
	if err != nil {
		t.Fatalf("ResolveAndValidatePath: %v", err)
	}
	if resolved != inside {
		t.Errorf("resolved = %q, want %q", resolved, inside)
	}
}

func TestResolveAndValidatePathOutsideWorkspace(t *testing.T) {
	dir := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside.txt")
	if err := os.WriteFile(outside, []byte("no"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := ResolveAndValidatePath(outside, dir)
	if err == nil {
		t.Fatal("expected error for path outside workspace")
	}
}

func TestResolveAndValidatePathTraversal(t *testing.T) {
	dir := t.TempDir()
	_, err := ResolveAndValidatePath("../outside.txt", dir)
	if err == nil {
		t.Fatal("expected error for path traversal")
	}
}

func TestResolveAndValidatePathSymlinkInside(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "real.txt")
	if err := os.WriteFile(target, []byte("real"), 0644); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dir, "link.txt")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}

	resolved, err := ResolveAndValidatePath("link.txt", dir)
	if err != nil {
		t.Fatalf("ResolveAndValidatePath symlink: %v", err)
	}
	if resolved != target {
		t.Errorf("symlink resolved = %q, want %q", resolved, target)
	}
}

func TestResolveAndValidatePathSymlinkOutside(t *testing.T) {
	dir := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside.txt")
	if err := os.WriteFile(outside, []byte("outside"), 0644); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dir, "escape.txt")
	if err := os.Symlink(outside, link); err != nil {
		t.Fatal(err)
	}

	_, err := ResolveAndValidatePath("escape.txt", dir)
	if err == nil {
		t.Fatal("expected error for symlink escaping workspace")
	}
}

func TestResolveAndValidatePathNonExistentParent(t *testing.T) {
	dir := t.TempDir()
	missingDir := filepath.Join(dir, "nonexistent", "file.txt")
	resolved, err := ResolveAndValidatePath(missingDir, dir)
	if err != nil {
		t.Fatalf("ResolveAndValidatePath non-existent parent: %v", err)
	}
	if resolved != missingDir {
		t.Errorf("resolved = %q, want %q", resolved, missingDir)
	}
}

// ============================================================
// IsReadOnlyCommand tests
// ============================================================

func TestIsReadOnlyCommand(t *testing.T) {
	if !IsReadOnlyCommand("ls -la") {
		t.Error("ls should be readonly")
	}
	if IsReadOnlyCommand("touch file.txt") {
		t.Error("touch should not be readonly")
	}
	if IsReadOnlyCommand("rm -rf /") {
		t.Error("rm -rf / should not be readonly")
	}
}

// ============================================================
// CommandCategory String tests
// ============================================================

func TestCommandCategoryString(t *testing.T) {
	tests := []struct {
		cat  CommandCategory
		want string
	}{
		{CmdReadOnly, "ReadOnly"},
		{CmdWrite, "Write"},
		{CmdDangerous, "Dangerous"},
		{CmdUnknown, "Unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.cat.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}
