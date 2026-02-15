// Global state
let currentDownload = null;

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

// Format file size
function formatFileSize(bytes) {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// Format speed
function formatSpeed(bytesPerSecond) {
    return formatFileSize(bytesPerSecond) + '/s';
}

// Format time
function formatTime(seconds) {
    if (seconds < 60) return seconds.toFixed(1) + 's';
    const minutes = Math.floor(seconds / 60);
    const secs = Math.floor(seconds % 60);
    return `${minutes}m ${secs}s`;
}

// Refresh files list
async function refreshFiles() {
    const path = document.getElementById('currentPath').value;
    const filesList = document.getElementById('filesList');

    filesList.innerHTML = '<div class="text-center text-gray-400 py-4"><i class="fas fa-spinner fa-spin text-2xl"></i><p class="mt-2">Loading...</p></div>';

    try {
        const response = await fetch(`/api/files/view?dir=${encodeURIComponent(path)}`);

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        displayFiles(data);
        showToast('Success', 'Files loaded successfully', 'success');
    } catch (error) {
        filesList.innerHTML = `<div class="text-center text-red-400 py-4"><i class="fas fa-exclamation-circle text-2xl"></i><p class="mt-2">Error: ${error.message}</p></div>`;
        showToast('Error', 'Failed to load files', 'error');
    }
}

// Display files in the list
function displayFiles(files) {
    const filesList = document.getElementById('filesList');

    if (!files || files.length === 0) {
        filesList.innerHTML = '<div class="text-center text-gray-400 py-4"><i class="fas fa-folder-open text-2xl"></i><p class="mt-2">Empty directory</p></div>';
        return;
    }

    filesList.innerHTML = files.map(file => {
        const icon = file.IsDir
            ? 'fa-folder text-yellow-400'
            : 'fa-file text-blue-400';
        const size = file.IsDir ? '-' : formatFileSize(file.Size);
        const action = file.IsDir
            ? `onclick="navigateToDir('${file.Path}')"`
            : `onclick="selectFile('${file.Path}')"`;

        return `
            <div ${action} class="flex items-center justify-between p-3 hover:bg-slate-700/50 rounded-lg cursor-pointer transition-all group">
                <div class="flex items-center gap-3 flex-1 min-w-0">
                    <i class="fas ${icon} text-xl"></i>
                    <div class="flex-1 min-w-0">
                        <div class="font-medium truncate group-hover:text-blue-400 transition-colors">${file.Name}</div>
                        <div class="text-xs text-gray-400">${size}</div>
                    </div>
                </div>
                <div class="flex items-center gap-2">
                    ${!file.IsDir ? `
                        <button onclick="event.stopPropagation(); selectFileForShare('${file.Path}')" 
                            class="bg-purple-600 hover:bg-purple-700 px-3 py-1 rounded text-sm transition-all opacity-0 group-hover:opacity-100">
                            <i class="fas fa-share-nodes"></i>
                        </button>
                    ` : ''}
                    <i class="fas fa-chevron-right text-gray-600 group-hover:text-blue-400 transition-colors"></i>
                </div>
            </div>
        `;
    }).join('');
}

// Navigate to directory
function navigateToDir(path) {
    document.getElementById('currentPath').value = path;
    refreshFiles();
}

// Navigate using the input button
function navigateTo() {
    refreshFiles();
}

// Select file for sharing
function selectFileForShare(path) {
    document.getElementById('shareFilePath').value = path;
    showToast('File Selected', 'Click "Generate Share Link" to create a temporary link', 'info');

    // Scroll to share section
    document.getElementById('shareFilePath').scrollIntoView({ behavior: 'smooth', block: 'center' });
}

// Generate share link
async function generateShareLink() {
    const filePath = document.getElementById('shareFilePath').value.trim();

    if (!filePath) {
        showToast('Error', 'Please enter a file path', 'error');
        return;
    }

    const linkResult = document.getElementById('shareLinkResult');
    linkResult.classList.add('hidden');

    try {
        const response = await fetch(`/files/d/r?file_path=${encodeURIComponent(filePath)}`);

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();

        // Build full URL
        const fullUrl = `${window.location.origin}${data.share_link}`;
        document.getElementById('generatedLink').value = fullUrl;

        linkResult.classList.remove('hidden');
        showToast('Success', 'Share link generated successfully!', 'success');

        // Scroll to result
        linkResult.scrollIntoView({ behavior: 'smooth', block: 'center' });
    } catch (error) {
        showToast('Error', `Failed to generate link: ${error.message}`, 'error');
    }
}

// Copy to clipboard
async function copyToClipboard() {
    const linkInput = document.getElementById('generatedLink');
    const text = linkInput.value;

    try {
        await navigator.clipboard.writeText(text);
        showToast('Copied!', 'Link copied to clipboard', 'success');
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
            showToast('Copied!', 'Link copied to clipboard', 'success');
        } catch (err) {
            // Show manual copy dialog
            showManualCopyDialog(text, 'Share Link');
        }

        document.body.removeChild(textarea);
    }
}

