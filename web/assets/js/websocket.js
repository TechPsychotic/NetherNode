const ws = new WebSocket(`wss://${window.location.host}/ws`);

ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    
    // Update server status in dashboard
    if (data.type === 'STATUS_UPDATE') {
        const serverElement = document.querySelector(`[data-server-id="${data.serverId}"]`);
        if (serverElement) {
            serverElement.querySelector('.server-status').className = 
                `server-status status-${data.status}`;
        }
    }
    
    // Update console output
    if (data.type === 'LOG_UPDATE' && window.location.pathname.includes('console.html')) {
        const consoleOutput = document.getElementById('console-output');
        consoleOutput.textContent += data.log + '\n';
        consoleOutput.scrollTop = consoleOutput.scrollHeight; // Auto-scroll
    }
};

// Reconnect on close
ws.onclose = () => {
    setTimeout(() => location.reload(), 1000);
};