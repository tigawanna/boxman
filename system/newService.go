package system

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type SystemdServiceConfig struct {
	// Unit section
	Unit struct {
		Description string
	}

	// Service section
	Service struct {
		Type           string
		User           string
		Group          string
		LimitNOFILE    int
		Restart        string
		RestartSec     string
		StandardOutput string
		StandardError  string
		ExecStart      string
	}

	// Install section
	Install struct {
		WantedBy string
	}
}

// ConfigOptions allows overriding default service settings
//
// Example usage:
//
//	package main
//
//	import (
//		"fmt"
//		"github.com/tigawanna/boxman/system"
//	)
//
//	func main() {
//		svc := system.SystemdServiceConfig{
//			Unit: system.Unit{
//				Description: "my service",
//			},
//			Service: system.Service{
//				Type:           "simple",
//				User:           "root",
//				Group:          "root",
//				LimitNOFILE:    4096,
//				Restart:        "always",
//				RestartSec:     "5s",
//				StandardOutput: "append:/root/pb/std.log",
//				StandardError:  "append:/root/pb/std.log",
//				ExecStart:      "/root/pb/pocketbase serve yourdomain.com",
//			},
//			Install: system.Install{
//				WantedBy: "multi-user.target",
//			},
//		}
//		fmt.Println(svc.ToString())
//	}
type ConfigOptions struct {
	Type        string
	User        string
	Group       string
	LimitNOFILE int
	Restart     string
	RestartSec  string
}

// NewServiceConfig generates a SystemdServiceConfig for the given service name,
// base directory and exec command. The opts parameter allows overriding default
// service settings. If opts is nil, default options are used.
//
// The base directory is expanded if it starts with ~ and is ensured to be an
// absolute path. The log file is created in the base directory under the
// "logs" directory.
//
// The generated service will be configured with the following defaults:
//
//   - Type: simple
//   - User: root
//   - Group: root
//   - LimitNOFILE: 4096
//   - Restart: always
//   - RestartSec: 5s
//   - StandardOutput: append:baseDir/logs/service.log
//   - StandardError: append:baseDir/logs/service.log
//   - ExecStart: baseDir/execCommand
//   - WantedBy: multi-user.target
//
// Example usage:
//
//		package main
//
//	config := NewServiceConfig(
//		"my-node-server",
//		"~/my-node-server",
//		"node /dist/index.js",
//		&ConfigOptions{
//			User:  "pocketbase",
//			Group: "pocketbase",
//		},
//	)
//	fmt.Println(config.ToString())
//
// )
func NewServiceConfig(serviceName, baseDir, execCommand string, opts *ConfigOptions) SystemdServiceConfig {
	// Default options
	if opts == nil {
		opts = &ConfigOptions{
			Type:        "simple",
			User:        "root",
			Group:       "root",
			LimitNOFILE: 4096,
			Restart:     "always",
			RestartSec:  "5s",
		}
	}

	// Expand home directory if path starts with ~
	if strings.HasPrefix(baseDir, "~/") {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			baseDir = filepath.Join(homeDir, baseDir[2:])
		}
	}

	// Ensure base directory is absolute
	baseDir, _ = filepath.Abs(baseDir)

	// Build paths
	logPath := filepath.Join(baseDir, "logs", "service.log")
	execPath := filepath.Join(baseDir, execCommand)

	return SystemdServiceConfig{
		Unit: struct{ Description string }{
			Description: fmt.Sprintf("%s service", serviceName),
		},
		Service: struct {
			Type           string
			User           string
			Group          string
			LimitNOFILE    int
			Restart        string
			RestartSec     string
			StandardOutput string
			StandardError  string
			ExecStart      string
		}{
			Type:           opts.Type,
			User:           opts.User,
			Group:          opts.Group,
			LimitNOFILE:    opts.LimitNOFILE,
			Restart:        opts.Restart,
			RestartSec:     opts.RestartSec,
			StandardOutput: "append:" + logPath,
			StandardError:  "append:" + logPath,
			ExecStart:      execPath,
		},
		Install: struct{ WantedBy string }{
			WantedBy: "multi-user.target",
		},
	}
}

func (c SystemdServiceConfig) ToString() string {
	var sb strings.Builder

	// [Unit] section
	sb.WriteString("[Unit]\n")
	sb.WriteString(fmt.Sprintf("Description=%s\n\n", c.Unit.Description))

	// [Service] section
	sb.WriteString("[Service]\n")
	sb.WriteString(fmt.Sprintf("Type=%s\n", c.Service.Type))
	sb.WriteString(fmt.Sprintf("User=%s\n", c.Service.User))
	sb.WriteString(fmt.Sprintf("Group=%s\n", c.Service.Group))
	sb.WriteString(fmt.Sprintf("LimitNOFILE=%d\n", c.Service.LimitNOFILE))
	sb.WriteString(fmt.Sprintf("Restart=%s\n", c.Service.Restart))
	sb.WriteString(fmt.Sprintf("RestartSec=%s\n", c.Service.RestartSec))
	sb.WriteString(fmt.Sprintf("StandardOutput=%s\n", c.Service.StandardOutput))
	sb.WriteString(fmt.Sprintf("StandardError=%s\n", c.Service.StandardError))
	sb.WriteString(fmt.Sprintf("ExecStart=%s\n\n", c.Service.ExecStart))

	// [Install] section
	sb.WriteString("[Install]\n")
	sb.WriteString(fmt.Sprintf("WantedBy=%s\n", c.Install.WantedBy))

	return sb.String()
}

// func main() {
// 	config := NewServiceConfig(
// 		"pocketbase",
// 		"~/pb",
// 		"pocketbase serve yourdomain.com",
// 		&ConfigOptions{
// 			User:  "pocketbase",
// 			Group: "pocketbase",
// 		},
// 	)
// 	fmt.Println(config.ToString())
// }
