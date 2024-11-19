// artist-details.js
let map;
let favorites = JSON.parse(localStorage.getItem('favorites')) || [];

mapboxgl.accessToken = 'pk.eyJ1Ijoic3RlbGxhYWNoYXJvaXJvIiwiYSI6ImNtMWhmZHNlODBlc3cybHF5OWh1MDI2dzMifQ.wk3v-v7IuiSiPwyq13qdHw';

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

function toggleFavorite(artistId) {
    const index = favorites.indexOf(artistId);
    if (index === -1) {
        favorites.push(artistId);
    } else {
        favorites.splice(index, 1);
    }
    localStorage.setItem('favorites', JSON.stringify(favorites));
    updateFavoriteButton(artistId);
}

function updateFavoriteButton(artistId) {
    const button = document.getElementById(`favorite-${artistId}`);
    if (button) {
        button.innerHTML = favorites.includes(artistId) 
            ? '<i class="fas fa-star"></i> Remove from Favorites' 
            : '<i class="far fa-star"></i> Add to Favorites';
    }
}

function shareArtist(artist) {
    if (navigator.share) {
        navigator.share({
            title: artist.name,
            text: `Check out ${artist.name} on Groupie Tracker!`,
            url: window.location.href
        }).then(() => console.log('Successful share'))
        .catch((error) => console.log('Error sharing:', error));
    } else {
        alert('Web Share API is not supported in your browser. You can copy the URL to share.');
    }
}

function displayMap(locations) {
    if (map) {
        map.remove();
    }

    map = new mapboxgl.Map({
        container: 'map',
        style: 'mapbox://styles/mapbox/streets-v11',
        center: [0, 0],
        zoom: 1
    });

    const bounds = new mapboxgl.LngLatBounds();

    locations.forEach(location => {
        if (!location.lon || !location.lat) {
            console.error('Invalid coordinates for location:', location);
            return;
        }

        const el = document.createElement('div');
        el.className = 'custom-marker';
        el.innerHTML = '<i class="fas fa-map-marker-alt"></i>';
        el.style.color = '#FF0000';
        el.style.fontSize = '24px';

        el.addEventListener('click', () => {
            const popup = document.createElement('div');
            popup.className = 'custom-popup';
            popup.innerHTML = `
                <h3 style="color: black; margin: 0; padding: 5px 0;">${location.address}</h3>
                <button class="popup-close">&times;</button>
            `;
            popup.style.position = 'absolute';
            popup.style.backgroundColor = 'white';
            popup.style.padding = '10px';
            popup.style.borderRadius = '4px';
            popup.style.boxShadow = '0 2px 4px rgba(0,0,0,0.2)';
            popup.style.zIndex = '1000';
            popup.style.minWidth = '100px';
            popup.style.textAlign = 'center';
        
            const existingPopup = document.querySelector('.custom-popup');
            if (existingPopup) {
                existingPopup.remove();
            }
        
            map.getCanvasContainer().appendChild(popup);
        
            const point = map.project([location.lon, location.lat]);
            popup.style.left = `${point.x - (popup.offsetWidth / 2)}px`;
            popup.style.top = `${point.y - popup.offsetHeight - 10}px`;

            const closeButton = popup.querySelector('.popup-close');
            closeButton.addEventListener('click', (e) => {
                e.stopPropagation();
                popup.remove();
            });
        });

        new mapboxgl.Marker({ element: el })
            .setLngLat([location.lon, location.lat])
            .addTo(map);

        bounds.extend([location.lon, location.lat]);
    });

    if (!bounds.isEmpty()) {
        map.fitBounds(bounds, { padding: 50 });
    }
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
                <button id="favorite-${details.artist.id}" onclick="toggleFavorite(${details.artist.id})">
                    ${favorites.includes(details.artist.id) ? '<i class="fas fa-star"></i> Remove from Favorites' : '<i class="far fa-star"></i> Add to Favorites'}
                </button>
                <button onclick="shareArtist(${JSON.stringify(details.artist)})">
                    <i class="fas fa-share-alt"></i> Share
                </button>
            </div>
        </div>
    `;

    displayMap(details.locations);
}

// Load artist details when the page loads
window.addEventListener('load', () => {
    const artistId = getArtistId();
    if (artistId) {
        showLoading();
        fetch(`/api/artist/${artistId}`)
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            })
            .then(data => {
                displayArtistDetails(data);
                hideLoading();
            })
            .catch(error => {
                console.error('Error:', error);
                showError('An error occurred while fetching artist details. Please try again later.');
                hideLoading();
            });
    }
});