// Show toast notification
function showToast(title, message, type = 'info') {
    const toast = document.getElementById('toast');
    const toastIcon = document.getElementById('toastIcon');
    const toastTitle = document.getElementById('toastTitle');
    const toastMessage = document.getElementById('toastMessage');

    const icons = {
        success: 'fa-check-circle text-green-400',
        error: 'fa-exclamation-circle text-red-400',
        warning: 'fa-exclamation-triangle text-yellow-400',
        info: 'fa-info-circle text-blue-400'
    };

    toastIcon.className = `fas ${icons[type]} text-2xl`;
    toastTitle.textContent = title;
    toastMessage.textContent = message;

    toast.classList.remove('hidden');
    setTimeout(() => {
        toast.classList.add('hidden');
    }, 4000);
}

// Show manual copy dialog
function showManualCopyDialog(text, label) {
    const overlay = document.createElement('div');
    overlay.className = 'fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4';
    overlay.innerHTML = `
        <div class="bg-slate-800 rounded-2xl shadow-2xl border border-slate-700 max-w-2xl w-full p-6">
            <div class="flex items-center justify-between mb-4">
                <h3 class="text-xl font-bold flex items-center">
                    <i class="fas fa-copy mr-2 text-blue-400"></i>Manual Copy Required
                </h3>
                <button onclick="this.closest('.fixed').remove()" class="text-gray-400 hover:text-white">
                    <i class="fas fa-times text-xl"></i>
                </button>
            </div>
            <p class="text-gray-400 mb-4">Please select and copy the ${label} below:</p>
            <textarea readonly class="w-full bg-slate-700/50 border border-slate-600 rounded-lg px-4 py-3 font-mono text-sm h-32 focus:outline-none focus:border-blue-500 focus:ring-2 focus:ring-blue-500/50" id="manualCopyText">${text}</textarea>
            <div class="flex gap-3 mt-4">
                <button onclick="document.getElementById('manualCopyText').select();" class="flex-1 bg-blue-600 hover:bg-blue-700 py-2 rounded-lg font-semibold transition-all">
                    <i class="fas fa-mouse-pointer mr-2"></i>Select All
                </button>
                <button onclick="this.closest('.fixed').remove()" class="px-6 bg-slate-700 hover:bg-slate-600 py-2 rounded-lg font-semibold transition-all">
                    Close
                </button>
            </div>
        </div>
    `;
    document.body.appendChild(overlay);
    
    // Auto-select the text
    setTimeout(() => {
        const textarea = document.getElementById('manualCopyText');
        textarea.focus();
        textarea.select();
    }, 100);
}

// Generate new token
async function generateToken() {
    const targetDir = document.getElementById('targetDir').value.trim();
    const maxUploads = parseInt(document.getElementById('maxUploads').value);
    const expiryHours = parseInt(document.getElementById('expiryHours').value);

    if (!targetDir) {
        showToast('Error', 'Please enter a target directory', 'error');
        return;
    }

    console.log("Generating token");

    try {
        const response = await authenticatedFetch(
            `/api/upload-tokens/generate?target_dir=${encodeURIComponent(targetDir)}&max_uploads=${maxUploads}&expiry_hours=${expiryHours}`,
            { method: 'POST' }
        );

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();

        // Display results
        const protocol = window.location.protocol;
        const host = window.location.host;
        const uploadURL = `${protocol}//${data.upload_url}`;

        document.getElementById('generatedToken').value = data.token;
        document.getElementById('generatedURL').value = uploadURL;
        document.getElementById('resultTargetDir').textContent = data.target_dir;
        document.getElementById('resultMaxUploads').textContent = data.max_uploads;

        document.getElementById('tokenResult').classList.remove('hidden');

        showToast('Success', 'Token generated successfully!', 'success');

        // Refresh token list
        setTimeout(refreshTokens, 500);
    } catch (error) {
        showToast('Error', `Failed to generate token: ${error.message}`, 'error');
    }
}

// Copy token to clipboard
async function copyToken() {
    const tokenInput = document.getElementById('generatedToken');
    const text = tokenInput.value;
    
    try {
        await navigator.clipboard.writeText(text);
        showToast('Copied!', 'Token copied to clipboard', 'success');
    } catch (error) {
        // Fallback for mobile devices
        const textarea = document.createElement('textarea');
        textarea.value = text;
        textarea.style.position = 'fixed';
        textarea.style.left = '-999999px';
        textarea.style.top = '-999999px';
        document.body.appendChild(textarea);
        textarea.focus();
        textarea.select();
        
        try {
            document.execCommand('copy');
            showToast('Copied!', 'Token copied to clipboard', 'success');
        } catch (err) {
            // Show manual copy dialog
            showManualCopyDialog(text, 'Token');
        }
        
        document.body.removeChild(textarea);
    }
}

// Copy URL to clipboard
async function copyURL() {
    const urlInput = document.getElementById('generatedURL');
    const text = urlInput.value;
    
    try {
        await navigator.clipboard.writeText(text);
        showToast('Copied!', 'URL copied to clipboard', 'success');
    } catch (error) {
        // Fallback for mobile devices
        const textarea = document.createElement('textarea');
        textarea.value = text;
        textarea.style.position = 'fixed';
        textarea.style.left = '-999999px';
        textarea.style.top = '-999999px';
        document.body.appendChild(textarea);
        textarea.focus();
        textarea.select();
        
        try {
            document.execCommand('copy');
            showToast('Copied!', 'URL copied to clipboard', 'success');
        } catch (err) {
            // Show manual copy dialog
            showManualCopyDialog(text, 'Upload URL');
        }
        
        document.body.removeChild(textarea);
    }
}

