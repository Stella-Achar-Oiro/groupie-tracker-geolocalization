// artist-details.js
let map;
let favorites = JSON.parse(localStorage.getItem('favorites')) || [];

mapboxgl.accessToken = 'pk.eyJ1Ijoic3RlbGxhYWNoYXJvaXJvIiwiYSI6ImNtMWhmZHNlODBlc3cybHF5OWh1MDI2dzMifQ.wk3v-v7IuiSiPwyq13qdHw';

// Utility Functions
function showLoading() {
    document.getElementById('loading').style.display = 'block';
}

function hideLoading() {
    document.getElementById('loading').style.display = 'none';
}

function showError(message) {
    const errorElement = document.getElementById('error-message');
    errorElement.textContent = message;
    errorElement.style.display = 'block';
    setTimeout(() => {
        errorElement.style.display = 'none';
    }, 5000);
}

function getArtistId() {
    const pathParts = window.location.pathname.split('/');
    return pathParts[pathParts.length - 1];
}

// Map Utility Functions
function calculateDistance(coords) {
    let totalDistance = 0;
    for (let i = 0; i < coords.length - 1; i++) {
        const [lon1, lat1] = coords[i];
        const [lon2, lat2] = coords[i + 1];
        
        // Haversine formula
        const R = 6371; // Earth's radius in km
        const dLat = (lat2 - lat1) * Math.PI / 180;
        const dLon = (lon2 - lon1) * Math.PI / 180;
        const a = 
            Math.sin(dLat/2) * Math.sin(dLat/2) +
            Math.cos(lat1 * Math.PI / 180) * Math.cos(lat2 * Math.PI / 180) * 
            Math.sin(dLon/2) * Math.sin(dLon/2);
        const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a));
        totalDistance += R * c;
    }
    return Math.round(totalDistance);
}

function addDistanceControl(map, distance) {
    const distanceControl = document.createElement('div');
    distanceControl.className = 'mapboxgl-ctrl mapboxgl-ctrl-group distance-control';
    distanceControl.innerHTML = `
        <div class="distance-info">
            <i class="fas fa-route"></i>
            <span>Tour Distance: ${distance} km</span>
        </div>
    `;
    map.getContainer().appendChild(distanceControl);
}

function addRouteAnimation(map) {
    map.addLayer({
        'id': 'route-animation',
        'type': 'line',
        'source': 'route',
        'layout': {
            'line-cap': 'round',
            'line-join': 'round'
        },
        'paint': {
            'line-color': '#1DB954',
            'line-width': 3,
            'line-opacity': 0.8,
            'line-gradient': [
                'interpolate',
                ['linear'],
                ['line-progress'],
                0, '#1DB954',
                0.5, '#ffffff',
                1, '#1DB954'
            ]
        }
    });

    // Animate the line-gradient
    let progress = 0;
    function animateLine() {
        progress = (progress + 0.005) % 1;
        map.setPaintProperty('route-animation', 'line-progress', progress);
        requestAnimationFrame(animateLine);
    }
    animateLine();
}

