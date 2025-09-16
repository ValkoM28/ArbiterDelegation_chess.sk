// Rounds Management JavaScript
let currentRounds = [];
let currentLeague = null;
let directorInfo = '';
let contactPerson = '';

// Load rounds data for a specific league
async function loadRoundsData(leagueId) {
    try {
        const response = await fetch('/get-rounds', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ leagueId: leagueId })
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        currentRounds = data.rounds;
        currentLeague = data.league;
        
        // Extract director info and contact person from league
        directorInfo = `${currentLeague.directorFirstName} ${currentLeague.directorSurname}`;
        if (currentLeague.directorEmail) {
            directorInfo += ` (${currentLeague.directorEmail})`;
        }
        contactPerson = ''; // This will be set by user

        console.log('Rounds loaded:', currentRounds);
        console.log('League:', currentLeague);
        
        // Display the rounds editor
        displayRoundsEditor();
        
        return data;
    } catch (error) {
        console.error('Error loading rounds data:', error);
        showStatus('Error loading rounds data: ' + error.message, 'error');
        throw error;
    }
}

// Display the rounds editor interface
function displayRoundsEditor() {
    const roundsContainer = document.getElementById('roundsEditor');
    if (!roundsContainer) {
        console.error('Rounds editor container not found');
        return;
    }

    // Show the rounds editor section
    roundsContainer.classList.remove('hidden');

    // Build the HTML for rounds editor
    let html = `
        <div class="bg-white rounded-lg shadow-md p-6">
            <h2 class="text-2xl font-semibold text-gray-700 mb-6">Rounds Data Editor</h2>
            
            <!-- Global Fields -->
            <div class="mb-8 p-4 bg-gray-50 rounded-lg">
                <h3 class="text-lg font-medium text-gray-700 mb-4">Global Information</h3>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Director Info</label>
                        <input 
                            type="text" 
                            id="globalDirectorInfo" 
                            value="${directorInfo}"
                            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                        />
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Contact Person</label>
                        <input 
                            type="text" 
                            id="globalContactPerson" 
                            value="${contactPerson}"
                            placeholder="Enter contact person"
                            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                        />
                    </div>
                </div>
            </div>

            <!-- Rounds List -->
            <div class="space-y-6">
    `;

    // Add each round
    currentRounds.forEach((round, roundIndex) => {
        html += `
            <div class="border border-gray-200 rounded-lg p-4">
                <div class="flex items-center justify-between mb-4">
                    <h4 class="text-lg font-medium text-gray-700">Round ${round.number}</h4>
                    <div class="flex items-center space-x-4">
                        <label class="text-sm font-medium text-gray-700">Date & Time:</label>
                        <input 
                            type="text" 
                            id="round_${roundIndex}_datetime" 
                            value="${round.dateTime}"
                            class="px-3 py-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                        />
                    </div>
                </div>
                
                <div class="space-y-3">
        `;

        // Add each match in the round
        round.matches.forEach((match, matchIndex) => {
            html += `
                <div class="grid grid-cols-1 md:grid-cols-4 gap-3 p-3 bg-gray-50 rounded">
                    <div>
                        <label class="block text-xs font-medium text-gray-600 mb-1">Home Team</label>
                        <input 
                            type="text" 
                            id="round_${roundIndex}_match_${matchIndex}_home" 
                            value="${match.homeTeam}"
                            class="w-full px-2 py-1 text-sm border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-blue-500"
                        />
                    </div>
                    <div>
                        <label class="block text-xs font-medium text-gray-600 mb-1">Guest Team</label>
                        <input 
                            type="text" 
                            id="round_${roundIndex}_match_${matchIndex}_guest" 
                            value="${match.guestTeam}"
                            class="w-full px-2 py-1 text-sm border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-blue-500"
                        />
                    </div>
                    <div>
                        <label class="block text-xs font-medium text-gray-600 mb-1">Date & Time</label>
                        <input 
                            type="text" 
                            id="round_${roundIndex}_match_${matchIndex}_datetime" 
                            value="${match.dateTime}"
                            class="w-full px-2 py-1 text-sm border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-blue-500"
                        />
                    </div>
                    <div>
                        <label class="block text-xs font-medium text-gray-600 mb-1">Address</label>
                        <input 
                            type="text" 
                            id="round_${roundIndex}_match_${matchIndex}_address" 
                            value="${match.address}"
                            placeholder="Enter address"
                            class="w-full px-2 py-1 text-sm border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-blue-500"
                        />
                    </div>
                </div>
            `;
        });

        html += `
                </div>
            </div>
        `;
    });

    html += `
            </div>
            
            <!-- Action Buttons -->
            <div class="mt-8 flex justify-end space-x-4">
                <button 
                    onclick="saveRoundsData()"
                    class="bg-green-500 hover:bg-green-600 text-white font-bold py-2 px-6 rounded transition duration-200"
                >
                    Save Changes
                </button>
                <button 
                    onclick="hideRoundsEditor()"
                    class="bg-gray-500 hover:bg-gray-600 text-white font-bold py-2 px-6 rounded transition duration-200"
                >
                    Close Editor
                </button>
            </div>
        </div>
    `;

    roundsContainer.innerHTML = html;
}

