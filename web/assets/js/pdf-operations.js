// PDF Operations Functions
async function preparePDFData() {
    const arbiterSelect = document.getElementById('arbiterSelect');
    const leagueSelect = document.getElementById('leagueSelect');
    const prepareStatus = document.getElementById('prepareStatus');
    const preparePdfBtn = document.getElementById('preparePdfBtn');
    
    if (!arbiterSelect.value || !leagueSelect.value) {
        prepareStatus.innerHTML = '<span class="text-red-600">Please select both an arbiter and a league first.</span>';
        return;
    }
    
    preparePdfBtn.disabled = true;
    preparePdfBtn.textContent = 'Preparing...';
    prepareStatus.textContent = 'Preparing PDF data...';
    
    try {
        const response = await fetch('/prepare-pdf-data', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                arbiterId: parseInt(arbiterSelect.value),
                leagueId: parseInt(leagueSelect.value)
            })
        });
        
        const result = await response.json();
        
        if (response.ok) {
            prepareStatus.innerHTML = `
                <span class="text-green-600">✓ ${result.message}</span><br>
                <span class="text-sm">Data has been printed to the server console.</span>
            `;
        } else {
            prepareStatus.innerHTML = `<span class="text-red-600">✗ Error: ${result.error}</span>`;
        }
    } catch (error) {
        prepareStatus.innerHTML = `<span class="text-red-600">✗ Network error: ${error.message}</span>`;
    } finally {
        preparePdfBtn.disabled = false;
        preparePdfBtn.textContent = 'Prepare PDF Data';
    }
}

async function listFields() {
    const status = document.getElementById('status');
    status.textContent = 'Listing PDF fields...';
    
    try {
        const response = await fetch('/list-fields');
        const result = await response.json();
        
        if (response.ok) {
            status.textContent = '✓ Fields listed to console. Check server logs for details.';
        } else {
            status.textContent = `✗ Error: ${result.error}`;
        }
    } catch (error) {
        status.textContent = `✗ Network error: ${error.message}`;
    }
}

async function generatePDF() {
    const status = document.getElementById('status');
    status.textContent = 'Generating PDF...';
    
    try {
        const response = await fetch('/generate', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                // Add your PDF field data here
                test_field: 'Test Value'
            })
        });
        
        if (response.ok) {
            // Download the PDF
            const blob = await response.blob();
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = 'delegacny.pdf';
            document.body.appendChild(a);
            a.click();
            window.URL.revokeObjectURL(url);
            document.body.removeChild(a);
            
            status.textContent = '✓ PDF generated and downloaded!';
        } else {
            const result = await response.json();
            status.textContent = `✗ Error: ${result.error}`;
        }
    } catch (error) {
        status.textContent = `✗ Network error: ${error.message}`;
    }
}
