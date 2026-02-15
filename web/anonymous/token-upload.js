// Global state
let uploadToken = null;
let tokenInfo = null;
const uploadQueue = new Map();
let completedUploads = 0;
let totalBytesUploaded = 0;

// Constants
const CHUNK_SIZE = 50 * 1024 * 1024; // 50MB chunks
const CHUNKED_THRESHOLD = 100 * 1024 * 1024; // 100MB

// Check for token in URL
document.addEventListener('DOMContentLoaded', () => {
    const urlParams = new URLSearchParams(window.location.search);
    const token = urlParams.get('token');

    if (token) {
        document.getElementById('tokenInput').value = token;
        validateToken();
    }

    setupUploadHandlers();
});

// Setup upload handlers
function setupUploadHandlers() {
    const dropZone = document.getElementById('dropZone');
    const fileInput = document.getElementById('fileInput');

    dropZone.addEventListener('click', () => {
        fileInput.click();
    });

    fileInput.addEventListener('change', (e) => {
        handleFiles(Array.from(e.target.files));
        fileInput.value = '';
    });

    dropZone.addEventListener('dragover', (e) => {
        e.preventDefault();
        dropZone.classList.add('border-blue-500', 'bg-blue-500/10');
    });

    dropZone.addEventListener('dragleave', (e) => {
        e.preventDefault();
        dropZone.classList.remove('border-blue-500', 'bg-blue-500/10');
    });

    dropZone.addEventListener('drop', (e) => {
        e.preventDefault();
        dropZone.classList.remove('border-blue-500', 'bg-blue-500/10');

        const files = Array.from(e.dataTransfer.files);
        handleFiles(files);
    });
}

// Validate token
async function validateToken() {
    const token = document.getElementById('tokenInput').value.trim();
    const errorDiv = document.getElementById('tokenError');

    if (!token) {
        errorDiv.textContent = 'Please enter a token';
        errorDiv.classList.remove('hidden');
        return;
    }

    try {
        const response = await fetch(`/upload/token/validate?token=${encodeURIComponent(token)}`);
        const data = await response.json();

        if (!data.valid) {
            errorDiv.textContent = data.error || 'Invalid token';
            errorDiv.classList.remove('hidden');
            return;
        }

        // Token is valid
        uploadToken = token;
        tokenInfo = data;

        // Show upload section
        document.getElementById('tokenSection').classList.add('hidden');
        document.getElementById('uploadSection').classList.remove('hidden');

        // Display token info
        document.getElementById('tokenTargetDir').textContent = data.target_dir;
        document.getElementById('tokenRemaining').textContent = data.uploads_remaining;
        document.getElementById('tokenExpiry').textContent = new Date(data.expires_at).toLocaleString();

        showToast('Success', 'Token validated! You can now upload files.', 'success');
    } catch (error) {
        errorDiv.textContent = 'Failed to validate token: ' + error.message;
        errorDiv.classList.remove('hidden');
    }
}

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
    }, 5000);
}

