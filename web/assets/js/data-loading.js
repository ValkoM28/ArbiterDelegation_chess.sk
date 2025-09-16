async function loadExternalData() {
    const btn = document.getElementById('loadDataBtn');
    const status = document.getElementById('loadStatus');
    const seasonYear = document.getElementById('seasonYear').value;
    
    btn.disabled = true;
    btn.textContent = 'Loading...';
    status.textContent = 'Loading external data...';
    
    try {
        const response = await fetch('/load-external-data', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                seasonStartYear: seasonYear
            })
        });
        
        const result = await response.json();
        
        if (response.ok) {
            status.innerHTML = `
                <span class="text-green-600">✓ ${result.message}</span><br>
            `;
            
            // Show data preview and populate dropdowns
            if (result.arbiters_loaded && result.leagues_loaded) {
                showDataPreview();
                populateLeagueDropdown();
            }
        } else {
            status.innerHTML = `<span class="text-red-600">✗ Error: ${result.error}</span>`;
        }
    } catch (error) {
        status.innerHTML = `<span class="text-red-600">✗ Network error: ${error.message}</span>`;
    } finally {
        btn.disabled = false;
        btn.textContent = 'Load External Data';
    }
}

async function showDataPreview() {
    const preview = document.getElementById('dataPreview');
    const arbitersPreview = document.getElementById('arbitersPreview');
    const leaguesPreview = document.getElementById('leaguesPreview');
    
    try {
        const arbitersResponse = await fetch('/arbiters');
        const arbitersData = await arbitersResponse.json();
        
        if (arbitersData.arbiters && arbitersData.arbiters.length > 0) {
            arbitersPreview.innerHTML = `<div class="text-green-600 font-medium">Loaded ${arbitersData.arbiters.length} active arbiters</div>`;
        } else {
            arbitersPreview.innerHTML = '<div class="text-gray-400">No arbiters data</div>';
        }

        const leaguesResponse = await fetch('/leagues');
        const leaguesData = await leaguesResponse.json();
        
        if (leaguesData.leagues && leaguesData.leagues.length > 0) {
            leaguesPreview.innerHTML = `<div class="text-green-600 font-medium">Loaded ${leaguesData.leagues.length} leagues</div>`;
        } else {
            leaguesPreview.innerHTML = '<div class="text-gray-400">No leagues data</div>';
        }

        // Show the preview
        preview.classList.remove('hidden');
    } catch (error) {
        console.error('Error loading data preview:', error);
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


async function onLeagueSelected() {
    const leagueSelect = document.getElementById('leagueSelect');
    const presetFields = document.getElementById('presetFields');
    
    if (leagueSelect.value) {
        // Show preset fields
        presetFields.classList.remove('hidden');
        
        // Automatically load rounds data
        try {
            await loadRoundsData(parseInt(leagueSelect.value));
            // The button will be enabled in the rounds editor after it's created
        } catch (error) {
            console.error('Error loading rounds data:', error);
        }
    } else {
        // Hide preset fields
        presetFields.classList.add('hidden');
    }
}