// Refresh tokens list
async function refreshTokens() {
    const tokensList = document.getElementById('tokensList');
    tokensList.innerHTML = '<div class="text-center text-gray-400 py-4"><i class="fas fa-spinner fa-spin text-2xl"></i><p class="mt-2">Loading...</p></div>';

    try {
        const response = await authenticatedFetch('/api/upload-tokens/list');

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        displayTokens(data.tokens);
        showToast('Success', 'Tokens refreshed', 'success');
    } catch (error) {
        tokensList.innerHTML = `<div class="text-center text-red-400 py-4"><i class="fas fa-exclamation-circle text-2xl"></i><p class="mt-2">Error: ${error.message}</p></div>`;
        showToast('Error', 'Failed to load tokens', 'error');
    }
}

// Display tokens
function displayTokens(tokens) {
    const tokensList = document.getElementById('tokensList');

    if (!tokens || tokens.length === 0) {
        tokensList.innerHTML = '<div class="text-center text-gray-400 py-4"><i class="fas fa-key text-2xl"></i><p class="mt-2">No active tokens</p></div>';
        return;
    }

    tokensList.innerHTML = tokens.map(token => {
        const isExpired = new Date(token.ExpiresAt) < new Date();
        const isRevoked = token.IsRevoked;
        const isActive = !isExpired && !isRevoked;
        const uploadsRemaining = token.MaxUploads - token.UploadCount;

        let statusClass = 'bg-green-600/20 text-green-400';
        let statusText = 'Active';
        let statusIcon = 'fa-check-circle';

        if (isRevoked) {
            statusClass = 'bg-red-600/20 text-red-400';
            statusText = 'Revoked';
            statusIcon = 'fa-ban';
        } else if (isExpired) {
            statusClass = 'bg-yellow-600/20 text-yellow-400';
            statusText = 'Expired';
            statusIcon = 'fa-clock';
        } else if (uploadsRemaining === 0) {
            statusClass = 'bg-gray-600/20 text-gray-400';
            statusText = 'Used';
            statusIcon = 'fa-check';
        }

        return `
            <div class="bg-slate-700/30 rounded-lg p-4 hover:bg-slate-700/50 transition-all">
                <div class="flex items-start justify-between mb-3">
                    <div class="flex-1 min-w-0">
                        <div class="text-xs text-gray-400 mb-1">Token</div>
                        <div class="font-mono text-sm truncate">${token.Token}</div>
                    </div>
                    <span class="px-3 py-1 rounded-full text-xs font-semibold ${statusClass} ml-2">
                        <i class="fas ${statusIcon} mr-1"></i>${statusText}
                    </span>
                </div>

                <div class="grid grid-cols-2 gap-3 text-sm mb-3">
                    <div>
                        <div class="text-xs text-gray-400">Target Directory</div>
                        <div class="font-mono text-blue-400">${token.TargetDir}</div>
                    </div>
                    <div>
                        <div class="text-xs text-gray-400">Uploads</div>
                        <div class="font-bold">${token.UploadCount}/${token.MaxUploads}</div>
                    </div>
                    <div>
                        <div class="text-xs text-gray-400">Created</div>
                        <div>${new Date(token.CreatedAt).toLocaleString()}</div>
                    </div>
                    <div>
                        <div class="text-xs text-gray-400">Expires</div>
                        <div>${new Date(token.ExpiresAt).toLocaleString()}</div>
                    </div>
                </div>

                <div class="flex gap-2">
                    ${isActive ? `
                        <button onclick="copyUploadLink('${token.Token}')" 
                            class="flex-1 bg-blue-600 hover:bg-blue-700 px-3 py-2 rounded text-sm transition-all">
                            <i class="fas fa-link mr-1"></i>Copy Link
                        </button>
                        <button onclick="revokeToken('${token.Token}')" 
                            class="bg-red-600 hover:bg-red-700 px-3 py-2 rounded text-sm transition-all">
                            <i class="fas fa-ban"></i>
                        </button>
                    ` : ''}
                </div>
            </div>
        `;
    }).join('');
}

// Copy upload link
async function copyUploadLink(token) {
    const protocol = window.location.protocol;
    const host = window.location.host;
    const url = `${protocol}//${host}/upload?token=${token}`;

    try {
        await navigator.clipboard.writeText(url);
        showToast('Copied!', 'Upload link copied to clipboard', 'success');
    } catch (error) {
        // Fallback for mobile devices
        const textarea = document.createElement('textarea');
        textarea.value = url;
        textarea.style.position = 'fixed';
        textarea.style.left = '-999999px';
        textarea.style.top = '-999999px';
        document.body.appendChild(textarea);
        textarea.focus();
        textarea.select();
        
        try {
            document.execCommand('copy');
            showToast('Copied!', 'Upload link copied to clipboard', 'success');
        } catch (err) {
            // Show manual copy dialog
            showManualCopyDialog(url, 'Upload Link');
        }
        
        document.body.removeChild(textarea);
    }
}

// Revoke token
async function revokeToken(token) {
    if (!confirm('Are you sure you want to revoke this token?')) {
        return;
    }

    try {
        const response = await authenticatedFetch(`/api/upload-tokens/revoke?token=${encodeURIComponent(token)}`, {
            method: 'POST'
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        showToast('Success', 'Token revoked successfully', 'success');
        refreshTokens();
    } catch (error) {
        showToast('Error', `Failed to revoke token: ${error.message}`, 'error');
    }
}

// Load tokens on page load
document.addEventListener('DOMContentLoaded', () => {
    refreshTokens();
});
