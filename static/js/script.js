async function fetchMessage() {
    try {
        const response = await fetch('/api/message');
        const data = await response.json();
        document.getElementById('api-message').textContent = data.message;
    } catch (error) {
        console.error('Error fetching message:', error);
    }
}