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
        contactPerson = '';
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

    roundsContainer.classList.remove('hidden');

    // JavaScript injected html, probably enough for the usecase
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
                <div class="p-3 bg-gray-50 rounded border">
                    <div class="grid grid-cols-1 md:grid-cols-4 gap-3 mb-3">
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
                    
                    <!-- Match Arbiter Selection -->
                    <div class="mt-3 p-2 bg-blue-50 rounded">
                        <label class="block text-xs font-medium text-gray-600 mb-1">Arbiter for this Match:</label>
                        <select 
                            id="round_${roundIndex}_match_${matchIndex}_arbiter" 
                            class="w-full px-2 py-1 text-sm border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-blue-500"
                            onchange="onMatchArbiterSelected(${roundIndex}, ${matchIndex})"
                        >
                            <option value="">Select an arbiter for this match...</option>
                        </select>
                        <div id="round_${roundIndex}_match_${matchIndex}_arbiter_details" class="mt-1 text-xs text-gray-600 hidden">
                            <!-- Arbiter details will be shown here -->
                        </div>
                    </div>
                </div>
            `;
        });
        // end the editor
        html += `
                </div>
            </div>
        `;
    });

    roundsContainer.innerHTML = html;
    
    // Populate arbiter dropdowns for all matches
    populateMatchArbiterDropdowns();
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



// Populate arbiter dropdowns for all matches
async function populateMatchArbiterDropdowns() {
    try {
        // Get arbiters data
        const response = await fetch('/arbiters');
        const data = await response.json();
        
        if (data.arbiters && data.arbiters.length > 0) {
            // Populate all match arbiter dropdowns
            currentRounds.forEach((round, roundIndex) => {
                round.matches.forEach((match, matchIndex) => {
                    const selectElement = document.getElementById(`round_${roundIndex}_match_${matchIndex}_arbiter`);
                    if (selectElement) {
                        // Clear existing options except the first one
                        selectElement.innerHTML = '<option value="">Select an arbiter for this match...</option>';
                        
                        // Add arbiter options
                        data.arbiters.forEach(arbiter => {
                            const option = document.createElement('option');
                            option.value = arbiter.ArbiterId;
                            option.textContent = `${arbiter.FirstName} ${arbiter.LastName} (${arbiter.ArbiterLevel})${arbiter.KlubName ? ` - ${arbiter.KlubName}` : ''}`;
                            selectElement.appendChild(option);
                        });
                    }
                });
            });
        }
    } catch (error) {
        console.error('Error loading arbiters for match dropdowns:', error);
    }
}

// Handle arbiter selection for a specific match
async function onMatchArbiterSelected(roundIndex, matchIndex) {
    const selectElement = document.getElementById(`round_${roundIndex}_match_${matchIndex}_arbiter`);
    const detailsElement = document.getElementById(`round_${roundIndex}_match_${matchIndex}_arbiter_details`);
    
    if (selectElement.value) {
        try {
            const response = await fetch(`/arbiters/${selectElement.value}`);
            const data = await response.json();
            
            if (data.arbiter) {
                const clubInfo = data.arbiter.KlubName ? ` - ${data.arbiter.KlubName}` : '';
                detailsElement.innerHTML = `
                    <strong>Selected:</strong> ${data.arbiter.FirstName} ${data.arbiter.LastName} (ID: ${data.arbiter.PlayerId})${clubInfo}
                `;
                detailsElement.classList.remove('hidden');
            }
        } catch (error) {
            console.error('Error fetching arbiter details:', error);
            detailsElement.innerHTML = '<span class="text-red-600">Error loading arbiter details</span>';
            detailsElement.classList.remove('hidden');
        }
    } else {
        detailsElement.classList.add('hidden');
    }
}

// Prepare PDFData array from current rounds data
function preparePDFDataArray() {
    const leagueSelect = document.getElementById('leagueSelect');
    
    // Get league info from the form
    const directorInfo = document.getElementById('directorField').value;
    const directorContact = document.getElementById('directorContactField').value;
    
    // Get global info from rounds editor if available
    const globalDirectorInfo = document.getElementById('globalDirectorInfo')?.value || directorInfo;
    const globalContactPerson = document.getElementById('globalContactPerson')?.value || '';
    
    // Get league name from the selected option
    const selectedLeagueOption = leagueSelect.options[leagueSelect.selectedIndex];
    const leagueName = selectedLeagueOption ? selectedLeagueOption.textContent.split(' (')[0] : '';
    const leagueYear = selectedLeagueOption ? selectedLeagueOption.textContent.match(/\((.+?)\)/)?.[1] || '' : '';
    
    const pdfDataArray = [];
    
    // Create PDFData for each match
    currentRounds.forEach((round, roundIndex) => {
        round.matches.forEach((match, matchIndex) => {
            // Get current form data (including any user edits)
            const homeTeam = document.getElementById(`round_${roundIndex}_match_${matchIndex}_home`)?.value || match.homeTeam;
            const guestTeam = document.getElementById(`round_${roundIndex}_match_${matchIndex}_guest`)?.value || match.guestTeam;
            const dateTime = document.getElementById(`round_${roundIndex}_match_${matchIndex}_datetime`)?.value || match.dateTime;
            const address = document.getElementById(`round_${roundIndex}_match_${matchIndex}_address`)?.value || match.address;
            
            // Get arbiter info from the match's arbiter dropdown
            const arbiterSelect = document.getElementById(`round_${roundIndex}_match_${matchIndex}_arbiter`);
            const selectedArbiterId = arbiterSelect ? arbiterSelect.value : '';
            
            // Get arbiter details from the details element
            const arbiterDetails = document.getElementById(`round_${roundIndex}_match_${matchIndex}_arbiter_details`);
            let arbiterName = '';
            let arbiterId = '';
            let arbiterClubName = '';
            
            if (arbiterDetails && arbiterDetails.textContent) {
                // Extract name, ID, and club from the details text
                const detailsText = arbiterDetails.textContent;
                const nameMatch = detailsText.match(/Selected: (.+?) \(ID: (.+?)\)(?: - (.+))?$/);
                if (nameMatch) {
                    arbiterName = nameMatch[1];
                    arbiterId = nameMatch[2];
                    arbiterClubName = nameMatch[3] || '';
                }
            }
            
            // If no arbiter details found, try to get from dropdown selection
            if (!arbiterName && selectedArbiterId) {
                const selectedArbiterOption = arbiterSelect.options[arbiterSelect.selectedIndex];
                if (selectedArbiterOption) {
                    const optionText = selectedArbiterOption.textContent;
                    // Updated regex to handle club name in the format: "Name (Level) - Club"
                    const arbiterMatch = optionText.match(/^(.+?) \((.+?)\)(?: - (.+))?$/);
                    if (arbiterMatch) {
                        arbiterName = arbiterMatch[1];
                        arbiterId = selectedArbiterId;
                    }
                }
            }
            
            const pdfData = {
                league: {
                    name: leagueName,
                    year: leagueYear
                },
                director: {
                    contact: globalDirectorInfo || `${directorInfo} (${directorContact})`
                },
                arbiter: {
                    firstName: arbiterName.split(' ')[0] || '',
                    lastName: arbiterName.split(' ').slice(1).join(' ') || '',
                    playerId: arbiterId,
                    clubName: '' // just because of the updated ArbiterInfo in backend it wont run without this line :D 
                },
                match: {
                    homeTeam: homeTeam,
                    guestTeam: guestTeam,
                    dateTime: dateTime,
                    address: address
                },
                contactPerson: globalContactPerson
            };
            
            pdfDataArray.push(pdfData);
        });
    });
    
    return pdfDataArray;
}

// Prepare delegation data and send to backend
async function prepareDelegationData() {
    const leagueSelect = document.getElementById('leagueSelect');
    const roundsStatus = document.getElementById('roundsStatus');
    
    if (!leagueSelect.value) {
        roundsStatus.innerHTML = '<span class="text-red-600">✗ Please select a league first</span>';
        return;
    }
    
    if (!currentRounds || currentRounds.length === 0) {
        roundsStatus.innerHTML = '<span class="text-red-600">✗ Please load rounds data first</span>';
        return;
    }
    
    try {
        // Prepare PDFData array from current rounds data
        const pdfDataArray = preparePDFDataArray();
        
        // Validate that we have data
        if (pdfDataArray.length === 0) {
            roundsStatus.innerHTML = '<span class="text-red-600">✗ No match data found</span>';
            return;
        }
        
        // Check if all matches have arbiters assigned (for testing, we'll allow missing arbiters)
        const missingArbiters = pdfDataArray.filter(data => !data.arbiter.firstName || !data.arbiter.lastName);
        if (missingArbiters.length > 0) {
            console.log(`Warning: ${missingArbiters.length} matches missing arbiter assignment - proceeding anyway for testing`);
        }
        
        // Send to backend
        const response = await fetch('/delegate-arbiters', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(pdfDataArray)
        });
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const result = await response.json();
        console.log('Delegation data sent:', result);
        
        roundsStatus.innerHTML = `
            <span class="text-green-600">✓ ${result.message}</span><br>
            <span class="text-sm text-gray-600">Count: ${result.count} items</span><br>
            <span class="text-sm text-gray-600">Check server console for detailed output</span>
        `;
        
    } catch (error) {
        console.error('Error preparing delegation data:', error);
        roundsStatus.innerHTML = `<span class="text-red-600">✗ Error: ${error.message}</span>`;
    }
}