// Map Display Function
function displayMap(locations) {
    if (map) {
        map.remove();
    }

    // Calculate center point of all locations
    let centerLon = 0;
    let centerLat = 0;
    let validLocations = locations.filter(loc => loc.lon && loc.lat);
    
    if (validLocations.length > 0) {
        centerLon = validLocations.reduce((sum, loc) => sum + loc.lon, 0) / validLocations.length;
        centerLat = validLocations.reduce((sum, loc) => sum + loc.lat, 0) / validLocations.length;
    }

    map = new mapboxgl.Map({
        container: 'map',
        style: 'mapbox://styles/mapbox/dark-v10',
        center: [centerLon, centerLat],
        zoom: 1,
        projection: 'mercator' // Changed from 'globe' to 'mercator' for better compatibility
    });

    // Add controls
    map.addControl(new mapboxgl.NavigationControl());
    map.addControl(new mapboxgl.FullscreenControl());

    const coordinates = [];

    // Add markers and collect coordinates
    locations.forEach((location, index) => {
        if (!location.lon || !location.lat) {
            console.error('Invalid coordinates for location:', location);
            return;
        }

        coordinates.push([location.lon, location.lat]);

        const el = document.createElement('div');
        el.className = 'custom-marker';
        el.innerHTML = `
            <div class="marker-content">
                <i class="fas fa-map-marker-alt"></i>
                <span class="marker-number">${index + 1}</span>
            </div>
        `;

        const popup = new mapboxgl.Popup({
            offset: 25,
            closeButton: true,
            closeOnClick: false
        })
        .setHTML(`
            <div class="custom-popup">
                <h3>${location.address}</h3>
                <p>Stop ${index + 1} on the tour</p>
            </div>
        `);

        new mapboxgl.Marker(el)
            .setLngLat([location.lon, location.lat])
            .setPopup(popup)
            .addTo(map);
    });

    // Calculate appropriate zoom level based on coordinates spread
    if (coordinates.length > 1) {
        const lons = coordinates.map(coord => coord[0]);
        const lats = coordinates.map(coord => coord[1]);
        const minLon = Math.min(...lons);
        const maxLon = Math.max(...lons);
        const minLat = Math.min(...lats);
        const maxLat = Math.max(...lats);
        
        const bounds = [
            [minLon - 10, minLat - 10], // Add padding
            [maxLon + 10, maxLat + 10]  // Add padding
        ];
        
        map.fitBounds(bounds, {
            padding: 50,
            duration: 1000
        });
    }

    map.on('style.load', () => {
        if (coordinates.length > 1) {
            map.addSource('route', {
                'type': 'geojson',
                'data': {
                    'type': 'Feature',
                    'properties': {},
                    'geometry': {
                        'type': 'LineString',
                        'coordinates': coordinates
                    }
                }
            });

            // Add static route
            map.addLayer({
                'id': 'route',
                'type': 'line',
                'source': 'route',
                'layout': {
                    'line-join': 'round',
                    'line-cap': 'round'
                },
                'paint': {
                    'line-color': '#1DB954',
                    'line-width': 2,
                    'line-opacity': 0.5,
                    'line-dasharray': [2, 2]
                }
            });

            // Add animated route layer
            map.addLayer({
                'id': 'route-animation',
                'type': 'line',
                'source': 'route',
                'layout': {
                    'line-cap': 'round',
                    'line-join': 'round'
                },
                'paint': {
                    'line-color': '#1DB954',
                    'line-width': 3,
                    'line-opacity': 0.8
                }
            });

            // Calculate and display total distance
            const totalDistance = calculateDistance(coordinates);
            addDistanceControl(map, totalDistance);

            // Start the animation
            let progress = 0;
            function animateLine() {
                if (!map.getLayer('route-animation')) return;
                progress = (progress + 0.005) % 1;
                map.setPaintProperty('route-animation', 'line-gradient', [
                    'interpolate',
                    ['linear'],
                    ['line-progress'],
                    0, '#1DB954',
                    progress, '#ffffff',
                    1, '#1DB954'
                ]);
                requestAnimationFrame(animateLine);
            }
            animateLine();
        }
    });
}

function showActionConfirmation(message, type = 'success') {
    const confirmationEl = document.createElement('div');
    confirmationEl.className = `feedback-message ${type}`;
    confirmationEl.innerHTML = `
        <i class="fas ${type === 'success' ? 'fa-check-circle' : 'fa-exclamation-circle'}"></i>
        <span>${message}</span>
    `;
    
    document.body.appendChild(confirmationEl);
    
    setTimeout(() => {
        confirmationEl.style.animation = 'slideOut 0.3s ease-out forwards';
        setTimeout(() => confirmationEl.remove(), 300);
    }, 3000);
}

function toggleFavorite(artistId) {
    // Convert artistId to string for consistent comparison
    artistId = artistId.toString();
    
    const button = document.getElementById(`favorite-${artistId}`);
    const index = favorites.indexOf(artistId);
    
    if (index === -1) {
        favorites.push(artistId);
        button.classList.add('favorited');
        button.innerHTML = '<i class="fas fa-star"></i> Remove from Favorites';
        showActionConfirmation('Added to favorites!');
    } else {
        favorites.splice(index, 1);
        button.classList.remove('favorited');
        button.innerHTML = '<i class="far fa-star"></i> Add to Favorites';
        showActionConfirmation('Removed from favorites!');
    }
    
    localStorage.setItem('favorites', JSON.stringify(favorites));
}

function shareArtist(artistId, artistName, creationDate = '', firstAlbum = '') {
    const shareData = {
        title: `${artistName} - Groupie Tracker`,
        text: `Check out ${artistName} on Groupie Tracker!${creationDate ? `\nCreation Date: ${creationDate}` : ''}${firstAlbum ? `\nFirst Album: ${firstAlbum}` : ''}`,
        url: window.location.href
    };

    try {
        if (navigator.share) {
            // Use Web Share API if available (mobile devices)
            navigator.share(shareData)
                .then(() => showActionConfirmation('Shared successfully!'))
                .catch(error => {
                    if (error.name !== 'AbortError') {
                        showActionConfirmation('Error sharing content', 'error');
                    }
                });
        } else {
            // Fallback to clipboard copy (desktop)
            navigator.clipboard.writeText(`${shareData.text}\n${shareData.url}`)
                .then(() => showActionConfirmation('Copied to clipboard!', 'clipboard'))
                .catch(() => showActionConfirmation('Error copying to clipboard', 'error'));
        }
    } catch (error) {
        showActionConfirmation('Error sharing content', 'error');
    }
}

function showNotification(message, type = 'success') {
    // Remove existing notification if any
    const existingNotification = document.querySelector('.notification');
    if (existingNotification) {
        existingNotification.remove();
    }

    const notification = document.createElement('div');
    notification.className = `notification notification-${type}`;
    
    // Add icon based on type
    let icon = 'check';
    if (type === 'clipboard') icon = 'clipboard';
    if (type === 'error') icon = 'exclamation-circle';
    
    notification.innerHTML = `
        <i class="fas fa-${icon}"></i>
        ${message}
    `;
    
    document.body.appendChild(notification);
    
    // Remove notification after 2 seconds
    setTimeout(() => {
        notification.style.animation = 'fadeOut 0.3s ease-out forwards';
        setTimeout(() => {
            notification.remove();
        }, 300);
    }, 2000);
}

