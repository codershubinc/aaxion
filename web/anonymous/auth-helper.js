// Authentication helper for admin pages

// Check if user is authenticated
function checkAuth() {
    const token = localStorage.getItem('authToken');
    if (!token) {
        // Redirect to login with current page as redirect target
        const currentPath = window.location.pathname + window.location.search;
        window.location.href = `/login?redirect=${encodeURIComponent(currentPath)}`;
        return false;
    }
    return true;
}

// Get auth token
function getAuthToken() {
    return localStorage.getItem('authToken');
}

// Make authenticated API request
async function authenticatedFetch(url, options = {}) {
    const token = getAuthToken();

    if (!token) {
        window.location.href = '/login';
        throw new Error('Not authenticated');
    }

    // Add auth header
    const headers = {
        ...options.headers,
        'Authorization': `Bearer ${token}`
    };

    const response = await fetch(url, {
        ...options,
        headers
    });

    // If unauthorized, redirect to login
    if (response.status === 401) {
        localStorage.removeItem('authToken');
        localStorage.removeItem('username');
        window.location.href = '/login';
        throw new Error('Session expired');
    }

    return response;
}

// Logout
function logout() {
    const token = getAuthToken();

    // Call logout endpoint
    if (token) {
        fetch('/auth/logout', {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        }).catch(() => { });
    }

    // Clear local storage
    localStorage.removeItem('authToken');
    localStorage.removeItem('username');

    // Redirect to login
    window.location.href = '/login';
}

// Get username
function getUsername() {
    return localStorage.getItem('username') || 'User';
}

// Add logout button to page
function addLogoutButton() {
    const username = getUsername();

    const logoutDiv = document.createElement('div');
    logoutDiv.className = 'fixed top-4 right-4 z-50';
    logoutDiv.innerHTML = `
        <div class="bg-slate-800/90 backdrop-blur-lg rounded-lg shadow-lg p-3 border border-slate-700/50">
            <div class="flex items-center gap-3">
                <div class="text-right">
                    <div class="text-xs text-gray-400">Logged in as</div>
                    <div class="font-semibold text-white">${username}</div>
                </div>
                <button onclick="logout()" 
                    class="bg-red-600 hover:bg-red-700 px-4 py-2 rounded-lg transition-all flex items-center gap-2">
                    <i class="fas fa-sign-out-alt"></i>
                    <span>Logout</span>
                </button>
            </div>
        </div>
    `;

    document.body.appendChild(logoutDiv);
}
