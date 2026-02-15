async function loadExternalData() {
    console.log('[DATA-LOADING] ===== START loadExternalData =====');
    const btn = document.getElementById('loadDataBtn');
    const status = document.getElementById('loadStatus');
    const seasonYear = document.getElementById('seasonYear').value;
    
    console.log('[DATA-LOADING] Season year:', seasonYear);
    console.log('[DATA-LOADING] Button element found:', !!btn);
    console.log('[DATA-LOADING] Status element found:', !!status);

    btn.disabled = true;
    btn.textContent = 'Loading...';
    status.textContent = 'Loading external data...';
    
    try {
        console.log('[DATA-LOADING] Preparing fetch request to /load-external-data');
        console.log('[DATA-LOADING] Request body:', JSON.stringify({ seasonStartYear: seasonYear }));

        const requestStartTime = performance.now();
        const response = await fetch('/load-external-data', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                seasonStartYear: seasonYear
            })
        });
        const requestEndTime = performance.now();
        const requestDuration = (requestEndTime - requestStartTime).toFixed(2);

        console.log('[DATA-LOADING] Response received in', requestDuration, 'ms');
        console.log('[DATA-LOADING] Response status:', response.status);
        console.log('[DATA-LOADING] Response ok:', response.ok);
        console.log('[DATA-LOADING] Response headers:', Object.fromEntries(response.headers.entries()));

        const result = await response.json();
        console.log('[DATA-LOADING] Response JSON parsed:', result);

        if (response.ok) {
            console.log('[DATA-LOADING] ✓ Request successful');
            console.log('[DATA-LOADING] Arbiters loaded:', result.arbiters_loaded);
            console.log('[DATA-LOADING] Leagues loaded:', result.leagues_loaded);

            status.innerHTML = `
                <span class="text-green-600">✓ ${result.message}</span><br>
            `;
            
            if (result.arbiters_loaded && result.leagues_loaded) {
                console.log('[DATA-LOADING] Both arbiters and leagues loaded, populating dropdown');
                populateLeagueDropdown();
            } else {
                console.log('[DATA-LOADING] ⚠ Not all data loaded - arbiters:', result.arbiters_loaded, 'leagues:', result.leagues_loaded);
            }
        } else {
            console.error('[DATA-LOADING] ✗ Request failed with status:', response.status);
            console.error('[DATA-LOADING] Error message:', result.error);
            status.innerHTML = `<span class="text-red-600">✗ Error: ${result.error}</span>`;
        }
    } catch (error) {
        console.error('[DATA-LOADING] ✗ Exception caught:', error);
        console.error('[DATA-LOADING] Error name:', error.name);
        console.error('[DATA-LOADING] Error message:', error.message);
        console.error('[DATA-LOADING] Error stack:', error.stack);
        status.innerHTML = `<span class="text-red-600">✗ Network error: ${error.message}</span>`;
    } finally {
        console.log('[DATA-LOADING] Resetting button state');
        btn.disabled = false;
        btn.textContent = 'Načítaj dáta z chess.sk';
        console.log('[DATA-LOADING] ===== END loadExternalData =====');
    }
}