function displayArtistDetails(details) {
    const container = document.getElementById('artist-details');
    
    container.innerHTML = `
        <h2>${details.artist.name}</h2>
        <img src="${details.artist.image}" alt="${details.artist.name}" class="artist-image">
        <div class="artist-info">
            <p><i class="fas fa-users"></i> Members: ${details.artist.members.join(', ')}</p>
            <p><i class="fas fa-calendar-alt"></i> Creation Date: ${details.artist.creationDate}</p>
            <p><i class="fas fa-compact-disc"></i> First Album: ${details.artist.firstAlbum}</p>
            
            <h3><i class="fas fa-map-marker-alt"></i> Locations:</h3>
            <ul class="locations-list">
                ${details.locations.map(loc => `<li>${loc.address}</li>`).join('')}
            </ul>
            
            <h3><i class="fas fa-calendar-check"></i> Dates:</h3>
            <ul class="dates-list">
                ${details.dates.map(date => `<li>${date}</li>`).join('')}
            </ul>
            
            <h3><i class="fas fa-link"></i> Relations:</h3>
            <ul class="relations-list">
                ${Object.entries(details.relations).map(([loc, dates]) => `
                    <li>${loc}: ${dates.join(', ')}</li>
                `).join('')}
            </ul>
            
            <div class="action-buttons">
                <button id="favorite-${details.artist.id}" 
                        class="action-btn favorite-btn ${favorites.includes(details.artist.id.toString()) ? 'favorited' : ''}"
                        onclick="toggleFavorite('${details.artist.id}')">
                    <i class="fas ${favorites.includes(details.artist.id.toString()) ? 'fa-star' : 'fa-star-o'}"></i>
                    ${favorites.includes(details.artist.id.toString()) ? 'Remove from Favorites' : 'Add to Favorites'}
                <button class="action-btn share-btn" 
                        onclick="shareArtist('${details.artist.id}', '${details.artist.name}', '${details.artist.creationDate}', '${details.artist.firstAlbum}')">
                    <i class="fas fa-share-alt"></i> Share
                </button>
            </div>
        </div>
    `;

    displayMap(details.locations);
}

function initializeKeyboardShortcuts() {
    document.addEventListener('keydown', (e) => {
        // Alt + Left Arrow for back navigation
        if (e.altKey && e.key === 'ArrowLeft') {
            window.location.href = '/';
        }
        
        // Alt + F to toggle favorite
        if (e.altKey && e.key === 'f') {
            const artistId = getArtistId();
            toggleFavorite(artistId);
        }

        // Alt + S to share
        if (e.altKey && e.key === 's') {
            const artistDetails = document.getElementById('artist-details');
            if (artistDetails) {
                // Extract information from the artist details container
                const name = artistDetails.querySelector('h2')?.textContent || '';
                const creationDate = artistDetails.querySelector('p:nth-child(3)')?.textContent.split(': ')[1] || '';
                const firstAlbum = artistDetails.querySelector('p:nth-child(4)')?.textContent.split(': ')[1] || '';
                
                if (name) {
                    shareArtist(getArtistId(), name, creationDate, firstAlbum);
                }
            }
        }
    });
}
// Event Listeners
window.addEventListener('load', () => {
    const artistId = getArtistId();
    if (artistId) {
        const container = document.getElementById('artist-details');
        showLoading();
        
        fetch(`/api/artist/${artistId}`)
            .then(response => {
                if (!response.ok) {
                    throw { response };
                }
                return response.json();
            })
            .then(data => {
                displayArtistDetails(data);
                hideLoading();
            })
            .catch(error => handleError(error, container));
    }
});

document.addEventListener('DOMContentLoaded', () => {
    // Initialize keyboard shortcuts
    initializeKeyboardShortcuts();

    // Initialize shortcuts guide
    const shortcutsGuide = document.querySelector('.shortcuts-guide');
    const shortcutsToggle = document.querySelector('.shortcuts-toggle');

    if (shortcutsGuide && shortcutsToggle) {
        shortcutsToggle.addEventListener('click', () => {
            shortcutsGuide.classList.toggle('active');
        });

        // Close shortcuts guide when clicking outside
        document.addEventListener('click', (e) => {
            if (!shortcutsGuide.contains(e.target)) {
                shortcutsGuide.classList.remove('active');
            }
        });

        // Close shortcuts guide when pressing Escape
        document.addEventListener('keydown', (e) => {
            if (e.key === 'Escape') {
                shortcutsGuide.classList.remove('active');
            }
        });
    }
});

window.addEventListener('online', () => {
    // Refresh the page when connection is restored
    window.location.reload();
});

window.addEventListener('offline', () => {
    const container = document.getElementById('results-container'); // or 'artist-details'
    createOfflineUI(container);
});