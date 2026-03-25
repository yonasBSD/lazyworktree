package config

import "fmt"

// ContainerMount represents a bind mount for a container.
type ContainerMount struct {
	Source   string // Host path (supports env var expansion)
	Target   string // Container path
	ReadOnly bool   // Mount as read-only
	Options  string // Extra mount options (e.g. "z", "Z", "shared")
}

// ContainerCommand configures OCI container execution for a custom command.
type ContainerCommand struct {
	Image       string            // Required: container image (e.g. "golang:1.22")
	Runtime     string            // Optional: "docker" or "podman" (auto-detected if empty)
	Mounts      []ContainerMount  // Optional: additional bind mounts
	Env         map[string]string // Optional: extra environment variables
	WorkingDir  string            // Optional: working directory inside container (default: /workspace)
	ExtraArgs   []string          // Optional: extra docker/podman run arguments
	Args        []string          // Optional: arguments passed after image (as CMD)
	Interactive bool              // Optional: allocate TTY for interactive use
}

func parseContainerCommand(data map[string]any) *ContainerCommand {
	cmd := &ContainerCommand{
		Image:       getString(data, "image"),
		Runtime:     getString(data, "runtime"),
		WorkingDir:  getString(data, "working_dir"),
		Interactive: coerceBool(data["interactive"], false),
	}
	if args, ok := data["args"].([]any); ok {
		for _, a := range args {
			cmd.Args = append(cmd.Args, fmt.Sprint(a))
		}
	}
	if mounts, ok := data["mounts"].([]any); ok {
		for _, m := range mounts {
			if mData, ok := m.(map[string]any); ok {
				cmd.Mounts = append(cmd.Mounts, ContainerMount{
					Source:   getString(mData, "source"),
					Target:   getString(mData, "target"),
					ReadOnly: coerceBool(mData["read_only"], false),
					Options:  getString(mData, "options"),
				})
			}
		}
	}
	if envData, ok := data["env"].(map[string]any); ok {
		cmd.Env = make(map[string]string)
		for k, v := range envData {
			cmd.Env[k] = fmt.Sprint(v)
		}
	}
	if args, ok := data["extra_args"].([]any); ok {
		for _, a := range args {
			cmd.ExtraArgs = append(cmd.ExtraArgs, fmt.Sprint(a))
		}
	}
	if cmd.Image == "" {
		return nil
	}
	return cmd
}
