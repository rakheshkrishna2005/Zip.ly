document.addEventListener('DOMContentLoaded', function() {
    const shortenForm = document.getElementById('shortenForm');
    const result = document.getElementById('result');
    const shortUrlInput = document.getElementById('shortUrl');
    const copyButton = document.getElementById('copyButton');
    const expiryNotice = document.getElementById('expiryNotice');
    const createNewBtn = document.getElementById('createNewBtn');
    
    // Smooth scrolling for anchor links
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function (e) {
            e.preventDefault();
            
            const targetId = this.getAttribute('href');
            const targetElement = document.querySelector(targetId);
            
            if (targetElement) {
                window.scrollTo({
                    top: targetElement.offsetTop - 80,
                    behavior: 'smooth'
                });
            }
        });
    });
    
    // Form submission
    shortenForm.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        // Get form values
        const originalUrl = document.getElementById('originalUrl').value;
        const customAlias = document.getElementById('customAlias').value;
        const expiresIn = parseInt(document.getElementById('expiresIn').value, 10);
        
        // Validate custom alias length
        if (customAlias && (customAlias.length < 3 || customAlias.length > 10)) {
            alert("Custom alias must be between 3 and 10 characters");
            return;
        }
        
        try {
            // Show loading state
            const submitBtn = shortenForm.querySelector('button[type="submit"]');
            const originalBtnText = submitBtn.innerHTML;
            submitBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Processing...';
            submitBtn.disabled = true;
            
            // Prepare request payload
            const payload = {
                original_url: originalUrl,
                expires_in: expiresIn
            };
            
            // Add custom alias if provided
            if (customAlias) {
                payload.custom_alias = customAlias;
            }
            
            // Send request to API
            const response = await fetch('/api/v1/urls', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(payload)
            });
            
            // Reset button state
            submitBtn.innerHTML = originalBtnText;
            submitBtn.disabled = false;
            
            if (!response.ok) {
                const errorText = await response.text();
                let errorMessage;
                try {
                    // Try to parse as JSON
                    const errorData = JSON.parse(errorText);
                    errorMessage = errorData.message || errorText;
                } catch (e) {
                    // If not JSON, use the text directly
                    errorMessage = errorText;
                }
                throw new Error(errorMessage || `Error: ${response.status} ${response.statusText}`);
            }
            
            const data = await response.json();
            
            // Display the shortened URL
            shortUrlInput.value = data.short_url;
            
            // Set expiry notice if applicable
            if (data.expires_at) {
                const expiryDate = new Date(data.expires_at);
                expiryNotice.textContent = `This link will expire on ${expiryDate.toLocaleDateString()} at ${expiryDate.toLocaleTimeString()}.`;
            } else {
                expiryNotice.textContent = 'This link does not expire.';
            }
            
            // Hide form and show result
            shortenForm.closest('.url-card').style.display = 'none';
            result.classList.remove('hidden');
            
        } catch (error) {
            alert(`Failed to shorten URL: ${error.message}`);
            console.error('Error shortening URL:', error);
        }
    });
    
    // Copy to clipboard functionality
    copyButton.addEventListener('click', function() {
        shortUrlInput.select();
        document.execCommand('copy');
        
        // Visual feedback
        const originalText = copyButton.innerHTML;
        copyButton.innerHTML = '<i class="fas fa-check"></i> Copied!';
        copyButton.style.backgroundColor = 'var(--success-color)';
        
        setTimeout(() => {
            copyButton.innerHTML = originalText;
            copyButton.style.backgroundColor = '';
        }, 2000);
    });
    
    // Create new link button
    createNewBtn.addEventListener('click', function() {
        // Show form and hide result
        shortenForm.closest('.url-card').style.display = 'block';
        result.classList.add('hidden');
        
        // Reset form
        shortenForm.reset();
    });
});