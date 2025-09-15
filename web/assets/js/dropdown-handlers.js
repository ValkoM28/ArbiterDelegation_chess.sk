// Dropdown Population and Selection Handlers
async function populateArbiterDropdown() {
    const arbiterSelect = document.getElementById('arbiterSelect');
    
    try {
        const response = await fetch('/arbiters');
        const data = await response.json();
        
        if (data.arbiters && data.arbiters.length > 0) {
            // Clear existing options
            arbiterSelect.innerHTML = '<option value="">Select an arbiter...</option>';
            
            // Add arbiter options
            data.arbiters.forEach(arbiter => {
                const option = document.createElement('option');
                option.value = arbiter.ArbiterId;
                option.textContent = `${arbiter.FirstName} ${arbiter.LastName} (${arbiter.ArbiterLevel})`;
                arbiterSelect.appendChild(option);
            });
            
            // Enable the dropdown
            arbiterSelect.disabled = false;
        } else {
            arbiterSelect.innerHTML = '<option value="">No arbiters available</option>';
        }
    } catch (error) {
        console.error('Error loading arbiters:', error);
        arbiterSelect.innerHTML = '<option value="">Error loading arbiters</option>';
    }
}

async function populateLeagueDropdown() {
    const leagueSelect = document.getElementById('leagueSelect');
    
    try {
        const response = await fetch('/leagues');
        const data = await response.json();
        
        if (data.leagues && data.leagues.length > 0) {
            // Clear existing options
            leagueSelect.innerHTML = '<option value="">Select a league...</option>';
            
            // Add league options
            data.leagues.forEach(league => {
                const option = document.createElement('option');
                option.value = league.leagueId;
                option.textContent = `${league.leagueName} (${league.saisonName})`;
                leagueSelect.appendChild(option);
            });
            
            // Enable the dropdown
            leagueSelect.disabled = false;
        } else {
            leagueSelect.innerHTML = '<option value="">No leagues available</option>';
        }
    } catch (error) {
        console.error('Error loading leagues:', error);
        leagueSelect.innerHTML = '<option value="">Error loading leagues</option>';
    }
}

function onArbiterSelected() {
    const arbiterSelect = document.getElementById('arbiterSelect');
    const arbiterDetails = document.getElementById('arbiterDetails');
    const arbiterNameField = document.getElementById('arbiterNameField');
    const arbiterIdField = document.getElementById('arbiterIdField');
    
    if (arbiterSelect.value) {
        // Show arbiter details
        arbiterDetails.classList.remove('hidden');
        
        // Fetch specific arbiter data
        fetchArbiterDetails(arbiterSelect.value);
    } else {
        // Hide arbiter details
        arbiterDetails.classList.add('hidden');
        arbiterNameField.value = '';
        arbiterIdField.value = '';
    }
    
    // Update prepare PDF button state
    updatePreparePdfButtonState();
}

function onLeagueSelected() {
    const leagueSelect = document.getElementById('leagueSelect');
    const presetFields = document.getElementById('presetFields');
    const directorField = document.getElementById('directorField');
    const directorContactField = document.getElementById('directorContactField');
    
    if (leagueSelect.value) {
        // Show preset fields
        presetFields.classList.remove('hidden');
        
        // Fetch specific league data
        fetchLeagueDetails(leagueSelect.value);
    } else {
        // Hide preset fields
        presetFields.classList.add('hidden');
        directorField.value = '';
        directorContactField.value = '';
    }
    
    // Update prepare PDF button state
    updatePreparePdfButtonState();
}

async function fetchArbiterDetails(arbiterId) {
    const arbiterNameField = document.getElementById('arbiterNameField');
    const arbiterIdField = document.getElementById('arbiterIdField');
    
    // Show loading state
    arbiterNameField.value = 'Loading...';
    arbiterIdField.value = 'Loading...';
    
    try {
        const response = await fetch(`/arbiters/${arbiterId}`);
        const data = await response.json();
        
        if (data.arbiter) {
            // Update fields with real arbiter data
            arbiterNameField.value = `${data.arbiter.FirstName} ${data.arbiter.LastName}`;
            arbiterIdField.value = data.arbiter.PlayerId || 'N/A';
        } else {
            arbiterNameField.value = 'Arbiter not found';
            arbiterIdField.value = 'N/A';
        }
    } catch (error) {
        console.error('Error fetching arbiter details:', error);
        arbiterNameField.value = 'Error loading data';
        arbiterIdField.value = 'Error loading data';
    }
}

async function fetchLeagueDetails(leagueId) {
    const directorField = document.getElementById('directorField');
    const directorContactField = document.getElementById('directorContactField');
    
    // Show loading state
    directorField.value = 'Loading...';
    directorContactField.value = 'Loading...';
    
    try {
        const response = await fetch(`/leagues/${leagueId}`);
        const data = await response.json();
        
        if (data.league) {
            // Update fields with real league data
            directorField.value = `${data.league.directorFirstName} ${data.league.directorSurname}`;
            directorContactField.value = data.league.directorEmail || 'Contact not specified';
        } else {
            directorField.value = 'League not found';
            directorContactField.value = 'Contact not available';
        }
    } catch (error) {
        console.error('Error fetching league details:', error);
        directorField.value = 'Error loading data';
        directorContactField.value = 'Error loading data';
    }
}

function updatePreparePdfButtonState() {
    const arbiterSelect = document.getElementById('arbiterSelect');
    const leagueSelect = document.getElementById('leagueSelect');
    const preparePdfBtn = document.getElementById('preparePdfBtn');
    const downloadExcelBtn = document.getElementById('downloadExcelBtn');
    
    // Enable buttons only if both arbiter and league are selected
    if (arbiterSelect.value && leagueSelect.value) {
        preparePdfBtn.disabled = false;
        if (downloadExcelBtn) {
            downloadExcelBtn.disabled = false;
        }
    } else {
        preparePdfBtn.disabled = true;
        if (downloadExcelBtn) {
            downloadExcelBtn.disabled = true;
        }
    }
}

async function downloadExcelFile() {
    const leagueSelect = document.getElementById('leagueSelect');
    const downloadBtn = document.getElementById('downloadExcelBtn');
    const downloadStatus = document.getElementById('downloadStatus');
    
    if (!leagueSelect.value) {
        downloadStatus.innerHTML = '<span class="text-red-600">✗ Please select a league first</span>';
        return;
    }
    
    // Update button state
    downloadBtn.disabled = true;
    downloadBtn.textContent = 'Downloading...';
    downloadStatus.textContent = 'Downloading Excel file...';
    
    try {
        const response = await fetch('/download-excel', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                leagueId: parseInt(leagueSelect.value)
            })
        });
        
        const result = await response.json();
        
        if (response.ok) {
            downloadStatus.innerHTML = `
                <span class="text-green-600">✓ ${result.message}</span><br>
                <span class="text-sm text-gray-600">File saved to: ${result.filePath}</span><br>
                <span class="text-sm text-gray-600">League: ${result.league}</span>
            `;
        } else {
            downloadStatus.innerHTML = `<span class="text-red-600">✗ Error: ${result.error}</span>`;
        }
    } catch (error) {
        downloadStatus.innerHTML = `<span class="text-red-600">✗ Network error: ${error.message}</span>`;
    } finally {
        downloadBtn.disabled = false;
        downloadBtn.textContent = 'Download Excel File';
    }
}
