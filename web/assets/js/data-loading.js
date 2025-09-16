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
            
            if (result.arbiters_loaded && result.leagues_loaded) {
                populateLeagueDropdown();
            }
        } else {
            status.innerHTML = `<span class="text-red-600">✗ Error: ${result.error}</span>`;
        }
    } catch (error) {
        status.innerHTML = `<span class="text-red-600">✗ Network error: ${error.message}</span>`;
    } finally {
        btn.disabled = false;
        btn.textContent = 'Načítaj dáta z chess.sk';
    }
}

async function populateLeagueDropdown() {
    const leagueSelect = document.getElementById('leagueSelect');
    
    try {
        const response = await fetch('/leagues');
        const data = await response.json();
        
        if (data.leagues && data.leagues.length > 0) {
            // Clear existing options
            leagueSelect.innerHTML = '<option value="">Vyberte ligu...</option>';
            
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
