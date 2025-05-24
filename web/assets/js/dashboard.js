document.addEventListener('DOMContentLoaded', async () => {
    // Check authentication
    const token = localStorage.getItem('token');
    if (!token) {
        window.location.href = '/login.html';
        return;
    }

    // Load user's servers
    await loadServers();

    // Setup event listeners
    document.getElementById('create-server').addEventListener('click', createServer);
    document.getElementById('logout-btn').addEventListener('click', logout);
});

async function fetchServers() {
    const response = await fetch('http://serveo.net:6969/servers', {
  headers: {
    Authorization: token // assuming `token` is set
  }
});

if (!response.ok) {
  const errorText = await response.text();
  throw new Error(`Failed to load servers: ${errorText}`);
}

const data = await response.json();

    console.log(data);
}
fetchServers(); // ✅ Call the async function

async function createServer() {
    const serverName = prompt("Enter a name for your new server:");
    console.log("button pressed");
    if (!serverName || serverName.trim() === '') {
        alert('Server name cannot be empty.');
        return;
    }

    try {
        const response = await fetch('/servers', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                
            },
            body: JSON.stringify({ name: serverName })
        });

        if (response.ok) {
            showNotification('Server created successfully');
            await loadServers();
        } else {
            const errMsg = await response.text();
            alert(`Failed to create server: ${errMsg}`);
        }
    } catch (err) {
        console.error('Create server error:', err);
    }
}

async function loadServers() {
    try {
        const response = await fetch('/servers', {
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            }
        });

        if (response.ok) {
            const servers = await response.json();
            renderServers(servers);
        } else if (response.status === 401) {
            logout();
        }
    } catch (err) {
        console.error('Failed to load servers:', err);
    }
}

function renderServers(servers) {
    const container = document.getElementById('server-list');
    container.innerHTML = servers.map(server => `
        <div class="server-card" data-server-id="${server.id}">
            <h3>${server.name}</h3>
            <p>Status: <span class="status">${server.status}</span></p>
            <div class="server-actions">
                <button class="start-btn" ${server.status === 'running' ? 'disabled' : ''}>
                    Start
                </button>
                <button class="stop-btn" ${server.status !== 'running' ? 'disabled' : ''}>
                    Stop
                </button>
                <button class="settings-btn">Settings</button>
                <button class="console-btn">Console</button>
            </div>
        </div>
    `).join('');

    // Add event listeners to all action buttons
    document.querySelectorAll('.start-btn').forEach(btn => {
        btn.addEventListener('click', startServer);
    });
    // Add similar listeners for other buttons
}

async function startServer(e) {
    const serverId = e.target.closest('.server-card').dataset.serverId;
    try {
        const response = await fetch(`/servers/${server.id}/start`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            }
        });

        if (response.ok) {
            await loadServers(); // Refresh server list
        }
    } catch (err) {
        console.error('Failed to start server:', err);
    }
}

function logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    window.location.href = '/login.html';
}
// Add these functions to your existing dashboard.js

async function restartServer(serverId) {
    try {
        const response = await fetch(`/servers/${server.id}/restart`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            }
        });
        
        if (response.ok) {
            showNotification('Server restart initiated');
            await loadServers();
        }
    } catch (err) {
        console.error('Restart failed:', err);
    }
}

async function viewConsole(serverId) {
    try {
        const response = await fetch(`/servers/${server.id}/console`, {
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            }
        });
        
        if (response.ok) {
            const data = await response.json();
            openConsoleModal(data.output);
        }
    } catch (err) {
        console.error('Failed to get console:', err);
    }
}

function openConsoleModal(content) {
    // Implement a modal to show console output
    const modal = document.createElement('div');
    modal.className = 'console-modal';
    modal.innerHTML = `
        <div class="modal-content">
            <pre>${content}</pre>
            <button class="close-btn">Close</button>
        </div>
    `;
    document.body.appendChild(modal);
    
    modal.querySelector('.close-btn').addEventListener('click', () => {
        modal.remove();
    });
}

// Update your server card rendering to include new buttons
function renderServerCard(server) {
    return `
        <div class="server-card" data-server-id="${server.id}">
            <h3>${server.name}</h3>
            <p>Status: <span class="status">${server.status}</span></p>
            <div class="server-actions">
                <button class="start-btn" ${server.status === 'running' ? 'disabled' : ''}>
                    Start
                </button>
                <button class="stop-btn" ${server.status !== 'running' ? 'disabled' : ''}>
                    Stop
                </button>
                <button class="restart-btn">Restart</button>
                <button class="console-btn">Console</button>
                <button class="settings-btn">Settings</button>
            </div>
        </div>
    `;
}

// Add event listeners for new buttons
document.addEventListener('click', async (e) => {
    if (e.target.classList.contains('restart-btn')) {
        const serverId = e.target.closest('.server-card').dataset.serverId;
        await restartServer(serverId);
    }
    
    if (e.target.classList.contains('console-btn')) {
        const serverId = e.target.closest('.server-card').dataset.serverId;
        await viewConsole(serverId);
    }
});
// Add to dashboard.js
async function openSettings(serverId) {
    try {
        // Get current properties
        const response = await fetch(`/servers/${server.id}/properties`, {
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            }
        });
        
        if (response.ok) {
            const properties = await response.json();
            showSettingsModal(serverId, properties);
        }
    } catch (err) {
        console.error('Failed to get properties:', err);
    }
}

function showSettingsModal(serverId, properties) {
    const modal = document.createElement('div');
    modal.className = 'settings-modal';
    
    // Create form inputs for properties
    let formHTML = '<form class="server-properties-form">';
    for (const [key, value] of Object.entries(properties)) {
        if (!key.startsWith("#")) {
            formHTML += `
                <div class="form-group">
                    <label>${key}</label>
                    <input type="text" name="${key}" value="${value}">
                </div>
            `;
        }
    }
    formHTML += `
        <button type="submit" class="btn-primary">Save</button>
        <button type="button" class="btn-outline close-btn">Cancel</button>
    </form>`;
    
    modal.innerHTML = `
        <div class="modal-content">
            <h3>Server Properties</h3>
            ${formHTML}
        </div>
    `;
    
    document.body.appendChild(modal);
    
    modal.querySelector('form').addEventListener('submit', async (e) => {
        e.preventDefault();
        const formData = new FormData(e.target);
        const properties = {};
        
        for (const [key, value] of formData.entries()) {
            properties[key] = value;
        }
        
        try {
            const response = await fetch(`/servers/${server.id}/properties`, {
                method: 'PUT',
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('token')}`,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(properties)
            });
            
            if (response.ok) {
                showNotification('Properties updated successfully');
                modal.remove();
            }
        } catch (err) {
            console.error('Failed to update properties:', err);
        }
    });
    
    modal.querySelector('.close-btn').addEventListener('click', () => {
        modal.remove();
    });
}
// In your dashboard.js
async function refreshConsole(serverId) {
    const response = await fetch(`/servers/${serverId}/console`);
    const data = await response.json();
    document.getElementById('console-output').textContent = data.output;
}