// Test download with speed monitoring
async function testDownload() {
    const url = document.getElementById('generatedLink').value;

    if (!url) {
        showToast('Error', 'No share link available', 'error');
        return;
    }

    const statsPanel = document.getElementById('downloadStats');
    statsPanel.classList.remove('hidden');
    statsPanel.scrollIntoView({ behavior: 'smooth', block: 'center' });

    // Reset stats
    document.getElementById('downloadSpeed').textContent = '0 MB/s';
    document.getElementById('downloadSize').textContent = '0 MB';
    document.getElementById('downloadProgress').textContent = '0%';
    document.getElementById('downloadTime').textContent = '0s';
    document.getElementById('progressBar').style.width = '0%';
    document.getElementById('downloadStatus').textContent = 'Initializing download...';

    const startTime = Date.now();
    let lastLoaded = 0;
    let lastTime = startTime;

    try {
        const response = await fetch(url);

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const contentLength = response.headers.get('content-length');
        const total = parseInt(contentLength, 10);

        document.getElementById('downloadStatus').textContent = `Downloading... (Total: ${formatFileSize(total)})`;

        const reader = response.body.getReader();
        let loaded = 0;
        const chunks = [];

        // Update interval for smooth stats
        const updateInterval = setInterval(() => {
            const now = Date.now();
            const elapsed = (now - startTime) / 1000;
            const timeDiff = (now - lastTime) / 1000;
            const loadedDiff = loaded - lastLoaded;

            // Calculate current speed
            const speed = timeDiff > 0 ? loadedDiff / timeDiff : 0;

            // Update UI
            document.getElementById('downloadSpeed').textContent = formatSpeed(speed);
            document.getElementById('downloadSize').textContent = formatFileSize(loaded);
            document.getElementById('downloadTime').textContent = formatTime(elapsed);

            if (total) {
                const progress = Math.round((loaded / total) * 100);
                document.getElementById('downloadProgress').textContent = progress + '%';
                document.getElementById('progressBar').style.width = progress + '%';
            }

            lastLoaded = loaded;
            lastTime = now;
        }, 100);

        while (true) {
            const { done, value } = await reader.read();

            if (done) break;

            chunks.push(value);
            loaded += value.length;
        }

        clearInterval(updateInterval);

        // Final update
        const totalTime = (Date.now() - startTime) / 1000;
        const avgSpeed = loaded / totalTime;

        document.getElementById('downloadSpeed').textContent = formatSpeed(avgSpeed);
        document.getElementById('downloadSize').textContent = formatFileSize(loaded);
        document.getElementById('downloadProgress').textContent = '100%';
        document.getElementById('progressBar').style.width = '100%';
        document.getElementById('downloadTime').textContent = formatTime(totalTime);
        document.getElementById('downloadStatus').innerHTML = `
            <div class="text-green-400 font-semibold">
                <i class="fas fa-check-circle mr-2"></i>Download completed successfully!
                <div class="text-sm text-gray-400 mt-2">
                    Average speed: ${formatSpeed(avgSpeed)} | Total time: ${formatTime(totalTime)}
                </div>
            </div>
        `;

        // Trigger actual download
        const blob = new Blob(chunks);
        const downloadUrl = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = downloadUrl;
        a.download = url.split('/').pop() || 'download';
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        window.URL.revokeObjectURL(downloadUrl);

        showToast('Complete!', 'Download finished successfully', 'success');
    } catch (error) {
        document.getElementById('downloadStatus').innerHTML = `
            <div class="text-red-400 font-semibold">
                <i class="fas fa-exclamation-circle mr-2"></i>Download failed: ${error.message}
            </div>
        `;
        showToast('Error', 'Download failed', 'error');
    }
}

// Initialize on page load
document.addEventListener('DOMContentLoaded', () => {
    // Load files from root by default
    refreshFiles();

    // Allow Enter key in path input
    document.getElementById('currentPath').addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            navigateTo();
        }
    });

    // Allow Enter key in share path input
    document.getElementById('shareFilePath').addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            generateShareLink();
        }
    });
});