function closeToast() {
    document.getElementById('toast').classList.add('hidden');
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

// Handle selected files
function handleFiles(files) {
    if (files.length === 0) return;

    if (!uploadToken) {
        showToast('Error', 'Please validate your token first', 'error');
        return;
    }

    // Check remaining uploads
    if (files.length > tokenInfo.uploads_remaining) {
        showToast('Error', `Token only allows ${tokenInfo.uploads_remaining} more upload(s)`, 'error');
        return;
    }

    document.getElementById('statsPanel').classList.remove('hidden');

    files.forEach(file => {
        const uploadId = Date.now() + '_' + file.name;
        const useChunked = document.getElementById('useChunked').checked || file.size > CHUNKED_THRESHOLD;

        const uploadItem = {
            id: uploadId,
            file: file,
            useChunked: useChunked,
            progress: 0,
            speed: 0,
            uploaded: 0,
            status: 'pending',
            startTime: null
        };

        uploadQueue.set(uploadId, uploadItem);
        createUploadCard(uploadItem);

        setTimeout(() => startUpload(uploadId), 100);
    });

    updateStats();
}

// Create upload card UI
function createUploadCard(item) {
    const queueContainer = document.getElementById('uploadQueue');

    const card = document.createElement('div');
    card.id = `upload-${item.id}`;
    card.className = 'bg-slate-800/50 backdrop-blur-lg rounded-xl shadow-lg border border-slate-700/50 p-6';

    card.innerHTML = `
        <div class="flex items-start justify-between mb-4">
            <div class="flex items-center gap-3 flex-1 min-w-0">
                <i class="fas fa-file text-3xl text-blue-400"></i>
                <div class="flex-1 min-w-0">
                    <div class="font-semibold truncate">${item.file.name}</div>
                    <div class="text-sm text-gray-400">${formatFileSize(item.file.size)} • ${item.useChunked ? 'Chunked Upload' : 'Direct Upload'}</div>
                </div>
            </div>
            <div class="flex items-center gap-2">
                <span id="status-${item.id}" class="px-3 py-1 rounded-full text-xs font-semibold bg-yellow-600/20 text-yellow-400">
                    <i class="fas fa-clock mr-1"></i>Pending
                </span>
            </div>
        </div>

        <div class="mb-4">
            <div class="flex justify-between text-sm mb-2">
                <span class="text-gray-400">Progress</span>
                <span id="progress-text-${item.id}" class="font-semibold">0%</span>
            </div>
            <div class="bg-slate-700 rounded-full h-3 overflow-hidden">
                <div id="progress-bar-${item.id}" class="bg-gradient-to-r from-blue-500 to-purple-500 h-full transition-all duration-300" style="width: 0%"></div>
            </div>
        </div>

        <div class="grid grid-cols-3 gap-4 text-center">
            <div>
                <div class="text-xs text-gray-400">Speed</div>
                <div id="speed-${item.id}" class="font-semibold text-blue-400">0 MB/s</div>
            </div>
            <div>
                <div class="text-xs text-gray-400">Uploaded</div>
                <div id="uploaded-${item.id}" class="font-semibold text-purple-400">0 MB</div>
            </div>
            <div>
                <div class="text-xs text-gray-400">Remaining</div>
                <div id="remaining-${item.id}" class="font-semibold text-gray-400">${formatFileSize(item.file.size)}</div>
            </div>
        </div>

        <div id="error-${item.id}" class="hidden mt-4 p-3 bg-red-900/30 border border-red-700 rounded-lg text-red-400 text-sm"></div>
    `;

    queueContainer.appendChild(card);
}

// Update upload card
function updateUploadCard(item) {
    const progressText = document.getElementById(`progress-text-${item.id}`);
    const progressBar = document.getElementById(`progress-bar-${item.id}`);
    const speed = document.getElementById(`speed-${item.id}`);
    const uploaded = document.getElementById(`uploaded-${item.id}`);
    const remaining = document.getElementById(`remaining-${item.id}`);
    const status = document.getElementById(`status-${item.id}`);

    if (progressText) progressText.textContent = item.progress + '%';
    if (progressBar) progressBar.style.width = item.progress + '%';
    if (speed) speed.textContent = formatSpeed(item.speed);
    if (uploaded) uploaded.textContent = formatFileSize(item.uploaded);
    if (remaining) remaining.textContent = formatFileSize(item.file.size - item.uploaded);

    if (status) {
        if (item.status === 'uploading') {
            status.className = 'px-3 py-1 rounded-full text-xs font-semibold bg-blue-600/20 text-blue-400';
            status.innerHTML = '<i class="fas fa-spinner fa-spin mr-1"></i>Uploading';
        } else if (item.status === 'completed') {
            status.className = 'px-3 py-1 rounded-full text-xs font-semibold bg-green-600/20 text-green-400';
            status.innerHTML = '<i class="fas fa-check-circle mr-1"></i>Completed';
        } else if (item.status === 'error') {
            status.className = 'px-3 py-1 rounded-full text-xs font-semibold bg-red-600/20 text-red-400';
            status.innerHTML = '<i class="fas fa-exclamation-circle mr-1"></i>Failed';
        }
    }
}

// Start upload
async function startUpload(uploadId) {
    const item = uploadQueue.get(uploadId);
    if (!item) return;

    item.status = 'uploading';
    item.startTime = Date.now();
    updateUploadCard(item);

    try {
        if (item.useChunked) {
            await uploadChunked(item);
        } else {
            await uploadDirect(item);
        }

        item.status = 'completed';
        item.progress = 100;
        completedUploads++;

        // Update token info
        if (tokenInfo) {
            tokenInfo.uploads_remaining--;
            document.getElementById('tokenRemaining').textContent = tokenInfo.uploads_remaining;
        }

        updateUploadCard(item);
        updateStats();
        showToast('Success', `${item.file.name} uploaded successfully!`, 'success');
    } catch (error) {
        item.status = 'error';
        updateUploadCard(item);

        const errorDiv = document.getElementById(`error-${item.id}`);
        if (errorDiv) {
            errorDiv.textContent = `Error: ${error.message}`;
            errorDiv.classList.remove('hidden');
        }

        showToast('Upload Failed', `${item.file.name}: ${error.message}`, 'error');
    }
}

// Direct upload
async function uploadDirect(item) {
    const formData = new FormData();
    formData.append('file', item.file);

    const xhr = new XMLHttpRequest();

    return new Promise((resolve, reject) => {
        let lastTime = Date.now();
        let lastLoaded = 0;

        xhr.upload.addEventListener('progress', (e) => {
            if (e.lengthComputable) {
                const now = Date.now();
                const timeDiff = (now - lastTime) / 1000;
                const loadedDiff = e.loaded - lastLoaded;

                item.uploaded = e.loaded;
                item.progress = Math.round((e.loaded / e.total) * 100);
                item.speed = timeDiff > 0 ? loadedDiff / timeDiff : 0;

                totalBytesUploaded += loadedDiff;
                lastLoaded = e.loaded;
                lastTime = now;

                updateUploadCard(item);
                updateStats();
            }
        });

        xhr.addEventListener('load', () => {
            if (xhr.status >= 200 && xhr.status < 300) {
                resolve();
            } else {
                reject(new Error(`Server returned ${xhr.status}`));
            }
        });

        xhr.addEventListener('error', () => reject(new Error('Network error')));
        xhr.addEventListener('abort', () => reject(new Error('Upload cancelled')));

        xhr.open('POST', `/upload/token/file?token=${encodeURIComponent(uploadToken)}`);
        xhr.send(formData);
    });
}

// Chunked upload
async function uploadChunked(item) {
    const file = item.file;
    const totalChunks = Math.ceil(file.size / CHUNK_SIZE);

    // Initialize
    const initResponse = await fetch(
        `/upload/token/chunk/start?token=${encodeURIComponent(uploadToken)}&filename=${encodeURIComponent(file.name)}`,
        { method: 'POST' }
    );

    if (!initResponse.ok) {
        throw new Error('Failed to initialize chunked upload');
    }

    // Upload chunks
    for (let chunkIndex = 0; chunkIndex < totalChunks; chunkIndex++) {
        const start = chunkIndex * CHUNK_SIZE;
        const end = Math.min(start + CHUNK_SIZE, file.size);
        const chunk = file.slice(start, end);

        const startTime = Date.now();

        const response = await fetch(
            `/upload/token/chunk?token=${encodeURIComponent(uploadToken)}&filename=${encodeURIComponent(file.name)}&chunk_index=${chunkIndex}`,
            {
                method: 'POST',
                body: chunk
            }
        );

        if (!response.ok) {
            throw new Error(`Failed to upload chunk ${chunkIndex + 1}/${totalChunks}`);
        }

        const elapsed = (Date.now() - startTime) / 1000;
        const chunkSize = end - start;

        item.uploaded = end;
        item.progress = Math.round((end / file.size) * 100);
        item.speed = elapsed > 0 ? chunkSize / elapsed : 0;

        totalBytesUploaded += chunkSize;
        updateUploadCard(item);
        updateStats();
    }

    // Complete
    const completeResponse = await fetch(
        `/upload/token/chunk/complete?token=${encodeURIComponent(uploadToken)}&filename=${encodeURIComponent(file.name)}`,
        { method: 'POST' }
    );

    if (!completeResponse.ok) {
        throw new Error('Failed to complete chunked upload');
    }
}

// Update global stats
function updateStats() {
    const queueCount = uploadQueue.size;
    const activeUploads = Array.from(uploadQueue.values()).filter(item => item.status === 'uploading');
    const totalSpeed = activeUploads.reduce((sum, item) => sum + item.speed, 0);

    document.getElementById('totalSpeed').textContent = formatSpeed(totalSpeed);
    document.getElementById('totalUploaded').textContent = formatFileSize(totalBytesUploaded);
    document.getElementById('queueCount').textContent = queueCount;
    document.getElementById('completedCount').textContent = completedUploads;
}
