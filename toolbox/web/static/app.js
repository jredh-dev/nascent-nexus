// Toolbox Frontend JavaScript

const conversationEl = document.getElementById('conversation');
const messageForm = document.getElementById('messageForm');
const messageInput = document.getElementById('messageInput');
const sendBtn = document.getElementById('sendBtn');
const resetBtn = document.getElementById('resetBtn');

// Add message to conversation
function addMessage(role, content) {
    const messageDiv = document.createElement('div');
    messageDiv.className = `message ${role}`;
    
    const roleLabel = document.createElement('div');
    roleLabel.className = 'message-role';
    roleLabel.textContent = role;
    
    const contentDiv = document.createElement('div');
    contentDiv.className = 'message-content';
    contentDiv.textContent = content;
    
    messageDiv.appendChild(roleLabel);
    messageDiv.appendChild(contentDiv);
    conversationEl.appendChild(messageDiv);
    
    // Scroll to bottom
    conversationEl.scrollTop = conversationEl.scrollHeight;
}

// Show loading indicator
function showLoading() {
    const loadingDiv = document.createElement('div');
    loadingDiv.id = 'loading-indicator';
    loadingDiv.className = 'message assistant';
    loadingDiv.innerHTML = `
        <div class="message-role">assistant</div>
        <div class="message-content">
            <div class="loading"></div> Thinking...
        </div>
    `;
    conversationEl.appendChild(loadingDiv);
    conversationEl.scrollTop = conversationEl.scrollHeight;
}

// Remove loading indicator
function removeLoading() {
    const loadingDiv = document.getElementById('loading-indicator');
    if (loadingDiv) {
        loadingDiv.remove();
    }
}

// Send message to server
async function sendMessage(message) {
    try {
        const response = await fetch('/api/message', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ message }),
        });
        
        const data = await response.json();
        
        if (!response.ok) {
            throw new Error(data.error || 'Failed to send message');
        }
        
        return data.response;
    } catch (error) {
        console.error('Error sending message:', error);
        throw error;
    }
}

// Reset conversation
async function resetConversation() {
    try {
        await fetch('/api/reset', {
            method: 'POST',
        });
        
        // Clear UI
        conversationEl.innerHTML = `
            <div class="message system">
                <div class="message-content">
                    Conversation reset. What would you like to do?
                </div>
            </div>
        `;
    } catch (error) {
        console.error('Error resetting conversation:', error);
        alert('Failed to reset conversation');
    }
}

// Handle form submission
messageForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const message = messageInput.value.trim();
    if (!message) return;
    
    // Disable input
    messageInput.disabled = true;
    sendBtn.disabled = true;
    
    // Add user message to UI
    addMessage('user', message);
    
    // Clear input
    messageInput.value = '';
    
    // Show loading
    showLoading();
    
    try {
        // Send to server
        const response = await sendMessage(message);
        
        // Remove loading
        removeLoading();
        
        // Add assistant response
        addMessage('assistant', response);
    } catch (error) {
        // Remove loading
        removeLoading();
        
        // Show error
        addMessage('system', `Error: ${error.message}`);
    } finally {
        // Re-enable input
        messageInput.disabled = false;
        sendBtn.disabled = false;
        messageInput.focus();
    }
});

// Handle reset button
resetBtn.addEventListener('click', async () => {
    if (confirm('Are you sure you want to reset the conversation?')) {
        await resetConversation();
        messageInput.focus();
    }
});

// Focus input on load
messageInput.focus();
