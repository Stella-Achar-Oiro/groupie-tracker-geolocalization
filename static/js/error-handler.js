// error-handler.js

function displayError(container, statusCode, message) {
    container.innerHTML = `
        <div class="error-container">
            <h1>${statusCode}</h1>
            <p class="error-message">${message}</p>
            <p class="error-description">${getStatusText(statusCode)}</p>
            <button class="button" onclick="window.location.href='/'">
                Back to Home
            </button>
        </div>
    `;
}

function getStatusText(code) {
    const statusTexts = {
        400: 'Bad Request',
        401: 'Unauthorized',
        403: 'Forbidden',
        404: 'Not Found',
        408: 'Request Timeout',
        429: 'Too Many Requests',
        500: 'Internal Server Error',
        502: 'Bad Gateway',
        503: 'Service Unavailable',
        504: 'Gateway Timeout'
    };
    return statusTexts[code] || 'Unknown Error';
}

function isOffline() {
    return !navigator.onLine;
}

function createOfflineUI(container) {
    container.innerHTML = `
        <div class="offline-container">
            <div class="offline-icon">
                <i class="fas fa-wifi-slash"></i>
            </div>
            <h2>No Internet Connection</h2>
            <p>Please check your internet connection and try again.</p>
            <button class="retry-button" onclick="window.location.reload()">
                <i class="fas fa-redo"></i> Try Again
            </button>
        </div>
    `;
}

function handleError(error, container) {
    if (isOffline()) {
        createOfflineUI(container);
        return;
    }

    let statusCode = 500;
    let message = 'An unexpected error occurred';

    if (error.response) {
        statusCode = error.response.status;
        switch (statusCode) {
            case 404:
                message = 'The requested resource was not found';
                break;
            case 429:
                message = 'Too many requests. Please try again later';
                break;
            case 400:
                message = 'Invalid request. Please check your input';
                break;
            case 401:
                message = 'You need to be authenticated to access this resource';
                break;
            case 403:
                message = 'You don\'t have permission to access this resource';
                break;
            case 408:
                message = 'Request timed out. Please try again';
                break;
            default:
                message = error.response.data?.message || 'Server error occurred';
        }
    } else if (error.request) {
        statusCode = 503;
        message = 'Unable to reach the server. Please check your connection';
    } else {
        message = error.message || 'An unexpected error occurred';
    }

    container.innerHTML = `
        <div class="error-container">
            <div class="error-icon">
                <i class="fas fa-exclamation-circle"></i>
            </div>
            <h2>Error ${statusCode}</h2>
            <p>${message}</p>
            <button class="retry-button" onclick="window.location.reload()">
                <i class="fas fa-redo"></i> Try Again
            </button>
        </div>
    `;

    hideLoading();
    displayError(container, statusCode, message);
}

function detectSlowConnection() {
    const connection = navigator.connection || 
                      navigator.mozConnection || 
                      navigator.webkitConnection;
                      
    if (connection) {
        if (connection.effectiveType === 'slow-2g' || connection.effectiveType === '2g') {
            showSlowConnectionWarning();
        }
    }
}

function showSlowConnectionWarning() {
    const warning = document.createElement('div');
    warning.className = 'loading-indicator';
    warning.innerHTML = `
        <i class="fas fa-spinner"></i>
        <span>Slow connection detected</span>
    `;
    document.body.appendChild(warning);
    
    setTimeout(() => {
        warning.remove();
    }, 3000);
}