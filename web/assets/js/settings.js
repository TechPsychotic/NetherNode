document.addEventListener('DOMContentLoaded', async () => {
    const serverId = new URLSearchParams(window.location.search).get('serverId');
    const form = document.getElementById('server-properties-form');
    
    // Load existing properties
    const properties = await fetch(`/api/servers/${serverId}/properties`, {
        headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
    }).then(res => res.json());
    
    // Populate form
    Object.entries(properties).forEach(([key, value]) => {
        const input = form.elements[key];
        if (input) input.value = value;
    });
    
    // Save changes
    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        const formData = new FormData(form);
        const properties = Object.fromEntries(formData.entries());
        
        await fetch(`/api/servers/${serverId}/properties`, {
            method: 'PUT',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(properties)
        });
        alert('Settings saved!');
    });
});