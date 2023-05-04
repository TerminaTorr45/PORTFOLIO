mapboxgl.accessToken = "pk.eyJ1IjoibGF0aW1lcyIsImEiOiJjajhvcXRraGUwNnlwMzNyczR3cTBsaWh1In0.0cPKLwe2A0ET4P5CtWSiLQ"
var map = new mapboxgl.Map({
    container: 'map',
    style: 'mapbox://styles/mapbox/streets-v11',
    center: [-98.5795, 39.8283],
    zoom: 3
});

map.addControl(new mapboxgl.NavigationControl());

var markers = [];

function addMarker(lng, lat, popupText) {
    var marker = new mapboxgl.Marker()
        .setLngLat([lng, lat])
        .setPopup(new mapboxgl.Popup({ offset: 25 }) // add popups
            .setHTML('<h3>' + popupText + '</h3>'))
        .addTo(map);
    markers.push(marker);
}

function removeMarkers() {
    markers.forEach(function(marker) {
        marker.remove();
    });
    markers = [];
}

function getLocations() {
    removeMarkers();
    fetch('/locations')
        .then(function(response) {
            return response.json();
        })
        .then(function(locations) {
            locations.forEach(function(location) {
                addMarker(location.longitude, location.latitude, location.location);
            });
        })
        .catch(function(error) {
            console.error(error);
        });
}

function generateGeojson() {
    fetch('/locations')
        .then(function(response) {
            return response.json();
        })
        .then(function(locations) {
            var geojson = {
                type: 'FeatureCollection',
                features: []
            };

            locations.forEach(function(location) {
                var feature = {
                    type: 'Feature',
                    geometry: {
                        type: 'Point',
                        coordinates: [location.longitude, location.latitude]
                    },
                    properties: {
                        title: location.location
                    }
                };
                geojson.features.push(feature);
            });

            var data = JSON.stringify(geojson);
            document.getElementById('data-input').value = data;
        })
        .catch(function(error) {
            console.error(error);
        });
}

document.getElementById('search-form').addEventListener('submit', function(event) {
    event.preventDefault();
    var query = document.getElementById('search-input').value;
    fetch('/search?query=' + query)
        .then(function(response) {
            return response.json();
        })
        .then(function(locations) {
            removeMarkers();
            console.log(locations)
            locations.forEach(function(location) {
                addMarker(location.longitude, location.latitude, location.location);
            });
        })
        .catch(function(error) {
            console.error(error);
        });
});

getLocations();
generateGeojson();
