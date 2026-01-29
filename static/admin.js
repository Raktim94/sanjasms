function showToast(message, type = 'success') {
    const container = document.getElementById('toast-container');
    if (!container) return;

    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    toast.innerHTML = `
        <span>${message}</span>
        <span style="cursor:pointer; margin-left:10px;" onclick="this.parentElement.remove()">×</span>
    `;

    container.appendChild(toast);

    setTimeout(() => {
        toast.style.animation = 'fadeOut 0.5s forwards';
        setTimeout(() => toast.remove(), 500);
    }, 4000);
}

// Global form submission handler for "Success" toasts
document.addEventListener('submit', function (e) {
    // We can't easily catch the redirect, but we can detect intent
    // In a real SPA we'd catch the response. For SSR, we might check URL params on load.
});

// Check for success/error in URL
window.addEventListener('load', () => {
    const params = new URLSearchParams(window.location.search);
    if (params.has('msg')) {
        showToast(params.get('msg'));
    }
    if (params.has('err')) {
        showToast(params.get('err'), 'error');
    }
});
