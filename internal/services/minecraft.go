package services

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strconv"
    "strings"
    "sync"

    "NetherNode/internal/storage"
    "NetherNode/internal/utils"
)

var (
    activeProcesses = make(map[uint]*exec.Cmd)
    processMutex    = &sync.Mutex{}
)

// Server management functions
func StartMinecraftServer(serverID uint) error {
    processMutex.Lock()
    defer processMutex.Unlock()

    server, err := storage.GetServer(serverID)
    if err != nil {
        return fmt.Errorf("get server error: %w", err)
    }

    serverPath := utils.GetServerPath(server.UserID, serverID)
    cmd := exec.Command(
        "java",
        "-Xmx1024M",
        "-jar", "server.jar",
        "--nogui",
        "--port", strconv.Itoa(server.Port),
    )
    cmd.Dir = serverPath

    stdout, err := cmd.StdoutPipe()
    if err != nil {
        return fmt.Errorf("stdout pipe error: %w", err)
    }

    go func() {
    scanner := bufio.NewScanner(stdout)
    for scanner.Scan() {
        BroadcastLog(serverID, scanner.Text()) // ✅ Correct usage
    }
}()


    if err := cmd.Start(); err != nil {
        return fmt.Errorf("start command error: %w", err)
    }

    activeProcesses[serverID] = cmd
    return nil
}

func StopMinecraftServer(serverID uint) error {
    processMutex.Lock()
    defer processMutex.Unlock()

    cmd, exists := activeProcesses[serverID]
    if !exists {
        return fmt.Errorf("server not running")
    }

    if err := cmd.Process.Kill(); err != nil {
        return fmt.Errorf("kill process error: %w", err)
    }

    delete(activeProcesses, serverID)
    return nil
}

func RestartMinecraftServer(serverID uint) error {
    if err := StopMinecraftServer(serverID); err != nil {
        return fmt.Errorf("stop failed: %w", err)
    }
    return StartMinecraftServer(serverID)
}

func GetServerConsole(serverID uint) (string, error) {
    server, err := storage.GetServer(serverID)
    if err != nil {
        return "", fmt.Errorf("get server error: %w", err)
    }

    logPath := filepath.Join(utils.GetServerPath(server.UserID, serverID), "logs", "latest.log")
    content, err := os.ReadFile(logPath)
    return string(content), err
}

func UpdateServerProperties(serverID uint, properties map[string]string) error {
    server, err := storage.GetServer(serverID)
    if err != nil {
        return err
    }

    propsPath := filepath.Join(utils.GetServerPath(server.UserID, serverID), "server.properties")
    content, err := os.ReadFile(propsPath)
    if err != nil {
        return err
    }

    lines := strings.Split(string(content), "\n")
    for i, line := range lines {
        parts := strings.SplitN(line, "=", 2)
        if len(parts) == 2 {
            if newVal, exists := properties[strings.TrimSpace(parts[0])]; exists {
                lines[i] = parts[0] + "=" + newVal
            }
        }
    }

    return os.WriteFile(propsPath, []byte(strings.Join(lines, "\n")), 0644)
}

