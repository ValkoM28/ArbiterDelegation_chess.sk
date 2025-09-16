// Data Loading Functions
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
                <span class="text-sm">Arbiters: ${result.arbiters_loaded ? 'Loaded' : 'Not loaded'}</span><br>
                <span class="text-sm">Leagues: ${result.leagues_loaded ? 'Loaded' : 'Not loaded'}</span>
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
        // Fetch arbiters data
        const arbitersResponse = await fetch('/arbiters');
        const arbitersData = await arbitersResponse.json();
        
        if (arbitersData.arbiters && arbitersData.arbiters.length > 0) {
            const firstFew = arbitersData.arbiters.slice(0, 3);
            arbitersPreview.innerHTML = firstFew.map(arbiter => 
                `<div>• ${arbiter.FirstName} ${arbiter.LastName} (${arbiter.ArbiterLevel})</div>`
            ).join('');
            if (arbitersData.arbiters.length > 3) {
                arbitersPreview.innerHTML += `<div class="text-gray-400">... and ${arbitersData.arbiters.length - 3} more</div>`;
            }
        } else {
            arbitersPreview.innerHTML = '<div class="text-gray-400">No arbiters data</div>';
        }

        // Fetch leagues data
        const leaguesResponse = await fetch('/leagues');
        const leaguesData = await leaguesResponse.json();
        
        if (leaguesData.leagues && leaguesData.leagues.length > 0) {
            const firstFew = leaguesData.leagues.slice(0, 3);
            leaguesPreview.innerHTML = firstFew.map(league => 
                `<div>• ${league.leagueName} (${league.saisonName})</div>`
            ).join('');
            if (leaguesData.leagues.length > 3) {
                leaguesPreview.innerHTML += `<div class="text-gray-400">... and ${leaguesData.leagues.length - 3} more</div>`;
            }
        } else {
            leaguesPreview.innerHTML = '<div class="text-gray-400">No leagues data</div>';
        }

        // Show the preview
        preview.classList.remove('hidden');
    } catch (error) {
        console.error('Error loading data preview:', error);
    }
}
