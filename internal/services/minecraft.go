package services

import (
    "os/exec"
    "your-project/internal/storage"
    "your-project/internal/utils"
)

var activeProcesses = make(map[uint]*exec.Cmd)

func StartMinecraftServer(serverID uint) {
    server, _ := storage.GetServer(serverID)
    serverPath := utils.GetServerPath(server.UserID, serverID)
    
    cmd := exec.Command(
        "java",
        "-Xmx1024M",
        "-jar", "server.jar",
        "--nogui",
        "--port", strconv.Itoa(server.Port),
    )
    cmd.Dir = serverPath
    
    // Capture and broadcast logs
    stdout, _ := cmd.StdoutPipe()
    go func() {
        scanner := bufio.NewScanner(stdout)
        for scanner.Scan() {
            BroadcastLog(serverID, scanner.Text())
        }
    }()
    
    cmd.Start()
    activeProcesses[serverID] = cmd
}