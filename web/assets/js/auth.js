document.addEventListener('DOMContentLoaded', () => {
    // Registration form handler
    const registerForm = document.getElementById('register-form');
    if (registerForm) {
        registerForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const formData = {
                username: registerForm.querySelector('[name="username"]').value,
                email: registerForm.querySelector('[name="email"]').value,
                password: registerForm.querySelector('[name="password"]').value,
                dob: registerForm.querySelector('[name="dob"]').value
            };

            try {
                const response = await fetch('/register', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(formData)
                });

                if (response.ok) {
                    window.location.href = '/login';
                } else {
                    const error = await response.json();
                    alert(error.error);
                }
            } catch (err) {
                console.error('Registration failed:', err);
                alert('Registration failed. Please try again.');
            }
        });
    }

    // Login form handler
    const loginForm = document.getElementById('login-form');
    if (loginForm) {
        loginForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const formData = {
                login: loginForm.querySelector('[name="login"]').value,
                password: loginForm.querySelector('[name="password"]').value
            };

            try {
                const response = await fetch('/login', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(formData)
                });

                // After successful login
if (response.ok) {
    const data = await response.json();
    localStorage.setItem('token', data.token);
    localStorage.setItem('user', JSON.stringify(data.user));
    
    // Redirect to dashboard instead of home
    window.location.href = '/dashboard'; 
}else {
                    const error = await response.json();
                    alert(error.error);
                }
            } catch (err) {
                console.error('Login failed:', err);
                alert('Login failed. Please try again.');
            }
        });
    }
});
