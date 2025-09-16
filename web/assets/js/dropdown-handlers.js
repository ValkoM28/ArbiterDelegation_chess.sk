async function onLeagueSelected() {
    const leagueSelect = document.getElementById('leagueSelect');
    const presetFields = document.getElementById('presetFields');
    const directorField = document.getElementById('directorField');
    const directorContactField = document.getElementById('directorContactField');
    const prepareDelegationBtn = document.getElementById('prepareDelegationBtn');
    
    if (leagueSelect.value) {
        // Show preset fields
        presetFields.classList.remove('hidden');
        
        // Fetch specific league data
        await fetchLeagueDetails(leagueSelect.value);
        
        // Automatically load rounds data
        try {
            await loadRoundsData(parseInt(leagueSelect.value));
            prepareDelegationBtn.disabled = false;
        } catch (error) {
            console.error('Error loading rounds data:', error);
        }
    } else {
        // Hide preset fields
        presetFields.classList.add('hidden');
        directorField.value = '';
        directorContactField.value = '';
        prepareDelegationBtn.disabled = true;
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
