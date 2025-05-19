document.addEventListener('DOMContentLoaded', async () => {
    // Load servers on page load
    const servers = await fetchServers();
    renderServers(servers);
});

async function fetchServers() {
    const response = await fetch('/api/servers', {
        headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
    });
    return await response.json();
}

function renderServers(servers) {
    const container = document.getElementById('server-list');
    container.innerHTML = servers.map(server => `
      <div class="server-card">
        <div class="server-header">
          <span class="server-status status-${server.status}"></span>
          <h3>${server.name}</h3>
        </div>
        <div class="server-meta">
          <span>Port: ${server.port}</span>
          <span>Players: 0/20</span>
        </div>
        <div class="server-actions">
          ${server.status === 'running' ? 
            `<button class="btn-primary stop-btn" onclick="stopServer(${server.id})">Stop</button>` :
            `<button class="btn-primary start-btn" onclick="startServer(${server.id})">Start</button>`
          }
          <a href="/web/server/console.html?serverId=${server.id}" class="btn-primary">Console</a>
          <a href="/web/server/settings.html?serverId=${server.id}" class="btn-primary">Settings</a>
        </div>
      </div>
    `).join('');
  }
async function startServer(serverId) {
    await fetch(`/api/servers/${serverId}/start`, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
    });
    // Refresh list after action
    const servers = await fetchServers();
    renderServers(servers);
}