async function populateLeagueDropdown() {
    console.log('[LEAGUE-DROPDOWN] ===== START populateLeagueDropdown =====');
    const leagueSelect = document.getElementById('leagueSelect');
    console.log('[LEAGUE-DROPDOWN] League select element found:', !!leagueSelect);

    try {
        console.log('[LEAGUE-DROPDOWN] Fetching leagues from /leagues');
        const requestStartTime = performance.now();
        const response = await fetch('/leagues');
        const requestEndTime = performance.now();
        const requestDuration = (requestEndTime - requestStartTime).toFixed(2);

        console.log('[LEAGUE-DROPDOWN] Response received in', requestDuration, 'ms');
        console.log('[LEAGUE-DROPDOWN] Response status:', response.status);
        console.log('[LEAGUE-DROPDOWN] Response ok:', response.ok);

        const data = await response.json();
        console.log('[LEAGUE-DROPDOWN] Response data:', data);
        console.log('[LEAGUE-DROPDOWN] Leagues count:', data.leagues ? data.leagues.length : 0);

        if (data.leagues && data.leagues.length > 0) {
            console.log('[LEAGUE-DROPDOWN] Processing', data.leagues.length, 'leagues');

            // Clear existing options
            leagueSelect.innerHTML = '<option value="">Vyberte ligu...</option>';
            console.log('[LEAGUE-DROPDOWN] Cleared existing options');

            // Add league options
            data.leagues.forEach((league, index) => {
                console.log(`[LEAGUE-DROPDOWN] Adding league ${index + 1}/${data.leagues.length}:`, {
                    leagueId: league.leagueId,
                    leagueName: league.leagueName,
                    saisonName: league.saisonName,
                    chessResultsLink: league.chessResultsLink
                });

                const option = document.createElement('option');
                option.value = league.leagueId;
                option.textContent = `${league.leagueName} (${league.saisonName})`;
                leagueSelect.appendChild(option);
            });
            
            // Enable the dropdown
            leagueSelect.disabled = false;
            console.log('[LEAGUE-DROPDOWN] ✓ Dropdown populated and enabled');
        } else {
            console.warn('[LEAGUE-DROPDOWN] ⚠ No leagues available');
            leagueSelect.innerHTML = '<option value="">No leagues available</option>';
        }
    } catch (error) {
        console.error('[LEAGUE-DROPDOWN] ✗ Error loading leagues:', error);
        console.error('[LEAGUE-DROPDOWN] Error name:', error.name);
        console.error('[LEAGUE-DROPDOWN] Error message:', error.message);
        console.error('[LEAGUE-DROPDOWN] Error stack:', error.stack);
        leagueSelect.innerHTML = '<option value="">Error loading leagues</option>';
    }

    console.log('[LEAGUE-DROPDOWN] ===== END populateLeagueDropdown =====');
}


async function onLeagueSelected() {
    console.log('[LEAGUE-SELECTED] ===== START onLeagueSelected =====');
    const leagueSelect = document.getElementById('leagueSelect');
    const presetFields = document.getElementById('presetFields');
    
    console.log('[LEAGUE-SELECTED] League select element found:', !!leagueSelect);
    console.log('[LEAGUE-SELECTED] Preset fields element found:', !!presetFields);
    console.log('[LEAGUE-SELECTED] Selected league ID:', leagueSelect.value);
    console.log('[LEAGUE-SELECTED] Selected option text:', leagueSelect.options[leagueSelect.selectedIndex]?.text);

    if (leagueSelect.value) {
        console.log('[LEAGUE-SELECTED] League selected, showing preset fields');
        // Show preset fields
        presetFields.classList.remove('hidden');
        
        // Automatically load rounds data
        try {
            console.log('[LEAGUE-SELECTED] Calling loadRoundsData with leagueId:', leagueSelect.value);
            const roundsStartTime = performance.now();
            await loadRoundsData(parseInt(leagueSelect.value));
            const roundsEndTime = performance.now();
            const roundsDuration = (roundsEndTime - roundsStartTime).toFixed(2);
            console.log('[LEAGUE-SELECTED] ✓ loadRoundsData completed in', roundsDuration, 'ms');
            // The button will be enabled in the rounds editor after it's created
        } catch (error) {
            console.error('[LEAGUE-SELECTED] ✗ Error loading rounds data:', error);
            console.error('[LEAGUE-SELECTED] Error name:', error.name);
            console.error('[LEAGUE-SELECTED] Error message:', error.message);
            console.error('[LEAGUE-SELECTED] Error stack:', error.stack);
        }
    } else {
        console.log('[LEAGUE-SELECTED] No league selected, hiding preset fields');
        // Hide preset fields
        presetFields.classList.add('hidden');
    }

    console.log('[LEAGUE-SELECTED] ===== END onLeagueSelected =====');
}