// Save rounds data
async function saveRoundsData() {
    try {
        // Collect all the data from the form
        const updatedRounds = [];
        
        currentRounds.forEach((round, roundIndex) => {
            const updatedRound = {
                number: round.number,
                dateTime: document.getElementById(`round_${roundIndex}_datetime`).value,
                matches: []
            };

            round.matches.forEach((match, matchIndex) => {
                const updatedMatch = {
                    homeTeam: document.getElementById(`round_${roundIndex}_match_${matchIndex}_home`).value,
                    guestTeam: document.getElementById(`round_${roundIndex}_match_${matchIndex}_guest`).value,
                    dateTime: document.getElementById(`round_${roundIndex}_match_${matchIndex}_datetime`).value,
                    address: document.getElementById(`round_${roundIndex}_match_${matchIndex}_address`).value
                };
                updatedRound.matches.push(updatedMatch);
            });

            updatedRounds.push(updatedRound);
        });

        // Update global info
        directorInfo = document.getElementById('globalDirectorInfo').value;
        contactPerson = document.getElementById('globalContactPerson').value;

        // Send to backend
        const response = await fetch('/save-rounds', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ rounds: updatedRounds })
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        console.log('Rounds saved:', data);
        
        // Update current rounds
        currentRounds = updatedRounds;
        
        showStatus('Rounds data saved successfully!', 'success');
        
    } catch (error) {
        console.error('Error saving rounds data:', error);
        showStatus('Error saving rounds data: ' + error.message, 'error');
    }
}

// Hide rounds editor
function hideRoundsEditor() {
    const roundsContainer = document.getElementById('roundsEditor');
    if (roundsContainer) {
        roundsContainer.classList.add('hidden');
    }
}

// Show status message
function showStatus(message, type = 'info') {
    const statusElement = document.getElementById('roundsStatus');
    if (statusElement) {
        statusElement.textContent = message;
        statusElement.className = `text-sm ${type === 'error' ? 'text-red-600' : type === 'success' ? 'text-green-600' : 'text-blue-600'}`;
    }
}

// Get current rounds data (for use by other scripts)
function getCurrentRounds() {
    return currentRounds;
}

// Get director info
function getDirectorInfo() {
    return directorInfo;
}

// Get contact person
function getContactPerson() {
    return contactPerson;
}

// Load rounds for the currently selected league
async function loadRoundsForCurrentLeague() {
    const leagueSelect = document.getElementById('leagueSelect');
    const loadRoundsBtn = document.getElementById('loadRoundsBtn');
    const roundsStatus = document.getElementById('roundsStatus');
    
    if (!leagueSelect.value) {
        roundsStatus.innerHTML = '<span class="text-red-600">✗ Please select a league first</span>';
        return;
    }
    
    // Update button state
    loadRoundsBtn.disabled = true;
    loadRoundsBtn.textContent = 'Loading...';
    roundsStatus.textContent = 'Loading rounds data...';
    
    try {
        await loadRoundsData(parseInt(leagueSelect.value));
        roundsStatus.innerHTML = '<span class="text-green-600">✓ Rounds data loaded successfully</span>';
    } catch (error) {
        roundsStatus.innerHTML = `<span class="text-red-600">✗ Error: ${error.message}</span>`;
    } finally {
        loadRoundsBtn.disabled = false;
        loadRoundsBtn.textContent = 'Load & Edit Rounds';
    }
}
